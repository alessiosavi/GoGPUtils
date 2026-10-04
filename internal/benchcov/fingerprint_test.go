package benchcov

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustFingerprint(t *testing.T, dir string, include func(string) bool) string {
	t.Helper()
	got, err := Fingerprint(dir, include)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 64 {
		t.Fatalf("fingerprint %q is not a SHA-256 hex digest", got)
	}

	return got
}

func TestFingerprint(t *testing.T) {
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q")
	writeFile(t, dir, ".gitignore", "ignored.txt\n")
	writeFile(t, dir, "a.go", "package a\n")
	writeFile(t, dir, "a_test.go", "package a\n")
	untracked := mustFingerprint(t, dir, nil)

	gitIn(t, dir, "add", "a.go")
	if mustFingerprint(t, dir, nil) != untracked {
		t.Error("staging a file without changing it changed the fingerprint")
	}
	writeFile(t, dir, "ignored.txt", "anything")
	if mustFingerprint(t, dir, nil) != untracked {
		t.Error("an ignored file changed the fingerprint")
	}
	suite := mustFingerprint(t, dir, SuiteFile)
	writeFile(t, dir, "a.go", "package a // edited\n")
	if mustFingerprint(t, dir, nil) == untracked {
		t.Error("editing a file did not change the fingerprint")
	}
	if mustFingerprint(t, dir, SuiteFile) != suite {
		t.Error("editing a non-suite file changed the suite fingerprint")
	}
	writeFile(t, dir, "a_test.go", "package a // edited\n")
	if mustFingerprint(t, dir, SuiteFile) == suite {
		t.Error("editing a test file did not change the suite fingerprint")
	}
}

func TestFingerprintSymlinks(t *testing.T) {
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q")
	writeFile(t, dir, ".gitignore", "fixture.txt\n")
	writeFile(t, dir, "fixture.txt", "first")
	if err := os.Symlink("fixture.txt", filepath.Join(dir, "source.go")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	before := mustFingerprint(t, dir, nil)
	writeFile(t, dir, "fixture.txt", "second")
	if mustFingerprint(t, dir, nil) == before {
		t.Error("editing a symlink's target did not change the fingerprint")
	}
	if err := os.Remove(filepath.Join(dir, "fixture.txt")); err != nil {
		t.Fatal(err)
	}
	if _, err := Fingerprint(dir, nil); err == nil {
		t.Error("a broken symlink did not fail the fingerprint")
	}
}

func TestSuiteFile(t *testing.T) {
	tests := map[string]bool{
		"a_test.go": true, "x/y/bench_test.go": true, "internal/benchkit/benchkit.go": true,
		"a.go": false, "internal/benchcov/check.go": false, "docs/x.md": false,
	}
	for path, want := range tests {
		if got := SuiteFile(path); got != want {
			t.Errorf("SuiteFile(%q) = %v, want %v", path, got, want)
		}
	}
}
