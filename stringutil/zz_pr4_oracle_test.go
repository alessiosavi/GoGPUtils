package stringutil

// Frozen bodies from 4c05787953a944b57478180f4c27100869db12f6.
// Only private names and references to other copied bodies are changed.
import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
	"unsafe"
	"weak"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

type zzBaseCleanOption func(*zzBaseCleanConfig)
type zzBaseCleanConfig struct {
	unicodeNorm, htmlStrip, dbSanitize bool
	dbMaxLen                           int
	dbReplaceNul                       bool
}

func zzBaseReverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func zzBaseTruncate(s string, maxLen int, suffix string) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	suffixRunes := []rune(suffix)
	if len(suffixRunes) >= maxLen {
		return string(suffixRunes[:maxLen])
	}

	truncateAt := maxLen - len(suffixRunes)

	return string(runes[:truncateAt]) + suffix
}

func zzBaseTruncateWords(s string, maxLen int, suffix string) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	suffixRunes := []rune(suffix)
	if len(suffixRunes) >= maxLen {
		return string(suffixRunes[:maxLen])
	}

	truncateAt := maxLen - len(suffixRunes)

	// Find last space before truncate point
	lastSpace := -1

	for i := truncateAt - 1; i >= 0; i-- {
		if unicode.IsSpace(runes[i]) {
			lastSpace = i

			break
		}
	}

	if lastSpace > 0 {
		return string(runes[:lastSpace]) + suffix
	}

	return string(runes[:truncateAt]) + suffix
}

func zzBasePadLeft(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	padding := length - len(runes)

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))

	for range padding {
		b.WriteRune(padChar)
	}

	b.WriteString(s)

	return b.String()
}

func zzBasePadRight(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	padding := length - len(runes)

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))
	b.WriteString(s)

	for range padding {
		b.WriteRune(padChar)
	}

	return b.String()
}

func zzBasePadCenter(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	totalPadding := length - len(runes)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))

	for range leftPadding {
		b.WriteRune(padChar)
	}

	b.WriteString(s)

	for range rightPadding {
		b.WriteRune(padChar)
	}

	return b.String()
}

func zzBaseCountLines(s string) int {
	if s == "" {
		return 0
	}

	count := 1

	for _, r := range s {
		if r == '\n' {
			count++
		}
	}
	// Don't count trailing newline as extra line
	if strings.HasSuffix(s, "\n") {
		count--
	}

	return count
}

func zzBaseCapitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func zzBasePascalCase(s string) string {
	result := zzBaseCamelCase(s)
	if result == "" {
		return ""
	}

	runes := []rune(result)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func zzBaseSafeSlice(s string, start, end int) string {
	runes := []rune(s)

	if start < 0 {
		start = 0
	}

	if end > len(runes) {
		end = len(runes)
	}

	if start >= end || start >= len(runes) {
		return ""
	}

	return string(runes[start:end])
}

func zzBaseCommonPrefix(strs ...string) string {
	if len(strs) < 2 {
		if len(strs) == 1 {
			return strs[0]
		}

		return ""
	}

	// Find shortest string to bound our search
	minLen := len(strs[0])
	for _, s := range strs[1:] {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	var prefix strings.Builder

	for i := range minLen {
		char := strs[0][i]
		for _, s := range strs[1:] {
			if s[i] != char {
				return prefix.String()
			}
		}

		prefix.WriteByte(char)
	}

	return prefix.String()
}

func zzBaseCommonSuffix(strs ...string) string {
	if len(strs) < 2 {
		if len(strs) == 1 {
			return strs[0]
		}

		return ""
	}

	// Reverse all strings, find prefix, then reverse result
	reversed := make([]string, len(strs))
	for i, s := range strs {
		reversed[i] = zzBaseReverse(s)
	}

	return zzBaseReverse(zzBaseCommonPrefix(reversed...))
}

func zzBaseBetween(s, start, end string) (string, bool) {
	startIdx := strings.Index(s, start)
	if startIdx == -1 {
		return "", false
	}

	startIdx += len(start)

	endIdx := strings.Index(s[startIdx:], end)
	if endIdx == -1 {
		return "", false
	}

	return s[startIdx : startIdx+endIdx], true
}

func zzBaseBetweenAll(s, start, end string) []string {
	if start == "" && end == "" {
		return nil
	}

	var results []string

	remaining := s

	for {
		result, ok := zzBaseBetween(remaining, start, end)
		if !ok {
			break
		}

		results = append(results, result)

		// Consume through both markers. At least one is nonempty, so this
		// strictly shortens remaining after each match.
		idx := strings.Index(remaining, start)
		remaining = remaining[idx+len(start):]
		idx = strings.Index(remaining, end)
		remaining = remaining[idx+len(end):]
	}

	return results
}

func zzBaseCamelCase(s string) string {
	var b strings.Builder

	b.Grow(len(s))

	capitalizeNext := false

	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			capitalizeNext = true

			continue
		}

		if capitalizeNext {
			b.WriteRune(unicode.ToUpper(r))

			capitalizeNext = false
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	return b.String()
}

func zzBaseTruncateRunes(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen])
}

