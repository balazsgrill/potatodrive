package win

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"
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
		exec, err = os.Getwd()
	} else {
		exec = filepath.Dir(exec)
	}

	if err != nil {
		panic(err)
	}
	return ToShortPath(filepath.Join(exec, relativename))
}

func CheckAlreadyRunning() error {
	name, err := syscall.UTF16PtrFromString(MUTEXNAME)
	if err != nil {
		return err
	}
	_, err = windows.CreateMutex(nil, true, name)
	return err
}

func ToLongPath(localpath string) string {
	shortpathp, err := windows.UTF16FromString(localpath)
	if err != nil {
		log.Printf("Failed to convert path '%s' to UTF16: %v", localpath, err)
		return localpath
	}
	longpathp := make([]uint16, windows.MAX_PATH)
	_, err = windows.GetLongPathName(&shortpathp[0], &longpathp[0], uint32(len(longpathp)))
	if err != nil {
		log.Printf("Failed to convert path '%s' to long path: %v", localpath, err)
		return localpath
	}
	return windows.UTF16ToString(longpathp)
}

func ToShortPath(localpath string) string {
	shortpathp, err := windows.UTF16FromString(localpath)
	if err != nil {
		log.Printf("Failed to convert path '%s' to UTF16: %v", localpath, err)
		return localpath
	}
	longpathp := make([]uint16, windows.MAX_PATH)
	_, err = windows.GetShortPathName(&shortpathp[0], &longpathp[0], uint32(len(longpathp)))
	if err != nil {
		log.Printf("Failed to convert path '%s' to short path: %v", localpath, err)
		return localpath
	}
	return windows.UTF16ToString(longpathp)
}
