package filesystem

import (
	"io"
	"log"
	"syscall"
	"unsafe"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/balazsgrill/potatodrive/win/cfapi"
)

func (instance *VirtualizationInstance) callback_getRemoteFilePath(info *cfapi.CF_CALLBACK_INFO) string {
	return instance.path_localToRemote(win.GetString(info.VolumeDosName) + win.GetString(info.NormalizedPath))
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
	log.Printf("Fetch data: %s %d bytes at %d", filename, length, byteOffset)
	log.Printf("Optional %d at %d", data.OptionalLength, data.OptionalFileOffset)
	file, err := instance.fs.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}
	defer file.Close()
	buffer := make([]byte, length)

	var n int
	var count int64
	for count < length {
		n, err = file.ReadAt(buffer[count:], byteOffset+count)
		count += int64(n)
		if err == io.EOF {
			err = nil
			break
		}
	}

	log.Printf("Read %d bytes", count)
	if err != nil {
		log.Printf("Error reading file %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}

	var transfer cfapi.CF_OPERATION_PARAMETERS_TransferData
	transfer.Buffer = uintptr(unsafe.Pointer(&buffer[0]))
	transfer.Length = count
	transfer.Offset = byteOffset
	transfer.ParamSize = uint32(unsafe.Sizeof(transfer))
	transfer.Flags = cfapi.CF_OPERATION_TRANSFER_DATA_FLAG_NONE
	hr := instance.transferData(info, &transfer)
	if hr != 0 {
		log.Printf("Error transferring data: %s", win.ErrorByCode(hr))
		return hr
	}
	return 0
}
