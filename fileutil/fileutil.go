package fileutil

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LineTerminator represents different line ending styles.
type LineTerminator int

const (
	// LF represents Unix-style line endings (\n).
	LF LineTerminator = iota
	// CRLF represents Windows-style line endings (\r\n).
	CRLF
	// CR represents old Mac-style line endings (\r).
	CR
	// Mixed indicates the file has inconsistent line endings.
	Mixed
	// Unknown indicates no line endings were found.
	Unknown
)

// String returns the string representation of the line terminator.
func (lt LineTerminator) String() string {
	switch lt {
	case LF:
		return "LF (\\n)"
	case CRLF:
		return "CRLF (\\r\\n)"
	case CR:
		return "CR (\\r)"
	case Mixed:
		return "Mixed"
	default:
		return "Unknown"
	}
}

// Bytes returns the byte sequence for the line terminator.
func (lt LineTerminator) Bytes() []byte {
	switch lt {
	case LF:
		return []byte{'\n'}
	case CRLF:
		return []byte{'\r', '\n'}
	case CR:
		return []byte{'\r'}
	default:
		return []byte{'\n'}
	}
}

// ListOption configures List behavior.
type ListOption int

const (
	// FilesOnly lists only files (not directories).
	FilesOnly ListOption = 1 << iota
	// DirsOnly lists only directories (not files).
	DirsOnly
	// Recursive lists contents recursively.
	Recursive
	// IncludeHidden includes hidden files (starting with .)
	IncludeHidden
)

// Common errors.
var (
	ErrNotFile      = errors.New("path is not a file")
	ErrNotDir       = errors.New("path is not a directory")
	ErrNotExist     = errors.New("path does not exist")
	ErrReadCanceled = errors.New("read operation canceled")
	ErrSameFile     = errors.New("source and destination are the same file")
)

// ============================================================================
// Existence and Type Checks
// ============================================================================

// Exists reports whether the path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// IsFile reports whether the path is a regular file.
func IsFile(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.Mode().IsRegular()
}

// IsDir reports whether the path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// IsSymlink reports whether the path is a symbolic link.
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)

	return err == nil && info.Mode()&os.ModeSymlink != 0
}

// IsExecutable reports whether the path is executable.
func IsExecutable(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.Mode()&0111 != 0
}

// IsEmpty reports whether the file or directory is empty.
// For files, returns true if size is 0.
// For directories, returns true if no entries exist.
func IsEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}
		defer f.Close()

		_, err = f.Readdirnames(1)
		if errors.Is(err, io.EOF) {
			return true, nil
		}

		return false, err
	}

	return info.Size() == 0, nil
}

// ============================================================================
// File Information
// ============================================================================

// Size returns the size of the file in bytes.
func Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		return 0, ErrNotFile
	}

	return info.Size(), nil
}

// ModTime returns the modification time of the file.
func ModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return info.ModTime(), nil
}

// Mode returns the file mode of the path.
func Mode(path string) (fs.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}

// Extension returns the file extension including the dot.
// Returns empty string if no extension.
func Extension(path string) string {
	return filepath.Ext(path)
}

// BaseName returns the file name without extension.
func BaseName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)

	return strings.TrimSuffix(base, ext)
}

// ============================================================================
// Reading Files
// ============================================================================

// ReadBytes reads the entire file and returns its contents as bytes.
// Respects context cancellation before starting the read.
// The read itself is not interruptible by context cancellation.
func ReadBytes(ctx context.Context, path string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return os.ReadFile(path)
}

