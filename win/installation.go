package win

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	POTATOICO = "potato.ico"
)

func InstalledFile(relativename string) string {
	exec, err := os.Executable()
	if strings.HasPrefix(exec, os.TempDir()) {
		// Detects development mode, where files are looked for in the working directory rather than along with the exe
		exec = "."
	}

	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Dir(exec), relativename)
}
