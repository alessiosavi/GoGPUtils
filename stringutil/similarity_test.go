package stringutil

import (
	"math"
	"testing"
	"unicode/utf8"
)

func TestLevenshteinEmptyInputs(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"empty", "", 0},
		{"ASCII", "hello", 5},
		{"multibyte letter", "é", 1},
		{"multibyte letters", "世界", 2},
		{"emoji", "🙂", 1},
		{"mixed", "aé🙂", 3},
		{"combining sequence", "e\u0301", 2},
		{"invalid byte", "\xff", 1},
		{"invalid bytes", "\xff\xfe", 2},
		{"multibyte and invalid", "é\xff", 2},
		{"truncated encoding", "\xf0\x9f\x99", 3},
		{"overlong encoding", "\xc0\x80", 2},
		{"surrogate encoding", "\xed\xa0\x80", 3},
		{"replacement character", "\ufffd", 1},
		{"replacement and invalid", "\ufffd\xff", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantSimilarity := 0.0
			if tt.s == "" {
				wantSimilarity = 1.0
			}
			for _, pair := range [][2]string{{"", tt.s}, {tt.s, ""}} {
				if got := LevenshteinDistance(pair[0], pair[1]); got != tt.want {
					t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", pair[0], pair[1], got, tt.want)
				}
				if got := LevenshteinSimilarity(pair[0], pair[1]); got != wantSimilarity {
					t.Errorf("LevenshteinSimilarity(%q, %q) = %g, want %g", pair[0], pair[1], got, wantSimilarity)
				}
			}
		})
	}
}

func TestLevenshteinRuneSemantics(t *testing.T) {
	tests := []struct {
		name           string
		s1, s2         string
		wantDistance   int
		wantSimilarity float64
	}{
		{"ASCII substitution", "ab", "ac", 1, 0.5},
		{"Unicode substitution", "x", "é", 1, 0},
		{"emoji substitution", "a🙂", "a🙃", 1, 0.5},
		{"combining mark deletion", "e\u0301", "e", 1, 0.5},
		{"no normalization", "é", "e\u0301", 2, 0},
		{"invalid bytes decode equally", "\xff", "\xfe", 0, 1},
		{"invalid byte and replacement", "\xff", "\ufffd", 0, 1},
		{"invalid byte deletion", "é\xff", "é", 1, 0.5},
		{"truncated encoding", "\xe2\x82", "\ufffd\ufffd", 0, 1},
		{"equal Unicode", "é🙂", "é🙂", 0, 1},
		{"equal malformed", "é\xff", "é\xff", 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, pair := range [][2]string{{tt.s1, tt.s2}, {tt.s2, tt.s1}} {
				if got := LevenshteinDistance(pair[0], pair[1]); got != tt.wantDistance {
					t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", pair[0], pair[1], got, tt.wantDistance)
				}
				if got := LevenshteinSimilarity(pair[0], pair[1]); got != tt.wantSimilarity {
					t.Errorf("LevenshteinSimilarity(%q, %q) = %g, want %g", pair[0], pair[1], got, tt.wantSimilarity)
				}
			}
		})
	}
}

func TestLevenshteinRuneBounds(t *testing.T) {
	inputs := []string{"", "a", "ab", "é", "世界", "🙂", "aé🙂", "e\u0301", "\xff", "\xfe", "é\xff", "\xe2\x82", "\ufffd"}
	for _, s1 := range inputs {
		for _, s2 := range inputs {
			n1, n2 := utf8.RuneCountInString(s1), utf8.RuneCountInString(s2)
			lower, upper := max(n1, n2)-min(n1, n2), max(n1, n2)
			distance := LevenshteinDistance(s1, s2)
			if distance < lower || distance > upper {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, outside rune bounds [%d, %d]", s1, s2, distance, lower, upper)
			}
			if reverse := LevenshteinDistance(s2, s1); distance != reverse {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, reverse = %d", s1, s2, distance, reverse)
			}
			similarity := LevenshteinSimilarity(s1, s2)
			if math.IsNaN(similarity) || similarity < 0 || similarity > 1 {
				t.Errorf("LevenshteinSimilarity(%q, %q) = %g, outside [0, 1]", s1, s2, similarity)
			}
			if reverse := LevenshteinSimilarity(s2, s1); similarity != reverse {
				t.Errorf("LevenshteinSimilarity(%q, %q) = %g, reverse = %g", s1, s2, similarity, reverse)
			}
		}
	}
}
