package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/balazsgrill/potatodrive/bindings/utils"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/core/cfapi"
)

func (instance *VirtualizationInstance) syncRemoteToLocal() error {
	return utils.Walk(instance.fs, "", func(path string, remoteinfo fs.FileInfo, err error) error {
		instance.Logger.Debug().Msgf("Syncing remote file '%s'", path)
		if os.IsNotExist(err) {
			instance.Logger.Error().Msgf("Not exists: %v", err)
			return nil
		}
		if err != nil {
			instance.Logger.Err(err).Send()
			return fmt.Errorf("syncRemoteToLocal.1 %w", err)
		}

		filename := instance.path_getNameRemote(path)
		if strings.HasPrefix(filename, ".") && len(filename) > 1 {
			// do not skip "."
			if remoteinfo.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		localpath := instance.path_remoteToLocal(path)
		placeholderstate, err := getPlaceholderState(localpath)
		instance.Logger.Debug().Msgf("Placeholder state for '%s' is %x", localpath, placeholderstate)
		if os.IsNotExist(err) {
			if remoteinfo.IsDir() {
				// local dir does not exist, create it
				instance.Logger.Debug().Msgf("Creating local dir '%s'", localpath)
				return os.MkdirAll(localpath, 0777)
			} else {
				localdir := filepath.Dir(localpath)
				// placeholder does not exists, create it
				placeholder := getPlaceholder(remoteinfo)
				var EntriesProcessed uint32
				hr := cfapi.CfCreatePlaceholders(core.GetPointer(localdir), &placeholder, 1, cfapi.CF_CREATE_FLAG_NONE, &EntriesProcessed)
				if hr != 0 {
					return core.ErrorByCodeWithContext("syncRemoteToLocal:CfCreatePlaceholders", hr)
				}
				if EntriesProcessed != 1 {
					return fmt.Errorf("syncRemoteToLocal: unexpected number of entries processed: %d", EntriesProcessed)
				}
				// done here, return
				return nil
			}
		}
		if err != nil {
			return fmt.Errorf("syncRemoteToLocal.2 %w", err)
		}

		insync := (placeholderstate & cfapi.CF_PLACEHOLDER_STATE_IN_SYNC) != 0
		isaplaceholder := (placeholderstate & cfapi.CF_PLACEHOLDER_STATE_PLACEHOLDER) != 0

		// check if remote is newer
		localinfo, _ := os.Stat(localpath)
		if localinfo.ModTime().UTC().Unix() < remoteinfo.ModTime().UTC().Unix() {
			return instance.markFileAsDirty(path, localpath, remoteinfo, localinfo, insync, isaplaceholder)
		}

		return nil
	})
}

func (instance *VirtualizationInstance) markFileAsDirty(path string, localpath string, remoteinfo fs.FileInfo, localinfo os.FileInfo, insync bool, isaplaceholder bool) error {
	instance.Logger.Debug().Msgf("Updating local file '%s'", path)

	var handle syscall.Handle
	hr := cfapi.CfOpenFileWithOplock(core.GetPointer(localpath), cfapi.CF_OPEN_FILE_FLAG_WRITE_ACCESS|cfapi.CF_OPEN_FILE_FLAG_EXCLUSIVE, &handle)
	if hr != 0 {
		return core.ErrorByCodeWithContext("syncRemoteToLocal:CfOpenFileWithOplock", hr)
	}
	defer cfapi.CfCloseHandle(handle)
	placeholder := getPlaceholder(remoteinfo)

	if !isaplaceholder {
		// setting in-sync state only works if it's a placeholder
		instance.Logger.Info().Msgf("Converting to placeholder '%s'", path)
		hr = cfapi.CfConvertToPlaceholder(handle, placeholder.FileIdentity, placeholder.FileIdentityLength, cfapi.CF_CONVERT_FLAG_NONE, 0, 0)
		if hr != 0 {
			return core.ErrorByCodeWithContext("syncRemoteToLocal:CfConvertToPlaceholder", hr)
		}
	}
	if !insync {
		// updating a placeholder only works if it is marked as in-sync
		hr = cfapi.CfSetInSyncState(handle, cfapi.CF_IN_SYNC_STATE_IN_SYNC, cfapi.CF_SET_IN_SYNC_FLAG_NONE, nil)
		if hr != 0 {
			return core.ErrorByCodeWithContext("syncRemoteToLocal:CfSetInSyncState", hr)
		}
	}
	var fileRange cfapi.CF_FILE_RANGE
	fileRange.StartingOffset = 0
	fileRange.Length = localinfo.Size()
	hr = cfapi.CfUpdatePlaceholder(handle, &placeholder.FsMetadata, placeholder.FileIdentity, placeholder.FileIdentityLength, &fileRange, 1, cfapi.CF_UPDATE_FLAG_CLEAR_IN_SYNC|cfapi.CF_UPDATE_FLAG_DEHYDRATE, nil, 0)
	if hr != 0 {
		return core.ErrorByCodeWithContext("syncRemoteToLocal:CfUpdatePlaceholder", hr)
	}
	instance.FileSynchronizing(localpath)
	return nil
}
