package fileutil

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// Helper to create temp directory for tests.
func setupTestDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "fileutil-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() { os.RemoveAll(dir) })

	return dir
}

// Helper to create test file.
func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return path
}

// ============================================================================
// Existence and Type Tests
// ============================================================================

func TestExists(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	if !Exists(file) {
		t.Error("Exists should return true for existing file")
	}

	if !Exists(dir) {
		t.Error("Exists should return true for existing directory")
	}

	if Exists(filepath.Join(dir, "nonexistent")) {
		t.Error("Exists should return false for nonexistent path")
	}
}

func TestIsFile(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	if !IsFile(file) {
		t.Error("IsFile should return true for regular file")
	}

	if IsFile(dir) {
		t.Error("IsFile should return false for directory")
	}

	if IsFile(filepath.Join(dir, "nonexistent")) {
		t.Error("IsFile should return false for nonexistent path")
	}
}

func TestIsDir(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	if !IsDir(dir) {
		t.Error("IsDir should return true for directory")
	}

	if IsDir(file) {
		t.Error("IsDir should return false for regular file")
	}
}

func TestIsEmpty(t *testing.T) {
	dir := setupTestDir(t)

	// Empty file
	emptyFile := createTestFile(t, dir, "empty.txt", "")

	empty, err := IsEmpty(emptyFile)
	if err != nil || !empty {
		t.Error("IsEmpty should return true for empty file")
	}

	// Non-empty file
	nonEmptyFile := createTestFile(t, dir, "nonempty.txt", "content")

	empty, err = IsEmpty(nonEmptyFile)
	if err != nil || empty {
		t.Error("IsEmpty should return false for non-empty file")
	}

	// Empty directory
	emptyDir := filepath.Join(dir, "emptydir")
	os.Mkdir(emptyDir, 0755)

	empty, err = IsEmpty(emptyDir)
	if err != nil || !empty {
		t.Error("IsEmpty should return true for empty directory")
	}

	// Non-empty directory
	empty, err = IsEmpty(dir)
	if err != nil || empty {
		t.Error("IsEmpty should return false for non-empty directory")
	}
}

// ============================================================================
// File Information Tests
// ============================================================================

func TestSize(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	file := createTestFile(t, dir, "test.txt", content)

	size, err := Size(file)
	if err != nil {
		t.Fatalf("Size() error: %v", err)
	}

	if size != int64(len(content)) {
		t.Errorf("Size() = %d, want %d", size, len(content))
	}

	// Directory should error
	_, err = Size(dir)
	if !errors.Is(err, ErrNotFile) {
		t.Errorf("Size(dir) error = %v, want ErrNotFile", err)
	}
}

func TestModTime(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	modTime, err := ModTime(file)
	if err != nil {
		t.Fatalf("ModTime() error: %v", err)
	}

	// Should be recent (within last minute)
	if time.Since(modTime) > time.Minute {
		t.Errorf("ModTime() = %v, expected recent time", modTime)
	}
}

