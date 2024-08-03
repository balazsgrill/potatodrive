//go:build windows

package projfs

type IProvider interface {
	CancelCommand(commandID int32)
	StartDirectoryEnumeration()
}
