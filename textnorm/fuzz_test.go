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
		"Ꮳꮳ Ᏸᏸ",
		"00000000000000000000000000000ꮒ0",
		"🙂🙂",
		string([]byte{'g', 'o', 0xff, 0x00, '!', 0xfe}),
		"0000000000\xf2Ă", // A malformed byte can block decomposition without sanitization.
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		pipe := New().SanitizeUTF8().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace()
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

var searchPresetSeeds = []string{
	"  Café,   go!  ",
	"Straße",
	"Ꮳꮳ Ᏸᏸ",
	"🙂 mixed CASE 🙂",
	"",
	"00000000000000000000000000000ꮒ0",
	"0000000000\xf2Ă",
	"go\xff\x00!\xfe",
	"Ｇｏ ＣＡＦÉ",
	// Filtering can expose Hangul L+V and LV+T composition boundaries.
	"ᄀ!ᅡ",
	"ᄀ-ᅡ-ᆨ",
	"가!ᆨ",
	"ᄒ,ᅵ,ᇂ",
	"ᄀ🙂ᅡ",
	"ᄀ+ᅡ🙂ᆨ",
	"가🙂ᆨ",
	"ᄀ\x00ᅡ",
	"ᄀ\x00ᅡ\x00ᆨ",
	"가\x00ᆨ",
	"ᄀ\xffᅡ",
	"ᄀ\xfeᅡ\xffᆨ",
	"가\xffᆨ",
	"ᄀ\u0903ᅡ\u200dᆨ",
	"  ＧＯ ᄀ!ᅡ 가🙂ᆨ ＣＡＦÉ  ",
}

func FuzzSearchPreset(f *testing.F) {
	for _, seed := range searchPresetSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		for _, widthFold := range []bool{false, true} {
			var opts []PresetOption
			if widthFold {
				opts = append(opts, WithWidthFold())
			}
			pipe := SearchPreset(opts...)
			out1, err := pipe.Run(input)
			if err != nil {
				t.Fatalf("Run() error = %v (widthFold=%t)", err, widthFold)
			}
			out2, err := pipe.Run(out1)
			if err != nil {
				t.Fatalf("second Run() error = %v (widthFold=%t)", err, widthFold)
			}
			if out1 != out2 {
				t.Fatalf("SearchPreset(%q) not idempotent (widthFold=%t): %q != %q", input, widthFold, out1, out2)
			}
		}
	})
}

func FuzzSearchPresetDifferential(f *testing.F) {
	for _, seed := range searchPresetSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		checkSearchPresetAgainstBase(t, input)
	})
}

func FuzzCanonicalPreset(f *testing.F) {
	addWidthFoldPresetSeeds(f)
	for _, seed := range []string{
		"  Hello,   World!  ",
		"Café",
		"Ꮳꮳ Ᏸᏸ",
		"mixed   whitespace",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		for _, widthFold := range []bool{false, true} {
			var opts []PresetOption
			if widthFold {
				opts = append(opts, WithWidthFold())
			}
			pipe := CanonicalPreset(opts...)
			out1, err := pipe.Run(input)
			if err != nil {
				t.Fatalf("Run() error = %v (widthFold=%t)", err, widthFold)
			}
			out2, err := pipe.Run(out1)
			if err != nil {
				t.Fatalf("second Run() error = %v (widthFold=%t)", err, widthFold)
			}
			if out1 != out2 {
				t.Fatalf("CanonicalPreset(%q) not idempotent (widthFold=%t): %q != %q", input, widthFold, out1, out2)
			}
		}
	})
}

func FuzzMeaningPreset(f *testing.F) {
	addWidthFoldPresetSeeds(f)
	f.Add("Galaxy S22+ 4.5\" 1,000 100%")
	f.Add("c++ a+b % off 1.000")
	f.Add("Café किताब مَكتَب")
	f.Add("Ꮳꮳ Ᏸᏸ")
	f.Fuzz(func(t *testing.T, in string) {
		// Meaning still has unrelated case-folding and punctuation idempotence
		// failures; compare its exact output with the width-normalization model.
		for _, widthFold := range []bool{false, true} {
			var opts []PresetOption
			if widthFold {
				opts = append(opts, WithWidthFold())
			}
			out, err := MeaningPreset(opts...).Run(in)
			if err != nil {
				t.Fatalf("MeaningPreset(%q, widthFold=%t): %v", in, widthFold, err)
			}
			model := widthFoldBaseMeaningPreset()
			if widthFold {
				model = widthFoldPresetModel("Meaning", true, true)
			}
			want, wantErr := model.Run(in)
			if out != want || err != wantErr { //nolint:errorlint // Compare oracle error identity.
				t.Fatalf("MeaningPreset(%+q, widthFold=%t) = %+q, %v; model = %+q, %v", in, widthFold, out, err, want, wantErr)
			}
			for _, r := range out {
				ok := unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsMark(r) || r == ' ' ||
					r == '.' || r == ',' || r == '+' || r == '%'
				if !ok {
					t.Fatalf("illegal rune %q in output %q for input %q (widthFold=%t)", r, out, in, widthFold)
				}
			}
			if strings.Contains(out, "  ") {
				t.Fatalf("uncollapsed whitespace in %q (widthFold=%t)", out, widthFold)
			}
		}
	})
}

func FuzzHygienePreset(f *testing.F) {
	f.Add("Café \u200bCrème &amp; Co!")
	f.Add("a\n\tb\x00c &#8203; &amp;amp;")
	f.Fuzz(func(t *testing.T, in string) {
		// No idempotency check: HTML-entity decoding is legitimately
		// non-idempotent ("&amp;amp;" → "&amp;" → "&").
		out, err := HygienePreset().Run(in)
		if err != nil {
			t.Fatalf("HygienePreset(%q): %v", in, err)
		}
		if !utf8.ValidString(out) {
			t.Fatalf("invalid UTF-8 in output %q for input %q", out, in)
		}
		for _, r := range out {
			if unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Cc, r) {
				t.Fatalf("format/control rune %U in output %q for input %q", r, out, in)
			}
			if unicode.IsSpace(r) && r != ' ' {
				t.Fatalf("non-collapsed whitespace %U in output %q for input %q", r, out, in)
			}
		}
		if strings.Contains(out, "  ") || out != strings.TrimSpace(out) {
			t.Fatalf("whitespace not collapsed/trimmed in %q", out)
		}
	})
}

func FuzzDBSafePreset(f *testing.F) {
	addWidthFoldPresetSeeds(f)
	for _, seed := range []string{
		string([]byte{'g', 'o', 0x00, 0xff, '!', 0xfe}),
		"valid text",
		"🙂 null \x00 mix",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		for _, widthFold := range []bool{false, true} {
			var opts []PresetOption
			if widthFold {
				opts = append(opts, WithWidthFold())
			}
			pipe := DBSafePreset(opts...)
			out1, err := pipe.Run(input)
			if err != nil {
				t.Fatalf("Run() error = %v (widthFold=%t)", err, widthFold)
			}
			if !utf8.ValidString(out1) {
				t.Fatalf("DBSafePreset produced invalid UTF-8: %q (widthFold=%t)", out1, widthFold)
			}
			out2, err := pipe.Run(out1)
			if err != nil {
				t.Fatalf("second Run() error = %v (widthFold=%t)", err, widthFold)
			}
			if out1 != out2 {
				t.Fatalf("DBSafePreset(%q) not idempotent (widthFold=%t): %q != %q", input, widthFold, out1, out2)
			}
		}
	})
}