func TestExtension(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/file.txt", ".txt"},
		{"/path/to/file.tar.gz", ".gz"},
		{"/path/to/file", ""},
		{"file.go", ".go"},
	}

	for _, tt := range tests {
		got := Extension(tt.path)
		if got != tt.want {
			t.Errorf("Extension(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestBaseName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/file.txt", "file"},
		{"/path/to/file.tar.gz", "file.tar"},
		{"/path/to/file", "file"},
		{"file.go", "file"},
	}

	for _, tt := range tests {
		got := BaseName(tt.path)
		if got != tt.want {
			t.Errorf("BaseName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

// ============================================================================
// Reading Tests
// ============================================================================

func TestReadBytes(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	file := createTestFile(t, dir, "test.txt", content)

	ctx := context.Background()

	data, err := ReadBytes(ctx, file)
	if err != nil {
		t.Fatalf("ReadBytes() error: %v", err)
	}

	if string(data) != content {
		t.Errorf("ReadBytes() = %q, want %q", string(data), content)
	}
}

func TestReadString(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	file := createTestFile(t, dir, "test.txt", content)

	ctx := context.Background()

	got, err := ReadString(ctx, file)
	if err != nil {
		t.Fatalf("ReadString() error: %v", err)
	}

	if got != content {
		t.Errorf("ReadString() = %q, want %q", got, content)
	}
}

func TestReadLines(t *testing.T) {
	dir := setupTestDir(t)
	content := "line1\nline2\nline3"
	file := createTestFile(t, dir, "test.txt", content)

	ctx := context.Background()

	lines, err := ReadLines(ctx, file)
	if err != nil {
		t.Fatalf("ReadLines() error: %v", err)
	}

	want := []string{"line1", "line2", "line3"}
	if !slices.Equal(lines, want) {
		t.Errorf("ReadLines() = %v, want %v", lines, want)
	}
}

func TestReadLinesN(t *testing.T) {
	dir := setupTestDir(t)
	content := "line1\nline2\nline3\nline4\nline5"
	file := createTestFile(t, dir, "test.txt", content)

	ctx := context.Background()

	lines, err := ReadLinesN(ctx, file, 3)
	if err != nil {
		t.Fatalf("ReadLinesN() error: %v", err)
	}

	want := []string{"line1", "line2", "line3"}
	if !slices.Equal(lines, want) {
		t.Errorf("ReadLinesN() = %v, want %v", lines, want)
	}
}

func TestReadWithCanceledContext(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := ReadBytes(ctx, file)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("ReadBytes with canceled context error = %v, want context.Canceled", err)
	}
}

func TestCountLines(t *testing.T) {
	dir := setupTestDir(t)
	content := "line1\nline2\nline3\nline4\n"
	file := createTestFile(t, dir, "test.txt", content)

	ctx := context.Background()

	count, err := CountLines(ctx, file)
	if err != nil {
		t.Fatalf("CountLines() error: %v", err)
	}

	if count != 4 {
		t.Errorf("CountLines() = %d, want 4", count)
	}
}

// ============================================================================
// Writing Tests
// ============================================================================

func TestWriteBytes(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "output.txt")
	content := []byte("hello world")

	err := WriteBytes(path, content, 0644)
	if err != nil {
		t.Fatalf("WriteBytes() error: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != string(content) {
		t.Errorf("File content = %q, want %q", string(got), string(content))
	}
}

func TestWriteLines(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "output.txt")
	lines := []string{"line1", "line2", "line3"}

	err := WriteLines(path, lines, LF, 0644)
	if err != nil {
		t.Fatalf("WriteLines() error: %v", err)
	}

	got, _ := os.ReadFile(path)

	want := "line1\nline2\nline3\n"

	if string(got) != want {
		t.Errorf("File content = %q, want %q", string(got), want)
	}
}

func TestWriteLinesCRLF(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "output.txt")
	lines := []string{"line1", "line2"}

	err := WriteLines(path, lines, CRLF, 0644)
	if err != nil {
		t.Fatalf("WriteLines() error: %v", err)
	}

	got, _ := os.ReadFile(path)

	want := "line1\r\nline2\r\n"

	if string(got) != want {
		t.Errorf("File content = %q, want %q", string(got), want)
	}
}

func TestAppendString(t *testing.T) {
	dir := setupTestDir(t)
	path := createTestFile(t, dir, "test.txt", "hello")

	err := AppendString(path, " world", 0644)
	if err != nil {
		t.Fatalf("AppendString() error: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "hello world" {
		t.Errorf("File content = %q, want %q", string(got), "hello world")
	}
}

// ============================================================================
// Directory Operations Tests
// ============================================================================

func TestEnsureDir(t *testing.T) {
	dir := setupTestDir(t)
	newDir := filepath.Join(dir, "a", "b", "c")

	err := EnsureDir(newDir, 0755)
	if err != nil {
		t.Fatalf("EnsureDir() error: %v", err)
	}

	if !IsDir(newDir) {
		t.Error("Directory was not created")
	}

	// Should succeed for existing directory
	err = EnsureDir(newDir, 0755)
	if err != nil {
		t.Errorf("EnsureDir() on existing dir error: %v", err)
	}

	// Should fail for existing file
	file := createTestFile(t, dir, "file.txt", "content")

	err = EnsureDir(file, 0755)
	if !errors.Is(err, ErrNotDir) {
		t.Errorf("EnsureDir() on file error = %v, want ErrNotDir", err)
	}
}

func TestList(t *testing.T) {
	dir := setupTestDir(t)
	createTestFile(t, dir, "file1.txt", "content")
	createTestFile(t, dir, "file2.txt", "content")
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	createTestFile(t, dir, filepath.Join("subdir", "file3.txt"), "content")

	ctx := context.Background()

	// Non-recursive
	entries, err := List(ctx, dir, 0)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("List() returned %d entries, want 3", len(entries))
	}

	// Recursive
	entries, err = List(ctx, dir, Recursive)
	if err != nil {
		t.Fatalf("List(Recursive) error: %v", err)
	}

	if len(entries) != 4 { // 2 files + 1 dir + 1 nested file
		t.Errorf("List(Recursive) returned %d entries, want 4", len(entries))
	}

	// Files only
	entries, err = List(ctx, dir, FilesOnly)
	if err != nil {
		t.Fatalf("List(FilesOnly) error: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("List(FilesOnly) returned %d entries, want 2", len(entries))
	}

	// Dirs only
	entries, err = List(ctx, dir, DirsOnly)
	if err != nil {
		t.Fatalf("List(DirsOnly) error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("List(DirsOnly) returned %d entries, want 1", len(entries))
	}
}

func TestFind(t *testing.T) {
	dir := setupTestDir(t)
	createTestFile(t, dir, "file1.go", "content")
	createTestFile(t, dir, "file2.go", "content")
	createTestFile(t, dir, "file3.txt", "content")
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	createTestFile(t, dir, filepath.Join("subdir", "file4.go"), "content")

	ctx := context.Background()

	matches, err := Find(ctx, dir, "*.go")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("Find(*.go) returned %d matches, want 3", len(matches))
	}
}

func TestFindByExtension(t *testing.T) {
	dir := setupTestDir(t)
	createTestFile(t, dir, "file1.txt", "content")
	createTestFile(t, dir, "file2.txt", "content")
	createTestFile(t, dir, "file3.go", "content")

	ctx := context.Background()

	matches, err := FindByExtension(ctx, dir, ".txt")
	if err != nil {
		t.Fatalf("FindByExtension() error: %v", err)
	}

	if len(matches) != 2 {
		t.Errorf("FindByExtension(.txt) returned %d matches, want 2", len(matches))
	}

	// Without dot prefix
	matches, err = FindByExtension(ctx, dir, "txt")
	if err != nil {
		t.Fatalf("FindByExtension() error: %v", err)
	}

	if len(matches) != 2 {
		t.Errorf("FindByExtension(txt) returned %d matches, want 2", len(matches))
	}
}

// ============================================================================
// File Operations Tests
// ============================================================================

func TestCopy(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	src := createTestFile(t, dir, "source.txt", content)
	dst := filepath.Join(dir, "destination.txt")

	ctx := context.Background()

	err := Copy(ctx, src, dst)
	if err != nil {
		t.Fatalf("Copy() error: %v", err)
	}

	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("Copied file content = %q, want %q", string(got), content)
	}
}

func TestCopyToNewDir(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	src := createTestFile(t, dir, "source.txt", content)
	dst := filepath.Join(dir, "newdir", "subdir", "destination.txt")

	ctx := context.Background()

	err := Copy(ctx, src, dst)
	if err != nil {
		t.Fatalf("Copy() error: %v", err)
	}

	if !Exists(dst) {
		t.Error("Destination file was not created")
	}
}

func TestMove(t *testing.T) {
	dir := setupTestDir(t)
	content := "hello world"
	src := createTestFile(t, dir, "source.txt", content)
	dst := filepath.Join(dir, "destination.txt")

	ctx := context.Background()

	err := Move(ctx, src, dst)
	if err != nil {
		t.Fatalf("Move() error: %v", err)
	}

	if Exists(src) {
		t.Error("Source file should not exist after move")
	}

	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("Moved file content = %q, want %q", string(got), content)
	}
}

