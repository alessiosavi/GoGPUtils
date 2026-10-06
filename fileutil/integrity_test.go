package fileutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
	"testing"
)

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Errorf("%s: got %d bytes, want %d bytes with unchanged contents", path, len(got), len(want))
	}
}

func assertNoCopyTemps(t *testing.T, dir string) {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join(dir, ".fileutil-copy-*"))
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 0 {
		t.Errorf("temporary copy files left behind: %v", paths)
	}
}

func TestCopy_SameFileLeavesContentsUntouched(t *testing.T) {
	for _, alias := range []string{"same path", "dot path", "hard link", "symlink"} {
		t.Run(alias, func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "must survive")
			dst := src

			switch alias {
			case "dot path":
				dst = dir + string(os.PathSeparator) + "." + string(os.PathSeparator) + "source"
			case "hard link":
				dst = filepath.Join(dir, "alias")
				if err := os.Link(src, dst); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				dst = filepath.Join(dir, "alias")
				if err := os.Symlink(src, dst); err != nil {
					t.Fatal(err)
				}
			}

			if err := Copy(context.Background(), src, dst); !errors.Is(err, ErrSameFile) {
				t.Errorf("Copy() error = %v, want ErrSameFile", err)
			}

			assertFileContent(t, src, "must survive")
			assertFileContent(t, dst, "must survive")
			assertNoCopyTemps(t, dir)
		})
	}
}

// Observe progress synchronously when Copy polls the context, avoiding timing races.
type checkingContext struct {
	context.Context
	check func()
}

func (c checkingContext) Err() error {
	c.check()

	return c.Context.Err()
}

func (c checkingContext) Done() <-chan struct{} {
	c.check()

	return c.Context.Done()
}

func TestCopy_CancellationPreservesDestination(t *testing.T) {
	for _, exists := range []bool{false, true} {
		t.Run(fmt.Sprintf("destination_exists=%t", exists), func(t *testing.T) {
			dir := t.TempDir()
			content := strings.Repeat("source data", 1<<17)
			src := createTestFile(t, dir, "source", content)
			dst := filepath.Join(dir, "destination")
			if exists {
				createTestFile(t, dir, "destination", "old destination")
			}

			base, cancel := context.WithCancel(context.Background())
			defer cancel()

			var copied int64
			ctx := checkingContext{Context: base, check: func() {
				paths, err := filepath.Glob(filepath.Join(dir, ".fileutil-copy-*"))
				if err != nil {
					t.Fatal(err)
				}

				for _, path := range paths {
					info, err := os.Stat(path)
					if err != nil {
						t.Fatal(err)
					}

					if info.Size() > 0 {
						copied = info.Size()
						cancel()
					}
				}
			}}

			if err := Copy(ctx, src, dst); !errors.Is(err, context.Canceled) {
				t.Errorf("Copy() error = %v, want context.Canceled", err)
			}

			if copied == 0 || copied >= int64(len(content)) {
				t.Errorf("cancellation observed after %d of %d bytes; want a partial copy", copied, len(content))
			}

			if exists {
				assertFileContent(t, dst, "old destination")
			} else if _, err := os.Stat(dst); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("destination should not exist: %v", err)
			}

			assertFileContent(t, src, content)
			assertNoCopyTemps(t, dir)
		})
	}
}

func TestCopy_ReplacesDestinationWithoutChangingItsOtherLinks(t *testing.T) {
	for _, link := range []string{"hard link", "symlink"} {
		t.Run(link, func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "new contents")
			other := createTestFile(t, dir, "other", "old contents")
			dst := filepath.Join(dir, "destination")
			if err := os.Chmod(src, 0750); err != nil {
				t.Fatal(err)
			}

			makeLink := os.Link
			if link == "symlink" {
				makeLink = os.Symlink
			}

			if err := makeLink(other, dst); err != nil {
				t.Fatal(err)
			}

			if err := Copy(context.Background(), src, dst); err != nil {
				t.Fatal(err)
			}

			assertFileContent(t, dst, "new contents")
			assertFileContent(t, other, "old contents")
			info, err := os.Lstat(dst)
			if err != nil {
				t.Fatal(err)
			}

			if !info.Mode().IsRegular() || info.Mode().Perm() != 0750 {
				t.Errorf("destination mode = %v, want regular file with mode 0750", info.Mode())
			}

			assertNoCopyTemps(t, dir)
		})
	}
}

