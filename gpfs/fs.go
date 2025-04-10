package gpfs

import (
	"context"
	"errors"
	"os"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/spf13/afero"
)

type fs struct {
	*gphotos.Client
}

// Chown implements afero.Fs.
func (f *fs) Chown(name string, uid int, gid int) error {
	return ErrReadOnlyFs
}

var ErrReadOnlyFs = errors.New("read-only file system")

func NewFs(client *gphotos.Client) afero.Fs {
	return &fs{client}
}

func (f *fs) Name() string {
	return "gpfs"
}

func (f *fs) Create(name string) (afero.File, error) {
	return nil, ErrReadOnlyFs
}

func (f *fs) Mkdir(name string, perm os.FileMode) error {
	return ErrReadOnlyFs
}

func (f *fs) MkdirAll(path string, perm os.FileMode) error {
	return ErrReadOnlyFs
}

func (f *fs) Open(name string) (afero.File, error) {
	albums, err := f.Client.Albums.List(context.Background())
	if err != nil {
		return nil, err
	}

	for _, album := range albums {
		if album.Title == name {
			// Logic to represent the album as a file can be implemented here
			return nil, nil // Replace with actual file representation
		}
	}

	return nil, os.ErrNotExist
}

func (f *fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, ErrReadOnlyFs
	}
	return f.Open(name)
}

func (f *fs) Remove(name string) error {
	return ErrReadOnlyFs
}

func (f *fs) RemoveAll(path string) error {
	return ErrReadOnlyFs
}

func (f *fs) Rename(oldname, newname string) error {
	return ErrReadOnlyFs
}

func (f *fs) Stat(name string) (os.FileInfo, error) {
	// Implement logic to fetch metadata for media or albums
	return nil, os.ErrNotExist
}

func (f *fs) Chmod(name string, mode os.FileMode) error {
	return ErrReadOnlyFs
}

func (f *fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return ErrReadOnlyFs
}
