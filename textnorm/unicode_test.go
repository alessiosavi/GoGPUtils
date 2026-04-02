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
