package filesystem

import (
	"io"
	"syscall"
	"unsafe"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/balazsgrill/potatodrive/win/cfapi"
)

const BUFFER_SIZE int64 = 1024 * 1024

func (instance *VirtualizationInstance) callback_getRemoteFilePath(info *cfapi.CF_CALLBACK_INFO) string {
	return instance.path_localToRemote(win.GetString(info.VolumeDosName) + win.GetString(info.NormalizedPath))
}

type transferBuffer struct {
	info       *cfapi.CF_CALLBACK_INFO
	instance   *VirtualizationInstance
	transfer   cfapi.CF_OPERATION_PARAMETERS_TransferData
	buffer     []byte
	count      int64
	byteOffset int64
}

func (tb *transferBuffer) send() error {
	if tb.count == 0 {
		return nil
	}
	tb.transfer.Buffer = uintptr(unsafe.Pointer(&tb.buffer[0]))
	tb.transfer.Length = tb.count
	tb.transfer.Offset = tb.byteOffset
	tb.transfer.ParamSize = uint32(unsafe.Sizeof(tb.transfer))
	tb.transfer.Flags = cfapi.CF_OPERATION_TRANSFER_DATA_FLAG_NONE
	hr := tb.instance.transferData(tb.info, &tb.transfer)
	tb.byteOffset += tb.count
	tb.count = 0
	return win.ErrorByCode(hr)
}

func (instance *VirtualizationInstance) fetchData(info *cfapi.CF_CALLBACK_INFO, data *cfapi.CF_CALLBACK_PARAMETERS_FetchData) uintptr {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	filename := instance.callback_getRemoteFilePath(info)
	length := data.RequiredLength
	byteOffset := data.RequiredFileOffset
	if length == 0 || length < 0 {
		length = info.FileSize
	}
	if data.OptionalLength > data.RequiredLength {
		length = data.OptionalLength
		byteOffset = data.OptionalFileOffset
	}
	instance.Logger.Debug().Msgf("Fetch data: %s %d bytes at %d", filename, length, byteOffset)
	instance.Logger.Debug().Msgf("Optional %d at %d", data.OptionalLength, data.OptionalFileOffset)
	file, err := instance.fs.Open(filename)
	if err != nil {
		instance.Logger.Debug().Msgf("Error opening file %s: %s", filename, err)
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
		n, err = file.ReadAt(tb.buffer[tb.count:], byteOffset+count)
		count += int64(n)
		if err == io.EOF {
			err = nil
			break
		}
		tb.count += int64(n)
		if tb.count >= BUFFER_SIZE {
			err = tb.send()
			if err != nil {
				return uintptr(syscall.EIO)
			}
		}
	}
	err = tb.send()
	instance.Logger.Debug().Msgf("Read %d bytes", count)
	if err != nil {
		instance.Logger.Debug().Msgf("Error reading file %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}

	return 0
}