func TestCopy_RenameFailureCleansUpTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	src := createTestFile(t, dir, "source", "source contents")
	dst := filepath.Join(dir, "destination")
	if err := os.Mkdir(dst, 0755); err != nil {
		t.Fatal(err)
	}

	child := createTestFile(t, dst, "child", "keep me")
	if err := Copy(context.Background(), src, dst); err == nil {
		t.Fatal("Copy() succeeded over a directory")
	}

	assertFileContent(t, src, "source contents")
	assertFileContent(t, child, "keep me")
	assertNoCopyTemps(t, dir)
}

func TestMove_RenamesIntoExistingOrMissingDirectory(t *testing.T) {
	for _, parent := range []string{".", "new/nested"} {
		t.Run(parent, func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "contents")
			dst := filepath.Join(dir, parent, "destination")
			before, err := os.Stat(src)
			if err != nil {
				t.Fatal(err)
			}

			if err := Move(context.Background(), src, dst); err != nil {
				t.Fatal(err)
			}

			after, err := os.Stat(dst)
			if err != nil {
				t.Fatal(err)
			}

			if !os.SameFile(before, after) {
				t.Error("Move copied the file instead of renaming it on the same filesystem")
			}

			if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("source should be removed: %v", err)
			}

			assertFileContent(t, dst, "contents")
		})
	}
}

func TestMove_CanceledContextHasNoSideEffects(t *testing.T) {
	for _, parent := range []string{".", "new/nested"} {
		t.Run(parent, func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "source contents")
			dst := filepath.Join(dir, parent, "destination")
			if parent == "." {
				createTestFile(t, dir, "destination", "destination contents")
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			if err := Move(ctx, src, dst); !errors.Is(err, context.Canceled) {
				t.Errorf("Move() error = %v, want context.Canceled", err)
			}

			assertFileContent(t, src, "source contents")
			if parent == "." {
				assertFileContent(t, dst, "destination contents")
			} else if _, err := os.Stat(filepath.Join(dir, "new")); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("canceled Move created a destination directory: %v", err)
			}
		})
	}
}

func TestMove_PermissionDeniedDoesNotOverwriteDestination(t *testing.T) {
	if os.Geteuid() == 0 || runtime.GOOS == "windows" {
		t.Skip("requires Unix directory permissions and a non-root user")
	}

	dir := t.TempDir()
	srcDir := filepath.Join(dir, "source-dir")
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	src := createTestFile(t, srcDir, "source", "source contents")
	dst := createTestFile(t, dir, "destination", "destination contents")
	if err := os.Chmod(srcDir, 0555); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.Chmod(srcDir, 0755); err != nil {
			t.Error(err)
		}
	})

	err := Move(context.Background(), src, dst)
	var linkErr *os.LinkError
	if !errors.Is(err, os.ErrPermission) || !errors.As(err, &linkErr) || linkErr.Op != "rename" {
		t.Errorf("Move() error = %v, want original permission-denied rename error", err)
	}

	assertFileContent(t, src, "source contents")
	assertFileContent(t, dst, "destination contents")
}

func TestMove_NonCrossDeviceErrorsDoNotCopy(t *testing.T) {
	for _, errno := range []syscall.Errno{syscall.EACCES, syscall.ENOENT, syscall.EIO} {
		t.Run(errno.Error(), func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "source contents")
			dst := createTestFile(t, dir, "destination", "destination contents")
			want := &os.LinkError{Op: "rename", Old: src, New: dst, Err: errno}
			err := moveFile(context.Background(), src, dst, func(string, string) error {
				return want
			}, func(string) error {
				t.Error("Move attempted source removal after a non-EXDEV rename error")

				return nil
			})
			if err != want {
				t.Errorf("Move() error = %v, want original error %v", err, want)
			}

			assertFileContent(t, src, "source contents")
			assertFileContent(t, dst, "destination contents")
		})
	}
}