func zzBaseCleanString(input string, options ...zzBaseCleanOption) (string, error) {
	if len(options) == 0 {
		return input, nil
	}

	cfg := &zzBaseCleanConfig{}
	for _, opt := range options {
		opt(cfg)
	}

	result := input

	// Step 1: HTML stripping (before text processing)
	if cfg.htmlStrip {
		result = StripHTMLEntities(result)
	}

	// Step 2: Unicode normalization
	if cfg.unicodeNorm {
		var err error

		result, err = NormalizeUnicode(result)
		if err != nil {
			return "", err
		}
	}

	// Step 3: Database sanitization (applied last)
	if cfg.dbSanitize {
		result = SanitizeUTF8(result)
		if cfg.dbMaxLen > 0 {
			result = zzBaseTruncateRunes(result, cfg.dbMaxLen)
		}
	}

	return result, nil
}

func zzBaseWithUnicodeNorm() zzBaseCleanOption {
	return func(c *zzBaseCleanConfig) {
		c.unicodeNorm = true
	}
}

func zzBaseWithHTMLStrip() zzBaseCleanOption {
	return func(c *zzBaseCleanConfig) {
		c.htmlStrip = true
	}
}

func zzBaseWithDBSanitize(maxLen int) zzBaseCleanOption {
	return func(c *zzBaseCleanConfig) {
		c.dbSanitize = true
		c.dbReplaceNul = true
		if maxLen > 0 {
			c.dbMaxLen = maxLen
		}
	}
}

type zzOutcome struct {
	Value                   any
	PanicType, PanicMessage string
}

