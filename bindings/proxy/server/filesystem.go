package server

import (
	"context"
	"os"
	"time"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
	"github.com/spf13/afero"
)

type FilesystemServer struct {
	openfiles map[proxy.FileHandle]afero.File
	count     int32
	fs        afero.Fs
}

func New(fs afero.Fs) proxy.Filesystem {
	return &FilesystemServer{
		fs:        fs,
		openfiles: make(map[proxy.FileHandle]afero.File),
		count:     0,
	}
}

var _ proxy.Filesystem = (*FilesystemServer)(nil)

func (fs *FilesystemServer) Chown(ctx context.Context, name string, uid int32, gid int32) (_err error) {
	return fs.fs.Chown(name, int(uid), int(gid))
}

func (fs *FilesystemServer) Chtimes(ctx context.Context, name string, atime proxy.Timestamp, mtime proxy.Timestamp) (_err error) {
	return fs.fs.Chtimes(name, time.UnixMicro(int64(atime)), time.UnixMicro(int64(mtime)))
}

func (fs *FilesystemServer) Create(ctx context.Context, name string) (_r proxy.FileHandle, _err error) {
	file, err := fs.fs.Create(name)
	if err != nil {
		return 0, err
	}
	fs.count++
	handle := proxy.FileHandle(fs.count)
	fs.openfiles[handle] = file
	return handle, nil
}

func (fs *FilesystemServer) Mkdir(ctx context.Context, path string, perm proxy.FileMode) (_err error) {
	return fs.fs.Mkdir(path, os.FileMode(perm))
}

func (fs *FilesystemServer) MkdirAll(ctx context.Context, path string, perm proxy.FileMode) (_err error) {
	return fs.fs.MkdirAll(path, os.FileMode(perm))
}

func (fs *FilesystemServer) Name(ctx context.Context) (_r string, _err error) {
	return fs.fs.Name(), nil
}

func (fs *FilesystemServer) Open(ctx context.Context, name string) (_r proxy.FileHandle, _err error) {
	file, err := fs.fs.Open(name)
	if err != nil {
		return 0, err
	}
	fs.count++
	handle := proxy.FileHandle(fs.count)
	fs.openfiles[handle] = file
	return handle, nil
}

func (fs *FilesystemServer) OpenFile(ctx context.Context, name string, flag int32, perm proxy.FileMode) (_r proxy.FileHandle, _err error) {
	file, err := fs.fs.OpenFile(name, int(flag), os.FileMode(perm))
	if err != nil {
		return 0, err
	}
	fs.count++
	handle := proxy.FileHandle(fs.count)
	fs.openfiles[handle] = file
	return handle, nil
}

func (fs *FilesystemServer) Remove(ctx context.Context, name string) (_err error) {
	return fs.fs.Remove(name)
}

func (fs *FilesystemServer) RemoveAll(ctx context.Context, name string) (_err error) {
	return fs.fs.RemoveAll(name)
}

func (fs *FilesystemServer) Rename(ctx context.Context, oldname string, newname string) (_err error) {
	return fs.fs.Rename(oldname, newname)
}

func wrapFileInfo(file os.FileInfo) *proxy.FileInfo {
	return &proxy.FileInfo{
		Fname:  file.Name(),
		Fmode:  proxy.FileMode(file.Mode()),
		Fsize:  file.Size(),
		Ftime:  proxy.Timestamp(file.ModTime().UnixMicro()),
		FisDir: file.IsDir(),
	}
}

func (fs *FilesystemServer) Stat(ctx context.Context, name string) (_r *proxy.FileInfo, _err error) {
	file, err := fs.fs.Stat(name)
	if err != nil {
		return nil, err
	}
	return wrapFileInfo(file), nil
}

func (fs *FilesystemServer) Chmod(ctx context.Context, name string, mode proxy.FileMode) error {
	return fs.fs.Chmod(name, os.FileMode(mode))
}
