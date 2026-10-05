package stringutil

import (
	"runtime"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
	"unsafe"
)

var commonAffixCases = []struct {
	name           string
	inputs         []string
	prefix, suffix string
}{
	{"zero", nil, "", ""},
	{"single empty", []string{""}, "", ""},
	{"single valid", []string{"é😀"}, "é😀", "é😀"},
	{"single malformed", []string{"a\xff\xe2\x82"}, "a\xff\xe2\x82", "a\xff\xe2\x82"},
	{"empty first", []string{"", "é"}, "", ""},
	{"empty last", []string{"é", ""}, "", ""},
	{"empty middle", []string{"é", "", "é"}, "", ""},
	{"ASCII", []string{"prefix-one-end", "prefix-two-end", "prefix-three-end"}, "prefix-", "-end"},
	{"identical", []string{"é😀", "é😀", "é😀"}, "é😀", "é😀"},
	{"split two-byte rune", []string{"é", "ê"}, "", ""},
	{"split three-byte rune", []string{"世", "丗"}, "", ""},
	{"split supplementary rune", []string{"😀", "😁"}, "", ""},
	{"shared continuation byte", []string{"é", "ĩ"}, "", ""},
	{"surrounding ASCII", []string{"AéZ", "AêZ"}, "A", "Z"},
	{"supplementary affixes", []string{"😀a𐀀", "😀b𐀀"}, "😀", "𐀀"},
	{"combining prefix", []string{"e\u0301", "e\u0300"}, "e", ""},
	{"combining suffix", []string{"a\u0301", "e\u0301"}, "", "\u0301"},
	{"no normalization", []string{"é", "e\u0301"}, "", ""},
	{"third prefix boundary", []string{"é😀x", "é😀y", "ê😀z"}, "", ""},
	{"third suffix boundary", []string{"x😀é", "y😀é", "z😀ê"}, "", ""},
	{"shortest valid", []string{"é", "é😀", "éx"}, "é", ""},
	{"different invalid bytes", []string{"\xff", "\xfe"}, "�", "�"},
	{"same invalid byte", []string{"\xff", "\xff"}, "�", "�"},
	{"invalid and replacement", []string{"\xff", "�"}, "�", "�"},
	{"replacement and invalid", []string{"�", "\xff"}, "�", "�"},
	{"different encoded lengths", []string{"\xffa\xfe", "�a�", "\x80a�"}, "�a�", "�a�"},
	{"three invalid inputs", []string{"\xffx\xfe", "�y�", "\x80z\xc0"}, "�", "�"},
	{"invalid outside prefix", []string{"abcX\xff", "abcY\xfe"}, "abc", "�"},
	{"invalid outside suffix", []string{"\xffXabc", "\xfeYabc"}, "�", "abc"},
	{"invalid before mismatch", []string{"\xffé", "\xffê"}, "�", ""},
	{"invalid after mismatch", []string{"é\xff", "ê\xff"}, "", "�"},
	{"truncated two-byte", []string{"\xc3", "\xc3"}, "�", "�"},
	{"truncated three-byte", []string{"\xe2\x82", "��"}, "��", "��"},
	{"truncated four-byte", []string{"\xf0\x9f\x92", "\xff\xfe\x80"}, "���", "���"},
	{"truncated versus complete", []string{"\xe2\x82", "€"}, "", ""},
	{"truncated versus complete reversed", []string{"€", "\xe2\x82"}, "", ""},
	{"valid versus continuation", []string{"é", "\xa9"}, "", ""},
	{"overlong", []string{"\xc0\xaf", "��"}, "��", "��"},
	{"surrogate", []string{"\xed\xa0\x80", "���"}, "���", "���"},
	{"out of range", []string{"\xf4\x90\x80\x80", "����"}, "����", "����"},
	{"invalid run lengths", []string{"\xff\xfe", "�"}, "�", "�"},
	{"invalid run lengths reversed", []string{"�", "\xff\xfe"}, "�", "�"},
	{"broken between valid runes", []string{"é\x80😀", "é�😀"}, "é�😀", "é�😀"},
	{"malformed prefix shortened later", []string{"\xffabc", "�abc", "�ax"}, "�a", ""},
	{"malformed suffix shortened later", []string{"abc\xff", "abc�", "xbc�"}, "", "bc�"},
}

