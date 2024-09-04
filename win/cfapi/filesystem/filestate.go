package filesystem

import "github.com/balazsgrill/potatodrive/win"

func (instance *VirtualizationInstance) NotifyFileState(path string, state win.FileSyncStateEnum) {
	if instance.handler == nil {
		return
	}
	instance.handler(win.FileSyncState{
		Path:  path,
		State: state,
	})
}

func (instance *VirtualizationInstance) NotifyFileError(path string, err error) {
	if instance.handler == nil {
		return
	}
	instance.handler(win.FileSyncState{
		Path:      path,
		State:     win.FileSyncStateError,
		LastError: err,
	})
}
