package filesystem

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/balazsgrill/potatodrive/win/cfapi"
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

func (instance *VirtualizationInstance) syncLocalToRemote() error {
	return filepath.Walk(instance.rootPath, func(localpath string, localinfo fs.FileInfo, err error) error {
		instance.Logger.Debug().Msgf("Syncing local file '%s'", localpath)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}

		path := instance.path_localToRemote(localpath)
		if localinfo.IsDir() {
			return instance.fs.MkdirAll(path, 0777)
		}
		if strings.HasPrefix(path, ".") {
			return nil
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
			localisnewer := os.IsNotExist(err) || (localinfo.ModTime().UTC().Unix() > remoteinfo.ModTime().UTC().Unix())

			if localisnewer {
				instance.Logger.Info().Msgf("Updating remote file '%s'", path)
				err = instance.streamLocalToRemote(path)
				if err != nil {
					return err
				}
			}
			// mark file as in-sync
			return instance.setInSync(localpath)

		}

		if deleted {
			return os.Remove(localpath)
		}
		return nil
	})
}
