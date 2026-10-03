//go:build windows

package fileutil

import (
	"errors"
	"syscall"
)

// errorNotSameDevice is Windows ERROR_NOT_SAME_DEVICE (17 / 0x11).
// https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes--0-499-
const errorNotSameDevice = syscall.Errno(17)

func isCrossDeviceError(err error) bool {
	// Verified in Go 1.27.1: os/file_windows.go calls windows.Rename in
	// internal/syscall/windows/syscall_windows.go, which uses MoveFileEx with
	// MOVEFILE_REPLACE_EXISTING only (without MOVEFILE_COPY_ALLOWED).
	// Cross-volume failure returns ERROR_NOT_SAME_DEVICE, not syscall.EXDEV.
	return errors.Is(err, errorNotSameDevice)
}
