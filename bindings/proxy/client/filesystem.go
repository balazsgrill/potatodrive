package client

import (
	"context"
	"io/fs"
	"time"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
	"github.com/spf13/afero"
)

type filesystemClient struct {
	client proxy.Filesystem
}

func New(fs proxy.Filesystem) afero.Fs {
	return &filesystemClient{
		client: fs,
	}
}

// Chmod implements afero.Fs.
func (f *filesystemClient) Chmod(name string, mode fs.FileMode) error {
	return eurap("chmod", f.client.Chmod(context.Background(), name, proxy.FileMode(mode)))
}

// Chown implements afero.Fs.
func (f *filesystemClient) Chown(name string, uid int, gid int) error {
	return eurap("chown", f.client.Chown(context.Background(), name, int32(uid), int32(gid)))
}

// Chtimes implements afero.Fs.
func (f *filesystemClient) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return eurap("chtimes", f.client.Chtimes(context.Background(), name, proxy.Timestamp(atime.UnixMicro()), proxy.Timestamp(mtime.UnixMicro())))
}

// Create implements afero.Fs.
func (f *filesystemClient) Create(name string) (afero.File, error) {
	h, err := f.client.Create(context.Background(), name)
	if err != nil {
		return nil, eurap("create", err)
	}
	return toFile(f, h), nil
}

// Mkdir implements afero.Fs.
func (f *filesystemClient) Mkdir(name string, perm fs.FileMode) error {
	return eurap("mkdir", f.client.Mkdir(context.Background(), name, proxy.FileMode(perm)))
}

// MkdirAll implements afero.Fs.
func (f *filesystemClient) MkdirAll(path string, perm fs.FileMode) error {
	return eurap("mkdirall", f.client.MkdirAll(context.Background(), path, proxy.FileMode(perm)))
}

// Name implements afero.Fs.
func (f *filesystemClient) Name() string {
	name, err := f.client.Name(context.Background())
	if err != nil {
		return ""
	}
	return name
}

// Open implements afero.Fs.
func (f *filesystemClient) Open(name string) (afero.File, error) {
	h, err := f.client.Open(context.Background(), name)
	if err != nil {
		return nil, eurap("open", err)
	}
	return toFile(f, h), nil
}

// OpenFile implements afero.Fs.
func (f *filesystemClient) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	h, err := f.client.OpenFile(context.Background(), name, int32(flag), proxy.FileMode(perm))
	if err != nil {
		return nil, eurap("openfile", err)
	}
	return toFile(f, h), nil
}

// Remove implements afero.Fs.
func (f *filesystemClient) Remove(name string) error {
	return eurap("remove", f.client.Remove(context.Background(), name))
}

// RemoveAll implements afero.Fs.
func (f *filesystemClient) RemoveAll(path string) error {
	return eurap("removeall", f.client.RemoveAll(context.Background(), path))
}

// Rename implements afero.Fs.
func (f *filesystemClient) Rename(oldname string, newname string) error {
	return eurap("rename", f.client.Rename(context.Background(), oldname, newname))
}

// Stat implements afero.Fs.
func (f *filesystemClient) Stat(name string) (fs.FileInfo, error) {
	fi, err := f.client.Stat(context.Background(), name)
	return fi, eurap("stat", err)
}
