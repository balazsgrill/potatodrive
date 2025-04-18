package gpfs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/spf13/afero"
)

func aError(context string, err error) error {
	return errors.New("gpfs: " + context + ": " + err.Error())
}

type fs struct {
	*gphotos.Client
	httpclient *http.Client
	albumcache map[string]*album
}

type item interface {
	afero.File
	os.FileInfo
}

// Chown implements afero.Fs.
func (f *fs) Chown(name string, uid int, gid int) error {
	return aError("Chown", ErrReadOnlyFs)
}

var ErrReadOnlyFs = errors.New("read-only file system")

func NewFs(client *http.Client) (afero.Fs, error) {
	gclient, err := gphotos.NewClient(client)
	if err != nil {
		return nil, aError("NewFs", err)
	}
	return &fs{
		Client:     gclient,
		httpclient: client,
	}, nil
}

func (f *fs) Name() string {
	return "gpfs"
}

func (f *fs) Create(name string) (afero.File, error) {
	return nil, aError("Create", ErrReadOnlyFs)
}

func (f *fs) Mkdir(name string, perm os.FileMode) error {
	return aError("Mkdir", ErrReadOnlyFs)
}

func (f *fs) MkdirAll(path string, perm os.FileMode) error {
	return aError("MkdirAll", ErrReadOnlyFs)
}

func (f *fs) ensureCachedAlbums() error {
	if f.albumcache != nil {
		return nil
	}
	albums, err := f.Client.Albums.List(context.Background(), false)
	if err != nil {
		return aError("ensureCachedAlbums", err)
	}
	f.albumcache = make(map[string]*album, len(albums))
	for _, album_ := range albums {
		f.albumcache[album_.Title] = &album{
			httpclient:        f.httpclient,
			MediaItemsService: f.Client.MediaItems,
			Album:             &album_,
		}
	}
	return nil
}

func (f *fs) openAlbum(name string) (*album, error) {
	err := f.ensureCachedAlbums()
	if err != nil {
		return nil, aError("openAlbum", err)
	}
	if album, ok := f.albumcache[name]; ok {
		return album, nil
	}

	return nil, aError(fmt.Sprintf("openAlbum(%s)", name), os.ErrNotExist)
}

func (f *fs) get(name string) (item, error) {
	sections := strings.Split(name, "/")
	var file item
	file = &rootDir{fs: f}
	filteredSections := []string{}
	for _, section := range sections {
		if section != "" && section != "." {
			filteredSections = append(filteredSections, section)
		}
	}
	if len(filteredSections) == 0 {
		return file, nil
	}
	if len(filteredSections) == 1 {
		album, err := f.openAlbum(filteredSections[0])
		if err != nil {
			return nil, aError("Open", err)
		}
		return album, nil
	}

	album, err := f.openAlbum(filteredSections[0])
	if err != nil {
		return nil, aError("Open", err)
	}
	file, err = album.findMediaItem(filteredSections[1])
	if err != nil {
		return nil, aError("Open", err)
	}

	return file, nil
}

func (f *fs) Open(name string) (afero.File, error) {
	return f.get(name)
}

func (f *fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, aError("OpenFile", ErrReadOnlyFs)
	}
	return f.Open(name)
}

func (f *fs) Remove(name string) error {
	return aError("Remove", ErrReadOnlyFs)
}

func (f *fs) RemoveAll(path string) error {
	return aError("RemoveAll", ErrReadOnlyFs)
}

func (f *fs) Rename(oldname, newname string) error {
	return aError("Rename", ErrReadOnlyFs)
}

func (f *fs) Stat(name string) (os.FileInfo, error) {
	return f.get(name)
}

func (f *fs) Chmod(name string, mode os.FileMode) error {
	return aError("Chmod", ErrReadOnlyFs)
}

func (f *fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return aError("Chtimes", ErrReadOnlyFs)
}
