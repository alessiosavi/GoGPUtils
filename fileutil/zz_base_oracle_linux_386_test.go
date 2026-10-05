package fileutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZZWriteLinesOverflow386(t *testing.T) {
	// The strings share 1 MiB of storage, but their joined length exceeds MaxInt.
	line := strings.Repeat("x", 1<<20)
	lines := make([]string, 2048)
	for i := range lines {
		lines[i] = line
	}
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing=%v", existing), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "must-not-be-opened")
			if existing {
				if err := os.WriteFile(path, []byte("unchanged"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			panicValue := func(write func(string, []string, LineTerminator, fs.FileMode) error) (value any) {
				defer func() { value = recover() }()
				if err := write(path, lines, LF, 0o600); err != nil {
					t.Fatalf("returned an error before the expected panic: %v", err)
				}
				return nil
			}
			want := panicValue(baseWriteLines)
			got := panicValue(WriteLines)
			if want != "strings: Join output length overflow" || got != want {
				t.Fatalf("panic = %#v, BASE = %#v", got, want)
			}
			data, err := os.ReadFile(path)
			if existing {
				if err != nil || string(data) != "unchanged" {
					t.Fatalf("destination changed: %q, %v", data, err)
				}
			} else if !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("destination created: %v", err)
			}
			t.Logf("BASE and HEAD panic = %q; destination unchanged", got)
		})
	}
}
