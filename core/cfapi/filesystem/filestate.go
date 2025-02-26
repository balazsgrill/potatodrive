package filesystem

import "github.com/balazsgrill/potatodrive/core"

var _ core.FileStateCallbacks = &VirtualizationInstance{}

func (vi *VirtualizationInstance) FileSynchronizing(path string) {
	if vi.callbacks != nil {
		vi.callbacks.FileSynchronizing(path)
	}
}

func (vi *VirtualizationInstance) FileDone(path string) {
	if vi.callbacks != nil {
		vi.callbacks.FileDone(path)
	}
}

func (vi *VirtualizationInstance) FileError(path string, err error) {
	if vi.callbacks != nil {
		vi.callbacks.FileError(path, err)
	}
}

func (vi *VirtualizationInstance) FileDownloading(path string, progress int) {
	if vi.callbacks != nil {
		vi.callbacks.FileDownloading(path, progress)
	}
}

func (vi *VirtualizationInstance) FileUploading(path string, progress int) {
	if vi.callbacks != nil {
		vi.callbacks.FileUploading(path, progress)
	}
}

func (vi *VirtualizationInstance) FileRemoved(path string) {
	if vi.callbacks != nil {
		vi.callbacks.FileRemoved(path)
	}
}
