package filesystem

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/balazsgrill/potatodrive/core/cfapi"
	"github.com/spf13/afero"
)

// isDeletedRemotely check whether file was deleted remotely
// if it was, it compares local hash with remote hash. Returns true only if the file has been deleted remotely and was not changed locally
func (instance *VirtualizationInstance) isDeletedRemotely(remotepath string, localpath string) (bool, error) {
	_, err := instance.fs.Stat(remotepath)
	if os.IsNotExist(err) {
		// chek if remote hash is known
		hash, err := instance.remoteCacheState.GetHash(remotepath)
		if err != nil {
			return false, err
		}
		if len(hash) > 0 {
			// on remote file existed before, upload only if hash is different
			localhash, err := instance.localHash(remotepath)
			if err != nil {
				return false, err
			}
			if localhash == nil {
				// local file does not exist, no need to upload
				// TODO is this a tombstone?
				return false, nil
			}
			if bytes.Equal(hash, localhash) {
				// hash is the same this file has been removed remotely, delete local file
				return true, nil
			}
		}

	}
	return false, nil
}

func (instance *VirtualizationInstance) syncLocalToRemote() ([]string, error) {
	uploads := []string{}
	return uploads, filepath.Walk(instance.rootPath, func(localpath string, localinfo fs.FileInfo, err error) error {
		instance.Logger.Debug().Msgf("Syncing local file '%s'", localpath)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}

		path := instance.path_localToRemote(localpath)
		if strings.HasPrefix(path, ".") {
			return filepath.SkipDir
		}
		if localinfo.IsDir() {
			if dir, err := afero.IsDir(instance.fs, path); dir {
				return err
			}
			return instance.fs.MkdirAll(path, 0777)
		}

		localstate, err := getPlaceholderState(localpath)
		if err != nil {
			return err
		}
		instance.Logger.Debug().Msgf("Local state %x", localstate)

		deleted, err := instance.isDeletedRemotely(path, localpath)
		if err != nil {
			return err
		}

		if ((localstate & cfapi.CF_PLACEHOLDER_STATE_IN_SYNC) == 0) && (!deleted) {
			// local file is a hydrated placeholder, but not in sync, upload it if local is newer

			remoteinfo, err := instance.fs.Stat(path)
			var localisnewer bool
			if os.IsNotExist(err) {
				localisnewer = true
			} else if err != nil {
				return fmt.Errorf("syncLocalToRemote.1 %w", err)
			} else if remoteinfo == nil {
				return fmt.Errorf("syncLocalToRemote.2 NPE")
			} else {
				localisnewer = (localinfo.ModTime().UTC().Unix() > remoteinfo.ModTime().UTC().Unix())
			}

			if localisnewer {
				uploads = append(uploads, path)
			}
			// mark file as in-sync
			return instance.setInSync(localpath)
		}

		if deleted {
			err := os.Remove(localpath)
			if err != nil {
				instance.FileError(localpath, err)
				return err
			} else {
				instance.FileRemoved(localpath)
			}
			return err
		}
		return nil
	})
}
