package utils

import (
	"os"
	"sync"
	"time"

	"github.com/spf13/afero"
)

type ConnectingFs struct {
	Connect func(onDisconnect func(error)) (afero.Fs, error)

	lock      sync.Mutex
	currentFs afero.Fs
}

var _ afero.Fs = (*ConnectingFs)(nil)

func (cfs *ConnectingFs) Chmod(name string, mode os.FileMode) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Chmod(name, mode)
	})
}

func (cfs *ConnectingFs) MkdirAll(path string, perm os.FileMode) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.MkdirAll(path, perm)
	})
}

func (cfs *ConnectingFs) Stat(name string) (os.FileInfo, error) {
	var fileInfo os.FileInfo
	err := cfs.withFs(func(fs afero.Fs) error {
		var err error
		fileInfo, err = fs.Stat(name)
		return err
	})
	return fileInfo, err
}
func (cfs *ConnectingFs) Rename(oldname, newname string) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Rename(oldname, newname)
	})
}
func (cfs *ConnectingFs) RemoveAll(path string) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.RemoveAll(path)
	})
}
func (cfs *ConnectingFs) Remove(name string) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Remove(name)
	})
}
func (cfs *ConnectingFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	var file afero.File
	err := cfs.withFs(func(fs afero.Fs) error {
		var err error
		file, err = fs.OpenFile(name, flag, perm)
		return err
	})
	return file, err
}
func (cfs *ConnectingFs) Open(name string) (afero.File, error) {
	var file afero.File
	err := cfs.withFs(func(fs afero.Fs) error {
		var err error
		file, err = fs.Open(name)
		return err
	})
	return file, err
}
func (cfs *ConnectingFs) Name() string {
	return "ConnectingFs"
}
func (cfs *ConnectingFs) Mkdir(name string, perm os.FileMode) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Mkdir(name, perm)
	})
}
func (cfs *ConnectingFs) Create(name string) (afero.File, error) {
	var file afero.File
	err := cfs.withFs(func(fs afero.Fs) error {
		var err error
		file, err = fs.Create(name)
		return err
	})
	return file, err
}
func (cfs *ConnectingFs) Chtimes(name string, atime, mtime time.Time) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Chtimes(name, atime, mtime)
	})
}
func (cfs *ConnectingFs) Chown(name string, uid, gid int) error {
	return cfs.withFs(func(fs afero.Fs) error {
		return fs.Chown(name, uid, gid)
	})
}
func (cfs *ConnectingFs) withFs(f func(fs afero.Fs) error) error {
	cfs.lock.Lock()
	defer cfs.lock.Unlock()
	if cfs.currentFs == nil {
		fs, err := cfs.Connect(func(error) {
			cfs.lock.Lock()
			defer cfs.lock.Unlock()
			cfs.currentFs = nil
		})
		if err != nil {
			return err
		}
		cfs.currentFs = fs
	}
	return f(cfs.currentFs)
}
