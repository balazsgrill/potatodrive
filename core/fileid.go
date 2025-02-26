//go:build windows

package core

import (
	"os"
	"syscall"
)

func GetFileID(filePath string) (uint64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileHandle := syscall.Handle(file.Fd())
	var fileInfo syscall.ByHandleFileInformation

	err = syscall.GetFileInformationByHandle(fileHandle, &fileInfo)
	if err != nil {
		return 0, err
	}

	fileID := (uint64(fileInfo.FileIndexHigh) << 32) + uint64(fileInfo.FileIndexLow)

	return fileID, nil
}
