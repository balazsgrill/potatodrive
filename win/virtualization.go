package win

import (
	"encoding/binary"
	"io"
	"syscall"
)

type Virtualization interface {
	io.Closer
	PerformSynchronization() error
	SetFileStateHandler(handler func(state FileSyncState))
}

func BytesToGuid(b []byte) *syscall.GUID {
	return &syscall.GUID{
		Data1: binary.LittleEndian.Uint32(b[0:4]),
		Data2: binary.LittleEndian.Uint16(b[4:6]),
		Data3: binary.LittleEndian.Uint16(b[6:8]),
		Data4: ([8]byte)(b[8:16]),
	}
}

type ConnectionState struct {
	ID             string
	SyncInProgress bool
	LastSyncError  error
}

type FileSyncStateEnum int

const (
	FileSyncStateUnknown     FileSyncStateEnum = 0
	FileSyncStatePending     FileSyncStateEnum = 1
	FileSyncStateUploading   FileSyncStateEnum = 2
	FileSyncStateDownloading FileSyncStateEnum = 3
	FileSyncStateDone        FileSyncStateEnum = 4
	FileSyncStateDeleted     FileSyncStateEnum = 5
	FileSyncStateDirty       FileSyncStateEnum = 6
	FileSyncStateError       FileSyncStateEnum = 7
)

type FileSyncState struct {
	Path      string
	State     FileSyncStateEnum
	LastError error
}
