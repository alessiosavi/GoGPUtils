package benchcov

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// enforce turns TestRepository from a progress report into the coverage gate.
// The final rollout task sets it to true.
const enforce = true

const repoModule = "github.com/alessiosavi/GoGPUtils"

func TestRepository(t *testing.T) {
	root, modulePath, err := FindModuleRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	if modulePath != repoModule {
		t.Fatalf("module path = %q, want %q", modulePath, repoModule)
	}
	rep, err := Run(root, modulePath, "internal/benchcov/exempt.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d packages, %d APIs, %d benchmarks, %d exemptions, %d problems",
		len(rep.Packages), len(rep.Inventory.APIs), len(rep.Inventory.Benchmarks), len(rep.Exemptions), len(rep.Problems))
	for _, p := range rep.Problems {
		if enforce {
			t.Error(p)
		} else {
			t.Log(p)
		}
	}
}

func TestFixtures(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			dir := filepath.Join("testdata", e.Name())
			rep := mustRun(t, dir)
			var problems []string
			for _, p := range rep.Problems {
				problems = append(problems, p.Brief())
			}
			compareGolden(t, filepath.Join(dir, "want.txt"), problems)
			if exists(filepath.Join(dir, "scope.txt")) {
				var dirs []string
				for _, p := range rep.Packages {
					dirs = append(dirs, p.Dir)
				}
				compareGolden(t, filepath.Join(dir, "scope.txt"), dirs)
			}
			if exists(filepath.Join(dir, "apis.txt")) {
				var apis []string
				for _, s := range rep.Inventory.APIs {
					apis = append(apis, s.Pkg+" "+s.Name)
				}
				compareGolden(t, filepath.Join(dir, "apis.txt"), apis)
			}
		})
	}
}

func TestProblemPositions(t *testing.T) {
	want := map[string]map[string]string{
		"basic": {
			"a Bar: missing benchmark BenchmarkBar":          "a/a.go:7",
			"- -: too many exemptions: 1 > 0 (5% of 7 APIs)": "exempt.txt",
		},
		"exemptions": {
			"p Two: duplicate exemption (first at line 3)":            "exempt.txt:6",
			"p One: stale exemption: now benchmarked by BenchmarkOne": "exempt.txt:4",
		},
	}
	for fixture, positions := range want {
		rep := mustRun(t, filepath.Join("testdata", fixture))
		for brief, pos := range positions {
			i := slices.IndexFunc(rep.Problems, func(p Problem) bool { return p.Brief() == brief })
			if i < 0 {
				t.Errorf("%s: no problem %q", fixture, brief)
				continue
			}
			if got := rep.Problems[i].Pos; got != pos {
				t.Errorf("%s: %q at %q, want %q", fixture, brief, got, pos)
			}
		}
	}
}

func TestFindModuleRoot(t *testing.T) {
	root, modulePath, err := FindModuleRoot(filepath.Join("testdata", "scope", "awsome"))
	if err != nil || filepath.Base(root) != "scope" || modulePath != "example.com/scope" {
		t.Fatalf("FindModuleRoot(awsome) = %q, %q, %v", root, modulePath, err)
	}
	_, modulePath, err = FindModuleRoot(filepath.Join("testdata", "scope", "nested"))
	if err != nil || modulePath != "example.com/nested" {
		t.Fatalf("FindModuleRoot(nested) = %q, %v", modulePath, err)
	}
}

func TestExcluded(t *testing.T) {
	tests := map[string]bool{
		".": false, "cache": false, "textnorm/stopwords": false, "awsome": false,
		"aws": true, "aws/s3": true, "internal": true, "internal/benchkit": true,
		"x/testdata": true, "x/vendor": true, ".git": true, "x/.cache": true, "_tools": true,
	}
	for dir, want := range tests {
		if got := excluded(dir); got != want {
			t.Errorf("excluded(%q) = %v, want %v", dir, got, want)
		}
	}
}

func TestFilenameConstraint(t *testing.T) {
	tests := map[string]bool{
		"a_linux.go": true, "a_linux_test.go": true, "a_amd64_test.go": true,
		"a_linux_amd64.go": true, "cross_device_windows.go": true,
		"linux.go": false, "a_default.go": false, "a_test.go": false, "bench_test.go": false,
	}
	for name, want := range tests {
		if got := hasFilenameConstraint(name); got != want {
			t.Errorf("hasFilenameConstraint(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestBuildConstraint(t *testing.T) {
	tests := []struct {
		name, src           string
		constrained, ignore bool
	}{
		{"none", "package p\n", false, false},
		{"go-build", "//go:build linux\n\npackage p\n", true, false},
		{"legacy-plus-build", "// +build linux darwin\n\npackage p\n", true, false},
		{"ignore", "//go:build ignore\n\npackage p\n", true, true},
		{"doc-comment", "// Package p does things.\npackage p\n", false, false},
		{"after-package", "package p\n\n//go:build linux\n", false, false},
	}
	for _, tt := range tests {
		f, err := parser.ParseFile(token.NewFileSet(), tt.name+".go", tt.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		x, err := buildConstraint(f)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if (x != nil) != tt.constrained || isIgnore(x) != tt.ignore {
			t.Errorf("%s: constraint %v, ignore %v; want constrained %v, ignore %v",
				tt.name, x, isIgnore(x), tt.constrained, tt.ignore)
		}
	}
}

func mustRun(t *testing.T, dir string) Report {
	t.Helper()
	root, modulePath, err := FindModuleRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root, modulePath, "exempt.txt")
	if err != nil {
		t.Fatal(err)
	}

	return rep
}

func compareGolden(t *testing.T, file string, got []string) {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	if s := strings.TrimRight(string(data), "\n"); s != "" {
		want = strings.Split(s, "\n")
	}
	if !slices.Equal(got, want) {
		t.Errorf("%s mismatch\ngot:\n  %s\nwant:\n  %s", file,
			strings.Join(got, "\n  "), strings.Join(want, "\n  "))
	}
}

func exists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}
