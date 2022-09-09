//go:build linux

package goutils

import "syscall"
import "fmt"

// GetUlimitValue return the current and max value for ulimit
func GetUlimitValue() (uint64, uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Printf("Error Getting Rlimit: %s\n", err)
		return 1024, 1024
	}
	return rLimit.Cur, rLimit.Max
}
