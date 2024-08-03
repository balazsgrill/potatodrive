//go:build windows

package projfs

import (
	"syscall"
	"unsafe"
)

func GetPointer(str string) uintptr {
	ptr, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return 0
	}
	return uintptr(unsafe.Pointer(ptr))
}

func GetString(str uintptr) string {
	p := (*uint16)(unsafe.Pointer(str))
	if p == nil {
		return ""
	}
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	return syscall.UTF16ToString(unsafe.Slice(p, n))
}
