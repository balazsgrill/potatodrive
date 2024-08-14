package filesystem

import (
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"golang.org/x/sys/windows"
)

func toLongPath(localpath string) string {
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

func toShortPath(localpath string) string {
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

func (instance *VirtualizationInstance) path_localToRemote(path string) string {
	p := toLongPath(path)
	p = strings.TrimPrefix(p, instance.shortprefix)
	p = strings.TrimPrefix(p, instance.longprefix)
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimPrefix(p, "/")
	return p
}

func (instance *VirtualizationInstance) path_remoteToLocal(path string) string {
	p := strings.TrimPrefix(path, "/")
	p = strings.ReplaceAll(p, "/", "\\")
	return filepath.Join(instance.rootPath, "\\", p)
}

func (instance *VirtualizationInstance) path_getNameRemote(path string) string {
	p := strings.TrimPrefix(path, "/")
	return filepath.Base(p)
}

func (instance *VirtualizationInstance) path_getNameLocal(path string) string {
	return filepath.Base(strings.ReplaceAll(path, "\\", "/"))
}

func (instance *VirtualizationInstance) path_hashFile(remotepath string) string {
	fname := filepath.Base(remotepath)
	dir := filepath.Dir(remotepath)
	return dir + "/.md5_" + fname
}
