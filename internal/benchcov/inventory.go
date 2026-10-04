package benchcov

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Symbol is an exported function ("Foo") or method ("T.M") of a package.
type Symbol struct {
	Pkg  string // module-relative package directory
	Name string
	Pos  string // module-relative "file:line"
}

// BenchmarkName returns the canonical benchmark name for s.
func (s Symbol) BenchmarkName() string {
	return "Benchmark" + strings.ReplaceAll(s.Name, ".", "_")
}

// Benchmark is a benchmark function found in a package's test files.
type Benchmark struct {
	Pkg  string
	Name string
	Pos  string
	// External is true for the external test package (package x_test).
	External bool
	// Constrained is true when the file has a build constraint or a
	// GOOS/GOARCH filename suffix, so it may not run by default.
	Constrained bool
}

// Inventory is the exported API and the benchmarks of a set of packages.
type Inventory struct {
	APIs       []Symbol
	Benchmarks []Benchmark
}

// Inspect inventories the exported API and the benchmarks of pkgs. Every
// eligible file is parsed regardless of build constraints, except files
// constrained exactly to the "ignore" tag; an identity declared in several
// files counts once.
func Inspect(root string, pkgs []Package) (Inventory, error) {
	var inv Inventory
	fset := token.NewFileSet()
	for _, p := range pkgs {
		dir := filepath.Join(root, filepath.FromSlash(p.Dir))
		entries, err := os.ReadDir(dir)
		if err != nil {
			return Inventory{}, err
		}
		seen := map[string]bool{}
		for _, e := range entries {
			if e.IsDir() || !eligible(e.Name()) {
				continue
			}
			file := filepath.Join(dir, e.Name())
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments|parser.SkipObjectResolution)
			if err != nil {
				return Inventory{}, err
			}
			expr, err := buildConstraint(f)
			if err != nil {
				return Inventory{}, fmt.Errorf("%s: %w", file, err)
			}
			if isIgnore(expr) {
				continue
			}
			constrained := expr != nil || hasFilenameConstraint(e.Name())
			if strings.HasSuffix(e.Name(), "_test.go") {
				inv.Benchmarks = append(inv.Benchmarks, benchmarks(fset, f, root, p.Dir, constrained)...)

				continue
			}
			for _, s := range exported(fset, f, root, p.Dir) {
				if !seen[s.Name] {
					seen[s.Name] = true
					inv.APIs = append(inv.APIs, s)
				}
			}
		}
	}
	slices.SortFunc(inv.APIs, func(a, b Symbol) int {
		return cmp.Or(cmp.Compare(a.Pkg, b.Pkg), cmp.Compare(a.Name, b.Name))
	})
	slices.SortFunc(inv.Benchmarks, func(a, b Benchmark) int {
		return cmp.Or(cmp.Compare(a.Pkg, b.Pkg), cmp.Compare(a.Name, b.Name), cmp.Compare(a.Pos, b.Pos))
	})

	return inv, nil
}

// buildConstraint returns the file's build constraint, or nil. A //go:build
// line takes precedence over legacy // +build lines, which are ANDed.
func buildConstraint(f *ast.File) (constraint.Expr, error) {
	var goBuild, plusBuild []constraint.Expr
	for _, cg := range f.Comments {
		if cg.Pos() >= f.Package {
			break
		}
		for _, c := range cg.List {
			if !constraint.IsGoBuild(c.Text) && !constraint.IsPlusBuild(c.Text) {
				continue
			}
			x, err := constraint.Parse(c.Text)
			if err != nil {
				return nil, err
			}
			if constraint.IsGoBuild(c.Text) {
				goBuild = append(goBuild, x)
			} else {
				plusBuild = append(plusBuild, x)
			}
		}
	}
	if len(goBuild) > 0 {
		return goBuild[0], nil
	}
	if len(plusBuild) == 0 {
		return nil, nil
	}
	x := plusBuild[0]
	for _, y := range plusBuild[1:] {
		x = &constraint.AndExpr{X: x, Y: y}
	}

	return x, nil
}

// isIgnore reports whether x is exactly the "ignore" tag.
func isIgnore(x constraint.Expr) bool {
	tag, ok := x.(*constraint.TagExpr)

	return ok && tag.Tag == "ignore"
}

var knownOS = setOf("aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris wasip1 windows zos")

