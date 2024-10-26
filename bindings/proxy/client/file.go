package client

import (
	"context"
	"io/fs"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
	"github.com/spf13/afero"
)

type file struct {
	fs     *filesystemClient
	handle proxy.FileHandle
}

// Close implements afero.File.
func (f *file) Close() error {
	return f.fs.client.Fclose(context.Background(), f.handle)
}

// Name implements afero.File.
func (f *file) Name() string {
	name, _ := f.fs.client.Fname(context.Background(), f.handle)
	return name
}

// Read implements afero.File.
func (f *file) Read(p []byte) (n int, err error) {
	data, err := f.fs.client.Fread(context.Background(), f.handle, int64(len(p)))
	if err != nil {
		return 0, err
	}
	return copy(p, data), nil
}

// ReadAt implements afero.File.
func (f *file) ReadAt(p []byte, off int64) (n int, err error) {
	data, err := f.fs.client.FreadAt(context.Background(), f.handle, int64(len(p)), off)
	return copy(p, data), eurap("readat", err)
}

// Readdir implements afero.File.
func (f *file) Readdir(count int) ([]fs.FileInfo, error) {
	dirs, err := f.fs.client.Freaddir(context.Background(), f.handle, int32(count))
	if err != nil {
		return nil, err
	}

	fileInfos := make([]fs.FileInfo, len(dirs))
	for i, dir := range dirs {
		fileInfos[i] = dir
	}

	return fileInfos, nil
}

// Readdirnames implements afero.File.
func (f *file) Readdirnames(n int) ([]string, error) {
	return f.fs.client.Freaddirnames(context.Background(), f.handle, int32(n))
}

// Seek implements afero.File.
func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.fs.client.Fseek(context.Background(), f.handle, offset, int32(whence))
}

// Stat implements afero.File.
func (f *file) Stat() (fs.FileInfo, error) {
	return f.fs.client.Fstat(context.Background(), f.handle)
}

// Sync implements afero.File.
func (f *file) Sync() error {
	return f.fs.client.Fsync(context.Background(), f.handle)
}

// Truncate implements afero.File.
func (f *file) Truncate(size int64) error {
	return f.fs.client.Ftruncate(context.Background(), f.handle, size)
}

// Write implements afero.File.
func (f *file) Write(p []byte) (n int, err error) {
	r, err := f.fs.client.Fwrite(context.Background(), f.handle, p)
	return int(r), err
}

// WriteAt implements afero.File.
func (f *file) WriteAt(p []byte, off int64) (n int, err error) {
	r, err := f.fs.client.FwriteAt(context.Background(), f.handle, p, off)
	return int(r), err
}

// WriteString implements afero.File.
func (f *file) WriteString(s string) (ret int, err error) {
	r, err := f.fs.client.FwriteString(context.Background(), f.handle, s)
	return int(r), err
}

func toFile(fs *filesystemClient, handle proxy.FileHandle) afero.File {
	return &file{
		fs:     fs,
		handle: handle,
	}
}
