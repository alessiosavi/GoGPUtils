//go:build !windows

package c

import "testing"

func benchCond(b *testing.B) { b.Skip("Cond is available only on windows") }
