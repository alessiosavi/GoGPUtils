package textnorm

import "testing"

func TestNormalizeUnicodeRemovesAccents(t *testing.T) {
	got, err := New().NormalizeUnicode().Run("caf\u00e9")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "cafe" {
		t.Fatalf("Run() = %q, want %q", got, "cafe")
	}
}

func TestRemoveAccentsMatchesNormalizeUnicode(t *testing.T) {
	got, err := New().RemoveAccents().Run("na\u00efve")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "naive" {
		t.Fatalf("Run() = %q, want %q", got, "naive")
	}
}

func TestNormalizeUnicodeLatin(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"latin accents stripped", "Caf\u00e9 Cr\u00e8me Herm\u00e8s", "Cafe Creme Hermes"},
		{"devanagari matras preserved", "\u0915\u093f\u0924\u093e\u092c", "\u0915\u093f\u0924\u093e\u092c"},
		{"arabic harakat preserved", "\u0645\u064e\u0643\u062a\u064e\u0628", "\u0645\u064e\u0643\u062a\u064e\u0628"},
		{"mixed latin+indic", "caf\u00e9 \u0915\u093f\u0924\u093e\u092c", "cafe \u0915\u093f\u0924\u093e\u092c"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := New().NormalizeUnicodeLatin().Run(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("NormalizeUnicodeLatin(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