func zzOutcomeOf(f func() any) (o zzOutcome) {
	defer func() {
		if p := recover(); p != nil {
			o.PanicType = fmt.Sprintf("%T", p)
			o.PanicMessage = fmt.Sprint(p)
		}
	}()
	o.Value = f()
	return
}
func zzCompare(t testing.TB, name string, got, want func() any) {
	t.Helper()
	g, w := zzOutcomeOf(got), zzOutcomeOf(want)
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("%s: got %#v; BASE %#v", name, g, w)
	}
}
func zzCheck(t testing.TB, s string, a, b int, r rune, suffix, start, end string) {
	t.Helper()
	for _, x := range []struct {
		name string
		g, w func(string) string
	}{
		{"Reverse", Reverse, zzBaseReverse}, {"Capitalize", Capitalize, zzBaseCapitalize}, {"PascalCase", PascalCase, zzBasePascalCase},
	} {
		zzCompare(t, x.name+fmt.Sprintf("(%q)", s), func() any { return x.g(s) }, func() any { return x.w(s) })
	}
	for _, x := range []struct {
		name string
		g, w func(string, int, string) string
	}{
		{"Truncate", Truncate, zzBaseTruncate}, {"TruncateWords", TruncateWords, zzBaseTruncateWords},
	} {
		zzCompare(t, x.name+fmt.Sprintf("(%q,%d,%q)", s, a, suffix), func() any { return x.g(s, a, suffix) }, func() any { return x.w(s, a, suffix) })
	}
	for _, x := range []struct {
		name string
		g, w func(string, int, rune) string
	}{
		{"PadLeft", PadLeft, zzBasePadLeft}, {"PadRight", PadRight, zzBasePadRight}, {"PadCenter", PadCenter, zzBasePadCenter},
	} {
		zzCompare(t, x.name+fmt.Sprintf("(%q,%d,%U)", s, a, r), func() any { return x.g(s, a, r) }, func() any { return x.w(s, a, r) })
	}
	zzCompare(t, "TruncateRunes", func() any { return TruncateRunes(s, a) }, func() any { return zzBaseTruncateRunes(s, a) })
	zzCompare(t, fmt.Sprintf("SafeSlice(%q,%d,%d)", s, a, b), func() any { return SafeSlice(s, a, b) }, func() any { return zzBaseSafeSlice(s, a, b) })
	zzCompare(t, "CountLines", func() any { return CountLines(s) }, func() any { return zzBaseCountLines(s) })
	zzCompare(t, "BetweenAll", func() any { return BetweenAll(s, start, end) }, func() any { return zzBaseBetweenAll(s, start, end) })
	for _, xs := range [][]string{nil, {s}, {s, suffix}, {s, suffix, start}, {s, suffix, start, end}} {
		zzCompare(t, "CommonPrefix", func() any { return CommonPrefix(xs...) }, func() any { return zzBaseCommonPrefix(xs...) })
		zzCompare(t, "CommonSuffix", func() any { return CommonSuffix(xs...) }, func() any { return zzBaseCommonSuffix(xs...) })
	}
	for mode := range 8 {
		var opts []CleanOption
		var old []zzBaseCleanOption
		if mode&1 != 0 {
			opts = append(opts, WithUnicodeNorm())
			old = append(old, zzBaseWithUnicodeNorm())
		}
		if mode&2 != 0 {
			opts = append(opts, WithHTMLStrip())
			old = append(old, zzBaseWithHTMLStrip())
		}
		if mode&4 != 0 {
			opts = append(opts, WithDBSanitize(a))
			old = append(old, zzBaseWithDBSanitize(a))
		}
		g, ge := CleanString(s, opts...)
		w, we := zzBaseCleanString(s, old...)
		if g != w || ge != we { //nolint:errorlint // The oracle contract requires identical errors, not just matching wrapped errors. //nolint:errorlint // The oracle contract requires identical errors, not just matching wrapped errors.
			t.Fatalf("CleanString mode=%d s=%q: got %q %v BASE %q %v", mode, s, g, ge, w, we)
		}
	}
}
func zzCorpus() []string {
	ss := []string{"", "a", " ", " a b c", "é a b c", "\u2003a b c", "a\xff c def", "a b\xffc def", "\xffabcde", "abcde\xff", "\r\n", "\n\n", "\xff\n", "a\n", "é", "è", "�", "Kİıſ", "ꮳᏣ", "👩🏽‍💻", "[one][two", "|a||b|", "aaaa", "<p>a\x00&amp;b</p>"}
	ss = append(ss, strings.Fields("café café naïve résumé Ærøskøbing straße Zürich ἀλήθεια λόγος Ελλάδα мир привет 東京 数据 日本語 キャッシュ 한국어 😀 👍🏽 ﬁne")...)
	for _, bad := range []string{"\x80", "\xc2", "\xe2\x82", "\xf0\x9f\x92", "\xc0\x80", "\xed\xa0\x80", "\xf4\x90\x80\x80", "\xff", "�"} {
		for i := range 11 {
			s := "a é b 世界"
			ss = append(ss, s[:i]+bad+s[i:])
		}
	}
	for n := range 41 {
		ss = append(ss, benchkit.Text(n, benchkit.ASCII), benchkit.Text(n, benchkit.Unicode))
	}
	for _, n := range []int{127, 128, 129, 1024} {
		ss = append(ss, benchkit.Text(n, benchkit.ASCII), benchkit.Text(n, benchkit.Unicode))
	}
	rng := rand.New(rand.NewPCG(4, 17))
	for range 150 {
		p := make([]byte, rng.IntN(48))
		for j := range p {
			p[j] = byte(rng.Uint32())
		}
		ss = append(ss, string(p))
	}
	return ss
}
func TestZZScansDifferential(t *testing.T) {
	for _, s := range zzCorpus() {
		n := utf8.RuneCountInString(s)
		if n != len([]rune(s)) {
			t.Fatal("rune count mismatch")
		}
		for _, suffix := range []string{"", "...", "…", "a very long suffix", "\xff", "a\xffxyz"} {
			k := utf8.RuneCountInString(suffix)
			for _, a := range []int{-1, 0, 1, n - 1, n, n + 1, k - 1, k, k + 1} {
				zzCheck(t, s, a, n/2, ' ', suffix, "[", "]")
			}
		}
		for _, r := range []rune{' ', 'é', '😀', '�', -1, 0xD800, 0x110000} {
			zzCheck(t, s, n+1, n-1, r, "", "", ",")
		}
		// Exhaustive index pairs on the small domain; sparse boundary pairs for large inputs.
		idx := []int{math.MinInt, -2, -1, 0, 1, n / 2, n - 1, n, n + 1, n + 2, math.MaxInt}
		if n <= 40 {
			for i := -2; i <= n+2; i++ {
				idx = append(idx, i)
			}
		}
		for _, a := range idx {
			for _, b := range idx {
				zzCompare(t, fmt.Sprintf("SafeSlice(%q,%d,%d)", s, a, b), func() any { return SafeSlice(s, a, b) }, func() any { return zzBaseSafeSlice(s, a, b) })
			}
		}
		for _, a := range []int{math.MinInt, math.MaxInt} {
			zzCompare(t, "extreme Truncate", func() any { return Truncate(s, a, "") }, func() any { return zzBaseTruncate(s, a, "") })
			zzCompare(t, "extreme TruncateWords", func() any { return TruncateWords(s, a, "") }, func() any { return zzBaseTruncateWords(s, a, "") })
			zzCompare(t, "extreme TruncateRunes", func() any { return TruncateRunes(s, a) }, func() any { return zzBaseTruncateRunes(s, a) })
		}
		for _, marker := range []string{"", "a", "aa", "[", "]", "\xff"} {
			for _, end := range []string{"", marker, "a", "]"} {
				zzCompare(t, "markers", func() any { return BetweenAll(s, marker, end) }, func() any { return zzBaseBetweenAll(s, marker, end) })
			}
		}
	}
	// Overflow before allocation: Grow must remain verbatim; avoid OOM inputs.
	for _, f := range []struct {
		g, w func(string, int, rune) string
	}{{PadLeft, zzBasePadLeft}, {PadRight, zzBasePadRight}, {PadCenter, zzBasePadCenter}} {
		zzCompare(t, "Grow overflow", func() any { return f.g("", math.MaxInt, 'é') }, func() any { return f.w("", math.MaxInt, 'é') })
	}
}
func zzShares(s, parent string) bool {
	if len(s) == 0 || len(parent) == 0 {
		return false
	}
	p, q := uintptr(unsafe.Pointer(unsafe.StringData(s))), uintptr(unsafe.Pointer(unsafe.StringData(parent)))
	return p >= q && p < q+uintptr(len(parent))
}
func TestZZScansOwnership(t *testing.T) {
	parent := strings.Repeat("abc def ", 128)
	s := parent[40:64]
	suffix := parent[72:90]
	for _, tc := range []struct {
		name string
		g, w func() string
	}{
		{"TruncateEmptySuffix", func() string { return Truncate(s, 3, "") }, func() string { return zzBaseTruncate(s, 3, "") }},
		{"WordsEmptySuffix", func() string { return TruncateWords(s, 7, "") }, func() string { return zzBaseTruncateWords(s, 7, "") }},
		{"TruncateSuffixOnly", func() string { return Truncate(s, 3, suffix) }, func() string { return zzBaseTruncate(s, 3, suffix) }},
		{"WordsSuffixOnly", func() string { return TruncateWords(s, 3, suffix) }, func() string { return zzBaseTruncateWords(s, 3, suffix) }},
		{"TruncateSuffix", func() string { return Truncate(s, 8, "...") }, func() string { return zzBaseTruncate(s, 8, "...") }},
		{"WordsSuffix", func() string { return TruncateWords(s, 8, "...") }, func() string { return zzBaseTruncateWords(s, 8, "...") }},
		{"TruncateNoop", func() string { return Truncate(s, 100, "") }, func() string { return zzBaseTruncate(s, 100, "") }},
		{"WordsNoop", func() string { return TruncateWords(s, 100, "") }, func() string { return zzBaseTruncateWords(s, 100, "") }},
		{"RunesNoop", func() string { return TruncateRunes(s, 100) }, func() string { return zzBaseTruncateRunes(s, 100) }},
		{"Runes", func() string { return TruncateRunes(s, 3) }, func() string { return zzBaseTruncateRunes(s, 3) }},
		{"SafeFull", func() string { return SafeSlice(s, 0, 100) }, func() string { return zzBaseSafeSlice(s, 0, 100) }},
		{"SafePart", func() string { return SafeSlice(s, 2, 7) }, func() string { return zzBaseSafeSlice(s, 2, 7) }},
		{"Prefix1", func() string { return CommonPrefix(s) }, func() string { return zzBaseCommonPrefix(s) }},
		{"Prefix2", func() string { return CommonPrefix(s, s) }, func() string { return zzBaseCommonPrefix(s, s) }},
		{"Reverse", func() string { return Reverse(s) }, func() string { return zzBaseReverse(s) }},
		{"Capitalize", func() string { return Capitalize(s) }, func() string { return zzBaseCapitalize(s) }},
		{"PascalCase", func() string { return PascalCase(s) }, func() string { return zzBasePascalCase(s) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g, w := tc.g(), tc.w()
			if g != w || zzShares(g, parent) != zzShares(w, parent) {
				t.Fatalf("got share=%v, BASE share=%v", zzShares(g, parent), zzShares(w, parent))
			}
		})
	}
	for _, f := range []struct {
		g, w func(string, int, rune) string
	}{{PadLeft, zzBasePadLeft}, {PadRight, zzBasePadRight}, {PadCenter, zzBasePadCenter}} {
		for _, n := range []int{0, 100} {
			g, w := f.g(s, n, ' '), f.w(s, n, ' ')
			if g != w || zzShares(g, parent) != zzShares(w, parent) {
				t.Fatal("pad ownership")
			}
		}
	}
	g, w := BetweenAll(s, "a", " "), zzBaseBetweenAll(s, "a", " ")
	for i := range g {
		if unsafe.StringData(g[i]) != unsafe.StringData(w[i]) {
			t.Fatal("BetweenAll ownership")
		}
	}
	runtime.KeepAlive(parent)
}

