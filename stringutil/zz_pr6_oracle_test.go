package stringutil

import (
	"errors"
	"math/rand/v2"
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unsafe"
	"weak"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// BASE oracles retain the production bodies from commit 4c05787.
func baseNormalizeUnicode(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	// Chain: NFKD decomposition → remove combining marks (diacritics)
	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

func baseToASCII(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mn)),
		runes.Map(func(r rune) rune {
			if r <= 127 {
				return r
			}
			// Non-ASCII characters that survived diacritic removal
			// get replaced with space
			return -1
		}),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
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

func TestZZStringOracles(t *testing.T) {
	zzInputs(func(in string) {
		for _, p := range []struct {
			name       string
			head, base func(string) (string, error)
		}{
			{"NormalizeUnicode", NormalizeUnicode, baseNormalizeUnicode}, {"ToASCII", ToASCII, baseToASCII},
		} {
			g, ge := p.head(in)
			w, we := p.base(in)
			zzEqual(t, p.name, in, g, ge, w, we)
		}
	})
}
func TestZZStringPointerIdentity(t *testing.T) {
	for _, n := range []int{1, 16, 127, 128, 129, 255, 256, 257, 1024, 4095, 4096, 4097, 65536} {
		s := strings.Clone(strings.Repeat("aZ\x00\r\n\t\x7f", n/7+1)[:n])
		for _, p := range []struct {
			name       string
			head, base func(string) (string, error)
		}{
			{"NormalizeUnicode", NormalizeUnicode, baseNormalizeUnicode}, {"ToASCII", ToASCII, baseToASCII},
		} {
			g, ge := p.head(s)
			w, we := p.base(s)
			zzEqual(t, p.name, s, g, ge, w, we)

			t.Logf("%s n=%d BASE aliases input=%v candidate=%v", p.name, n, unsafe.StringData(w) == unsafe.StringData(s), unsafe.StringData(g) == unsafe.StringData(s))
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
			}{{"NormalizeUnicode", baseNormalizeUnicode}, {"ToASCII", baseToASCII}} {
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

//go:noinline
func zzStringLiveResult(stage func(string) (string, error), n int) (string, weak.Pointer[byte]) {
	parent := strings.Repeat("a", 1<<20)
	w := weak.Make(unsafe.StringData(parent))
	out, err := stage(parent[64 : 64+n])
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(parent)
	return out, w
}
func TestZZStringLiveResult(t *testing.T) {
	for _, n := range []int{16, 127, 128, 129, 1024} {
		for _, p := range []struct {
			name       string
			head, base func(string) (string, error)
		}{
			{"NormalizeUnicode", NormalizeUnicode, baseNormalizeUnicode}, {"ToASCII", ToASCII, baseToASCII},
		} {
			for _, fn := range []func(string) (string, error){p.base, p.head} {
				out, w := zzStringLiveResult(fn, n)
				runtime.GC()
				runtime.GC()
				if retained := w.Value() != nil; retained != (n <= 128) {
					t.Errorf("%s n=%d parent retained=%v want=%v", p.name, n, retained, n <= 128)
				}
				runtime.KeepAlive(out)
			}
		}
	}
}

// Diagnostics use the matrix's actual fixtures; they do not execute benchmarks.
func TestZZMatrixAllocationDiagnostics(t *testing.T) {
	for _, n := range benchkit.Sizes {
		in := benchkit.Text(n, benchkit.ASCII)
		for _, p := range []struct {
			name string
			run  func(string) (string, error)
		}{
			{"NormalizeUnicode", NormalizeUnicode}, {"RemoveAccents", RemoveAccents},
			{"ToASCII", ToASCII}, {"Slugify", Slugify},
		} {
			var sink string
			allocs := testing.AllocsPerRun(100, func() {
				var err error
				sink, err = p.run(in)
				if err != nil {
					panic(err)
				}
			})
			runtime.KeepAlive(sink)
			target := 1.0
			if n == 16 {
				target = 0
			}
			if p.name == "Slugify" {
				target++
			}
			t.Logf("%s case=ascii n=%d allocs/run=%g target<=%g meets=%v", p.name, n, allocs, target, allocs <= target)
		}
	}
}
