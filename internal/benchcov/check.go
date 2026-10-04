package benchcov

import (
	"bufio"
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// maxExemptPercent caps the exemptions at 5% of the in-scope API.
const maxExemptPercent = 5

// Exemption is one entry of the exemption file.
type Exemption struct {
	Pkg    string // module-relative package directory
	Name   string // "Foo" or "T.M"
	Reason string
	Line   int
}

// Problem is one coverage-gate failure.
type Problem struct {
	Pkg    string // module-relative package directory, or "" for file-level problems
	Symbol string // API identity or benchmark name, or ""
	Pos    string // "file:line" relative to the module root, or a file name
	Msg    string
}

// Brief formats p without its position: "<pkg> <symbol>: <msg>".
func (p Problem) Brief() string {
	return cmp.Or(p.Pkg, "-") + " " + cmp.Or(p.Symbol, "-") + ": " + p.Msg
}

// String formats p with its position.
func (p Problem) String() string { return p.Pos + ": " + p.Brief() }

// Report is the result of one coverage run.
type Report struct {
	Packages   []Package
	Inventory  Inventory
	Exemptions []Exemption
	Problems   []Problem // sorted by package, symbol, and message
}

// Run checks the module at root: it selects the in-scope packages,
// inventories their API and benchmarks, and applies the exemption file at
// exemptRel (slash-separated and relative to root; a missing file means no
// exemptions).
func Run(root, modulePath, exemptRel string) (Report, error) {
	pkgs, err := Scope(root, modulePath)
	if err != nil {
		return Report{}, err
	}
	inv, err := Inspect(root, pkgs)
	if err != nil {
		return Report{}, err
	}
	var ex []Exemption
	var probs []Problem
	data, err := os.ReadFile(filepath.Clean(filepath.Join(root, filepath.FromSlash(exemptRel))))
	switch {
	case err == nil:
		ex, probs = parseExemptions(exemptRel, bytes.NewReader(data))
	case !errors.Is(err, fs.ErrNotExist):
		return Report{}, err
	}
	probs = append(probs, check(inv, ex, exemptRel)...)
	slices.SortFunc(probs, func(a, b Problem) int {
		return cmp.Or(cmp.Compare(a.Pkg, b.Pkg), cmp.Compare(a.Symbol, b.Symbol), cmp.Compare(a.Msg, b.Msg))
	})

	return Report{Packages: pkgs, Inventory: inv, Exemptions: ex, Problems: probs}, nil
}

// parseExemptions reads "<pkgdir> <Symbol> # <reason>" lines. Blank lines and
// lines starting with "#" are ignored.
func parseExemptions(name string, r io.Reader) ([]Exemption, []Problem) {
	var ex []Exemption
	var probs []Problem
	first := map[[2]string]int{}
	sc := bufio.NewScanner(r)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		pos := fmt.Sprintf("%s:%d", name, line)
		head, reason, hasReason := strings.Cut(text, "#")
		reason = strings.TrimSpace(reason)
		fields := strings.Fields(head)
		if len(fields) != 2 {
			probs = append(probs, Problem{Pos: pos, Msg: `malformed exemption: want "<pkgdir> <Symbol> # <reason>"`})

			continue
		}
		p := Problem{Pkg: fields[0], Symbol: fields[1], Pos: pos}
		key := [2]string{fields[0], fields[1]}
		switch at, dup := first[key]; {
		case !hasReason || reason == "":
			p.Msg = "exemption has no reason"
		case !canonical(fields[0], fields[1]):
			p.Msg = "exemption is not canonical"
		case dup:
			p.Msg = fmt.Sprintf("duplicate exemption (first at line %d)", at)
		default:
			first[key] = line
			ex = append(ex, Exemption{Pkg: fields[0], Name: fields[1], Reason: reason, Line: line})

			continue
		}
		probs = append(probs, p)
	}
	if err := sc.Err(); err != nil {
		probs = append(probs, Problem{Pos: name, Msg: "reading exemptions: " + err.Error()})
	}

	return ex, probs
}

