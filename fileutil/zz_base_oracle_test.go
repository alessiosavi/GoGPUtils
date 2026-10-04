package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
	"weak"
)

// Private copies from BASE a8b7b93b997b8986200995fdf2c65a5efbf2053b.
// Bodies are verbatim except for names and calls to private BASE dependencies.
// WriteBytes, AppendBytes, Exists and IsDir remain unchanged in this PR.
func baseNormalizeLineTerminators(data []byte, target LineTerminator) []byte {
	// First normalize to LF
	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\n'})
	data = bytes.ReplaceAll(data, []byte{'\r'}, []byte{'\n'})

	// Then convert to target
	if target == LF {
		return data
	}

	return bytes.ReplaceAll(data, []byte{'\n'}, target.Bytes())
}

//nolint:intrange // Preserve the verbatim BASE loop.
func baseDetectLineTerminator(data []byte) LineTerminator {
	hasLF := false
	hasCRLF := false
	hasCR := false

	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '\r':
			if i+1 < len(data) && data[i+1] == '\n' {
				hasCRLF = true
				i++ // Skip the \n
			} else {
				hasCR = true
			}
		case '\n':
			hasLF = true
		}
	}

	// Count how many types we found
	count := 0
	if hasLF {
		count++
	}

	if hasCRLF {
		count++
	}

	if hasCR {
		count++
	}

	if count == 0 {
		return Unknown
	}

	if count > 1 {
		return Mixed
	}

	if hasCRLF {
		return CRLF
	}

	if hasCR {
		return CR
	}

	return LF
}

func baseWriteLines(path string, lines []string, terminator LineTerminator, perm fs.FileMode) error {
	term := string(terminator.Bytes())

	content := strings.Join(lines, term)

	if len(lines) > 0 {
		content += term
	}

	return baseWriteString(path, content, perm)
}

func baseWriteString(path, content string, perm fs.FileMode) error {
	return WriteBytes(path, []byte(content), perm)
}

func baseAppendString(path, content string, perm fs.FileMode) error {
	return AppendBytes(path, []byte(content), perm)
}

func baseAppendLine(path, line string, terminator LineTerminator, perm fs.FileMode) error {
	return baseAppendString(path, line+string(terminator.Bytes()), perm)
}

func baseEnsureDir(path string, perm fs.FileMode) error {
	if Exists(path) {
		if !IsDir(path) {
			return ErrNotDir
		}

		return nil
	}

	return os.MkdirAll(path, perm)
}

func zzTerminatorTargets() []LineTerminator {
	return []LineTerminator{LF, CRLF, CR, Mixed, Unknown, 99, -1}
}

func zzCheckLineTerminators(t testing.TB, data []byte, target LineTerminator) {
	t.Helper()
	original := bytes.Clone(data)
	wantKind := baseDetectLineTerminator(data)
	gotKind := DetectLineTerminator(data)
	if gotKind != wantKind {
		t.Fatalf("DetectLineTerminator(%q) = %v, BASE = %v", data, gotKind, wantKind)
	}
	want := baseNormalizeLineTerminators(data, target)
	got := NormalizeLineTerminators(data, target)
	if !bytes.Equal(got, want) || (got == nil) != (want == nil) {
		t.Fatalf("NormalizeLineTerminators(%q, %d) = %q (nil=%v), BASE = %q (nil=%v)",
			data, target, got, got == nil, want, want == nil)
	}
	if !bytes.Equal(data, original) {
		t.Fatal("line terminator operation changed its input")
	}
	// Mutate every result byte: even an interior input view must not alias it.
	for i := range got {
		got[i] ^= 0xff
	}
	if !bytes.Equal(data, original) {
		t.Fatal("normalized result aliases its input")
	}
	for i := range want {
		want[i] ^= 0xff
	}
	if !bytes.Equal(data, original) {
		t.Fatal("BASE normalized result aliases its input")
	}
}

func TestZZLineTerminatorsExhaustive(t *testing.T) {
	inputs := 0
	var visit func([]byte)
	visit = func(data []byte) {
		inputs++
		for _, target := range zzTerminatorTargets() {
			zzCheckLineTerminators(t, data, target)
		}
		if len(data) < 10 {
			for _, ch := range []byte{'a', '\r', '\n'} {
				visit(append(data, ch))
			}
		}
	}
	visit([]byte{})
	for _, target := range zzTerminatorTargets() {
		zzCheckLineTerminators(t, nil, target)
		zzCheckLineTerminators(t, []byte("\x00\xff\xc0\x80\xef\xbf\xbd\r\r\n"), target)
		// A small view of a large input must still return a detached result.
		parent := bytes.Repeat([]byte{'a'}, 1<<20)
		zzCheckLineTerminators(t, parent[len(parent)/2:len(parent)/2+1], target)
	}
	if inputs != 88573 {
		t.Fatalf("visited %d inputs, want 88573", inputs)
	}
	t.Logf("%d exhaustive inputs, %d normalization comparisons", inputs, inputs*len(zzTerminatorTargets()))
}

