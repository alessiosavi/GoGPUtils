---
title: fileutil
parent: Packages
nav_order: 5
---

# fileutil

File system operations with proper error handling.
{: .fs-6 .fw-300 }

## Overview

The `fileutil` package provides file system operations with error returns for fallible operations, `context.Context` support for reads and transfers, and automatic file handle cleanup. Writers report close errors when writing succeeds.

```go
import "github.com/alessiosavi/GoGPUtils/fileutil"
```

## Design Philosophy

This package adheres to the following principles:

- **Errors over panics**: All functions return errors instead of panicking, allowing callers to handle failures gracefully
- **Context-aware**: Long-running operations accept `context.Context` for cancellation and timeout support
- **Resource safety**: File handles are closed, and writers check close errors
- **Zero external dependencies**: Uses only the Go standard library

---

## Types

### LineTerminator

Represents different line ending styles.

```go
type LineTerminator int
```

Constants:

| Constant  | Description                         |
| --------- | ----------------------------------- |
| `LF`      | Unix-style line endings (`\n`)      |
| `CRLF`    | Windows-style line endings (`\r\n`) |
| `CR`      | Old Mac-style line endings (`\r`)   |
| `Mixed`   | Inconsistent line endings detected  |
| `Unknown` | No line endings found               |

Methods:

```go
func (lt LineTerminator) String() string
func (lt LineTerminator) Bytes() []byte
```

### ListOption

Configures `List` behavior. Use as bit flags.

```go
type ListOption int
```

Constants:

| Constant        | Description                              |
| --------------- | ---------------------------------------- |
| `FilesOnly`     | List only files (not directories)        |
| `DirsOnly`      | List only directories (not files)        |
| `Recursive`     | List contents recursively                |
| `IncludeHidden` | Include hidden files (starting with `.`) |

---

## Common Errors

```go
var (
    ErrNotFile      = errors.New("path is not a file")
    ErrNotDir       = errors.New("path is not a directory")
    ErrNotExist     = errors.New("path does not exist")
    ErrReadCanceled = errors.New("read operation canceled")
    ErrSameFile     = errors.New("source and destination are the same file")
)
```

---

## Existence and Type Checks

### Exists

Reports whether the path exists.

```go
func Exists(path string) bool
```

### IsFile

Reports whether the path is a regular file.

```go
func IsFile(path string) bool
```

### IsDir

Reports whether the path is a directory.

```go
func IsDir(path string) bool
```

### IsSymlink

Reports whether the path is a symbolic link.

```go
func IsSymlink(path string) bool
```

### IsExecutable

Reports whether the path is executable.

```go
func IsExecutable(path string) bool
```

### IsEmpty

Reports whether the file or directory is empty. For files, returns true if size is 0. For directories, returns true if no entries exist.

```go
func IsEmpty(path string) (bool, error)
```

---

## File Information

### Size

Returns the size of the file in bytes. Returns `ErrNotFile` if the path is a directory.

```go
func Size(path string) (int64, error)
```

### ModTime

Returns the modification time of the file.

```go
func ModTime(path string) (time.Time, error)
```

### Mode

Returns the file mode of the path.

```go
func Mode(path string) (fs.FileMode, error)
```

### Extension

Returns the file extension including the dot. Returns empty string if no extension.

```go
func Extension(path string) string
```

### BaseName

Returns the file name without extension.

```go
func BaseName(path string) string
```

---

## Reading Files

### ReadBytes

Reads the entire file and returns its contents as bytes. Respects context cancellation before starting the read. The read itself is not interruptible by context cancellation; this also applies to `ReadString`.

```go
func ReadBytes(ctx context.Context, path string) ([]byte, error)
```

### ReadString

Reads the entire file and returns its contents as a string.

```go
func ReadString(ctx context.Context, path string) (string, error)
```

### ReadLines

Reads the file and returns its contents as a slice of lines, with no fixed line length limit. Strips the ending newline and one trailing carriage return. Includes a final unterminated line; an empty file returns a nil slice. Checks context cancellation before opening and between reads, including while reading a long line. An operating-system read already in progress cannot be interrupted by the context.

```go
func ReadLines(ctx context.Context, path string) ([]string, error)
```

### ReadLinesN

Reads the first n lines from a file. If n <= 0, reads all lines. Line lengths, terminators, and cancellation behave as in `ReadLines`. For an empty file, a positive limit returns an empty non-nil slice.

