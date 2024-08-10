package win

import (
	"encoding/binary"
	"io"
	"syscall"
)

type Virtualization interface {
	io.Closer
	PerformSynchronization() error
}

func BytesToGuid(b []byte) *syscall.GUID {
	return &syscall.GUID{
		Data1: binary.LittleEndian.Uint32(b[0:4]),
		Data2: binary.LittleEndian.Uint16(b[4:6]),
		Data3: binary.LittleEndian.Uint16(b[6:8]),
		Data4: ([8]byte)(b[8:16]),
	}
}
