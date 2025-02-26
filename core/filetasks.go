package core

import "github.com/balazsgrill/potatodrive/core/tasks"

type FileStateCallbacks interface {
	FileSynchronizing(path string)
	FileDone(path string)
	FileRemoved(path string)
	FileError(path string, err error)
	FileDownloading(path string, progress int)
	FileUploading(path string, progress int)
}

type fileStatesAsTasks struct {
	listener tasks.TaskStateListener
}

func (f *fileStatesAsTasks) FileSynchronizing(path string) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:    id,
			Name:  path,
			State: "Synchronizing",
		})
	}
}

func (f *fileStatesAsTasks) FileDone(path string) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:       id,
			Name:     path,
			State:    "Done",
			Progress: 100,
		})
	}

}

func (f *fileStatesAsTasks) FileRemoved(path string) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:       id,
			Name:     path,
			State:    "Removed",
			Progress: 100,
		})
	}
}

func (f *fileStatesAsTasks) FileError(path string, e error) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:    id,
			Name:  path,
			State: "Error",
			Error: e,
		})
	}
}

func (f *fileStatesAsTasks) FileDownloading(path string, progress int) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:       id,
			Name:     path,
			State:    "Downloading",
			Progress: progress,
		})
	}
}

func (f *fileStatesAsTasks) FileUploading(path string, progress int) {
	id, err := GetFileID(path)
	if err != nil {
		f.listener(tasks.TaskState{
			ID:       id,
			Name:     path,
			State:    "Uploading",
			Progress: progress,
		})
	}
}

func AsCallbacks(listener tasks.TaskStateListener) FileStateCallbacks {
	return &fileStatesAsTasks{listener}
}