```go
func ReadLinesN(ctx context.Context, path string, n int) ([]string, error)
```

### CountLines

Counts the number of lines in a file, including a final unterminated line. Returns the same count as `len(ReadLines(...))`: `"a"` has one line, `"a\nb"` has two, `"a\n"` has one, and an empty file has zero. Uses a fixed-size buffer for better memory efficiency than `ReadLines` on large files.

```go
func CountLines(ctx context.Context, path string) (int, error)
```

---

## Writing Files

### WriteBytes

Writes data to a file, creating it if necessary. Truncates existing files.

```go
func WriteBytes(path string, data []byte, perm fs.FileMode) error
```

### WriteString

Writes a string to a file.

```go
func WriteString(path, content string, perm fs.FileMode) error
```

### WriteLines

Writes lines to a file, joining with the specified terminator.

```go
func WriteLines(path string, lines []string, terminator LineTerminator, perm fs.FileMode) error
```

### AppendBytes

Appends data to a file, creating it if necessary. Returns a close error when the write succeeded; if writing fails, that error takes precedence. This also applies to `AppendString` and `AppendLine`.

```go
func AppendBytes(path string, data []byte, perm fs.FileMode) error
```

### AppendString

Appends a string to a file.

```go
func AppendString(path, content string, perm fs.FileMode) error
```

### AppendLine

Appends a line to a file with the specified terminator.

```go
func AppendLine(path, line string, terminator LineTerminator, perm fs.FileMode) error
```

---

## Directory Operations

### EnsureDir

Creates a directory if it doesn't exist. Creates parent directories as needed (like `mkdir -p`).

```go
func EnsureDir(path string, perm fs.FileMode) error
```

### List

Returns entries in a directory. Use `ListOption` flags to control behavior.

```go
func List(ctx context.Context, dir string, opts ListOption) ([]string, error)
```

### Find

Returns files matching the glob pattern. Pattern syntax follows `filepath.Match`.

```go
func Find(ctx context.Context, dir, pattern string) ([]string, error)
```

### FindByExtension

Returns files with the specified extension. Extension should include the dot (e.g., `".go"`).

```go
func FindByExtension(ctx context.Context, dir, ext string) ([]string, error)
```

### DirSize

Calculates the total size of a directory and its contents.

```go
func DirSize(ctx context.Context, dir string) (int64, error)
```

---

## File Operations

### Copy

Copies a file from src to dst, creating the destination directory if needed. Returns `ErrSameFile` without changing the file when the paths refer to the same file, including hard links and symbolic links.

Writes to a temporary file in the destination directory, applies the source mode, checks the close error, and renames the completed file over dst. Errors before the rename leave any existing destination intact. On any failure the temporary file is removed; if that removal fails, its error is joined to the returned error, which still matches the original cause with `errors.Is`. The rename has the platform guarantees of `os.Rename`; no crash-durability guarantee is made.

Replacement changes the destination's file identity: a destination symbolic link is replaced, and other hard links to the old destination retain their contents. Replacement requires write permission on the destination directory, and existing destination metadata is replaced by the new file's metadata, with mode copied from the source.

Checks context cancellation before starting, between reads, and before the final rename. Cancellation cannot interrupt a read or write already blocked in the operating system, or undo a rename already in progress. A destination directory created before cancellation can remain.

```go
func Copy(ctx context.Context, src, dst string) error
```

### Move

Moves a file from src to dst. Checks context cancellation before any side effect, creates the destination directory if needed, then attempts rename. Falls back to `Copy` followed by source removal only for a cross-device (`EXDEV`, or `ERROR_NOT_SAME_DEVICE` on Windows) rename error. Other rename errors are returned unchanged without overwriting the destination. On Plan 9, rename cannot move files across directories, so every rename error triggers the legacy copy fallback, including unrelated errors.

On Windows, renaming over a destination held open by another process can fail with an access-denied error. `Move` returns that error without copying; callers should retry or close the destination first.

If copying succeeds but source removal fails, both files remain. The returned error describes the successful copy and wraps the removal error for inspection with `errors.Is` or `errors.As`.

```go
func Move(ctx context.Context, src, dst string) error
```

### Touch

Creates an empty file or updates its modification time.

```go
func Touch(path string) error
```

---

## Line Terminator Operations

### DetectLineTerminator

