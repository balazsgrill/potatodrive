package client

import (
	"io"
	"io/fs"
	"os"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
)

func eurap(op string, err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*proxy.FilesystemException); ok {
		if e.EOF {
			return io.EOF
		}
		if e.Isnotexists {
			return &fs.PathError{
				Op:   op,
				Path: e.Message,
				Err:  os.ErrNotExist,
			}
		}
	}
	return err
}