//go:noinline
func zzRetained(f func(string) string) (string, weak.Pointer[byte]) {
	p := new([1 << 20]byte)
	copy(p[100:], "abc def ghi jkl mnop")
	s := unsafe.String(&p[100], 19)
	w := weak.Make(&p[0])
	out := f(s)
	runtime.KeepAlive(p)
	return out, w
}
func TestZZScansRetention(t *testing.T) {
	for _, tc := range []struct {
		name string
		f    func(string) string
	}{
		{"Runes", func(s string) string { return TruncateRunes(s, 3) }},
		{"Safe", func(s string) string { return SafeSlice(s, 1, 4) }},
		{"SafeFull", func(s string) string { return SafeSlice(s, 0, 100) }},
		{"Truncate", func(s string) string { return Truncate(s, 3, "") }},
		{"TruncateWords", func(s string) string { return TruncateWords(s, 7, "") }},
		{"TruncateSuffix", func(s string) string { return Truncate(s, 8, "...") }},
		{"WordsSuffix", func(s string) string { return TruncateWords(s, 8, "...") }},
		{"PrefixFull", func(s string) string { return CommonPrefix(s, s) }},
		{"TruncateSuffixOnly", func(s string) string { return Truncate("abcdefghijklmnopqrstu", 3, s) }},
		{"WordsSuffixOnly", func(s string) string { return TruncateWords("abcdefghijklmnopqrstu", 3, s) }},
		{"Prefix", func(s string) string { return CommonPrefix(s, "abcXYZ") }},
		{"Reverse", Reverse}, {"Capitalize", Capitalize}, {"PascalCase", PascalCase},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, w := zzRetained(tc.f)
			runtime.GC()
			runtime.GC()
			if w.Value() != nil {
				t.Error("live result retains 1 MiB parent")
			}
			runtime.KeepAlive(out)
		})
	}
}
func FuzzZZScansDifferential(f *testing.F) {
	for _, s := range zzCorpus()[:55] {
		f.Add(s, 3, 7, rune(' '))
	}
	f.Fuzz(func(t *testing.T, data string, a, b int, r rune) {
		// Bound allocation/time without abs(MinInt) overflow. Preserve negative indices.
		if len(data) > 2048 {
			data = data[:2048]
		}
		a %= utf8.RuneCountInString(data) + 5
		b %= utf8.RuneCountInString(data) + 5
		split := len(data) / 2
		suffix := data[split:]
		start := data[:min(2, len(data))]
		end := data[max(0, len(data)-2):]
		zzCheck(t, data, a, b, r, suffix, start, end)
		zzCheck(t, data, a, b, r, "", "", end)
	})
}

