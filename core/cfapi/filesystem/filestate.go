package filesystem

import "github.com/balazsgrill/potatodrive/core"

func (instance *VirtualizationInstance) NotifyFileState(path string, state core.FileSyncStateEnum) {
	if instance.handler == nil {
		return
	}
	instance.handler(core.FileSyncState{
		Path:  path,
		State: state,
	})
}

func (instance *VirtualizationInstance) NotifyFileError(path string, err error) {
	if instance.handler == nil {
		return
	}
	instance.handler(core.FileSyncState{
		Path:      path,
		State:     core.FileSyncStateError,
		LastError: err,
	})
}
