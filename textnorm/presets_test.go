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