func TestZZCutBoundaries(t *testing.T) {
	for _, s := range []string{" a b c", "\u2003a b c", "é a b c", "\xffabcde", "abcde\xff", "a\xff c def", "a b\xffc def", "� é x y", "a\xed\xa0\x80 b c", "a\xf4\x90\x80\x80 b c"} {
		for k := -2; k <= utf8.RuneCountInString(s)+2; k++ {
			for _, suffix := range []string{"", ".", "é", "\xff", "abc\xff", "�"} {
				zzCompare(t, fmt.Sprintf("Truncate(%q,%d,%q)", s, k, suffix), func() any { return Truncate(s, k, suffix) }, func() any { return zzBaseTruncate(s, k, suffix) })
				zzCompare(t, fmt.Sprintf("Words(%q,%d,%q)", s, k, suffix), func() any { return TruncateWords(s, k, suffix) }, func() any { return zzBaseTruncateWords(s, k, suffix) })
			}
		}
	}
	for i := range 256 {
		s := string([]byte{'a', byte(i), 'b', ' ', 'c'})
		if !utf8.ValidString(strings.ToLower(s)) || !utf8.ValidString(CamelCase(s)) {
			t.Fatalf("invalid case output: %q", s)
		}
		for k := range 7 {
			zzCompare(t, "byte TruncateRunes", func() any { return TruncateRunes(s, k) }, func() any { return zzBaseTruncateRunes(s, k) })
		}
	}
}
func TestZZPlanLayout(t *testing.T) {
	t.Logf("cleanConfig=%d frozen=%d dbMaxLen offset=%d", unsafe.Sizeof(cleanConfig{}), unsafe.Sizeof(zzBaseCleanConfig{}), unsafe.Offsetof(cleanConfig{}.dbMaxLen))
	// A custom option can retain the config; field order must not alter semantics.
	var cfg *cleanConfig
	out, err := CleanString("abc", func(c *cleanConfig) { cfg = c; c.dbMaxLen = 7 }, func(c *cleanConfig) {
		if c != cfg || c.dbMaxLen != 7 {
			t.Fatal("option identity/order")
		}
	})
	if out != "abc" || err != nil || cfg == nil {
		t.Fatal("custom options")
	}
}

