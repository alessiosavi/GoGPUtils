package stringutil

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
)

// Added with Task 4: these production helpers do not exist on BASE.
func TestPR3CosineGuard(t *testing.T) {
	const limit = 94_906_265
	for _, n := range []int{0, limit - 1, limit, limit + 1, math.MaxInt} {
		for _, pair := range [][2]int{{n, 0}, {0, n}, {n, limit}, {limit, n}} {
			want := n <= limit
			if got := cosineExactLengths(pair[0], pair[1]); got != want {
				t.Errorf("lengths=%v: got=%v want=%v", pair, got, want)
			}
		}
	}
	if cosineExactLengths(limit+1, limit+1) {
		t.Error("both oversized lengths accepted")
	}
}

// Artificial 64-bit counts exceed the public guarded domain to isolate
// conversion ordering without multi-GB strings. 50,000 is reachable on 32-bit.
func TestPR3CosineProducts(t *testing.T) {
	counts := []int64{50_000}
	if strconv.IntSize == 64 {
		counts = append(counts, 4_000_000_000)
	}
	for _, c := range counts {
		m1, m2 := map[string]int{"a": int(c)}, map[string]int{"a": int(c)}
		v := float64(c)
		want := (v * v) / (math.Sqrt(v*v) * math.Sqrt(v*v))
		if got := cosineFast(m1, m2); math.Float64bits(got) != math.Float64bits(want) {
			t.Errorf("count=%d: got=%g want=%g", c, got, want)
		}
	}
}

// Run instrumented copies through Go overlays: production stays unchanged.
// Route probes replace only the lengths and mark entry into the real branches;
// length probes record the arguments passed to the actual production predicate.
func TestPR3CosineDispatch(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(dir, "similarity.go")
	testPath := filepath.Join(dir, "zz_cosine_guard_test.go")
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	testSource, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"route", "lowered-lengths"} {
		t.Run(mode, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, sourcePath, source, 0)
			if err != nil {
				t.Fatal(err)
			}
			type edit struct {
				start, end int
				text       string
			}
			var edits []edit
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if mode == "route" {
					switch fn.Name.Name {
					case "CosineSimilarity":
						ast.Inspect(fn.Body, func(node ast.Node) bool {
							call, ok := node.(*ast.CallExpr)
							if !ok {
								return true
							}
							name, ok := call.Fun.(*ast.Ident)
							if ok && name.Name == "cosineExactLengths" && len(call.Args) == 2 {
								edits = append(edits, edit{fset.Position(call.Args[0].Pos()).Offset, fset.Position(call.Args[1].End()).Offset, "pr3ProbeLengths[0], pr3ProbeLengths[1]"})
							}
							return true
						})
					case "cosineFast", "cosineLegacy":
						pos := fset.Position(fn.Body.Lbrace).Offset + 1
						edits = append(edits, edit{pos, pos, "\npanic(" + strconv.Quote(fn.Name.Name) + ")\n"})
					}
				} else if fn.Name.Name == "cosineExactLengths" {
					edits = append(edits, edit{fset.Position(fn.Body.Pos()).Offset, fset.Position(fn.Body.End()).Offset, "{ panic([2]int{len1, len2}) }"})
				}
			}
			wantEdits := 1
			if mode == "route" {
				wantEdits = 3
			}
			if len(edits) != wantEdits {
				t.Fatalf("instrumentation edits=%d want=%d", len(edits), wantEdits)
			}
			sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
			instrumented := string(source)
			for _, e := range edits {
				instrumented = instrumented[:e.start] + e.text + instrumented[e.end:]
			}
			probe := pr3LoweredLengthProbe
			if mode == "route" {
				probe = pr3RouteProbe
			}
			tmp := t.TempDir()
			write := func(name string, data []byte) string {
				path := filepath.Join(tmp, name)
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
				return path
			}
			replacements := map[string]string{
				sourcePath: write("similarity.go", []byte(instrumented)),
				testPath:   write("guard_test.go", []byte(string(testSource)+probe)),
			}
			config, err := json.Marshal(map[string]any{"Replace": replacements})
			if err != nil {
				t.Fatal(err)
			}
			// Instrumentation intentionally makes the original function body unreachable.
			cmd := exec.Command("go", "test", "-vet=off", "-overlay", write("overlay.json", config), "-count=1", "-run", "^TestPR3DispatchProbe$", "-v", ".")
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			t.Logf("%s", out)
			if err != nil {
				t.Fatalf("dispatch probe failed: %v", err)
			}
		})
	}
}

const pr3RouteProbe = `
var pr3ProbeLengths [2]int
func TestPR3DispatchProbe(t *testing.T) {
 const limit = 94_906_265
 for _, n := range []int{0,limit-1,limit,limit+1,math.MaxInt} {
  for _, p := range [][2]int{{n,0},{0,n},{n,limit},{limit,n}} {
   pr3ProbeLengths=p
   want:="cosineFast"
   if n>limit {want="cosineLegacy"}
   func(){
    defer func(){if got:=recover();got!=want {t.Errorf("lengths=%v route=%v want=%v",p,got,want)}}()
    CosineSimilarity("ab","ac",1)
   }()
  }
 }
}
`

const pr3LoweredLengthProbe = `
func TestPR3DispatchProbe(t *testing.T) {
 for _, c:=range []struct{a,b string; want [2]int}{
  {"İẞA","b",[2]int{4,1}}, {"b","İẞA",[2]int{1,4}},
  {"\xffA","b",[2]int{4,1}}, {"b","\xffA",[2]int{1,4}},
  {"İẞA","\xffA",[2]int{4,4}},
 } {
  func(){
   defer func(){if got:=recover();got!=c.want {t.Errorf("inputs=%q,%q guard args=%v want=%v",c.a,c.b,got,c.want)}}()
   CosineSimilarity(c.a,c.b,1)
  }()
 }
}
`
