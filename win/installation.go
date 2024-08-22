package win

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	POTATOICO = "potato.ico"
	MUTEXNAME = "7122117b-3a2e-4c5d-bb21-2c9d0d2243bd"
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

func CheckAlreadyRunning() error {
	name, err := syscall.UTF16PtrFromString(MUTEXNAME)
	if err != nil {
		return err
	}
	_, err = windows.CreateMutex(nil, true, name)
	return err
}
