package textnorm

import (
	"testing"
	"unicode/utf8"
)

func FuzzPipeline(f *testing.F) {
	for _, seed := range []string{
		"",
		"  Café   Straße  ",
		"🙂🙂",
		string([]byte{'g', 'o', 0xff, 0x00, '!', 0xfe}),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pipe := New().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace()
		out1, err := pipe.Run(input)
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		out2, err := pipe.Run(out1)
		if err != nil {
			t.Fatalf("second Run() error = %v", err)
		}
		if out1 != out2 {
			t.Fatalf("pipeline not idempotent: %q != %q", out1, out2)
		}
	})
}

func FuzzSearchPreset(f *testing.F) {
	for _, seed := range []string{
		"  Café,   go!  ",
		"Straße",
		"🙂 mixed CASE 🙂",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pipe := SearchPreset()
		out1, err := pipe.Run(input)
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		out2, err := pipe.Run(out1)
		if err != nil {
			t.Fatalf("second Run() error = %v", err)
		}
		if out1 != out2 {
			t.Fatalf("SearchPreset not idempotent: %q != %q", out1, out2)
		}
	})
}

func FuzzCanonicalPreset(f *testing.F) {
	for _, seed := range []string{
		"  Hello,   World!  ",
		"Café",
		"mixed   whitespace",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pipe := CanonicalPreset()
		out1, err := pipe.Run(input)
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		out2, err := pipe.Run(out1)
		if err != nil {
			t.Fatalf("second Run() error = %v", err)
		}
		if out1 != out2 {
			t.Fatalf("CanonicalPreset not idempotent: %q != %q", out1, out2)
		}
	})
}

func FuzzDBSafePreset(f *testing.F) {
	for _, seed := range []string{
		string([]byte{'g', 'o', 0x00, 0xff, '!', 0xfe}),
		"valid text",
		"🙂 null \x00 mix",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pipe := DBSafePreset()
		out1, err := pipe.Run(input)
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		if !utf8.ValidString(out1) {
			t.Fatalf("DBSafePreset produced invalid UTF-8: %q", out1)
		}
		out2, err := pipe.Run(out1)
		if err != nil {
			t.Fatalf("second Run() error = %v", err)
		}
		if out1 != out2 {
			t.Fatalf("DBSafePreset not idempotent: %q != %q", out1, out2)
		}
	})
}