func TestTouch(t *testing.T) {
	dir := setupTestDir(t)

	// Create new file
	path := filepath.Join(dir, "newfile.txt")

	err := Touch(path)
	if err != nil {
		t.Fatalf("Touch() error: %v", err)
	}

	if !Exists(path) {
		t.Error("Touch should create new file")
	}

	// Update existing file
	file := createTestFile(t, dir, "existing.txt", "content")
	oldTime, _ := ModTime(file)

	time.Sleep(10 * time.Millisecond)

	err = Touch(file)
	if err != nil {
		t.Fatalf("Touch() error: %v", err)
	}

	newTime, _ := ModTime(file)
	if !newTime.After(oldTime) {
		t.Error("Touch should update modification time")
	}
}

// ============================================================================
// Line Terminator Tests
// ============================================================================

func TestDetectLineTerminator(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  LineTerminator
	}{
		{"LF", []byte("line1\nline2\n"), LF},
		{"CRLF", []byte("line1\r\nline2\r\n"), CRLF},
		{"CR", []byte("line1\rline2\r"), CR},
		{"Mixed", []byte("line1\nline2\r\n"), Mixed},
		{"Unknown", []byte("no newlines"), Unknown},
		{"Empty", []byte{}, Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectLineTerminator(tt.input)
			if got != tt.want {
				t.Errorf("DetectLineTerminator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeLineTerminators(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		target LineTerminator
		want   []byte
	}{
		{"CRLF to LF", []byte("a\r\nb\r\n"), LF, []byte("a\nb\n")},
		{"LF to CRLF", []byte("a\nb\n"), CRLF, []byte("a\r\nb\r\n")},
		{"CR to LF", []byte("a\rb\r"), LF, []byte("a\nb\n")},
		{"Mixed to LF", []byte("a\r\nb\rc\n"), LF, []byte("a\nb\nc\n")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeLineTerminators(tt.input, tt.target)
			if string(got) != string(tt.want) {
				t.Errorf("NormalizeLineTerminators() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLineTerminatorString(t *testing.T) {
	tests := []struct {
		lt   LineTerminator
		want string
	}{
		{LF, "LF (\\n)"},
		{CRLF, "CRLF (\\r\\n)"},
		{CR, "CR (\\r)"},
		{Mixed, "Mixed"},
		{Unknown, "Unknown"},
	}

	for _, tt := range tests {
		got := tt.lt.String()
		if got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.lt, got, tt.want)
		}
	}
}

// ============================================================================
// Path Operations Tests
// ============================================================================

func TestAbs(t *testing.T) {
	got := Abs("relative/path")
	if !filepath.IsAbs(got) {
		t.Errorf("Abs() returned relative path: %s", got)
	}
}

func TestClean(t *testing.T) {
	got := Clean("/path/./to/../file")

	want := "/path/file"

	if got != want {
		t.Errorf("Clean() = %q, want %q", got, want)
	}
}

func TestJoin(t *testing.T) {
	got := Join("path", "to", "file")

	want := filepath.Join("path", "to", "file")

	if got != want {
		t.Errorf("Join() = %q, want %q", got, want)
	}
}

func TestSplit(t *testing.T) {
	dir, file := Split("/path/to/file.txt")
	if dir != "/path/to/" || file != "file.txt" {
		t.Errorf("Split() = (%q, %q)", dir, file)
	}
}

// ============================================================================
// Utility Tests
// ============================================================================

func TestTempFile(t *testing.T) {
	path, err := TempFile("", "test-*")
	if err != nil {
		t.Fatalf("TempFile() error: %v", err)
	}
	defer os.Remove(path)

	if !Exists(path) {
		t.Error("TempFile did not create file")
	}
}

func TestTempDir(t *testing.T) {
	path, err := TempDir("", "test-*")
	if err != nil {
		t.Fatalf("TempDir() error: %v", err)
	}
	defer os.RemoveAll(path)

	if !IsDir(path) {
		t.Error("TempDir did not create directory")
	}
}

func TestSameFile(t *testing.T) {
	dir := setupTestDir(t)
	file := createTestFile(t, dir, "test.txt", "content")

	if !SameFile(file, file) {
		t.Error("SameFile should return true for same path")
	}

	file2 := createTestFile(t, dir, "test2.txt", "content")
	if SameFile(file, file2) {
		t.Error("SameFile should return false for different files")
	}
}

func TestDirSize(t *testing.T) {
	dir := setupTestDir(t)
	createTestFile(t, dir, "file1.txt", "hello") // 5 bytes
	createTestFile(t, dir, "file2.txt", "world") // 5 bytes
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	createTestFile(t, dir, filepath.Join("subdir", "file3.txt"), "test") // 4 bytes

	ctx := context.Background()

	size, err := DirSize(ctx, dir)
	if err != nil {
		t.Fatalf("DirSize() error: %v", err)
	}

	if size != 14 {
		t.Errorf("DirSize() = %d, want 14", size)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkReadLines(b *testing.B) {
	dir, _ := os.MkdirTemp("", "bench-*")
	defer os.RemoveAll(dir)

	// Create a file with 1000 lines
	var content string
	var contentSb740 strings.Builder
	for range 1000 {
		contentSb740.WriteString("This is a test line for benchmarking purposes.\n")
	}
	content += contentSb740.String()

	path := filepath.Join(dir, "bench.txt")
	os.WriteFile(path, []byte(content), 0644)

	ctx := context.Background()

	b.ResetTimer()

	for range b.N {
		ReadLines(ctx, path)
	}
}

func BenchmarkCountLines(b *testing.B) {
	dir, _ := os.MkdirTemp("", "bench-*")
	defer os.RemoveAll(dir)

	// Create a file with 10000 lines
	var content string
	var contentSb760 strings.Builder
	for range 10000 {
		contentSb760.WriteString("This is a test line for benchmarking purposes.\n")
	}
	content += contentSb760.String()

	path := filepath.Join(dir, "bench.txt")
	os.WriteFile(path, []byte(content), 0644)

	ctx := context.Background()

	b.ResetTimer()

	for range b.N {
		CountLines(ctx, path)
	}
}

func BenchmarkList(b *testing.B) {
	dir, _ := os.MkdirTemp("", "bench-*")
	defer os.RemoveAll(dir)

	// Create 100 files
	for i := range 100 {
		path := filepath.Join(dir, string(rune('a'+i%26))+".txt")
		os.WriteFile(path, []byte("content"), 0644)
	}

	ctx := context.Background()

	b.ResetTimer()

	for range b.N {
		List(ctx, dir, 0)
	}
}