func TestZZKnownBehaviorUnchanged(t *testing.T) {
	if got := CommonPrefix("é", "è"); got != "\xc3" {
		t.Fatalf("C-04 must remain unchanged: %q", got)
	}
	if got := CommonSuffix("é", "è"); got != "�" {
		t.Fatalf("C-04 CommonSuffix must remain unchanged: %q", got)
	}

	if got := Truncate("abcdef", 3, "\xff"); got != "ab\xff" {
		t.Fatalf("raw suffix must not be repaired: %q", got)
	}
	if got := Truncate("abcdef", 2, "xy\xff"); got != "xy" {
		t.Fatalf("discarded suffix bytes: %q", got)
	}
	if got := TruncateWords("ab cd\xffZZ", 6, ""); got != "ab" {
		t.Fatalf("discarded word bytes: %q", got)
	}
	if got := SafeSlice("\xffabc\xff", 1, 4); got != "abc" {
		t.Fatalf("discarded slice bytes: %q", got)
	}
}

func TestZZScansReencoding(t *testing.T) {
	// Explicit expected bytes ensure the oracle is not the only specification.
	cases := []struct{ name, raw, want string }{
		{"continuation", "\x80", "�"},
		{"truncated2", "\xc2", "�"},
		{"truncated3", "\xe2\x82", "��"},
		{"truncated4", "\xf0\x9f\x92", "���"},
		{"overlong2", "\xc0\x80", "��"},
		{"overlong3", "\xe0\x80\xaf", "���"},
		{"overlong4", "\xf0\x80\x80\xaf", "����"},
		{"surrogate", "\xed\xa0\x80", "���"},
		{"aboveMax", "\xf4\x90\x80\x80", "����"},
		{"ff", "\xff", "�"},
		{"genuine", "�", "�"},
		{"mixed", "�\xffé\xed\xa0\x80", "��é���"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if string([]rune(tc.raw)) != tc.want {
				t.Fatal("incorrect expected bytes")
			}
			raw := "ab" + tc.raw
			want := "ab" + tc.want
			k := utf8.RuneCountInString(raw)
			for _, suffix := range []string{"", ".", "\xff", "�\xff"} {
				n := k + utf8.RuneCountInString(suffix)
				input := raw + " " + strings.Repeat("x", n+1)
				for _, fn := range []func(string, int, string) string{Truncate, TruncateWords} {
					if got := fn(input, n, suffix); got != want+suffix {
						t.Fatalf("emitted %q want %q", got, want+suffix)
					}
					if got := fn(strings.Repeat("x", k+1), k, raw+"\xff"); got != want {
						t.Fatalf("suffix-only %q want %q", got, want)
					}
				}
			}
			if got := TruncateRunes(raw+"tail", k); got != want {
				t.Fatalf("runes %q want %q", got, want)
			}
			if got := SafeSlice("\xff"+raw+"\xff", 1, k+1); got != want {
				t.Fatalf("slice %q want %q", got, want)
			}
			if got := SafeSlice("\xffabc"+tc.raw, 1, 4); got != "abc" {
				t.Fatalf("outside %q", got)
			}
			if got := TruncateWords("ab "+tc.raw+"tail", k+1, ""); got != "ab" {
				t.Fatalf("word discarded %q", got)
			}
			for _, pad := range []string{"", strings.Repeat("z", 31), strings.Repeat("é", 33)} {
				s := pad + tc.raw + pad
				zzCompare(t, "Reverse long/short", func() any { return Reverse(s) }, func() any { return zzBaseReverse(s) })
			}
			// Every rune boundary, including cuts inside malformed byte sequences.
			s := tc.raw + "é � \xff tail"
			for i := 0; i <= utf8.RuneCountInString(s); i++ {
				zzCompare(t, "rune boundary", func() any { return TruncateRunes(s, i) }, func() any { return zzBaseTruncateRunes(s, i) })
				for j := i; j <= utf8.RuneCountInString(s); j++ {
					zzCompare(t, "slice boundary", func() any { return SafeSlice(s, i, j) }, func() any { return zzBaseSafeSlice(s, i, j) })
				}
			}
		})
	}
}