Detects the line terminator style in data. Returns `Unknown` if no line terminators are found. Returns `Mixed` if multiple styles are detected.

```go
func DetectLineTerminator(data []byte) LineTerminator
```

### NormalizeLineTerminators

Converts all line terminators to the specified style.

```go
func NormalizeLineTerminators(data []byte, target LineTerminator) []byte
```

### DetectFileLineTerminator

Detects line terminator style in a file. Checks context cancellation before opening and reads only the beginning of the file for efficiency.

```go
func DetectFileLineTerminator(ctx context.Context, path string) (LineTerminator, error)
```

---

## Path Operations

### Abs

Returns the absolute path. Returns the path unchanged if it's already absolute or on error.

```go
func Abs(path string) string
```

### Clean

Returns the cleaned path.

```go
func Clean(path string) string
```

### Join

Joins path elements.

```go
func Join(elem ...string) string
```

### Dir

Returns the directory component of the path.

```go
func Dir(path string) string
```

### Base

Returns the last element of the path.

```go
func Base(path string) string
```

### Rel

Returns a relative path from base to target.

```go
func Rel(base, target string) (string, error)
```

### Split

Splits a path into directory and file components.

```go
func Split(path string) (dir, file string)
```

---

## Utility Functions

### TempFile

Creates a temporary file and returns its path. The caller is responsible for removing the file. If closing fails, removes the temporary file and returns the close error, joined with the removal error if the removal also fails.

```go
func TempFile(dir, pattern string) (string, error)
```

### TempDir

Creates a temporary directory and returns its path. The caller is responsible for removing the directory.

```go
func TempDir(dir, pattern string) (string, error)
```

### SameFile

Reports whether two paths describe the same file.

```go
func SameFile(path1, path2 string) bool
```

---

## Usage Examples

### Context-Aware File Reading

Read operations check context cancellation before opening. Line readers also check between reads; `ReadBytes` and `ReadString` cannot interrupt a read once started:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/alessiosavi/GoGPUtils/fileutil"
)

func main() {
    // Read with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    lines, err := fileutil.ReadLines(ctx, "large-file.txt")
    if err != nil {
        if err == context.DeadlineExceeded {
            fmt.Println("Read timed out")
            return
        }
        // Handle other errors
        return
    }

    for _, line := range lines {
        fmt.Println(line)
    }
}
```

### File Existence and Type Checking

```go
path := "/path/to/file.txt"

if fileutil.Exists(path) {
    fmt.Println("Path exists")
}

if fileutil.IsFile(path) {
    fmt.Println("Is a regular file")
}

if fileutil.IsDir(path) {
    fmt.Println("Is a directory")
}

if fileutil.IsExecutable(path) {
    fmt.Println("Is executable")
}
```

### Reading Files

```go
ctx := context.Background()

// Read as bytes
data, err := fileutil.ReadBytes(ctx, "data.bin")
if err != nil {
    return err
}

// Read as string
content, err := fileutil.ReadString(ctx, "config.txt")
if err != nil {
    return err
}

// Read as lines
lines, err := fileutil.ReadLines(ctx, "lines.txt")
if err != nil {
    return err
}

// Read first 10 lines
head, err := fileutil.ReadLinesN(ctx, "large.txt", 10)
if err != nil {
    return err
}

// Count lines efficiently
count, err := fileutil.CountLines(ctx, "huge.txt")
if err != nil {
    return err
}
fmt.Printf("File has %d lines\n", count)
```

### Writing Files

```go
// Write bytes
err := fileutil.WriteBytes("output.bin", []byte{0x01, 0x02}, 0644)

// Write string
err := fileutil.WriteString("output.txt", "Hello, World!", 0644)

// Write lines with specific line endings
lines := []string{"line1", "line2", "line3"}
err := fileutil.WriteLines("output.txt", lines, fileutil.LF, 0644)

// Append to file
err := fileutil.AppendString("log.txt", "New log entry\n", 0644)

// Append a line
err := fileutil.AppendLine("log.txt", "Another entry", fileutil.LF, 0644)
```

### Directory Operations

```go
ctx := context.Background()

// Create directory tree
err := fileutil.EnsureDir("path/to/nested/dir", 0755)

// List all entries
entries, err := fileutil.List(ctx, "/path/to/dir", 0)

// List only files recursively
files, err := fileutil.List(ctx, "/path/to/dir", fileutil.FilesOnly|fileutil.Recursive)