// The oracle deliberately materializes every input's runes and compares them
// directly, independently of the production byte offsets and decoding scans.
func commonAffixOracle(inputs []string, suffix bool) string {
	if len(inputs) == 0 {
		return ""
	}
	if len(inputs) == 1 {
		return inputs[0]
	}
	common := []rune(inputs[0])
	for _, s := range inputs[1:] {
		runes := []rune(s)
		n := 0
		for n < len(common) && n < len(runes) {
			i, j := n, n
			if suffix {
				i, j = len(common)-1-n, len(runes)-1-n
			}
			if common[i] != runes[j] {
				break
			}
			n++
		}
		if suffix {
			common = common[len(common)-n:]
		} else {
			common = common[:n]
		}
	}
	return string(common)
}

func checkCommonAffixes(t testing.TB, inputs []string) {
	t.Helper()
	for _, tc := range []struct {
		name   string
		fn     func(...string) string
		suffix bool
	}{{"prefix", CommonPrefix, false}, {"suffix", CommonSuffix, true}} {
		got, want := tc.fn(inputs...), commonAffixOracle(inputs, tc.suffix)
		if got != want {
			t.Errorf("%s(%q) = %q, rune oracle = %q", tc.name, inputs, got, want)
		}
		if len(inputs) > 1 && !utf8.ValidString(got) {
			t.Errorf("%s(%q) returned invalid UTF-8: %q", tc.name, inputs, got)
		}
		for _, s := range inputs {
			if !utf8.ValidString(s) {
				continue
			}
			if (!tc.suffix && !strings.HasPrefix(s, got)) || (tc.suffix && !strings.HasSuffix(s, got)) {
				t.Errorf("%s result %q is not a byte affix of valid input %q", tc.name, got, s)
			}
		}
	}
}

func TestCommonAffixRunes(t *testing.T) {
	for _, tc := range commonAffixCases {
		t.Run(tc.name, func(t *testing.T) {
			if p, s := commonAffixOracle(tc.inputs, false), commonAffixOracle(tc.inputs, true); p != tc.prefix || s != tc.suffix {
				t.Fatalf("oracle = (%q, %q), explicit expectations = (%q, %q)", p, s, tc.prefix, tc.suffix)
			}
			checkCommonAffixes(t, tc.inputs)
			// Input order must not change multi-input decoded results.
			reversed := slices.Clone(tc.inputs)
			slices.Reverse(reversed)
			checkCommonAffixes(t, reversed)
		})
	}
}

func checkCommonAffixSegmentation(t testing.TB, s string) {
	t.Helper()
	type segment struct {
		r     rune
		width int
	}
	var forward, backward []segment
	var decoded []rune
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		forward = append(forward, segment{r, width})
		decoded = append(decoded, r)
		i += width
	}
	for i := len(s); i > 0; {
		r, width := utf8.DecodeLastRuneInString(s[:i])
		backward = append(backward, segment{r, width})
		i -= width
	}
	slices.Reverse(backward)
	if !slices.Equal(forward, backward) || !slices.Equal(decoded, []rune(s)) {
		t.Fatalf("segmentation differs for %q: forward=%v backward=%v runes=%U", s, forward, backward, []rune(s))
	}
}

func TestCommonAffixRuneSegmentation(t *testing.T) {
	for _, tc := range commonAffixCases {
		for _, s := range tc.inputs {
			checkCommonAffixSegmentation(t, s)
		}
	}
}

