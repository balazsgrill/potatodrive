package filesystem

import (
	"crypto/md5"
	"hash"
	"io"
	"syscall"
	"unsafe"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/core/cfapi"
)

const BUFFER_SIZE int64 = 100 * 1024

func (instance *VirtualizationInstance) callback_getFilePath(info *cfapi.CF_CALLBACK_INFO) string {
	return core.GetString(info.VolumeDosName) + core.GetString(info.NormalizedPath)
}

func (instance *VirtualizationInstance) callback_getRemoteFilePath(info *cfapi.CF_CALLBACK_INFO) string {
	return instance.path_localToRemote(instance.callback_getFilePath(info))
}

type transferBuffer struct {
	info       *cfapi.CF_CALLBACK_INFO
	instance   *VirtualizationInstance
	transfer   cfapi.CF_OPERATION_PARAMETERS_TransferData
	buffer     []byte
	count      int64
	byteOffset int64
}

func (tb *transferBuffer) send(updatehash hash.Hash) error {
	if tb.count == 0 {
		return nil
	}
	tb.instance.Logger.Debug().Msgf("Sending %d bytes", tb.count)
	tb.transfer.Buffer = uintptr(unsafe.Pointer(&tb.buffer[0]))
	tb.transfer.Length = tb.count
	tb.transfer.Offset = tb.byteOffset
	tb.transfer.ParamSize = uint32(unsafe.Sizeof(tb.transfer))
	tb.transfer.Flags = cfapi.CF_OPERATION_TRANSFER_DATA_FLAG_NONE
	hr := tb.instance.transferData(tb.info, &tb.transfer)

	if updatehash != nil {
		updatehash.Write(tb.buffer[:tb.transfer.Length])
	}

	tb.byteOffset += tb.count
	tb.count = 0

	return core.ErrorByCode(hr)
}

func (instance *VirtualizationInstance) fetchData(info *cfapi.CF_CALLBACK_INFO, data *cfapi.CF_CALLBACK_PARAMETERS_FetchData) uintptr {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	localpath := instance.callback_getFilePath(info)
	instance.FileDownloading(localpath, 0)

	filename := instance.path_localToRemote(localpath)
	length := data.RequiredLength
	byteOffset := data.RequiredFileOffset
	remoteinfo, err := instance.fs.Stat(filename)
	if err != nil {
		instance.Logger.Error().Msgf("Remote file is inaccessible %s: %s", filename, err)
		instance.FileError(localpath, err)
		return uintptr(syscall.EIO)
	}
	if length == 0 || length < 0 {
		length = remoteinfo.Size()
	}
	if data.OptionalLength > data.RequiredLength {
		length = data.OptionalLength
		byteOffset = data.OptionalFileOffset
	}

	wholeFileRequested := (byteOffset == 0) && (length == remoteinfo.Size())
	var updatehash hash.Hash
	if wholeFileRequested {
		// If the whole file is to be downloaded (a full hydration), it is a chance to update our cache of remote state
		updatehash = md5.New()
	}
	instance.Logger.Debug().Msgf("Fetch data: %s %d bytes at %d", filename, length, byteOffset)
	instance.Logger.Debug().Msgf("Optional %d at %d", data.OptionalLength, data.OptionalFileOffset)
	file, err := instance.fs.Open(filename)
	if err != nil {
		instance.Logger.Error().Msgf("Error opening file %s: %s", filename, err)
		instance.FileError(localpath, err)
		return uintptr(syscall.EIO)
	}
	defer file.Close()
	tb := &transferBuffer{
		info:       info,
		instance:   instance,
		buffer:     make([]byte, min(length, BUFFER_SIZE)),
		byteOffset: byteOffset,
		count:      0,
	}

	var n int
	var count int64
	for count < length {
		//instance.Logger.Debug().Msgf("Reading %d bytes", length-count)
		// last read may be partial
		n, err = file.ReadAt(tb.buffer[tb.count:min(BUFFER_SIZE, tb.count+length-count)], byteOffset+count)
		//instance.Logger.Debug().Msgf("Received %d bytes (%v)", n, err)
		count += int64(n)
		tb.count += int64(n)
		if err == io.EOF {
			instance.Logger.Debug().Msgf("Stream ended at %d bytes", count)
			err = nil
			break
		}
		if tb.count >= BUFFER_SIZE {
			err = tb.send(updatehash)
			if err != nil {
				instance.Logger.Error().Msgf("Error computing file hash %s: %s", filename, err)
				instance.FileError(localpath, err)
				return uintptr(syscall.EIO)
			}
			instance.FileDownloading(localpath, int(100*(float32(count)/float32(length))))
		}
	}
	err = tb.send(updatehash)
	instance.Logger.Debug().Msgf("Read %d bytes", count)
	if err != nil {
		instance.Logger.Error().Msgf("Error reading file %s: %s", filename, err)
		instance.FileError(localpath, err)
		return uintptr(syscall.EIO)
	}
	if updatehash != nil {
		err := instance.remoteCacheState.UpdateHash(filename, updatehash.Sum(nil))
		if err != nil {
			instance.Logger.Warn().Msgf("Error updating state cache %s: %s", filename, err)
		}
	}
	instance.FileDone(localpath)

	return 0
}
