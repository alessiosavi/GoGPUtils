package fileutil

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

func TestZZWriteFullDifferential(t *testing.T) {
	if _, err := os.Stat("/dev/full"); err != nil {
		t.Skipf("/dev/full unavailable: %v", err)
	}
	for _, pair := range zzWriters() {
		t.Run(pair.name, func(t *testing.T) {
			want := pair.base("/dev/full")
			got := pair.current("/dev/full")
			zzCheckError(t, got, want)
			var pathErr *os.PathError
			if !errors.As(got, &pathErr) || pathErr.Op != "write" || pathErr.Path != "/dev/full" || !errors.Is(got, syscall.ENOSPC) {
				t.Fatalf("error = %#v, want write /dev/full: ENOSPC", got)
			}
		})
	}
}