// ReadString reads the entire file and returns its contents as a string.
func ReadString(ctx context.Context, path string) (string, error) {
	data, err := ReadBytes(ctx, path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadLines reads the file and returns its contents as a slice of lines.
// Lines have no fixed length limit. The ending newline and one trailing carriage
// return are stripped; a final unterminated line is included.
// Checks context cancellation before opening and between reads.
func ReadLines(ctx context.Context, path string) ([]string, error) {
	return readLines(ctx, path, 0)
}

// ReadLinesN reads the first n lines from a file.
// If n <= 0, reads all lines.
// Line lengths, terminators, and cancellation behave as in ReadLines.
func ReadLinesN(ctx context.Context, path string, n int) ([]string, error) {
	return readLines(ctx, path, n)
}

func readLines(ctx context.Context, path string, n int) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	if n > 0 {
		lines = make([]string, 0, n)
	}

	reader := bufio.NewReader(contextReader{ctx: ctx, reader: f})
	for n <= 0 || len(lines) < n {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}

		if len(line) > 0 {
			line = strings.TrimSuffix(line, "\n")
			lines = append(lines, strings.TrimSuffix(line, "\r"))
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return lines, nil
}

// CountLines counts the number of lines in a file.
// More memory-efficient than ReadLines for large files.
// Includes a final unterminated line, so the count equals len(ReadLines(...)).
func CountLines(ctx context.Context, path string) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count := 0
	buf := make([]byte, 32*1024) // 32KB buffer
	lineSep := []byte{'\n'}
	last := byte('\n')

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		n, err := f.Read(buf)
		count += bytes.Count(buf[:n], lineSep)
		if n > 0 {
			last = buf[n-1]
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	if last != '\n' {
		count++
	}

	return count, nil
}

// contextReader checks cancellation between reads; it cannot interrupt a read
// already blocked in the underlying reader.
type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}

	return r.reader.Read(p)
}

// ============================================================================
// Writing Files
// ============================================================================

// WriteBytes writes data to a file, creating it if necessary.
// Truncates existing files.
func WriteBytes(path string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// WriteString writes a string to a file.
func WriteString(path, content string, perm fs.FileMode) error {
	return WriteBytes(path, []byte(content), perm)
}

// WriteLines writes lines to a file, joining with the specified terminator.
func WriteLines(path string, lines []string, terminator LineTerminator, perm fs.FileMode) error {
	term := string(terminator.Bytes())

	content := strings.Join(lines, term)

	if len(lines) > 0 {
		content += term
	}

	return WriteString(path, content, perm)
}

// AppendBytes appends data to a file, creating it if necessary.
// Returns a close error if the write succeeded.
func AppendBytes(path string, data []byte, perm fs.FileMode) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(data)

	return finishWrite(f, err)
}

func finishWrite(f io.Closer, err error) error {
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return err
}

// AppendString appends a string to a file.
func AppendString(path, content string, perm fs.FileMode) error {
	return AppendBytes(path, []byte(content), perm)
}

// AppendLine appends a line to a file with the specified terminator.
func AppendLine(path, line string, terminator LineTerminator, perm fs.FileMode) error {
	return AppendString(path, line+string(terminator.Bytes()), perm)
}

// ============================================================================
// Directory Operations
// ============================================================================

// EnsureDir creates a directory if it doesn't exist.
// Creates parent directories as needed (like mkdir -p).
func EnsureDir(path string, perm fs.FileMode) error {
	if Exists(path) {
		if !IsDir(path) {
			return ErrNotDir
		}

		return nil
	}

	return os.MkdirAll(path, perm)
}

// List returns entries in a directory.
// Use ListOption flags to control behavior.
func List(ctx context.Context, dir string, opts ListOption) ([]string, error) {
	if !IsDir(dir) {
		return nil, ErrNotDir
	}

	recursive := opts&Recursive != 0
	filesOnly := opts&FilesOnly != 0
	dirsOnly := opts&DirsOnly != 0
	includeHidden := opts&IncludeHidden != 0

	var entries []string

	walkFn := func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		// Skip root directory
		if path == dir {
			return nil
		}

		// Handle hidden files
		name := d.Name()
		if !includeHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		// Filter by type
		if filesOnly && d.IsDir() {
			if recursive {
				return nil // Continue into directory but don't include it
			}

			return filepath.SkipDir
		}

		if dirsOnly && !d.IsDir() {
			return nil
		}

		entries = append(entries, path)

		// Don't recurse if not requested
		if !recursive && d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	}

	if recursive {
		err := filepath.WalkDir(dir, walkFn)
		if err != nil {
			return nil, err
		}
	} else {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range dirEntries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			name := entry.Name()
			if !includeHidden && strings.HasPrefix(name, ".") {
				continue
			}

			if filesOnly && entry.IsDir() {
				continue
			}

			if dirsOnly && !entry.IsDir() {
				continue
			}

			entries = append(entries, filepath.Join(dir, name))
		}
	}

	return entries, nil
}