func TestCommonAffixDetachment(t *testing.T) {
	for _, tc := range []struct {
		name, input, other string
		fn                 func(...string) string
	}{
		{"prefix short", "é😀abc", "é😀xyz", CommonPrefix},
		{"prefix full", "é😀abc", "é😀abc", CommonPrefix},
		{"prefix malformed", "\xffabc", "�xyz", CommonPrefix},
		{"prefix malformed full", "\xffabc", "�abc", CommonPrefix},
		{"prefix discarded invalid", "abcX\xff", "abcY\xfe", CommonPrefix},
		{"suffix short", "abcé😀", "xyzé😀", CommonSuffix},
		{"suffix full", "abcé😀", "abcé😀", CommonSuffix},
		{"suffix malformed", "abc\xff", "xyz�", CommonSuffix},
		{"suffix malformed full", "abc\xff", "abc�", CommonSuffix},
		{"suffix discarded invalid", "\xffXabc", "\xfeYabc", CommonSuffix},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, first := range []bool{true, false} {
				out, w, shares := zzRetainedInput(tc.input, func(s string) string {
					if first {
						return tc.fn(s, tc.other)
					}
					return tc.fn(tc.other, s)
				})
				runtime.GC()
				runtime.GC()
				if out == "" || shares || w.Value() != nil {
					t.Errorf("result %q retains 1 MiB input or fixture is empty (input first=%v)", out, first)
				}
				runtime.KeepAlive(out)
			}
		})
	}
	for _, fn := range []func(...string) string{CommonPrefix, CommonSuffix} {
		parent := strings.Repeat("a\xff", 2048)
		s := parent[19:40]
		if got := fn(s); got != s || unsafe.StringData(got) != unsafe.StringData(s) {
			t.Fatal("single malformed input must be returned unchanged at its original address")
		}
		runtime.KeepAlive(parent)
	}
}

func TestCommonAffixAllocations(t *testing.T) {
	// Allocation counts only: timing benchmarks are run separately.
	for _, tc := range []struct {
		name   string
		inputs []string
		fn     func(...string) string
	}{
		{"prefix ASCII", []string{"shared-abc", "shared-def", "shared-xyz"}, CommonPrefix},
		{"prefix Unicode", []string{"é😀a", "é😀b", "é😀c"}, CommonPrefix},
		{"prefix malformed", []string{"é\xffa", "é�b", "é\xfec"}, CommonPrefix},
		{"suffix ASCII", []string{"abc-shared", "def-shared", "xyz-shared"}, CommonSuffix},
		{"suffix Unicode", []string{"aé😀", "bé😀", "cé😀"}, CommonSuffix},
		{"suffix malformed", []string{"a\xff😀", "b�😀", "c\xfe😀"}, CommonSuffix},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out string
			allocs := testing.AllocsPerRun(100, func() { out = tc.fn(tc.inputs...) })
			if out == "" || allocs > 1 {
				t.Errorf("nonempty result %q uses %g allocations, want at most one", out, allocs)
			}
			t.Logf("allocations per call: %g", allocs)
			runtime.KeepAlive(out)
		})
	}
}

func FuzzCommonAffixRunes(f *testing.F) {
	for _, tc := range commonAffixCases {
		var inputs [3]string
		for i := range inputs {
			if len(tc.inputs) > 0 {
				inputs[i] = tc.inputs[i%len(tc.inputs)]
			}
		}
		f.Add(inputs[0], inputs[1], inputs[2])
	}
	f.Fuzz(func(t *testing.T, a, b, c string) {
		// Bound oracle allocations without excluding malformed byte boundaries.
		a, b, c = a[:min(len(a), 4096)], b[:min(len(b), 4096)], c[:min(len(c), 4096)]
		for _, s := range []string{a, b, c} {
			checkCommonAffixSegmentation(t, s)
		}
		inputs := []string{a, b, c, a}
		for n := range len(inputs) + 1 {
			checkCommonAffixes(t, inputs[:n])
		}
		checkCommonAffixes(t, []string{c, b, a})
		checkCommonAffixes(t, []string{a, string([]rune(a))})
	})
}
