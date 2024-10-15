package proxy

import (
	"io/fs"
	"time"
)

func (f *FileInfo) Name() string {
	return f.Fname
}

func (f *FileInfo) Size() int64 {
	return f.Fsize
}
func (f *FileInfo) Mode() fs.FileMode {
	return fs.FileMode(f.Fmode)
}

func (f *FileInfo) ModTime() time.Time {
	return time.UnixMicro(int64(f.Ftime))
}
func (f *FileInfo) Sys() interface{} {
	return nil
}
func (f *FileInfo) IsDir() bool {
	return f.FisDir
}

var _ fs.FileInfo = (*FileInfo)(nil)
