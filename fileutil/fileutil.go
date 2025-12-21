package fileutil

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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
// Line terminators are stripped from each line.
func ReadLines(ctx context.Context, path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// ReadLinesN reads the first n lines from a file.
// If n <= 0, reads all lines.
func ReadLinesN(ctx context.Context, path string, n int) ([]string, error) {
	if n <= 0 {
		return ReadLines(ctx, path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, 0, n)

	scanner := bufio.NewScanner(f)
	for i := 0; i < n && scanner.Scan(); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// CountLines counts the number of lines in a file.
// More memory-efficient than ReadLines for large files.
func CountLines(ctx context.Context, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count := 0
	buf := make([]byte, 32*1024) // 32KB buffer
	lineSep := []byte{'\n'}

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		n, err := f.Read(buf)
		count += bytes.Count(buf[:n], lineSep)

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	return count, nil
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
func AppendBytes(path string, data []byte, perm fs.FileMode) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

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
// Creates destination directory if needed.
func Copy(ctx context.Context, src, dst string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return ErrNotFile
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := EnsureDir(dstDir, 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// Move moves a file from src to dst.
// Attempts rename first, falls back to copy+delete.
func Move(ctx context.Context, src, dst string) error {
	// Try rename first (fast path for same filesystem)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fall back to copy + delete
	if err := Copy(ctx, src, dst); err != nil {
		return err
	}

	return os.Remove(src)
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
func DetectFileLineTerminator(ctx context.Context, path string) (LineTerminator, error) {
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
func TempFile(dir, pattern string) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}

	name := f.Name()
	f.Close()

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
