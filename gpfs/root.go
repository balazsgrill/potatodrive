package gpfs

import (
	"maps"
	"os"
	"slices"
	"time"

	"github.com/spf13/afero"
)

type rootDir struct {
	*fs
}

// Sync implements afero.File.
func (r *rootDir) Sync() error {
	return nil
}

// Truncate implements afero.File.
func (r *rootDir) Truncate(size int64) error {
	return ErrReadOnlyFs
}

// WriteAt implements afero.File.
func (r *rootDir) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, ErrReadOnlyFs
}

// WriteString implements afero.File.
func (r *rootDir) WriteString(s string) (ret int, err error) {
	return 0, ErrReadOnlyFs
}

var _ = (afero.File)((*rootDir)(nil))
var _ = (os.FileInfo)((*rootDir)(nil))

func (r *rootDir) Name() string {
	return "/"
}

func (r *rootDir) Size() int64 {
	return 0
}

func (r *rootDir) Mode() os.FileMode {
	return os.ModeDir | 0555
}

func (r *rootDir) ModTime() time.Time {
	return time.Time{}
}

func (r *rootDir) IsDir() bool {
	return true
}

func (r *rootDir) Sys() interface{} {
	return nil
}

func (r *rootDir) Read(p []byte) (n int, err error) {
	return 0, os.ErrInvalid
}

func (r *rootDir) Write(p []byte) (n int, err error) {
	return 0, os.ErrInvalid
}

func (r *rootDir) Close() error {
	return nil
}

func (r *rootDir) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, os.ErrInvalid
}

func (r *rootDir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (r *rootDir) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (r *rootDir) Stat() (os.FileInfo, error) {
	return r, nil
}

// Readdirnames returns the names of files in the directory.
func (r *rootDir) Readdirnames(n int) ([]string, error) {
	err := r.fs.ensureCachedAlbums()
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Keys(r.fs.albumcache)), nil
}
