package benchcov

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Fingerprint returns the SHA-256 of a manifest of every tracked or untracked,
// non-ignored file resolving to a regular file in the checkout at dir that
// include accepts (nil accepts all): one "path NUL sha256(content) LF" line
// per file, sorted by path. Staging a file does not change the result;
// editing one does. Symlinks to regular files are read through the link and
// hashed under the link's path; broken symlinks and other non-regular entries
// are errors. A tracked path deleted from the working tree is omitted.
func Fingerprint(dir string, include func(path string) bool) (string, error) {
	cmd := exec.Command("git", "ls-files", "-z", "-co", "--exclude-standard")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("benchcov: git ls-files in %s: %w", dir, err)
	}
	paths := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	slices.Sort(paths)
	paths = slices.Compact(paths)
	h := sha256.New()
	for _, p := range paths {
		if p == "" || (include != nil && !include(p)) {
			continue
		}
		full := filepath.Join(dir, filepath.FromSlash(p))
		info, err := os.Lstat(full)
		if errors.Is(err, fs.ErrNotExist) {
			continue // a tracked file deleted from the working tree
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			info, err = os.Stat(full)
			if err != nil {
				return "", err
			}
		}
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("benchcov: cannot fingerprint non-regular path %q", p)
		}
		data, err := os.ReadFile(filepath.Clean(full))
		if err != nil {
			return "", err
		}
		if _, err := fmt.Fprintf(h, "%s\x00%x\n", p, sha256.Sum256(data)); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// SuiteFile reports whether a module-relative path belongs to the benchmark
// suite: a test file or the benchmark kit.
func SuiteFile(path string) bool {
	return strings.HasSuffix(path, "_test.go") || strings.HasPrefix(path, "internal/benchkit/")
}