// Find returns files matching the glob pattern.
// Pattern syntax follows filepath.Match.
func Find(ctx context.Context, dir, pattern string) ([]string, error) {
	if !IsDir(dir) {
		return nil, ErrNotDir
	}

	var matches []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, d.Name())
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// FindByExtension returns files with the specified extension.
// Extension should include the dot (e.g., ".go").
func FindByExtension(ctx context.Context, dir, ext string) ([]string, error) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return Find(ctx, dir, "*"+ext)
}

// ============================================================================
// File Operations
// ============================================================================

// Copy copies a file from src to dst.
// Creates the destination directory if needed and preserves the source mode.
// Returns ErrSameFile without changing either file if the paths refer to the
// same file, including through hard links or symbolic links.
// Copies to a temporary file in the destination directory and closes it before
// renaming it over dst. Errors before the rename leave any existing dst intact.
// Replacement changes dst's file identity: a destination symlink is replaced,
// and other hard links to the previous destination keep their original contents.
// Checks context cancellation between reads and before the final rename, but
// cannot interrupt a read or write already blocked in the operating system.
// Rename atomicity follows os.Rename's platform guarantees.
//
// Example:
//
//	err := Copy(ctx, "source.txt", "backup/source.txt")
func Copy(ctx context.Context, src, dst string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return ErrNotFile
	}

	dstInfo, err := os.Stat(dst)
	if err == nil && os.SameFile(srcInfo, dstInfo) {
		return ErrSameFile
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := EnsureDir(dstDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.CreateTemp(dstDir, ".fileutil-copy-*")
	if err != nil {
		return err
	}
	defer os.Remove(dstFile.Name())

	_, err = io.Copy(dstFile, contextReader{ctx: ctx, reader: srcFile})
	if err == nil {
		err = dstFile.Chmod(srcInfo.Mode())
	}

	if err := finishWrite(dstFile, err); err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return os.Rename(dstFile.Name(), dst)
}

// Move moves a file from src to dst.
// Checks context cancellation before any side effect and creates the destination
// directory if needed. Falls back to Copy followed by removal only when rename
// fails across filesystems. Other rename errors are returned unchanged.
// On Plan 9, every rename error triggers the fallback because rename cannot
// move files across directories, even within a filesystem.
// On Windows, renaming over a destination held open by another process can fail
// with an access-denied error. Move returns that error without copying; callers
// should retry or close the destination first.
// If source removal fails after a successful copy, both files remain and the
// returned error wraps the removal error and describes the completed copy.
//
// Example:
//
//	err := Move(ctx, "source.txt", "archive/source.txt")
func Move(ctx context.Context, src, dst string) error {
	return moveFile(ctx, src, dst, os.Rename, os.Remove)
}

func moveFile(ctx context.Context, src, dst string, rename func(string, string) error, remove func(string) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := EnsureDir(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	// Try rename first (fast path for same filesystem)
	err := rename(src, dst)
	if err == nil {
		return nil
	}

	if !isCrossDeviceError(err) {
		return err
	}

	// Fall back to copy + delete
	if err := Copy(ctx, src, dst); err != nil {
		return err
	}

	if err := remove(src); err != nil {
		return fmt.Errorf("copied %q to %q, but failed to remove source: %w", src, dst, err)
	}

	return nil
}

// Touch creates an empty file or updates its modification time.
func Touch(path string) error {
	if Exists(path) {
		now := time.Now()

		return os.Chtimes(path, now, now)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := EnsureDir(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return f.Close()
}

// ============================================================================
// Line Terminator Operations
// ============================================================================

// DetectLineTerminator detects the line terminator style in data.
// Returns Unknown if no line terminators are found.
// Returns Mixed if multiple styles are detected.
func DetectLineTerminator(data []byte) LineTerminator {
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

// NormalizeLineTerminators converts all line terminators to the specified style.
func NormalizeLineTerminators(data []byte, target LineTerminator) []byte {
	if len(data) == 0 {
		return nil
	}

	term := "\n"
	switch target {
	case CRLF:
		term = "\r\n"
	case CR:
		term = "\r"
	}

	// LF output without CR only needs a detached copy.
	cr := bytes.Count(data, []byte{'\r'})
	if cr == 0 && term == "\n" {
		return append([]byte(nil), data...)
	}

	lf := bytes.Count(data, []byte{'\n'})
	crlf := 0
	firstCR := bytes.IndexByte(data, '\r')
	for i := firstCR; i >= 0; {
		if i+1 < len(data) && data[i+1] == '\n' {
			crlf++
		}
		next := bytes.IndexByte(data[i+1:], '\r')
		if next < 0 {
			break
		}
		i += next + 1
	}

	// CR and LF occupy distinct bytes, so their sum cannot exceed len(data).
	breaks := cr + lf - crlf
	if breaks == 0 {
		return append([]byte(nil), data...)
	}
	const maxInt = int(^uint(0) >> 1)
	size := len(data) - cr - lf
	if breaks > (maxInt-size)/len(term) {
		return normalizeLineTerminatorsOverflow(data, target)
	}
	size += breaks * len(term)

	result := make([]byte, 0, size)
	pos := 0
	nextCR, nextLF := firstCR, bytes.IndexByte(data, '\n')
	for nextCR >= 0 || nextLF >= 0 {
		next := nextLF
		if nextCR >= 0 && (nextLF < 0 || nextCR < nextLF) {
			next = nextCR
		}
		result = append(result, data[pos:next]...)
		result = append(result, term...)
		pos = next + 1
		if next == nextCR && pos < len(data) && data[pos] == '\n' {
			pos++
		}

		// Refresh only consumed positions; retain absent bytes as -1 so a
		// missing or distant terminator never causes repeated suffix scans.
		if nextCR >= 0 && nextCR < pos {
			nextCR = bytes.IndexByte(data[pos:], '\r')
			if nextCR >= 0 {
				nextCR += pos
			}
		}
		if nextLF >= 0 && nextLF < pos {
			nextLF = bytes.IndexByte(data[pos:], '\n')
			if nextLF >= 0 {
				nextLF += pos
			}
		}
	}

	return append(result, data[pos:]...)
}

// Preserve BASE's allocation and panic behavior if the output length overflows.
func normalizeLineTerminatorsOverflow(data []byte, target LineTerminator) []byte {
	// First normalize to LF
	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\n'})
	data = bytes.ReplaceAll(data, []byte{'\r'}, []byte{'\n'})

	// Then convert to target
	if target == LF {
		return data
	}

	return bytes.ReplaceAll(data, []byte{'\n'}, target.Bytes())
}

// DetectFileLineTerminator detects line terminator style in a file.
// Reads only the beginning of the file for efficiency.
// Checks context cancellation before opening the file.
func DetectFileLineTerminator(ctx context.Context, path string) (LineTerminator, error) {
	if err := ctx.Err(); err != nil {
		return Unknown, err
	}

	f, err := os.Open(path)
	if err != nil {
		return Unknown, err
	}
	defer f.Close()

	// Read first 4KB to detect
	buf := make([]byte, 4096)

	n, err := f.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return Unknown, err
	}

	return DetectLineTerminator(buf[:n]), nil
}

// ============================================================================
// Path Operations
// ============================================================================

// Abs returns the absolute path.
// Returns the path unchanged if it's already absolute or on error.
func Abs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return abs
}

// Clean returns the cleaned path.
func Clean(path string) string {
	return filepath.Clean(path)
}

// Join joins path elements.
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Dir returns the directory component of the path.
func Dir(path string) string {
	return filepath.Dir(path)
}

// Base returns the last element of the path.
func Base(path string) string {
	return filepath.Base(path)
}

// Rel returns a relative path from base to target.
func Rel(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// Split splits a path into directory and file components.
func Split(path string) (dir, file string) {
	return filepath.Split(path)
}

// ============================================================================
// Utility Functions
// ============================================================================

// TempFile creates a temporary file and returns its path.
// The caller is responsible for removing the file.
// If closing the file fails, removes it and returns the close error.
func TempFile(dir, pattern string) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}

	name := f.Name()
	if err := finishWrite(f, nil); err != nil {
		os.Remove(name)

		return "", err
	}

	return name, nil
}

// TempDir creates a temporary directory and returns its path.
// The caller is responsible for removing the directory.
func TempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// SameFile reports whether fi1 and fi2 describe the same file.
func SameFile(path1, path2 string) bool {
	fi1, err := os.Stat(path1)
	if err != nil {
		return false
	}

	fi2, err := os.Stat(path2)
	if err != nil {
		return false
	}

	return os.SameFile(fi1, fi2)
}

// DirSize calculates the total size of a directory and its contents.
func DirSize(ctx context.Context, dir string) (int64, error) {
	var total int64

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			total += info.Size()
		}

		return nil
	})

	return total, err
}
