package textnorm

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSearchPreset(t *testing.T) {
	got, err := SearchPreset().Run("  Café,   go!  gophers ")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "cafe go gophers" {
		t.Fatalf("Run() = %q, want %q", got, "cafe go gophers")
	}
}

func TestCanonicalPreset(t *testing.T) {
	got, err := CanonicalPreset().Run("  Hello,   World!  ")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "hello, world!" {
		t.Fatalf("Run() = %q, want %q", got, "hello, world!")
	}
}

func TestDBSafePreset(t *testing.T) {
	input := string([]byte{'g', 'o', 0x00, 0xff, '!'})
	got, err := DBSafePreset().Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !utf8.ValidString(got) {
		t.Fatalf("Run() = %q, want valid UTF-8", got)
	}
	if strings.ContainsRune(got, '\x00') {
		t.Fatalf("Run() = %q, want no NUL bytes", got)
	}
}

func TestMeaningPresetGoldens(t *testing.T) {
	run := func(in string) string {
		out, err := MeaningPreset().Run(in)
		if err != nil {
			t.Fatalf("MeaningPreset(%q): %v", in, err)
		}
		return out
	}
	cases := []struct{ name, in, want string }{
		{"plus preserved", "Galaxy S22+", "galaxy s22+"},
		{"decimal preserved", "size 4.5", "size 4.5"},
		{"thousands preserved", "1,000", "1,000"},
		{"percent preserved", "100% Cotton", "100% cotton"},
		{"no token dedup", "Nike Air Air Max", "nike air air max"},
		{"latin accents folded", "Hermès Café", "hermes cafe"},
		{"indic preserved", "किताब", "किताब"},
		{"entity decoded then spaced", "Dolce &amp; Gabbana", "dolce gabbana"},
		{"punct to space collapses", "iPhone 13 | 128GB | Black", "iphone 13 128gb black"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := run(c.in); got != c.want {
				t.Fatalf("MeaningPreset(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
	if run("Galaxy S22+") == run("Galaxy S22") {
		t.Fatal("s22+ must differ from s22")
	}
	if run("size 4.5") == run("size 45") {
		t.Fatal("4.5 must differ from 45")
	}
	if run("1,000") == run("1 000") {
		t.Fatal("1,000 must differ from 1 000")
	}
}

func TestHygienePresetGoldens(t *testing.T) {
	run := func(in string) string {
		out, err := HygienePreset().Run(in)
		if err != nil {
			t.Fatalf("HygienePreset(%q): %v", in, err)
		}
		return out
	}
	if got := run("Café \u200bCrème!"); got != "Café Crème!" {
		t.Fatalf("hygiene must keep case/punct/diacritics, drop ZWSP: %q", got)
	}
	if got := run("A &amp; B"); got != "A & B" {
		t.Fatalf("hygiene must decode entities: %q", got)
	}
	if got := run("  spaced\n\nout\t "); got != "spaced out" {
		t.Fatalf("hygiene must collapse whitespace: %q", got)
	}
}

func TestWidthFoldOption(t *testing.T) {
	defaultOut, err := SearchPreset().Run("Ｇｏ")
	if err != nil {
		t.Fatalf("default Run() error = %v", err)
	}

	foldedOut, err := SearchPreset(WithWidthFold()).Run("Ｇｏ")
	if err != nil {
		t.Fatalf("folded Run() error = %v", err)
	}

	if foldedOut != "go" {
		t.Fatalf("folded Run() = %q, want %q", foldedOut, "go")
	}
	if defaultOut == foldedOut {
		t.Fatalf("default and folded outputs match: %q", defaultOut)
	}
}
