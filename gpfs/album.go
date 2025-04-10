package gpfs

import (
	"context"
	"net/http"
	"os"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
)

type album struct {
	httpclient *http.Client
	gphotos.MediaItemsService
	*albums.Album
}

func (a *album) Name() string {
	return a.Title
}

func (a *album) Readdir(count int) ([]os.FileInfo, error) {
	items, err := a.MediaItemsService.ListByAlbum(context.Background(), a.ID)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]os.FileInfo, len(items))
	for i, item := range items {
		fileInfos[i] = &media{
			httpclient: a.httpclient,
			MediaItem:  item,
		}
	}
	return fileInfos, nil
}

func (a *album) Stat() (os.FileInfo, error) {
	return a, nil
}

func (a *album) Close() error {
	// No resources to close for an album
	return nil
}

func (a *album) Read(p []byte) (n int, err error) {
	// Albums are not readable as files
	return 0, os.ErrInvalid
}

func (a *album) Write(p []byte) (n int, err error) {
	// Albums are not writable as files
	return 0, os.ErrInvalid
}

func (a *album) Seek(offset int64, whence int) (int64, error) {
	// Albums are not seekable
	return 0, os.ErrInvalid
}

func (a *album) Mode() os.FileMode {
	return os.ModeDir | 0755
}

func (a *album) ModTime() time.Time {
	return time.Time{}
}

func (a *album) IsDir() bool {
	return true
}

func (a *album) Sys() interface{} {
	return nil
}

func (a *album) Size() int64 {
	// Albums don't have a size, returning 0
	return 0
}