// canonical reports whether pkg is a clean module-relative directory and
// name is an exported "Foo" or "T.M" identity.
func canonical(pkg, name string) bool {
	if pkg == "" || pkg == ".." || path.Clean(pkg) != pkg || strings.HasPrefix(pkg, "../") || path.IsAbs(pkg) {
		return false
	}
	parts := strings.Split(name, ".")
	if len(parts) > 2 {
		return false
	}
	for _, part := range parts {
		if !token.IsIdentifier(part) || !token.IsExported(part) {
			return false
		}
	}

	return true
}

// check matches the inventory against the exemptions.
func check(inv Inventory, ex []Exemption, exemptRel string) []Problem {
	type key struct{ pkg, name string }
	var probs []Problem

	benches := map[key][]Benchmark{}
	for _, b := range inv.Benchmarks {
		k := key{b.Pkg, b.Name}
		benches[k] = append(benches[k], b)
	}
	runs := func(k key) bool {
		return slices.ContainsFunc(benches[k], func(b Benchmark) bool { return !b.Constrained })
	}
	for k, bs := range benches {
		var roots []Benchmark
		for _, b := range bs {
			if !b.Constrained {
				roots = append(roots, b)
			}
		}
		if len(roots) > 1 {
			probs = append(probs, Problem{Pkg: k.pkg, Symbol: k.name, Pos: roots[1].Pos,
				Msg: "duplicate benchmark in internal and external test packages"})
		}
	}

	apis := map[key]Symbol{}
	byBench := map[key][]string{}
	for _, s := range inv.APIs {
		apis[key{s.Pkg, s.Name}] = s
		k := key{s.Pkg, s.BenchmarkName()}
		byBench[k] = append(byBench[k], s.Name)
	}
	for k, names := range byBench {
		if len(names) > 1 {
			slices.Sort(names)
			probs = append(probs, Problem{Pkg: k.pkg, Symbol: k.name, Pos: apis[key{k.pkg, names[0]}].Pos,
				Msg: "ambiguous benchmark identity: " + strings.Join(names, ", ")})
		}
	}

	exempt := map[key]bool{}
	for _, e := range ex {
		exempt[key{e.Pkg, e.Name}] = true
	}
	for _, s := range inv.APIs {
		k := key{s.Pkg, s.BenchmarkName()}
		switch {
		case runs(k), exempt[key{s.Pkg, s.Name}]:
		case len(benches[k]) > 0:
			probs = append(probs, Problem{Pkg: s.Pkg, Symbol: s.Name, Pos: s.Pos,
				Msg: "benchmark " + s.BenchmarkName() + " exists only in build-constrained test files (not run by default)"})
		default:
			probs = append(probs, Problem{Pkg: s.Pkg, Symbol: s.Name, Pos: s.Pos,
				Msg: "missing benchmark " + s.BenchmarkName()})
		}
	}
	for _, e := range ex {
		pos := fmt.Sprintf("%s:%d", exemptRel, e.Line)
		s, ok := apis[key{e.Pkg, e.Name}]
		switch {
		case !ok:
			probs = append(probs, Problem{Pkg: e.Pkg, Symbol: e.Name, Pos: pos, Msg: "stale exemption: no such exported symbol"})
		case runs(key{s.Pkg, s.BenchmarkName()}):
			probs = append(probs, Problem{Pkg: e.Pkg, Symbol: e.Name, Pos: pos, Msg: "stale exemption: now benchmarked by " + s.BenchmarkName()})
		}
	}
	if limit := maxExemptPercent * len(inv.APIs) / 100; len(ex) > limit {
		probs = append(probs, Problem{Pos: exemptRel, Msg: fmt.Sprintf("too many exemptions: %d > %d (%d%% of %d APIs)",
			len(ex), limit, maxExemptPercent, len(inv.APIs))})
	}

	return probs
}
