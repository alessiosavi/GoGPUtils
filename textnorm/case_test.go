package textnorm

import (
	"testing"
	"unicode"

	"golang.org/x/text/runes"
)

func TestFoldCaseDiffersFromLower(t *testing.T) {
	folded, err := New().FoldCase().Run("Straße")
	if err != nil {
		t.Fatalf("FoldCase() error = %v", err)
	}
	if folded != "strasse" {
		t.Fatalf("FoldCase() = %q, want %q", folded, "strasse")
	}

	lowered, err := New().Lower().Run("Straße")
	if err != nil {
		t.Fatalf("Lower() error = %v", err)
	}
	if lowered != "straße" {
		t.Fatalf("Lower() = %q, want %q", lowered, "straße")
	}
}

func TestRuneTransformsMapAndFilter(t *testing.T) {
	got, err := New().MapRunes(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '-'
	}).FilterRunes(runes.Predicate(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-'
	})).Run("Go! 2")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "Go--2" {
		t.Fatalf("Run() = %q, want %q", got, "Go--2")
	}
}
