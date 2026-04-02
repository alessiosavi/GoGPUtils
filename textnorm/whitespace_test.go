package textnorm

import "testing"

func TestWhitespaceStagesCollapseAndTrim(t *testing.T) {
	got, err := New().TrimSpace().CollapseWhitespace().Run("\t  hello   world \n")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "hello world" {
		t.Fatalf("Run() = %q, want %q", got, "hello world")
	}
}

func TestSanitizeUTF8RemovesNulAndReplacesInvalidBytes(t *testing.T) {
	input := string([]byte{'g', 'o', 0xff, 0x00, '!', 0xfe})

	got, err := New().SanitizeUTF8().Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "go\uFFFD!\uFFFD" {
		t.Fatalf("Run() = %q, want %q", got, "go\uFFFD!\uFFFD")
	}
}
