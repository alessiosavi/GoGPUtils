package textnorm

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode"
	"unsafe"
	"weak"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// BASE oracles retain the production bodies from commit 4c05787.
func baseNormalizeUnicodeStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

func baseNormalizeUnicodeLatinStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	decomposed := norm.NFD.String(s)
	var b strings.Builder
	b.Grow(len(decomposed))
	lastBaseLatin := false
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			if lastBaseLatin {
				continue
			}
			b.WriteRune(r)
			continue
		}
		lastBaseLatin = unicode.Is(unicode.Latin, r)
		b.WriteRune(r)
	}
	return norm.NFC.String(b.String()), nil
}

func baseFoldCaseStage() Stage {
	return func(s string) (string, error) {
		return cases.Fold().String(s), nil
	}
}

func baseLowerStage() Stage {
	return func(s string) (string, error) {
		return cases.Lower(language.Und).String(s), nil
	}
}

func baseMapRunesStage(fn func(rune) rune) Stage {
	if fn == nil {
		return nil
	}
	return func(s string) (string, error) {
		result, _, err := transform.String(transform.Chain(runes.Map(fn)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	}
}

func baseFilterRunesStage(keep runes.Set) Stage {
	if keep == nil {
		return nil
	}
	return func(s string) (string, error) {
		remove := runes.Predicate(func(r rune) bool {
			return !keep.Contains(r)
		})
		result, _, err := transform.String(transform.Chain(runes.Remove(remove)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	}
}

// Deterministic, per-class PR-6 input distribution. The small-domain loops and
// explicit chunk boundaries supplement (rather than depend on) random coverage.
func zzInputs(visit func(string)) {
	for _, s := range []string{"", "A", "a", "a\x00Z\r\n\t\x7f", "\x80", "\xff", "\xc0\xaf", "\xed\xa0\x80", "\xef\xbf\xbd", "0000000000\xf2Ă", "00000000000000000000000000000ꮒ0", "Ꮳꮳ", "ΟΣ Σς", "café Å e\u0301 ﬁ Ａ", "क़ مَكتَب 中文 😀"} {
		visit(s)
	}
	for i := range 256 {
		visit(string([]byte{byte(i)}))
	}
	for _, n := range []int{127, 128, 129, 255, 256, 257, 4095, 4096, 4097, 65536} {
		for _, unit := range []string{"a", "AZ\x00", "é", "\xff", "\u0301", "ꮒ"} {
			visit(strings.Repeat(unit, n/len(unit)+1)[:n])
		}
	}
	rng := rand.New(rand.NewPCG(6, 1))
	alphabet := []rune("aAZéÅﬁＡİẞΟΣςᏣꮳक़مَ中文😀\u0301\x00\r\n\t\ufffd")
	for class := range 5 {
		for range 10000 {
			n := rng.IntN(80)
			var b strings.Builder
			for j := range n {
				switch class {
				case 0, 1:
					b.WriteByte(byte(rng.IntN(128)))
				case 2:
					b.WriteRune(alphabet[rng.IntN(len(alphabet))])
				case 3:
					b.WriteByte(byte(rng.IntN(256)))
				case 4:
					if j%3 == 0 {
						b.WriteRune('\ufffd')
					} else {
						b.WriteByte(byte(rng.IntN(128)))
					}
				}
			}
			s := b.String()
			if class == 1 {
				s += "é"
			}
			if class == 3 {
				s += "\x80"
			}
			visit(s)
		}
	}
}

func zzEqual(t testing.TB, name, in, got string, ge error, want string, we error) {
	t.Helper()
	if got != want || (!errors.Is(ge, we) || !errors.Is(we, ge)) {
		t.Fatalf("%s input=%q got=%q err=%v want=%q err=%v", name, in, got, ge, want, we)
	}
	if len(in) > 0 && len(want) > 0 {
		ba := unsafe.StringData(want) == unsafe.StringData(in)
		ca := unsafe.StringData(got) == unsafe.StringData(in)
		if ba != ca {
			t.Fatalf("%s input=%q pointer identity: BASE=%v candidate=%v", name, in, ba, ca)
		}
	}
	runtime.KeepAlive(in)
}

func zzPairs() []struct {
	name       string
	head, base Stage
} {
	return []struct {
		name       string
		head, base Stage
	}{
		{"NormalizeUnicode", normalizeUnicodeStage, baseNormalizeUnicodeStage},
		{"NormalizeUnicodeLatin", normalizeUnicodeLatinStage, baseNormalizeUnicodeLatinStage},
		{"FoldCase", New().FoldCase().Run, baseFoldCaseStage()},
		{"Lower", New().Lower().Run, baseLowerStage()},
		{"MapRunes", New().MapRunes(unicode.ToLower).Run, baseMapRunesStage(unicode.ToLower)},
		{"FilterRunes", New().FilterRunes(runes.In(unicode.Letter)).Run, baseFilterRunesStage(runes.In(unicode.Letter))},
	}
}
func zzCompareStages(t testing.TB, in string) {
	t.Helper()
	for _, p := range zzPairs() {
		g, ge := p.head(in)
		w, we := p.base(in)
		zzEqual(t, p.name, in, g, ge, w, we)
	}
}
func TestZZStageOracles(t *testing.T) { zzInputs(func(in string) { zzCompareStages(t, in) }) }
func TestZZStagePointerIdentity(t *testing.T) {
	for _, n := range []int{1, 16, 127, 128, 129, 255, 256, 257, 1024, 4095, 4096, 4097, 65536} {
		s := strings.Clone(strings.Repeat("az\x00\r\n\t\x7f", n/7+1)[:n])
		for _, p := range zzPairs()[:4] {
			g, ge := p.head(s)
			w, we := p.base(s)
			zzEqual(t, p.name, s, g, ge, w, we)
			alias := unsafe.StringData(w) == unsafe.StringData(s)

			t.Logf("%s n=%d BASE aliases input=%v candidate=%v", p.name, n, alias, unsafe.StringData(g) == unsafe.StringData(s))
		}
	}
}

// No inline: make parent liveness independent of the caller's frame/temporaries.
//
//go:noinline
func zzLiveResult(stage Stage, upper bool, n int) (string, weak.Pointer[byte]) {
	parent := strings.Repeat("a", 1<<20)
	if upper {
		parent = strings.Repeat("A", 1<<20)
	}
	w := weak.Make(unsafe.StringData(parent))
	out, err := stage(parent[64 : 64+n])
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(parent)
	return out, w
}
func TestZZLiveResult(t *testing.T) {
	for _, n := range []int{16, 127, 128, 129, 1024} {
		for _, p := range zzPairs()[:4] {
			for _, upper := range []bool{false, true} {
				for _, side := range []struct {
					name  string
					stage Stage
				}{{"BASE", p.base}, {"candidate", p.head}} {
					out, w := zzLiveResult(side.stage, upper, n)
					runtime.GC()
					runtime.GC()
					retained := w.Value() != nil
					want := (p.name == "NormalizeUnicode" && n <= 128) || ((p.name == "FoldCase" || p.name == "Lower") && !upper)
					if retained != want {
						t.Errorf("%s %s upper=%v n=%d parent retained=%v want=%v", side.name, p.name, upper, n, retained, want)
					}
					runtime.KeepAlive(out)
				}
			}
		}
	}
}
func TestZZCallbacks(t *testing.T) {
	for _, n := range []int{0, 1, 127, 128, 129, 255, 256, 257, 4095, 4096, 4097, 65536} {
		for _, unit := range []string{"aAaA", "é😀\xffa\ufffd", "\x80", "0000000000\xf2Ă"} {
			in := strings.Repeat(unit, n/len(unit)+1)[:n]
			for _, kind := range []string{"map", "filter"} {
				var gc, wc []rune
				mf := func(log *[]rune) func(rune) rune {
					return func(r rune) rune {
						*log = append(*log, r)
						switch len(*log) % 5 {
						case 0:
							return '😀'
						case 1:
							return r
						case 2:
							return -1
						case 3:
							return 'a'
						default:
							return '\ufffd'
						}
					}
				}
				keep := func(log *[]rune) runes.Set {
					return runes.Predicate(func(r rune) bool { *log = append(*log, r); return len(*log)%3 != 0 })
				}
				var h, b Stage
				if kind == "map" {
					h = New().MapRunes(mf(&gc)).Run
					b = baseMapRunesStage(mf(&wc))
				} else {
					h = New().FilterRunes(keep(&gc)).Run
					b = baseFilterRunesStage(keep(&wc))
				}
				for run := range 2 {
					g, ge := h(in)
					w, we := b(in)
					zzEqual(t, kind, in, g, ge, w, we)
					if !slices.Equal(gc, wc) {
						t.Fatalf("%s n=%d run=%d callback sequences differ: got=%d want=%d", kind, n, run, len(gc), len(wc))
					}
				}
			}
		}
	}
}
func TestZZNilAndImmutability(t *testing.T) {
	p := New().Then(func(s string) (string, error) { return s + "!", nil })
	for _, q := range []Pipeline{p.MapRunes(nil), p.FilterRunes(nil)} {
		if &q.stages[0] != &p.stages[0] {
			t.Fatal("nil callback/set changed pipeline")
		}
	}
	a := p.FoldCase()
	b := p.Lower()
	c := p.NormalizeUnicodeLatin()
	for _, q := range []Pipeline{a, b, c, p} {
		if _, e := q.Run("aAé"); e != nil {
			t.Fatal(e)
		}
	}
	if out, _ := p.Run("aAé"); out != "aAé!" {
		t.Fatal("parent pipeline mutated")
	}
}
func TestZZConcurrentStages(t *testing.T) {
	pairs := zzPairs()
	inputs := []string{strings.Repeat("Aé\xff😀ꮳ", 1000), strings.Repeat("aA\x00", 1000), "0000000000\xf2Ă"}
	wants := make([][]string, len(pairs))
	for i, p := range pairs {
		for _, s := range inputs {
			w, _ := p.base(s)
			wants[i] = append(wants[i], w)
		}
	}
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() {
			for range 20 {
				for i, p := range pairs {
					for k, in := range inputs {
						got, err := p.head(in)
						if err != nil || got != wants[i][k] {
							t.Errorf("concurrent %s mismatch", p.name)
							return
						}
					}
				}
			}
		})
	}
	wg.Wait()
}

// Every changed stage is exercised by the fuzzer. Stateful callback sequences
// and concurrency have dedicated deterministic tests above.
func FuzzZZStagesDifferential(f *testing.F) {
	for _, in := range []string{"", "A\x80", "ABC abc", "éक़مَ", "\xff", "0000000000\xf2Ă", "00000000000000000000000000000ꮒ0"} {
		f.Add(in)
	}
	f.Fuzz(func(t *testing.T, in string) { zzCompareStages(t, in) })
}
func TestZZSavedOutcomes(t *testing.T) {
	inputs := []string{"0000000000\xf2Ă", "00000000000000000000000000000ꮒ0"}
	for _, s := range inputs {
		for _, p := range []struct {
			name       string
			head, base Pipeline
		}{
			{"Pipeline", New().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace(), New().Then(baseNormalizeUnicodeStage).Then(baseFoldCaseStage()).TrimSpace().CollapseWhitespace()},
			{"CanonicalPreset", CanonicalPreset(), New().SanitizeUTF8().Then(baseNormalizeUnicodeStage).Then(baseFoldCaseStage()).TrimSpace().CollapseWhitespace()},
		} {
			h, he := p.head.Run(s)
			h2, he2 := p.head.Run(h)
			b, be := p.base.Run(s)
			b2, be2 := p.base.Run(b)
			if h != b || h2 != b2 || !errors.Is(he, be) || !errors.Is(be, he) || !errors.Is(he2, be2) || !errors.Is(be2, he2) {
				t.Fatalf("saved input mismatch %s %q", p.name, s)
			}
			msg := ""
			if h != h2 {
				prefix := "pipeline"
				if p.name == "CanonicalPreset" {
					prefix = p.name
				}
				msg = fmt.Sprintf("%s not idempotent: %q != %q", prefix, h, h2)
			}
			t.Logf("%s input=%q out1=%q out2=%q err1=%v err2=%v failure=%q", p.name, s, h, h2, he, he2, msg)
		}
	}
}

func TestZZASCIIBoundaryGrid(t *testing.T) {
	for ch := range 128 {
		for _, n := range []int{1, 2, 16, 127, 128, 129, 130, 255, 256, 257, 1024, 65536} {
			in := strings.Repeat(string(byte(ch)), n)
			for _, p := range []struct {
				name string
				base func(string) (string, error)
			}{{"NormalizeUnicode", baseNormalizeUnicodeStage}} {
				out, err := p.base(in)
				if err != nil || out != in {
					t.Fatalf("%s ASCII changed", p.name)
				}
				aliases := unsafe.StringData(in) == unsafe.StringData(out)
				if aliases != (n <= 128) {
					t.Fatalf("%s byte=%d n=%d aliases=%v expected=%v", p.name, ch, n, aliases, n <= 128)
				}
			}
		}
	}
}

// Allocation diagnostics only. These are not timing measurements or acceptance benchmarks.
func TestZZAllocationDiagnostics(t *testing.T) {
	for _, n := range []int{16, 1024, 65536} {
		in := strings.Repeat("a", n)
		for _, p := range zzPairs()[:4] {
			var sink string
			ba := testing.AllocsPerRun(20, func() { sink, _ = p.base(in) })
			ha := testing.AllocsPerRun(20, func() { sink, _ = p.head(in) })
			runtime.KeepAlive(sink)
			t.Logf("%s lowercase n=%d BASE=%g candidate=%g allocs/run", p.name, n, ba, ha)
		}
	}
	for _, p := range []struct {
		name       string
		head, base func() Pipeline
	}{
		{"MapRunes", func() Pipeline { return New().MapRunes(unicode.ToLower) }, func() Pipeline { return New().zzBaseMapRunes(unicode.ToLower) }},
		{"FilterRunes", func() Pipeline { return New().FilterRunes(runes.In(unicode.Letter)) }, func() Pipeline { return New().zzBaseFilterRunes(runes.In(unicode.Letter)) }},
	} {
		var sink Pipeline
		ba := testing.AllocsPerRun(20, func() { sink = p.base() })
		ha := testing.AllocsPerRun(20, func() { sink = p.head() })
		runtime.KeepAlive(sink)
		t.Logf("%s build (escaping result) BASE=%g candidate=%g allocs/run", p.name, ba, ha)
	}
}
func (p Pipeline) zzBaseMapRunes(fn func(rune) rune) Pipeline {
	if fn == nil {
		return p
	}

	return p.Then(func(s string) (string, error) {
		result, _, err := transform.String(transform.Chain(runes.Map(fn)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}

func (p Pipeline) zzBaseFilterRunes(keep runes.Set) Pipeline {
	if keep == nil {
		return p
	}

	return p.Then(func(s string) (string, error) {
		remove := runes.Predicate(func(r rune) bool {
			return !keep.Contains(r)
		})
		result, _, err := transform.String(transform.Chain(runes.Remove(remove)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}

func TestZZBuildDiscardedAllocationDiagnostics(t *testing.T) {
	ba := testing.AllocsPerRun(100, func() { New().zzBaseFilterRunes(benchmarkRuneSet{}) })
	ha := testing.AllocsPerRun(100, func() { New().FilterRunes(benchmarkRuneSet{}) })
	t.Logf("FilterRunes build (discarded result, existing fixture) BASE=%g candidate=%g allocs/run", ba, ha)
}

// Use the existing benchmark fixtures without invoking the benchmarks.
func TestZZRuneStageAllocationDiagnostics(t *testing.T) {
	for _, n := range benchkit.Sizes {
		in := benchkit.Text(n, benchkit.Unicode)
		for _, p := range []struct {
			name       string
			head, base Pipeline
		}{
			{"MapRunes", New().MapRunes(unicode.ToLower), New().zzBaseMapRunes(unicode.ToLower)},
			{"FilterRunes", New().FilterRunes(benchmarkRuneSet{}), New().zzBaseFilterRunes(benchmarkRuneSet{})},
		} {
			var sink string
			ba := testing.AllocsPerRun(100, func() { sink, _ = p.base.Run(in) })
			ha := testing.AllocsPerRun(100, func() { sink, _ = p.head.Run(in) })
			runtime.KeepAlive(sink)
			t.Logf("%s mode=run n=%d BASE=%g candidate=%g allocs/run", p.name, n, ba, ha)
		}
	}
	ba := testing.AllocsPerRun(100, func() { New().zzBaseMapRunes(unicode.ToLower) })
	ha := testing.AllocsPerRun(100, func() { New().MapRunes(unicode.ToLower) })
	t.Logf("MapRunes build (discarded result, existing fixture) BASE=%g candidate=%g allocs/run", ba, ha)
}