// Find Go files
goFiles, err := fileutil.Find(ctx, "/project", "*.go")

// Find by extension
txtFiles, err := fileutil.FindByExtension(ctx, "/project", ".txt")

// Get directory size
totalSize, err := fileutil.DirSize(ctx, "/path/to/dir")
```

### File Operations

```go
ctx := context.Background()

// Copy file
err := fileutil.Copy(ctx, "source.txt", "dest.txt")

// Copy to new directory (creates dirs automatically)
err := fileutil.Copy(ctx, "source.txt", "new/dir/dest.txt")

// Move file
err := fileutil.Move(ctx, "old.txt", "new.txt")

// Create or update file timestamp
err := fileutil.Touch("file.txt")
```

### Line Terminator Handling

```go
// Detect line endings in data
data := []byte("line1\r\nline2\r\n")
terminator := fileutil.DetectLineTerminator(data)
fmt.Println(terminator.String()) // "CRLF (\r\n)"

// Normalize to LF
normalized := fileutil.NormalizeLineTerminators(data, fileutil.LF)

// Detect in file
term, err := fileutil.DetectFileLineTerminator(ctx, "file.txt")
if err != nil {
    return err
}
```

### Path Manipulation

```go
// Get absolute path
abs := fileutil.Abs("relative/path")

// Clean path
clean := fileutil.Clean("/path/./to/../file") // "/path/file"

// Join paths
full := fileutil.Join("path", "to", "file") // "path/to/file"

// Get directory and base
dir := fileutil.Dir("/path/to/file.txt")   // "/path/to"
base := fileutil.Base("/path/to/file.txt") // "file.txt"

// Get extension and basename
ext := fileutil.Extension("file.txt")    // ".txt"
name := fileutil.BaseName("file.txt")    // "file"

// Split path
d, f := fileutil.Split("/path/to/file.txt") // d="/path/to/", f="file.txt"

// Relative path
rel, err := fileutil.Rel("/base", "/base/sub/file.txt") // "sub/file.txt"
```

### Temporary Files

```go
// Create temp file
tmpFile, err := fileutil.TempFile("", "prefix-*.txt")
if err != nil {
    return err
}
defer os.Remove(tmpFile) // Clean up

// Create temp directory
tmpDir, err := fileutil.TempDir("", "prefix-*")
if err != nil {
    return err
}
defer os.RemoveAll(tmpDir) // Clean up
```

---

## Error Handling Philosophy

The `fileutil` package follows the Go idiom of **errors over panics**. Every function that can fail returns an error as its last return value, allowing callers to make informed decisions about how to handle failures.

### Why Errors Over Panics?

1. **Caller control**: Callers decide whether a failure is fatal or recoverable
2. **Composability**: Functions can be chained and errors propagated up the call stack
3. **Testability**: Error paths can be tested explicitly
4. **Predictability**: No unexpected program termination

### Common Error Types

| Error                      | When Returned                                    |
| -------------------------- | ------------------------------------------------ |
| `ErrNotFile`               | Operation expects a file but path is a directory |
| `ErrNotDir`                | Operation expects a directory but path is a file |
| `ErrNotExist`              | Path does not exist (where applicable)           |
| `ErrSameFile`              | Copy source and destination identify the same file |
| `context.Canceled`         | Operation was canceled via context               |
| `context.DeadlineExceeded` | Operation timed out via context                  |

### Best Practices

```go
// Always check errors
content, err := fileutil.ReadString(ctx, "file.txt")
if err != nil {
    // Handle error - don't ignore it
    return fmt.Errorf("reading file: %w", err)
}

// Use errors.Is for specific error checking
size, err := fileutil.Size("somedir")
if errors.Is(err, fileutil.ErrNotFile) {
    fmt.Println("Expected a file, got a directory")
}

// Wrap errors with context
if err := fileutil.Copy(ctx, src, dst); err != nil {
    return fmt.Errorf("copying %s to %s: %w", src, dst, err)
}
```

---

## Performance Considerations

- `CountLines` is more memory-efficient than `ReadLines` for large files as it doesn't allocate a slice for all lines
- `ReadLinesN` is efficient for reading just the beginning of large files
- `List` with `Recursive` uses `filepath.WalkDir` which is efficient for deep directory trees
- `Copy` preserves file permissions from the source file
- `Move` attempts a fast rename first and falls back to copy+delete only across filesystems (or after any rename error on Plan 9)