func TestMove_CrossDeviceCopyReportsSourceRemovalFailure(t *testing.T) {
	renameErr := syscall.EXDEV
	if runtime.GOOS == "windows" {
		renameErr = syscall.Errno(17) // ERROR_NOT_SAME_DEVICE returned by MoveFileEx.
	}

	for _, removeFails := range []bool{false, true} {
		t.Run(fmt.Sprintf("remove_fails=%t", removeFails), func(t *testing.T) {
			dir := t.TempDir()
			src := createTestFile(t, dir, "source", "source contents")
			dst := createTestFile(t, dir, "destination", "old contents")
			removeErr := &os.PathError{Op: "remove", Path: src, Err: syscall.EACCES}
			removeCalled := false
			err := moveFile(context.Background(), src, dst, func(string, string) error {
				return &os.LinkError{Op: "rename", Old: src, New: dst, Err: renameErr}
			}, func(path string) error {
				removeCalled = true
				assertFileContent(t, dst, "source contents")
				if path != src {
					t.Errorf("removing %q, want %q", path, src)
				}

				if removeFails {
					return removeErr
				}

				return os.Remove(path)
			})
			if !removeCalled {
				t.Error("Move did not attempt removal after a successful copy")
			}

			if removeFails {
				if !errors.Is(err, removeErr) || err == removeErr || !strings.Contains(err.Error(), "copied") {
					t.Errorf("Move() error = %v, want wrapped removal failure explaining successful copy", err)
				}

				assertFileContent(t, src, "source contents")
			} else if err != nil {
				t.Fatal(err)
			} else if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("source should be removed: %v", err)
			}

			assertFileContent(t, dst, "source contents")
			assertNoCopyTemps(t, dir)
		})
	}
}

func TestReadLines_LongLinesAndLimits(t *testing.T) {
	long := strings.Repeat("x", 1<<20)
	path := createTestFile(t, t.TempDir(), "lines", long+"\r\nsecond\n"+long)
	for _, limit := range []int{-1, 0, 1, 2, 3, 4} {
		t.Run(fmt.Sprintf("limit=%d", limit), func(t *testing.T) {
			want := []string{long, "second", long}
			if limit > 0 && limit < len(want) {
				want = want[:limit]
			}

			got, err := ReadLinesN(context.Background(), path, limit)
			if err != nil || !slices.Equal(got, want) {
				t.Errorf("ReadLinesN() returned %d lines, error = %v; want %d complete lines", len(got), err, len(want))
			}
		})
	}

	got, err := ReadLines(context.Background(), path)
	if err != nil || !slices.Equal(got, []string{long, "second", long}) {
		t.Errorf("ReadLines() returned %d lines, error = %v; want 3 complete lines", len(got), err)
	}
}

func TestReadLines_PreservesLineSemantics(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{"empty", "", nil},
		{"unterminated", "a", []string{"a"}},
		{"terminated", "a\n", []string{"a"}},
		{"blank", "\n\n", []string{"", ""}},
		{"CRLF", "a\r\nb\r\n", []string{"a", "b"}},
		{"final CR", "a\r", []string{"a"}},
		{"only CR", "\r", []string{""}},
		{"embedded CR", "a\rb\r\r\n", []string{"a\rb\r"}},
		{"mixed", "\na\r\n\nb\n\r", []string{"", "a", "", "b", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestFile(t, t.TempDir(), "lines", tt.content)
			for _, limit := range []int{-1, 0, 1, 10} {
				want := tt.want
				if limit > 0 && limit < len(want) {
					want = want[:limit]
				}

				got, err := ReadLinesN(context.Background(), path, limit)
				if err != nil || !slices.Equal(got, want) {
					t.Errorf("ReadLinesN(%d) = %q, %v; want %q", limit, got, err, want)
				}

				if tt.content == "" && (got == nil) != (limit <= 0) {
					t.Errorf("ReadLinesN(%d) changed the empty-file nil slice behavior", limit)
				}
			}
		})
	}
}

func TestReaders_CanceledContextBeforeOpening(t *testing.T) {
	readers := []struct {
		name string
		read func(context.Context, string) error
	}{
		{"ReadBytes", func(ctx context.Context, path string) error { _, err := ReadBytes(ctx, path); return err }},
		{"ReadLines", func(ctx context.Context, path string) error { _, err := ReadLines(ctx, path); return err }},
		{"ReadLinesN", func(ctx context.Context, path string) error { _, err := ReadLinesN(ctx, path, 1); return err }},
		{"CountLines", func(ctx context.Context, path string) error { _, err := CountLines(ctx, path); return err }},
		{"DetectFileLineTerminator", func(ctx context.Context, path string) error {
			_, err := DetectFileLineTerminator(ctx, path)
			return err
		}},
	}
	dir := t.TempDir()
	empty := createTestFile(t, dir, "empty", "")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, reader := range readers {
		t.Run(reader.name, func(t *testing.T) {
			for _, path := range []string{empty, filepath.Join(dir, "missing")} {
				if err := reader.read(ctx, path); !errors.Is(err, context.Canceled) {
					t.Errorf("reading %s: error = %v, want context.Canceled", filepath.Base(path), err)
				}
			}
		})
	}
}

