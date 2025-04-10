package gpfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
)

type media struct {
	httpclient *http.Client
	*media_items.MediaItem
}

func (m *media) getUrl() string {
	return fmt.Sprintf("$s=w%d-h%d", m.MediaItem.BaseURL, m.MediaItem.MediaMetadata.Width, m.MediaItem.MediaMetadata.Height)
}

// Size implements fs.FileInfo.
func (m *media) Size() int64 {
	resp, err := m.httpclient.Head(m.getUrl())
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.ContentLength > 0 {
		return resp.ContentLength
	}

	return 0
}

func (m *media) Name() string {
	return m.Filename
}

func (m *media) Readdir(count int) ([]os.FileInfo, error) {
	return nil, io.EOF
}

func (m *media) Stat() (os.FileInfo, error) {
	return m, nil
}

func (m *media) Close() error {
	return nil
}

func (m *media) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *media) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}

func (m *media) Seek(offset int64, whence int) (int64, error) {
	return 0, io.EOF
}

func (m *media) Mode() os.FileMode {
	return 0444
}

func (m *media) ModTime() time.Time {
	return time.Now()
}

func (m *media) IsDir() bool {
	return false
}

func (m *media) Sys() interface{} {
	return nil
}
