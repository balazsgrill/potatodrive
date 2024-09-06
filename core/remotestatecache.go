package core

import (
	"path"

	"github.com/spf13/afero"
)

// RemoteStateCache is the contract for keeping track files seen by the client on the remote side.
type RemoteStateCache interface {
	UpdateHash(remotepath string, hash []byte) error
	GetHash(remotepath string) ([]byte, error)
}

type remoteHashFiles struct {
	fs afero.Fs
}

func HashFilesRemotely(fs afero.Fs) RemoteStateCache {
	return &remoteHashFiles{fs: fs}
}

// GetHash implements RemoteStateCache.
func (instance *remoteHashFiles) GetHash(remotepath string) ([]byte, error) {
	hashpath := instance.path_hashFile(remotepath)
	exists, err := afero.Exists(instance.fs, hashpath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return afero.ReadFile(instance.fs, hashpath)
}

// UpdateHash implements RemoteStateCache.
func (instance *remoteHashFiles) UpdateHash(remotepath string, hash []byte) error {
	return afero.WriteFile(instance.fs, instance.path_hashFile(remotepath), hash, 0666)
}

var _ RemoteStateCache = (*remoteHashFiles)(nil)

func (instance *remoteHashFiles) path_hashFile(remotepath string) string {
	fname := path.Base(remotepath)
	dir := path.Dir(remotepath)
	return dir + "/.md5_" + fname
}
