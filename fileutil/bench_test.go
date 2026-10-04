package fileutil

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

func benchmarkFile(b *testing.B, data []byte) string {
	b.Helper()
	path := filepath.Join(b.TempDir(), "fixture.bin")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		b.Fatal(err)
	}

	return path
}

// benchmarkLines creates n LF-terminated lines of exactly 40 bytes each.
func benchmarkLines(b *testing.B, n int) string {
	b.Helper()
	content := strings.Repeat(benchkit.Text(39, benchkit.ASCII)+"\n", n)

	return benchmarkFile(b, []byte(content))
}

// benchmarkTree creates files 1 KiB files: half at the root and half in four
// subdirectories. Each level contains equal numbers of .txt and .bin files.
func benchmarkTree(b *testing.B, files int) string {
	b.Helper()
	dir := b.TempDir()
	for i := range 4 {
		if err := os.Mkdir(filepath.Join(dir, "sub"+strconv.Itoa(i)), 0o700); err != nil {
			b.Fatal(err)
		}
	}
	data := benchkit.Bytes(1024)
	for i := range files {
		parent := dir
		if i >= files/2 {
			parent = filepath.Join(dir, "sub"+strconv.Itoa(i%4))
		}
		ext := ".txt"
		if i%2 != 0 {
			ext = ".bin"
		}
		path := filepath.Join(parent, "file"+strconv.Itoa(i)+ext)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			b.Fatal(err)
		}
	}

	return dir
}

// BenchmarkReadLines measures reading n LF-terminated lines of 40 bytes each
// from an unchanged file; n is the number of source lines.
func BenchmarkReadLines(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			lines, err := ReadLines(ctx, path)
			if err != nil {
				b.Fatal(err)
			}
			if len(lines) != n {
				b.Fatalf("read %d lines, want %d", len(lines), n)
			}
		}
	})
}

// BenchmarkCountLines measures counting n LF-terminated lines of 40 bytes
// each in an unchanged file; n is the number of source lines.
func BenchmarkCountLines(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			count, err := CountLines(ctx, path)
			if err != nil {
				b.Fatal(err)
			}
			if count != n {
				b.Fatalf("counted %d lines, want %d", count, n)
			}
		}
	})
}

// BenchmarkList measures listing files in an unchanged tree, half at the
// root and half in four subdirectories. files=16 and files=256 are fixed
// filesystem-cardinality overrides to avoid 65536-file setup costs.
// case=flat lists the root files; case=recursive lists all files.
func BenchmarkList(b *testing.B) {
	for _, files := range []int{16, 256} {
		b.Run("files="+strconv.Itoa(files), func(b *testing.B) {
			dir := benchmarkTree(b, files)
			ctx := context.Background()
			b.Run("case=flat", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					entries, err := List(ctx, dir, FilesOnly)
					if err != nil {
						b.Fatal(err)
					}
					if len(entries) != files/2 {
						b.Fatalf("listed %d files, want %d", len(entries), files/2)
					}
				}
			})
			b.Run("case=recursive", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					entries, err := List(ctx, dir, FilesOnly|Recursive)
					if err != nil {
						b.Fatal(err)
					}
					if len(entries) != files {
						b.Fatalf("listed %d files, want %d", len(entries), files)
					}
				}
			})
		})
	}
}

// BenchmarkLineTerminator_Bytes measures byte conversion for each enum value;
// case identifies the immutable scalar terminator, with no size sweep.
func BenchmarkLineTerminator_Bytes(b *testing.B) {
	for _, tc := range []struct {
		name string
		term LineTerminator
	}{
		{"LF", LF},
		{"CRLF", CRLF},
		{"CR", CR},
		{"Mixed", Mixed},
		{"Unknown", Unknown},
	} {
		b.Run("case="+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				tc.term.Bytes()
			}
		})
	}
}

// BenchmarkExists measures checking existence of an unchanged 1 KiB regular file.
// The fixture is fixed; there is no size sweep or measured restoration.
func BenchmarkExists(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		if !Exists(path) {
			b.Fatal("Exists returned false for its fixture")
		}
	}
}

