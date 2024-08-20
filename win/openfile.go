package win

import "golang.org/x/sys/windows"

func OpenFile(handle windows.Handle, path string) error {
	pathp, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	verpp, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	return windows.ShellExecute(handle, verpp, pathp, nil, nil, 1)
}
