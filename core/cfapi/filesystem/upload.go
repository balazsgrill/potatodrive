package filesystem

import (
	"crypto/md5"
	"io"
	"os"
)

func (instance *VirtualizationInstance) streamLocalToRemote(filename string) error {
	localpath := instance.path_remoteToLocal(filename)
	file, err := os.Open(localpath)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, 1024*1024)
	targetfile, err := instance.fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer targetfile.Close()

	hash := md5.New()
	done := false
	for !done {
		instance.Logger.Debug().Msgf("reading")
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				return err
			}
		}
		_, err = hash.Write(data[:n])
		if err != nil {
			return err
		}
		instance.Logger.Debug().Msgf("Uploading %d bytes", n)
		n2, err := targetfile.Write(data[:n])
		if err != nil {
			return err
		}
		instance.Logger.Debug().Msgf("uploaded chunk %d", n2)
	}
	instance.Logger.Debug().Msg("Done uploading")

	return instance.remoteCacheState.UpdateHash(filename, hash.Sum(nil))
}
