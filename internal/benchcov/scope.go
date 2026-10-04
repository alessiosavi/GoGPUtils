// Package benchcov checks that every exported function and method of the
// packages in the benchmark scope has a canonical benchmark: Foo needs
// BenchmarkFoo and T.M needs BenchmarkT_M. Exemptions are listed, with
// reasons, in exempt.txt. The scope excludes aws/..., internal/...,
// testdata, vendor, nested modules, and directories whose names start with
// "." or "_".
package benchcov

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Package is an in-scope package directory.
type Package struct {
	// Dir is the module-relative, slash-separated directory ("." for the root).
	Dir string
	// ImportPath is the package's import path.
	ImportPath string
}

// FindModuleRoot walks up from start to the nearest go.mod and returns its
// directory and module path.
func FindModuleRoot(start string) (root, modulePath string, err error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", "", err
	}
	for {
		data, err := os.ReadFile(filepath.Clean(filepath.Join(dir, "go.mod")))
		if err == nil {
			mp := modulePathOf(string(data))
			if mp == "" {
				return "", "", fmt.Errorf("benchcov: %s has no module directive", filepath.Join(dir, "go.mod"))
			}

			return dir, mp, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", fmt.Errorf("benchcov: no go.mod above %s", start)
		}
		dir = parent
	}
}

func modulePathOf(gomod string) string {
	for line := range strings.Lines(gomod) {
		if f := strings.Fields(line); len(f) >= 2 && f[0] == "module" {
			return strings.Trim(f[1], `"`)
		}
	}

	return ""
}

// excluded reports whether the module-relative directory rel is outside the
// benchmark scope. Matching is on whole path elements, so "awsome" is in
// scope while "aws" and "aws/s3" are not.
func excluded(rel string) bool {
	if rel == "." {
		return false
	}
	parts := strings.Split(rel, "/")
	if parts[0] == "aws" || parts[0] == "internal" {
		return true
	}
	for _, p := range parts {
		if p == "testdata" || p == "vendor" || strings.HasPrefix(p, ".") || strings.HasPrefix(p, "_") {
			return true
		}
	}

	return false
}

// eligible reports whether the go command would read a file with this name:
// a .go file whose name does not start with "." or "_".
func eligible(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "_")
}

// Scope returns the in-scope packages under root, sorted by Dir: every
// directory with an eligible non-test .go file, minus excluded directories
// and nested modules.
func Scope(root, modulePath string) ([]Package, error) {
	var pkgs []Package
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel != "." {
			if excluded(rel) {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				return filepath.SkipDir // nested module
			}
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if !e.IsDir() && eligible(e.Name()) && !strings.HasSuffix(e.Name(), "_test.go") {
				ip := modulePath
				if rel != "." {
					ip += "/" + rel
				}
				pkgs = append(pkgs, Package{Dir: rel, ImportPath: ip})

				break
			}
		}

		return nil
	})
	slices.SortFunc(pkgs, func(a, b Package) int { return cmp.Compare(a.Dir, b.Dir) })

	return pkgs, err
}