var knownArch = setOf("386 amd64 amd64p32 arm armbe arm64 arm64be loong64 mips mipsle mips64 mips64le mips64p32 mips64p32le ppc ppc64 ppc64le riscv riscv64 s390 s390x sparc sparc64 wasm")

func setOf(words string) map[string]bool {
	m := map[string]bool{}
	for w := range strings.FieldsSeq(words) {
		m[w] = true
	}

	return m
}

// hasFilenameConstraint mirrors go/build: name_GOOS, name_GOARCH, and
// name_GOOS_GOARCH (each optionally followed by _test) constrain a file.
func hasFilenameConstraint(name string) bool {
	name, _, _ = strings.Cut(name, ".")
	i := strings.Index(name, "_")
	if i < 0 {
		return false
	}
	l := strings.Split(name[i:], "_")
	if n := len(l); n > 0 && l[n-1] == "test" {
		l = l[:n-1]
	}
	n := len(l)

	return (n >= 2 && knownOS[l[n-2]] && knownArch[l[n-1]]) || (n >= 1 && (knownOS[l[n-1]] || knownArch[l[n-1]]))
}

func exported(fset *token.FileSet, f *ast.File, root, pkg string) []Symbol {
	var out []Symbol
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || !fd.Name.IsExported() {
			continue
		}
		name := fd.Name.Name
		if fd.Recv != nil {
			if len(fd.Recv.List) != 1 {
				continue
			}
			recv := receiverBase(fd.Recv.List[0].Type)
			if !token.IsExported(recv) {
				continue
			}
			name = recv + "." + name
		}
		out = append(out, Symbol{Pkg: pkg, Name: name, Pos: position(fset, fd.Pos(), root)})
	}

	return out
}

// receiverBase strips pointers, parentheses, and type arguments from a
// receiver type and returns its base type name ("" if there is none).
func receiverBase(x ast.Expr) string {
	for {
		switch t := x.(type) {
		case *ast.StarExpr:
			x = t.X
		case *ast.ParenExpr:
			x = t.X
		case *ast.IndexExpr:
			x = t.X
		case *ast.IndexListExpr:
			x = t.X
		case *ast.Ident:
			return t.Name
		default:
			return ""
		}
	}
}

func benchmarks(fset *token.FileSet, f *ast.File, root, pkg string, constrained bool) []Benchmark {
	alias := testingAlias(f)
	if alias == "" {
		return nil
	}
	var out []Benchmark
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Type.TypeParams != nil || !isBenchmarkName(fd.Name.Name) {
			continue
		}
		if fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
			continue
		}
		params := fd.Type.Params.List
		if len(params) != 1 || len(params[0].Names) > 1 || !isTestingB(params[0].Type, alias) {
			continue
		}
		out = append(out, Benchmark{
			Pkg:         pkg,
			Name:        fd.Name.Name,
			Pos:         position(fset, fd.Pos(), root),
			External:    strings.HasSuffix(f.Name.Name, "_test"),
			Constrained: constrained,
		})
	}

	return out
}

// isBenchmarkName mirrors go test: "Benchmark" followed by nothing or by a
// character that is not a lowercase letter.
func isBenchmarkName(name string) bool {
	rest, ok := strings.CutPrefix(name, "Benchmark")
	if !ok {
		return false
	}
	if rest == "" {
		return true
	}
	r, _ := utf8.DecodeRuneInString(rest)

	return !unicode.IsLower(r)
}

// testingAlias returns the name the file uses for package testing: "" if it
// is not imported, "." for a dot import.
func testingAlias(f *ast.File) string {
	for _, imp := range f.Imports {
		if path, err := strconv.Unquote(imp.Path.Value); err != nil || path != "testing" {
			continue
		}
		if imp.Name == nil {
			return "testing"
		}
		if imp.Name.Name != "_" {
			return imp.Name.Name
		}
	}

	return ""
}

func isTestingB(x ast.Expr, alias string) bool {
	star, ok := x.(*ast.StarExpr)
	if !ok {
		return false
	}
	if alias == "." {
		id, ok := star.X.(*ast.Ident)

		return ok && id.Name == "B"
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)

	return ok && pkg.Name == alias && sel.Sel.Name == "B"
}

func position(fset *token.FileSet, pos token.Pos, root string) string {
	p := fset.Position(pos)
	rel, err := filepath.Rel(root, p.Filename)
	if err != nil {
		rel = p.Filename
	}

	return filepath.ToSlash(rel) + ":" + strconv.Itoa(p.Line)
}
