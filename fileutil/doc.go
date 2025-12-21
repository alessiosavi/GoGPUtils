// Package fileutil provides file system operations with proper error handling.
//
// All functions in this package follow Go best practices:
// - Return errors instead of panicking
// - Accept context.Context for cancellable operations
// - Close resources properly using defer
//
// # Reading Files
//
// Functions for reading file contents:
//
//	lines, err := fileutil.ReadLines(ctx, "file.txt")
//	content, err := fileutil.ReadString(ctx, "file.txt")
//	data, err := fileutil.ReadBytes(ctx, "file.txt")
//
// # File Information
//
// Functions for checking file properties:
//
//	if fileutil.Exists(path) { ... }
//	if fileutil.IsDir(path) { ... }
//	size, err := fileutil.Size(path)
//	modTime, err := fileutil.ModTime(path)
//
// # Directory Operations
//
// Functions for working with directories:
//
//	files, err := fileutil.List(ctx, dir, fileutil.Recursive)
//	files, err := fileutil.Find(ctx, dir, "*.go")
//	err := fileutil.EnsureDir(path, 0755)
//
// # Line Terminator Handling
//
// Functions for handling different line endings:
//
//	terminator := fileutil.DetectLineTerminator(data)
//	normalized := fileutil.NormalizeLineTerminators(data, fileutil.LF)
//
// # Context Support
//
// Long-running operations accept context.Context for cancellation:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	lines, err := fileutil.ReadLines(ctx, "large-file.txt")
package fileutil