// BenchmarkIsFile measures checking an unchanged 1 KiB regular file.
// The fixture is fixed; there is no size sweep or measured restoration.
func BenchmarkIsFile(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		if !IsFile(path) {
			b.Fatal("IsFile returned false for its fixture")
		}
	}
}

// BenchmarkIsDir measures checking an existing temporary directory.
// The fixture is fixed; there is no size sweep or measured restoration.
func BenchmarkIsDir(b *testing.B) {
	path := b.TempDir()
	b.ReportAllocs()
	for b.Loop() {
		if !IsDir(path) {
			b.Fatal("IsDir returned false for its fixture")
		}
	}
}

// BenchmarkIsSymlink measures checking a symlink to an unchanged 1 KiB file.
// The fixture is fixed; there is no size sweep or measured restoration.
func BenchmarkIsSymlink(b *testing.B) {
	target := benchmarkFile(b, benchkit.Bytes(1024))
	path := filepath.Join(filepath.Dir(target), "link.bin")
	if err := os.Symlink(target, path); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if !IsSymlink(path) {
			b.Fatal("IsSymlink returned false for its fixture")
		}
	}
}

// BenchmarkIsExecutable measures checking an unchanged 1 KiB regular file with mode 0700.
// The fixture is fixed; there is no size sweep or measured restoration.
func BenchmarkIsExecutable(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	if err := os.Chmod(path, 0o700); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if !IsExecutable(path) {
			b.Fatal("IsExecutable returned false for its fixture")
		}
	}
}

// BenchmarkIsEmpty measures checking an unchanged nonempty 1 KiB regular
// file; the fixture is fixed and no restoration is measured.
func BenchmarkIsEmpty(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		empty, err := IsEmpty(path)
		if err != nil {
			b.Fatal(err)
		}
		if empty {
			b.Fatal("fixture unexpectedly empty")
		}
	}
}

// BenchmarkSize measures obtaining the size of an unchanged 1 KiB regular
// file; the fixture is fixed and no restoration is measured.
func BenchmarkSize(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		size, err := Size(path)
		if err != nil {
			b.Fatal(err)
		}
		if size != 1024 {
			b.Fatalf("size = %d, want 1024", size)
		}
	}
}

// BenchmarkModTime measures obtaining the modification time of an unchanged
// 1 KiB regular file with mode 0600; there is no size sweep or restoration.
func BenchmarkModTime(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := ModTime(path); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMode measures obtaining the mode of an unchanged
// 1 KiB regular file with mode 0600; there is no size sweep or restoration.
func BenchmarkMode(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Mode(path); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReadBytes measures reading the entire file as bytes.
// n is the number of source lines, each 40 bytes including LF; the file is
// unchanged and setup is outside the loop.
func BenchmarkReadBytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			data, err := ReadBytes(ctx, path)
			if err != nil {
				b.Fatal(err)
			}
			if len(data) != n*40 {
				b.Fatalf("result = %v, want %v", len(data), n*40)
			}
		}
	})
}

// BenchmarkReadString measures reading the entire file as a string.
// n is the number of source lines, each 40 bytes including LF; the file is
// unchanged and setup is outside the loop.
func BenchmarkReadString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			data, err := ReadString(ctx, path)
			if err != nil {
				b.Fatal(err)
			}
			if len(data) != n*40 {
				b.Fatalf("result = %v, want %v", len(data), n*40)
			}
		}
	})
}

// BenchmarkReadLinesN measures reading at most the first 100 lines.
// n is the number of source lines, each 40 bytes including LF; the file is
// unchanged and setup is outside the loop.
func BenchmarkReadLinesN(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			lines, err := ReadLinesN(ctx, path, 100)
			if err != nil {
				b.Fatal(err)
			}
			if len(lines) != min(n, 100) {
				b.Fatalf("result = %v, want %v", len(lines), min(n, 100))
			}
		}
	})
}

