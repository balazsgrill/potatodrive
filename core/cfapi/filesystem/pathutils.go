package filesystem

import (
	"path/filepath"
	"strings"

	"github.com/balazsgrill/potatodrive/core"
)

func (instance *VirtualizationInstance) path_localToRemote(path string) string {
	p := core.ToLongPath(path)
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