// A detached result must not retain a large backing array behind a tiny view.
func TestZZNormalizeDoesNotRetainInput(t *testing.T) {
	for _, normalize := range []struct {
		name string
		fn   func([]byte, LineTerminator) []byte
	}{{"BASE", baseNormalizeLineTerminators}, {"HEAD", NormalizeLineTerminators}} {
		for _, target := range zzTerminatorTargets() {
			for _, input := range []string{"a", "\r\n", "\r\r\n"} {
				t.Run(fmt.Sprintf("%s/target=%d/input=%q", normalize.name, target, input), func(t *testing.T) {
					result, parent := func() ([]byte, weak.Pointer[byte]) {
						data := bytes.Repeat([]byte{'a'}, 1<<20)
						start := len(data) / 2
						copy(data[start:], input)
						return normalize.fn(data[start:start+len(input)], target), weak.Make(&data[0])
					}()
					runtime.GC()
					if parent.Value() != nil {
						t.Fatal("live normalized result retains the input's large backing array")
					}
					runtime.KeepAlive(result)
				})
			}
		}
	}
}

func FuzzZZLineTerminatorsDifferential(f *testing.F) {
	for _, data := range [][]byte{nil, {}, []byte("plain"), []byte("\r\r\n"), []byte("\r\n\r\n"), []byte("\n\r"), []byte("\x00\xff\xc0\x80\xef\xbf\xbd\r\n")} {
		for _, target := range zzTerminatorTargets() {
			f.Add(data, int(target))
		}
	}
	f.Fuzz(func(t *testing.T, data []byte, target int) {
		zzCheckLineTerminators(t, data, LineTerminator(target))
	})
}

func zzCheckError(t *testing.T, got, want error) {
	t.Helper()
	// DeepEqual includes PathError's exact type, operation, path and errno.
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("error = %#v (%T), BASE = %#v (%T)", got, got, want, want)
	}
	for _, sentinel := range []error{ErrNotDir, fs.ErrNotExist, fs.ErrPermission, fs.ErrExist} {
		if errors.Is(got, sentinel) != errors.Is(want, sentinel) {
			t.Fatalf("error identity differs for %v: got %v, BASE %v", sentinel, got, want)
		}
	}
	if want == ErrNotDir && got != ErrNotDir { //nolint:errorlint // Exact sentinel identity is the contract.
		t.Fatalf("error = %v, want the ErrNotDir sentinel", got)
	}
}

// Each side uses the same path and starting state; paths in errors stay exact.
func zzCompareWrite(t *testing.T, base, current func(string) error, existing bool) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "output")
	before := []byte("existing\x00\xff\ncontents long enough to expose a missing truncate")
	var want []byte
	var wantMode fs.FileMode
	for i, write := range []func(string) error{base, current} {
		if existing {
			if err := os.WriteFile(path, before, 0o640); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, 0o640); err != nil {
				t.Fatal(err)
			}
		}
		if err := write(path); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			want, wantMode = data, info.Mode()
		} else if !bytes.Equal(data, want) || info.Mode() != wantMode {
			t.Fatalf("contents/mode = %q/%v, BASE = %q/%v", data, info.Mode(), want, wantMode)
		}
		if existing && runtime.GOOS != "windows" && info.Mode().Perm() != 0o640 {
			t.Fatalf("existing permissions changed to %v", info.Mode())
		}
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
}

func TestZZWriteLinesDifferential(t *testing.T) {
	inputs := [][]string{nil, {}, {""}, {"a"}, {"a", ""}, {"", ""}, {"a\r\nb", "\r", "\n", "\x00\xff\xc0\x80", "世界\ufffd"}}
	for i, lines := range inputs {
		for _, target := range zzTerminatorTargets() {
			for _, existing := range []bool{false, true} {
				t.Run(fmt.Sprintf("input=%d/target=%d/existing=%v", i, target, existing), func(t *testing.T) {
					before := slices.Clone(lines)
					zzCompareWrite(t,
						func(path string) error { return baseWriteLines(path, lines, target, 0o600) },
						func(path string) error { return WriteLines(path, lines, target, 0o600) }, existing)
					if !reflect.DeepEqual(lines, before) {
						t.Fatal("WriteLines changed its input")
					}
				})
			}
		}
	}
}

func TestZZStringWritesDifferential(t *testing.T) {
	for _, input := range []string{"", "abc", "\r\n\r\n\x00\xff\xc0\x80世界\ufffd", strings.Repeat("x", 65536)} {
		for _, existing := range []bool{false, true} {
			t.Run(fmt.Sprintf("bytes=%d/existing=%v", len(input), existing), func(t *testing.T) {
				zzCompareWrite(t,
					func(path string) error { return baseWriteString(path, input, 0o600) },
					func(path string) error { return WriteString(path, input, 0o600) }, existing)
				zzCompareWrite(t,
					func(path string) error { return baseAppendString(path, input, 0o600) },
					func(path string) error { return AppendString(path, input, 0o600) }, existing)
				for _, target := range zzTerminatorTargets() {
					zzCompareWrite(t,
						func(path string) error { return baseAppendLine(path, input, target, 0o600) },
						func(path string) error { return AppendLine(path, input, target, 0o600) }, existing)
				}
			})
		}
	}
}