func TestCountLines_EqualsReadLinesLength(t *testing.T) {
	inputs := []string{"", "a", "a\nb", "a\n", "\n", "\n\n", "a\n\n", "\r", "\r\n", "a\rb", "a\r\nb\r", "\x00\n\x00"}
	for _, size := range []int{32767, 32768, 32769, 1 << 20} {
		line := strings.Repeat("x", size)
		inputs = append(inputs, line, line+"\n", line+"\nb", line+"\r\n\n\r")
	}

	for i, input := range inputs {
		t.Run(fmt.Sprintf("input=%d", i), func(t *testing.T) {
			path := createTestFile(t, t.TempDir(), "lines", input)
			lines, err := ReadLines(context.Background(), path)
			if err != nil {
				t.Fatal(err)
			}

			got, err := CountLines(context.Background(), path)
			if err != nil || got != len(lines) {
				t.Errorf("CountLines() = %d, %v; len(ReadLines()) = %d", got, err, len(lines))
			}
		})
	}
}

type failingCloser struct {
	err   error
	calls int
}

func (f *failingCloser) Close() error {
	f.calls++

	return f.err
}

func TestWriters_ReportCloseErrorAfterSuccessfulWrite(t *testing.T) {
	closeErr := errors.New("close failed")
	for _, writeErr := range []error{nil, io.ErrShortWrite} {
		t.Run(fmt.Sprintf("write_error=%v", writeErr), func(t *testing.T) {
			f := &failingCloser{err: closeErr}
			want := writeErr
			if want == nil {
				want = closeErr
			}

			if got := finishWrite(f, writeErr); !errors.Is(got, want) {
				t.Errorf("finishWrite() = %v, want %v", got, want)
			}

			if f.calls != 1 {
				t.Errorf("Close() called %d times, want 1", f.calls)
			}
		})
	}
}

func TestRemoveTemp_ReportsCleanupFailure(t *testing.T) {
	// A nonempty directory fails removal without relying on Unix permissions
	// or an open file, whose removal behavior differs on Windows.
	dir := t.TempDir()
	child := createTestFile(t, dir, "child", "keep me")
	primary := fmt.Errorf("copy failed: %w", context.Canceled)
	err := removeTemp(dir, primary)
	if !errors.Is(err, primary) || !errors.Is(err, context.Canceled) {
		t.Errorf("removeTemp() error = %v, want preserved primary error", err)
	}

	var removeErr *os.PathError
	if !errors.As(err, &removeErr) {
		t.Fatalf("removeTemp() error = %v, want wrapped removal error", err)
	}

	if removeErr.Op != "remove" || removeErr.Path != dir || !errors.Is(err, removeErr.Err) {
		t.Errorf("removeTemp() error = %v, want removal failure for %q with matchable cause", err, dir)
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("remove temporary file %q:", dir)) {
		t.Errorf("removeTemp() error = %v, want temporary-file cleanup context", err)
	}

	assertFileContent(t, child, "keep me")
}

func TestRemoveTemp_PreservesPrimaryError(t *testing.T) {
	for _, exists := range []bool{false, true} {
		t.Run(fmt.Sprintf("file_exists=%t", exists), func(t *testing.T) {
			dir := t.TempDir()
			name := filepath.Join(dir, "temporary")
			if exists {
				createTestFile(t, dir, "temporary", "contents")
			}

			primary := fmt.Errorf("write failed: %w", io.ErrShortWrite)
			if err := removeTemp(name, primary); err != primary { //nolint:errorlint // Successful cleanup must return the identical error value.
				t.Errorf("removeTemp() error = %v, want original error %v", err, primary)
			}

			if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("temporary file should not exist: %v", err)
			}
		})
	}
}