func TestZZScansBytePairs(t *testing.T) {
	for a := range 256 {
		for b := range 256 {
			s := "a" + string([]byte{byte(a), byte(b)}) + " é tail"
			for _, k := range []int{2, 3, 4, 5} {
				zzCompare(t, fmt.Sprintf("pair %x %x cut %d", a, b, k), func() any { return TruncateRunes(s, k) }, func() any { return zzBaseTruncateRunes(s, k) })
				zzCompare(t, "words pair", func() any { return TruncateWords(s, k, "") }, func() any { return zzBaseTruncateWords(s, k, "") })
				zzCompare(t, "slice pair", func() any { return SafeSlice(s, 1, k) }, func() any { return zzBaseSafeSlice(s, 1, k) })
			}
		}
	}
}

//go:noinline
func zzRetainedInput(input string, f func(string) string) (string, weak.Pointer[byte], bool) {
	p := new([1 << 20]byte)
	copy(p[100:], input)
	parent := unsafe.String(&p[0], len(p))
	s := parent[100 : 100+len(input)]
	w := weak.Make(&p[0])
	out := f(s)
	shared := zzShares(out, parent)
	runtime.KeepAlive(p)
	return out, w, shared
}
func TestZZScansSpanRetention(t *testing.T) {
	cases := []struct {
		name, input string
		fn          func(string) string
	}{
		{"RunesInvalid", "ab\xff tail", func(s string) string { return TruncateRunes(s, 3) }},
		{"SliceInvalid", "x ab\xff tail", func(s string) string { return SafeSlice(s, 2, 5) }},
		{"SliceFullInvalid", "ab\xff tail", func(s string) string { return SafeSlice(s, 0, 100) }},
		{"TruncateInvalid", "ab\xff tail", func(s string) string { return Truncate(s, 3, "") }},
		{"TruncateInvalidSuffix", "ab\xff tail", func(s string) string { return Truncate(s, 4, ".") }},
		{"WordsInvalid", "ab\xff tail", func(s string) string { return TruncateWords(s, 6, "") }},
		{"WordsInvalidSuffix", "ab\xff tail", func(s string) string { return TruncateWords(s, 6, ".") }},
		{"TruncateInvalidSuffixOnly", "ab\xff tail", func(s string) string { return Truncate("0123456789", 3, s) }},
		{"WordsInvalidSuffixOnly", "ab\xff tail", func(s string) string { return TruncateWords("0123456789", 3, s) }},
		{"ReverseInvalidSmall", "ab\xff tail", Reverse},
		{"ReverseInvalidLong", strings.Repeat("ab", 30) + "\xff", Reverse},
		{"ReverseValidLong", strings.Repeat("ab", 30), Reverse},
		{"CapitalizeInvalid", "Ab\xff tail", Capitalize},
		{"PascalInvalid", "Ab\xff tail", PascalCase},
		{"PascalShrinks", strings.Repeat("_", 4096) + "a", PascalCase},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, w, shares := zzRetainedInput(tc.input, tc.fn)
			runtime.GC()
			runtime.GC()
			if shares || w.Value() != nil {
				t.Fatal("result retains parent")
			}
			runtime.KeepAlive(out)
		})
	}
	// No-op paths preserve the actual address, not only membership in a parent.
	parent := strings.Repeat("x", 4096)
	s := parent[27:40]
	for _, fn := range []func() string{
		func() string { return Truncate(s, 99, "\xff") }, func() string { return TruncateWords(s, 99, "\xff") },
		func() string { return TruncateRunes(s, 99) }, func() string { return CommonPrefix(s) },
		func() string { return PadLeft(s, 0, ' ') }, func() string { return PadRight(s, 0, ' ') },
		func() string { return PadCenter(s, 0, ' ') },
	} {
		if got := fn(); unsafe.StringData(got) != unsafe.StringData(s) {
			t.Fatal("no-op must retain original address")
		}
	}
	runtime.KeepAlive(parent)
}