// BenchmarkDetectFileLineTerminator measures detecting LF from at most the first 4096 bytes.
// n is the number of source lines, each 40 bytes including LF; the file is
// unchanged and setup is outside the loop.
func BenchmarkDetectFileLineTerminator(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		path := benchmarkLines(b, n)
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			term, err := DetectFileLineTerminator(ctx, path)
			if err != nil {
				b.Fatal(err)
			}
			if term != LF {
				b.Fatalf("result = %v, want %v", term, LF)
			}
		}
	})
}

// BenchmarkWriteBytes measures overwriting the same existing file on every
// iteration; n counts bytes of deterministic binary input. The API truncates the
// destination, so no separate restoration is measured and state stays bounded.
func BenchmarkWriteBytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Bytes(n)
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		for b.Loop() {
			if err := WriteBytes(path, data, 0o600); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkWriteString measures overwriting the same existing file on every
// iteration; n counts bytes of ASCII text. The API truncates the
// destination, so no separate restoration is measured and state stays bounded.
func BenchmarkWriteString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Text(n, benchkit.ASCII)
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		for b.Loop() {
			if err := WriteString(path, data, 0o600); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkWriteLines measures overwriting the same existing file on every
// iteration; n counts lines of 39 ASCII letters plus LF. The API truncates the
// destination, so no separate restoration is measured and state stays bounded.
func BenchmarkWriteLines(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Strings(n, 39, "abcdefghijklmnopqrstuvwxyz")
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		for b.Loop() {
			if err := WriteLines(path, data, LF, 0o600); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkAppendBytes measures appending n deterministic bytes to an
// initially empty file. Every 1024 appends it truncates the file to zero;
// that restoration is measured, bounding live file size at 1024*n bytes.
func BenchmarkAppendBytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Bytes(n)
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		i := 0
		for b.Loop() {
			if err := AppendBytes(path, data, 0o600); err != nil {
				b.Fatal(err)
			}
			i++
			if i%1024 == 0 {
				if err := os.Truncate(path, 0); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// BenchmarkAppendString measures appending n bytes of ASCII text to an
// initially empty file. Every 1024 appends it truncates the file to zero;
// that restoration is measured, bounding live file size at 1024*n bytes.
func BenchmarkAppendString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Text(n, benchkit.ASCII)
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		i := 0
		for b.Loop() {
			if err := AppendString(path, data, 0o600); err != nil {
				b.Fatal(err)
			}
			i++
			if i%1024 == 0 {
				if err := os.Truncate(path, 0); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// BenchmarkAppendLine measures appending n bytes of ASCII line content plus LF to an
// initially empty file. Every 1024 appends it truncates the file to zero;
// that restoration is measured, bounding live file size at 1024*(n+1) bytes.
func BenchmarkAppendLine(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Text(n, benchkit.ASCII)
		path := benchmarkFile(b, nil)
		b.ReportAllocs()
		i := 0
		for b.Loop() {
			if err := AppendLine(path, data, LF, 0o600); err != nil {
				b.Fatal(err)
			}
			i++
			if i%1024 == 0 {
				if err := os.Truncate(path, 0); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// BenchmarkEnsureDir measures ensuring one directory exists. case=existing
// keeps the same directory. case=new creates and removes the same empty child
// directory every iteration; removal is measured and confirms creation.
func BenchmarkEnsureDir(b *testing.B) {
	b.Run("case=existing", func(b *testing.B) {
		path := b.TempDir()
		b.ReportAllocs()
		for b.Loop() {
			if err := EnsureDir(path, 0o700); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("case=new", func(b *testing.B) {
		path := filepath.Join(b.TempDir(), "child")
		b.ReportAllocs()
		for b.Loop() {
			if err := EnsureDir(path, 0o700); err != nil {
				b.Fatal(err)
			}
			if err := os.Remove(path); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkFind measures matching *.txt against an unchanged tree.
// files=16 and files=256 are fixed filesystem-cardinality overrides to avoid
// 65536-file setup costs. Files are 1 KiB each, half .txt and half .bin, split
// equally between the root and four subdirectories; no restoration is measured.
func BenchmarkFind(b *testing.B) {
	for _, files := range []int{16, 256} {
		b.Run("files="+strconv.Itoa(files), func(b *testing.B) {
			dir := benchmarkTree(b, files)
			ctx := context.Background()
			b.ReportAllocs()
			for b.Loop() {
				matches, err := Find(ctx, dir, "*.txt")
				if err != nil {
					b.Fatal(err)
				}
				if len(matches) != files/2 {
					b.Fatalf("result = %d, want %d", len(matches), files/2)
				}
			}
		})
	}
}

// BenchmarkFindByExtension measures finding .txt files in an unchanged tree.
// files=16 and files=256 are fixed filesystem-cardinality overrides to avoid
// 65536-file setup costs. Files are 1 KiB each, half .txt and half .bin, split
// equally between the root and four subdirectories; no restoration is measured.
func BenchmarkFindByExtension(b *testing.B) {
	for _, files := range []int{16, 256} {
		b.Run("files="+strconv.Itoa(files), func(b *testing.B) {
			dir := benchmarkTree(b, files)
			ctx := context.Background()
			b.ReportAllocs()
			for b.Loop() {
				matches, err := FindByExtension(ctx, dir, ".txt")
				if err != nil {
					b.Fatal(err)
				}
				if len(matches) != files/2 {
					b.Fatalf("result = %d, want %d", len(matches), files/2)
				}
			}
		})
	}
}

// BenchmarkDirSize measures summing file sizes in an unchanged tree.
// files=16 and files=256 are fixed filesystem-cardinality overrides to avoid
// 65536-file setup costs. Files are 1 KiB each, half .txt and half .bin, split
// equally between the root and four subdirectories; no restoration is measured.
func BenchmarkDirSize(b *testing.B) {
	for _, files := range []int{16, 256} {
		b.Run("files="+strconv.Itoa(files), func(b *testing.B) {
			dir := benchmarkTree(b, files)
			ctx := context.Background()
			b.ReportAllocs()
			for b.Loop() {
				size, err := DirSize(ctx, dir)
				if err != nil {
					b.Fatal(err)
				}
				if size != int64(files)*1024 {
					b.Fatalf("result = %d, want %d", size, int64(files)*1024)
				}
			}
		})
	}
}

// BenchmarkCopy measures copying n bytes (1 KiB or 1 MiB) from an unchanged
// source over the same existing destination in one directory. The destination
// is replaced each iteration, so no separate restoration is measured.
func BenchmarkCopy(b *testing.B) {
	benchkit.Run(b, []int{1024, 1024 * 1024}, func(b *testing.B, n int) {
		data := benchkit.Bytes(n)
		src := benchmarkFile(b, data)
		dst := filepath.Join(filepath.Dir(src), "destination.bin")
		if err := os.WriteFile(dst, data, 0o600); err != nil {
			b.Fatal(err)
		}
		ctx := context.Background()
		b.ReportAllocs()
		for b.Loop() {
			if err := Copy(ctx, src, dst); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMove measures renaming a 1 KiB file inside one directory. The file
// alternates between two names, so every iteration moves an existing file to
// a free name.
func BenchmarkMove(b *testing.B) {
	ctx := context.Background()
	dir := b.TempDir()
	paths := [2]string{filepath.Join(dir, "a.bin"), filepath.Join(dir, "b.bin")}
	if err := os.WriteFile(paths[0], benchkit.Bytes(1024), 0o600); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		if err := Move(ctx, paths[i%2], paths[(i+1)%2]); err != nil {
			b.Fatal(err)
		}
		i++
	}
}

// BenchmarkTouch measures updating mtime on the same existing 1 KiB file;
// the fixed file remains present, with no size sweep or restoration.
func BenchmarkTouch(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	b.ReportAllocs()
	for b.Loop() {
		if err := Touch(path); err != nil {
			b.Fatal(err)
		}
	}
}

// benchmarkCRLF returns exactly n bytes of deterministic ASCII text with
// a CRLF pair at the end of every eight bytes.
func benchmarkCRLF(n int) []byte {
	data := []byte(benchkit.Text(n, benchkit.ASCII))
	for i := 6; i+1 < len(data); i += 8 {
		data[i], data[i+1] = '\r', '\n'
	}

	return data
}

// BenchmarkDetectLineTerminator measures detecting CRLF in n bytes of
// unchanged ASCII text with one CRLF pair every eight bytes.
func BenchmarkDetectLineTerminator(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchmarkCRLF(n)
		b.ReportAllocs()
		for b.Loop() {
			if term := DetectLineTerminator(data); term != CRLF {
				b.Fatalf("terminator = %v, want CRLF", term)
			}
		}
	})
}

// BenchmarkNormalizeLineTerminators measures converting CRLF to LF in n
// bytes of ASCII text with one CRLF pair every eight bytes. The input remains
// unchanged, so no restoration is measured.
func BenchmarkNormalizeLineTerminators(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchmarkCRLF(n)
		b.ReportAllocs()
		for b.Loop() {
			NormalizeLineTerminators(data, LF)
		}
	})
}

// BenchmarkExtension measures extracting .txt from a fixed four-segment path.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkExtension(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Extension(path)
	}
}

// BenchmarkBaseName measures extracting the stem of a fixed four-segment path.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkBaseName(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		BaseName(path)
	}
}

// BenchmarkAbs measures resolving a fixed relative four-segment path.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkAbs(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Abs(path)
	}
}

// BenchmarkClean measures cleaning a fixed path containing . and .. segments.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkClean(b *testing.B) {
	path := "workspace/reports/./drafts/../annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Clean(path)
	}
}

// BenchmarkDir measures extracting the directory of a fixed four-segment path.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkDir(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Dir(path)
	}
}

// BenchmarkBase measures extracting the file name of a fixed four-segment path.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkBase(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Base(path)
	}
}

// BenchmarkSplit measures splitting a fixed four-segment path into directory and file.
// The immutable path is a scalar input; no filesystem fixture or restoration.
func BenchmarkSplit(b *testing.B) {
	path := "workspace/reports/annual/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		Split(path)
	}
}

// BenchmarkJoin measures joining four fixed path segments; the immutable
// scalar input needs no filesystem fixture or restoration.
func BenchmarkJoin(b *testing.B) {
	parts := []string{"workspace", "reports", "annual", "report.txt"}
	b.ReportAllocs()
	for b.Loop() {
		Join(parts...)
	}
}

// BenchmarkRel measures computing a path between sibling directories. The
// fixed target contains .. and the relative result starts with ..; both
// inputs are immutable and no filesystem fixture or restoration is needed.
func BenchmarkRel(b *testing.B) {
	base := "workspace/reports/annual"
	target := "workspace/reports/drafts/../archive/report.txt"
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Rel(base, target); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTempFile measures creating one temporary file inside
// b.TempDir and removing it each iteration. Removal is measured, checks the
// returned path exists, and bounds live resources independently of iteration count.
func BenchmarkTempFile(b *testing.B) {
	dir := b.TempDir()
	b.ReportAllocs()
	for b.Loop() {
		path, err := TempFile(dir, "entry-*")
		if err != nil {
			b.Fatal(err)
		}
		if err := os.Remove(path); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTempDir measures creating one temporary empty directory inside
// b.TempDir and removing it each iteration. Removal is measured, checks the
// returned path exists, and bounds live resources independently of iteration count.
func BenchmarkTempDir(b *testing.B) {
	dir := b.TempDir()
	b.ReportAllocs()
	for b.Loop() {
		path, err := TempDir(dir, "entry-*")
		if err != nil {
			b.Fatal(err)
		}
		if err := os.Remove(path); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSameFile measures comparing two paths to the same unchanged
// 1 KiB file: its regular path and a symlink. The fixture is fixed; no restoration.
func BenchmarkSameFile(b *testing.B) {
	path := benchmarkFile(b, benchkit.Bytes(1024))
	link := filepath.Join(filepath.Dir(path), "link.bin")
	if err := os.Symlink(path, link); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if !SameFile(path, link) {
			b.Fatal("paths no longer refer to the same file")
		}
	}
}
