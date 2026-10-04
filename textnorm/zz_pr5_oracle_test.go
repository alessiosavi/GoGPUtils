package textnorm

import (
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unsafe"
	"weak"
)

func pr5BaseRemoveFormat(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r != '\n' && r != '\t' && (unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Cc, r)) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String(), nil
}

// pr5Alias reports a nonempty result's byte offset in an input, or -1.
func pr5Alias(out, in string) int {
	if len(out) == 0 || len(in) == 0 {
		return -1
	}
	p, q := uintptr(unsafe.Pointer(unsafe.StringData(out))), uintptr(unsafe.Pointer(unsafe.StringData(in)))
	if p >= q && p-q < uintptr(len(in)) {
		return int(p - q)
	}
	return -1
}
func pr5String(t testing.TB, name, in, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s(%q): got %q, want %q", name, in, got, want)
	}
	if pr5Alias(got, in) != pr5Alias(want, in) {
		t.Fatalf("%s(%q): alias got %d, want %d", name, in, pr5Alias(got, in), pr5Alias(want, in))
	}
}
func pr5Inputs() []string {
	out := []string{"", "a b", "a  b", " a ", "\r\r\n", "\r\n\r", "\xff", "\ufffd", "0000000000\xf2Ă", "ꮳ", "Ꮳ"}
	ws := []rune{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000}
	for r := rune(0x2000); r <= 0x200a; r++ {
		ws = append(ws, r)
	}
	for _, r := range ws {
		for _, a := range []string{"", "x", "\xff", "界"} {
			for _, b := range []string{"", "y", "\xfe", "界"} {
				out = append(out, a+string(r)+b, a+string(r)+string(r)+b)
			}
		}
	}
	controls := []rune{0x200b, 0x200c, 0x200d, 0x200e, 0x200f, 0x202a, 0x202b, 0x202c, 0x202d, 0x202e, 0x2066, 0x2067, 0x2068, 0x2069, 0xfeff, 0xad, 0xe0001, 0xe0020, 0xe007f, 0xe000, 0x378, 0xfffd}
	for r := rune(0); r <= 0x9f; r++ {
		controls = append(controls, r)
	}
	for _, r := range controls {
		out = append(out, string(r)+"x", "x"+string(r)+"y", "x"+string(r))
	}
	for _, s := range []string{"\xff", "\xc0\xaf", "\xed\xa0\x80", "\xe2\x82", "\xf4\x90\x80\x80"} {
		out = append(out, s+"x", "x"+s+"y", "x"+s, " \t"+s+"\u0085y")
	}
	atoms := []string{"\r", "\n", "\r\n", "\r\r\n", "a"}
	var visit func(string, int)
	visit = func(s string, n int) {
		out = append(out, s)
		if n > 0 {
			for _, a := range atoms {
				visit(s+a, n-1)
			}
		}
	}
	visit("", 4)
	for n := range 41 {
		for _, s := range []string{"x", "é \r\n", "\xff\t", "a b", "\u200b"} {
			out = append(out, strings.Repeat(s, n))
		}
	}
	for _, s := range []string{"x ", "\r\n", "é\t", "\xff ", "x\x07"} {
		out = append(out, strings.Repeat(s, 32768))
	}
	return out
}

type pr5Parent [1 << 20]byte

//go:noinline
func pr5Live(fn func(string) string, s string) (string, weak.Pointer[pr5Parent]) {
	p := new(pr5Parent)
	copy(p[97:], s)
	view := unsafe.String(&p[97], len(s))
	w := weak.Make(p)
	out := fn(view)
	runtime.KeepAlive(p)
	return out, w
}
func pr5Retention(t *testing.T, name string, fn func(string) string, s string, retained bool) {
	t.Helper()
	out, w := pr5Live(fn, s)
	runtime.GC()
	runtime.GC()
	got := w.Value() != nil
	runtime.KeepAlive(out)
	if got != retained {
		t.Fatalf("%s retained=%v want %v", name, got, retained)
	}
}

var pr5RemoveFormat = New().RemoveFormatChars()

func pr5Check(t testing.TB, s string) {
	for _, tc := range []struct {
		name      string
		got, base func(string) (string, error)
	}{{"RemoveFormatChars", pr5RemoveFormat.Run, pr5BaseRemoveFormat}} {
		a, e := tc.got(s)
		b, w := tc.base(s)
		if e != w { //nolint:errorlint // The oracle compares exact error identity, not wrapping.
			t.Fatalf("%s: error %v != %v", tc.name, e, w)
		}
		pr5String(t, tc.name, s, a, b)
	}
	// Compare the affected preset twice without asserting idempotence, which
	// is not guaranteed by the BASE stages.
	for _, tc := range []struct{ got, base Pipeline }{
		{HygienePreset(), New().SanitizeUTF8().DecodeHTMLEntities().Then(pr5BaseRemoveFormat).TrimSpace().CollapseWhitespace()},
	} {
		a, e := tc.got.Run(s)
		b, w := tc.base.Run(s)
		if a != b || e != w { //nolint:errorlint // The oracle compares exact error identity, not wrapping.
			t.Fatalf("pipeline %q: %q/%v != %q/%v", s, a, e, b, w)
		}
		a, e = tc.got.Run(a)
		b, w = tc.base.Run(b)
		if a != b || e != w { //nolint:errorlint // The oracle compares exact error identity, not wrapping.
			t.Fatal("second pipeline run differs")
		}
	}
}
func TestZZPR5Differential(t *testing.T) {
	for _, s := range pr5Inputs() {
		pr5Check(t, s)
	}
	for i := range 256 {
		pr5Check(t, string([]byte{byte(i)}))
		for j := range 256 {
			s := string([]byte{byte(i), byte(j)})
			a, e := pr5RemoveFormat.Run(s)
			b, w := pr5BaseRemoveFormat(s)
			if e != w { //nolint:errorlint // The oracle compares exact error identity, not wrapping.
				t.Fatal("errors")
			}
			pr5String(t, "Filter", s, a, b)
		}
	}
}
func TestZZPR5Retention(t *testing.T) {
	pr5Retention(t, "view control", func(s string) string { return s }, "clean", true)
	filter := func(s string) string { r, _ := pr5RemoveFormat.Run(s); return r }
	for _, s := range []string{"clean", "clean\u200b", "clean\xff", "\u200bclean", "\u200b", "\x00\u200b"} {
		pr5Retention(t, "filter", filter, s, false)
	}
}
func FuzzZZFormatCharsDifferential(f *testing.F) {
	for _, s := range []string{"", "a b", "\xff\u0085x", "a\u200db", "ꮳ", "0000000000\xf2Ă"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) { pr5Check(t, s) })
}
