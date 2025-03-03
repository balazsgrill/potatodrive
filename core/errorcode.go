package core

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func ErrorByCode(result uintptr) error {
	if result == 0 {
		return nil
	} else {
		message := make([]uint16, 256)
		_, err := windows.FormatMessage(windows.FORMAT_MESSAGE_IGNORE_INSERTS|windows.FORMAT_MESSAGE_FROM_SYSTEM, 0, uint32(result), 0, message, nil)
		if err != nil {
			return fmt.Errorf("can't extract message of %x: %w", result, err)
		}
		return fmt.Errorf("error result: %x - %s", result, windows.UTF16ToString(message))
	}
}

func ErrorByCodeWithContext(context string, result uintptr) error {
	err := ErrorByCode(result)
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}
