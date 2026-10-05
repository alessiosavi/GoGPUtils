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

func TestNormalizationRawMalformedUTF8(t *testing.T) {
	const input = "0000000000\xf2Ă"
	for _, tt := range []struct {
		name       string
		pipe       Pipeline
		wantFirst  string
		wantSecond string
	}{
		{"NormalizeUnicode", New().NormalizeUnicode(), "0000000000\ufffdĂ", "0000000000\ufffdA"},
		{"RemoveAccents", New().RemoveAccents(), "0000000000\ufffdĂ", "0000000000\ufffdA"},
		{"NormalizeUnicodeLatin", New().NormalizeUnicodeLatin(), "0000000000\ufffdĂ", "0000000000\ufffdA"},
		{"Pipeline", New().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace(), "0000000000\ufffdă", "0000000000\ufffda"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Raw normalization succeeds but leaves the accent on the first pass.
			first, err := tt.pipe.Run(input)
			if err != nil || first != tt.wantFirst {
				t.Fatalf("first Run(%q) = %q, %v; want %q, nil", input, first, err, tt.wantFirst)
			}
			second, err := tt.pipe.Run(first)
			if err != nil || second != tt.wantSecond {
				t.Fatalf("second Run(%q) = %q, %v; want %q, nil", first, second, err, tt.wantSecond)
			}
		})
	}
}

func TestSanitizeBeforeNormalization(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  string
	}{
		{"saved fuzz input", "0000000000\xf2Ă", "0000000000\ufffdA"},
		{"invalid lead", "\xffĂ", "\ufffdA"},
		{"stray continuation", "\x80Ă", "\ufffdA"},
		{"truncated sequence", "\xe2\x82Ă", "\ufffd\ufffdA"},
		{"overlong encoding", "\xc0\xafĂ", "\ufffd\ufffdA"},
		{"surrogate encoding", "\xed\xa0\x80Ă", "\ufffd\ufffd\ufffdA"},
		{"out of range", "\xf4\x90\x80\x80Ă", "\ufffd\ufffd\ufffd\ufffdA"},
		{"NUL removed", "\x00Ă", "A"},
		{"replacement character preserved", "\ufffdĂ", "\ufffdA"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, stage := range []struct {
				name string
				pipe Pipeline
			}{
				{"NormalizeUnicode", New().SanitizeUTF8().NormalizeUnicode()},
				{"RemoveAccents", New().SanitizeUTF8().RemoveAccents()},
				{"NormalizeUnicodeLatin", New().SanitizeUTF8().NormalizeUnicodeLatin()},
			} {
				t.Run(stage.name, func(t *testing.T) {
					got, err := stage.pipe.Run(tt.input)
					if err != nil || got != tt.want {
						t.Fatalf("Run(%q) = %q, %v; want %q, nil", tt.input, got, err, tt.want)
					}
				})
			}
		})
	}
}
