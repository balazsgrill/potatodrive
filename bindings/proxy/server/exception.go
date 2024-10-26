package server

import (
	"io"
	"os"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
)

func ewrap(err error) error {
	if err == nil {
		return nil
	}
	return &proxy.FilesystemException{
		Message:     err.Error(),
		Isnotexists: os.IsNotExist(err),
		EOF:         err == io.EOF,
	}
}
