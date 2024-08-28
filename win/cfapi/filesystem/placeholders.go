package filesystem

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/balazsgrill/potatodrive/win/cfapi"
	"github.com/spf13/afero"
)

func (instance *VirtualizationInstance) fetchPlaceholders(info *cfapi.CF_CALLBACK_INFO, data *cfapi.CF_CALLBACK_PARAMETERS_FetchPlaceholders) uintptr {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	name := getFileNameFromIdentity(info)
	instance.Logger.Debug().Msgf("Fetch placeholders: %s / %s", win.GetString(info.NormalizedPath), name)
	remotepath := instance.path_localToRemote(win.GetString(info.NormalizedPath))
	files, err := afero.ReadDir(instance.fs, remotepath)
	if err != nil {
		instance.Logger.Debug().Msgf("Error reading directory %s: %s", remotepath, err)
		return uintptr(syscall.EIO)
	}
	transfer := cfapi.CF_OPERATION_PARAMETERS_TransferPlaceholders{}
	transfer.ParamSize = uint32(unsafe.Sizeof(transfer))
	transfer.CompletionStatus = 0 //success
	transfer.Flags = cfapi.CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_NONE
	count := 0
	placeholders := make([]cfapi.CF_PLACEHOLDER_CREATE_INFO, len(files))
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			placeholders[count] = getPlaceholder(f)

			count += 1
		}
	}

	for i := 0; i < count; i++ {
		instance.Logger.Debug().Msgf("Sending %d", i)
		var placeholder cfapi.CF_PLACEHOLDER_CREATE_INFO
		transfer.PlaceholderTotalCount = int64(count)
		transfer.EntriesProcessed = 0
		transfer.PlaceholderCount = 1
		placeholder = placeholders[i]
		transfer.PlaceholderArray = &placeholder
		if i == count-1 {
			transfer.Flags = cfapi.CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_DISABLE_ON_DEMAND_POPULATION
		}

		hr := instance.transferPlaceholders(info, &transfer)

		if hr != 0 {
			instance.Logger.Debug().Msgf("Error transferring placeholders: %s", win.ErrorByCode(hr))
			return hr
		}
	}

	if count == 0 {
		// send empty placeholder array
		transfer.PlaceholderTotalCount = 0
		transfer.EntriesProcessed = 0
		transfer.PlaceholderCount = 0
		transfer.PlaceholderArray = nil
		transfer.Flags = cfapi.CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_DISABLE_ON_DEMAND_POPULATION
		hr := instance.transferPlaceholders(info, &transfer)

		if hr != 0 {
			instance.Logger.Debug().Msgf("Error transferring placeholders: %s", win.ErrorByCode(hr))
			return hr
		}
	}
	instance.Logger.Debug().Msgf("Sent %d entries", count)
	return 0
}
