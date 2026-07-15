package textnorm

import (
	"strings"
	"testing"
	"unicode"
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

func FuzzMeaningPreset(f *testing.F) {
	f.Add("Galaxy S22+ 4.5\" 1,000 100%")
	f.Add("c++ a+b % off 1.000")
	f.Add("Café किताब مَكتَب")
	f.Fuzz(func(t *testing.T, in string) {
		out, err := MeaningPreset().Run(in)
		if err != nil {
			t.Fatalf("MeaningPreset(%q): %v", in, err)
		}
		for _, r := range out {
			ok := unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsMark(r) || r == ' ' ||
				r == '.' || r == ',' || r == '+' || r == '%'
			if !ok {
				t.Fatalf("illegal rune %q in output %q for input %q", r, out, in)
			}
		}
		if strings.Contains(out, "  ") {
			t.Fatalf("uncollapsed whitespace in %q", out)
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
