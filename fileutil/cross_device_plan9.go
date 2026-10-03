//go:build plan9

package fileutil

func isCrossDeviceError(err error) bool {
	// Plan 9 rename rejects moves across directories even on the same filesystem
	// (Go's os/file_plan9.go). Preserve the legacy copy fallback for every rename
	// failure, at the cost of also attempting a copy for unrelated rename errors.
	// moveFile calls this helper only after rename returns a non-nil error.
	return true
}
