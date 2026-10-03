//go:build !windows && !plan9

package fileutil

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestIsCrossDeviceError(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil"},
		{name: "cross-device", err: syscall.EXDEV, want: true},
		{name: "rename error", err: &os.LinkError{Op: "rename", Old: "src", New: "dst", Err: syscall.EXDEV}, want: true},
		{name: "wrapped error", err: fmt.Errorf("move: %w", syscall.EXDEV), want: true},
		{name: "permission denied", err: &os.LinkError{Op: "rename", Old: "src", New: "dst", Err: syscall.EACCES}},
		{name: "not found", err: syscall.ENOENT},
		{name: "I/O error", err: syscall.EIO},
		{name: "same text", err: errors.New(syscall.EXDEV.Error())},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := isCrossDeviceError(tc.err); got != tc.want {
				t.Errorf("isCrossDeviceError(%v) = %t, want %t", tc.err, got, tc.want)
			}
		})
	}
}
