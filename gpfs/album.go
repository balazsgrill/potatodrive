package gpfs

import (
	"context"
	"maps"
	"net/http"
	"os"
	"slices"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/spf13/afero"
)

type album struct {
	httpclient *http.Client
	gphotos.MediaItemsService
	*albums.Album
	mediacache map[string]*media
}

var _ = (afero.File)((*album)(nil))
var _ = (os.FileInfo)((*album)(nil))

func enuserMediacached(a *album) error {
	if a.mediacache != nil {
		return nil
	}
	if a.mediacache == nil {
		a.mediacache = make(map[string]*media)
	}
	items, err := listMediaItemsByAlbum(a)
	if err != nil {
		return err
	}
	for _, item := range items {
		a.mediacache[item.Filename] = &media{
			httpclient: a.httpclient,
			MediaItem:  item,
		}
	}
	return nil
}

func listMediaItemsByAlbum(a *album) ([]*media_items.MediaItem, error) {
	return a.MediaItemsService.ListByAlbum(context.Background(), a.ID)
}

func (a *album) findMediaItem(name string) (*media, error) {
	err := enuserMediacached(a)
	if err != nil {
		return nil, err
	}
	item, ok := a.mediacache[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return item, nil
}

// ReadAt implements afero.File.
func (a *album) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, os.ErrInvalid
}

// Readdirnames implements afero.File.
func (a *album) Readdirnames(n int) ([]string, error) {
	err := enuserMediacached(a)
	if err != nil {
		return nil, err
	}
	return slices.Collect(maps.Keys(a.mediacache)), nil
}

// Sync implements afero.File.
func (a *album) Sync() error {
	return nil
}

// Truncate implements afero.File.
func (a *album) Truncate(size int64) error {
	return os.ErrInvalid
}

// WriteAt implements afero.File.
func (a *album) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, os.ErrInvalid
}

// WriteString implements afero.File.
func (a *album) WriteString(s string) (ret int, err error) {
	return 0, os.ErrInvalid
}

func (a *album) Name() string {
	return a.Title
}

func (a *album) Readdir(count int) ([]os.FileInfo, error) {
	err := enuserMediacached(a)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(a.mediacache))
	for _, item := range a.mediacache {
		infos = append(infos, item)
	}
	return infos, nil
}

func (a *album) Stat() (os.FileInfo, error) {
	return a, nil
}

func (a *album) Close() error {
	a.mediacache = nil
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
	return os.ModeDir | 0555
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
