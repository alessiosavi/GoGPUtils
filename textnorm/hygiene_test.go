package textnorm

import "testing"

func TestDecodeHTMLEntities(t *testing.T) {
	got, err := New().DecodeHTMLEntities().Run("A &amp; B &#39;quoted&#39;")
	if err != nil || got != "A & B 'quoted'" {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestRemoveFormatChars(t *testing.T) {
	in := "zero\u200bwidth bidi\u202e mark" // ZWSP + RLO
	want := "zerowidth bidi mark"           // Cf runes dropped, letters untouched
	got, err := New().RemoveFormatChars().Run(in)
	if err != nil || got != want {
		t.Fatalf("got %q err %v", got, err)
	}
	// \n and \t survive (whitespace collapse is a separate stage's job)
	got, err = New().RemoveFormatChars().Run("a\n\tb\x00c")
	if err != nil || got != "a\n\tbc" {
		t.Fatalf("control handling: got %q err %v", got, err)
	}
}