func TestZZScansCompactResult(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		fn          func(string) string
	}{
		{"PascalShrinks", strings.Repeat("_", 1<<20) + "1", PascalCase},
		{"WordsDiscardInvalid", "a " + strings.Repeat("\xff", 1<<19) + "xyz", func(s string) string { return TruncateWords(s, len(s)-2, "") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Eight live one-byte outputs should not retain eight ~1 MiB intermediates.
			// The 2 MiB tolerance exceeds runtime bookkeeping by a wide margin.
			out := make([]string, 8)
			runtime.GC()
			runtime.GC()
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			for i := range out {
				out[i] = tc.fn(tc.input)
				if len(out[i]) != 1 {
					t.Fatal("unexpected fixture result")
				}
			}
			runtime.GC()
			runtime.GC()
			runtime.ReadMemStats(&after)
			retained := int64(after.HeapAlloc) - int64(before.HeapAlloc)
			t.Logf("retained delta for 8 compact results: %d B", retained)
			if retained > 2<<20 {
				t.Fatal("compact result retains oversized intermediate backing")
			}
			runtime.KeepAlive(out)
			runtime.KeepAlive(tc.input)
		})
	}
}

func TestZZScansLineThresholds(t *testing.T) {
	sizes := []int{0, 1, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129, 1024, 65536}
	for _, n := range sizes {
		text := strings.Repeat("x", n)
		for _, s := range []string{text, "\n" + text, text + "\n", text + "\n\n", "\xff" + text + "\r\n", strings.Repeat("\n", n)} {
			zzCompare(t, fmt.Sprintf("CountLines %d bytes", len(s)), func() any { return CountLines(s) }, func() any { return zzBaseCountLines(s) })
		}
	}
}