type zzWriterPair struct {
	name          string
	base, current func(string) error
}

func zzWriters() []zzWriterPair {
	return []zzWriterPair{
		{"WriteString", func(p string) error { return baseWriteString(p, "\x00x\xff", 0o600) }, func(p string) error { return WriteString(p, "\x00x\xff", 0o600) }},
		{"WriteLines", func(p string) error { return baseWriteLines(p, []string{"a", "b"}, CRLF, 0o600) }, func(p string) error { return WriteLines(p, []string{"a", "b"}, CRLF, 0o600) }},
		{"AppendString", func(p string) error { return baseAppendString(p, "\x00x\xff", 0o600) }, func(p string) error { return AppendString(p, "\x00x\xff", 0o600) }},
		{"AppendLine", func(p string) error { return baseAppendLine(p, "\x00x\xff", CRLF, 0o600) }, func(p string) error { return AppendLine(p, "\x00x\xff", CRLF, 0o600) }},
	}
}

func TestZZWriteErrorsDifferential(t *testing.T) {
	for _, pair := range zzWriters() {
		for _, kind := range []string{"missing-parent", "directory", "read-only"} {
			t.Run(pair.name+"/"+kind, func(t *testing.T) {
				if kind == "read-only" && (os.Geteuid() == 0 || runtime.GOOS == "windows") {
					t.Skip("requires unprivileged Unix permissions")
				}
				path := filepath.Join(t.TempDir(), "target")
				switch kind {
				case "missing-parent":
					path = filepath.Join(path, "file")
				case "directory":
					if err := os.Mkdir(path, 0o700); err != nil {
						t.Fatal(err)
					}
				case "read-only":
					if err := os.WriteFile(path, []byte("unchanged"), 0o400); err != nil {
						t.Fatal(err)
					}
				}
				want := pair.base(path)
				got := pair.current(path)
				if want == nil {
					t.Fatal("BASE did not exercise the error path")
				}
				zzCheckError(t, got, want)
				if kind == "read-only" {
					data, err := os.ReadFile(path)
					if err != nil || string(data) != "unchanged" {
						t.Fatalf("read-only target changed: %q, %v", data, err)
					}
				}
			})
		}
	}
}

func TestZZEnsureDirDifferential(t *testing.T) {
	for _, kind := range []string{"directory", "file", "missing", "symlink-dir", "symlink-file", "dangling-symlink", "not-dir-component", "unreadable-parent"} {
		t.Run(kind, func(t *testing.T) {
			if kind == "unreadable-parent" && (os.Geteuid() == 0 || runtime.GOOS == "windows") {
				t.Skip("requires unprivileged Unix permissions")
			}
			root := filepath.Join(t.TempDir(), "fixture")
			path := filepath.Join(root, "target")
			var wantErr error
			var wantMode fs.FileMode
			for i, ensure := range []func(string, fs.FileMode) error{baseEnsureDir, EnsureDir} {
				if err := os.Mkdir(root, 0o700); err != nil {
					t.Fatal(err)
				}
				switch kind {
				case "directory", "symlink-dir":
					if err := os.Mkdir(path, 0o750); err != nil {
						t.Fatal(err)
					}
				case "file", "symlink-file", "not-dir-component":
					if err := os.WriteFile(path, nil, 0o600); err != nil {
						t.Fatal(err)
					}
				case "unreadable-parent":
					if err := os.Chmod(root, 0); err != nil {
						t.Fatal(err)
					}
				}
				callPath := path
				if strings.Contains(kind, "symlink") {
					callPath = filepath.Join(root, "link")
					if err := os.Symlink(path, callPath); err != nil {
						t.Skipf("symlinks unavailable: %v", err)
					}
				} else if kind == "not-dir-component" || kind == "missing" {
					callPath = filepath.Join(path, "child")
				}
				err := ensure(callPath, 0o700)
				var mode fs.FileMode
				if err == nil {
					info, statErr := os.Stat(callPath)
					if statErr != nil {
						t.Fatal(statErr)
					}
					mode = info.Mode()
				}
				if kind == "unreadable-parent" {
					if chmodErr := os.Chmod(root, 0o700); chmodErr != nil {
						t.Fatal(chmodErr)
					}
					if !errors.Is(err, fs.ErrPermission) {
						t.Fatalf("unreadable parent error = %v, want permission denied", err)
					}
				}
				if i == 0 {
					wantErr, wantMode = err, mode
				} else {
					zzCheckError(t, err, wantErr)
					if mode != wantMode {
						t.Fatalf("mode = %v, BASE = %v", mode, wantMode)
					}
				}
				if err := os.RemoveAll(root); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
