package gpfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/spf13/afero"
)

type media struct {
	httpclient *http.Client
	*media_items.MediaItem
	cache  []byte
	offset int64
}

var _ = (afero.File)((*media)(nil))
var _ = (os.FileInfo)((*media)(nil))

func (m *media) ensureCached() error {
	if m.cache != nil {
		return nil
	}
	resp, err := m.httpclient.Get(m.getUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	m.cache, err = io.ReadAll(resp.Body)
	return err
}

// ReadAt implements afero.File.
func (m *media) ReadAt(p []byte, off int64) (n int, err error) {
	if err := m.ensureCached(); err != nil {
		return 0, err
	}
	if off >= int64(len(m.cache)) {
		return 0, io.EOF
	}
	if off+int64(len(p)) > int64(len(m.cache)) {
		p = p[:int64(len(m.cache))-off]
	}
	copy(p, m.cache[off:])
	return 0, fmt.Errorf("cannot write to file: read-only file system")
}

// Readdirnames implements afero.File.
func (m *media) Readdirnames(n int) ([]string, error) {
	return nil, io.EOF
}

// Sync implements afero.File.
func (m *media) Sync() error {
	panic("unimplemented")
}

// Truncate implements afero.File.
func (m *media) Truncate(size int64) error {
	return fmt.Errorf("cannot truncate file: read-only file system")
}

// WriteAt implements afero.File.
func (m *media) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, fmt.Errorf("cannot write to file: read-only file system")
}

// WriteString implements afero.File.
func (m *media) WriteString(s string) (ret int, err error) {
	return 0, fmt.Errorf("cannot write to file: read-only file system")
}

func (m *media) getUrl() string {
	return fmt.Sprintf("%s=w%d-h%d", m.MediaItem.BaseURL, m.MediaItem.MediaMetadata.Width, m.MediaItem.MediaMetadata.Height)
}

// Size implements fs.FileInfo.
func (m *media) Size() int64 {
	if m.cache != nil {
		return int64(len(m.cache))
	}
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
	m.cache = nil
	return nil
}

func (m *media) Read(p []byte) (n int, err error) {
	if err := m.ensureCached(); err != nil {
		return 0, err
	}
	if len(m.cache) == 0 {
		return 0, io.EOF
	}

	// Ensure the offset is within bounds
	if m.offset >= int64(len(m.cache)) {
		return 0, io.EOF
	}

	// Calculate how much data can be read
	remaining := int64(len(m.cache)) - m.offset
	toRead := int64(len(p))
	if toRead > remaining {
		toRead = remaining
	}

	// Copy data from the cache to the buffer
	copy(p[:toRead], m.cache[m.offset:m.offset+toRead])

	// Update the offset
	m.offset += toRead

	return int(toRead), nil
}

func (m *media) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}

func (m *media) Seek(offset int64, whence int) (int64, error) {
	if err := m.ensureCached(); err != nil {
		return 0, err
	}

	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = int64(len(m.cache)) + offset
	case io.SeekEnd:
		newOffset = int64(len(m.cache)) + offset
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}

	if newOffset < 0 || newOffset > int64(len(m.cache)) {
		return 0, fmt.Errorf("invalid offset: %d", newOffset)
	}

	m.offset = newOffset
	return m.offset, nil
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
