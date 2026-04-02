package textnorm

import "testing"

func TestFoldWidthIsOptIn(t *testing.T) {
	defaultOut, err := New().Run("Ｇｏ")
	if err != nil {
		t.Fatalf("default Run() error = %v", err)
	}
	if defaultOut != "Ｇｏ" {
		t.Fatalf("default Run() = %q, want %q", defaultOut, "Ｇｏ")
	}

	folded, err := New().FoldWidth().Run("Ｇｏ")
	if err != nil {
		t.Fatalf("FoldWidth() error = %v", err)
	}
	if folded != "Go" {
		t.Fatalf("FoldWidth() = %q, want %q", folded, "Go")
	}
}
