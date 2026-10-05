package textnorm

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

var widthFoldPresets = []struct {
	name       string
	head, base func(...PresetOption) Pipeline
}{
	{"Canonical", CanonicalPreset, widthFoldBaseCanonicalPreset},
	{"DBSafe", DBSafePreset, widthFoldBaseDBSafePreset},
	{"Meaning", MeaningPreset, widthFoldBaseMeaningPreset},
	{"Search", SearchPreset, widthFoldBaseSearchPreset},
	{"Hygiene", HygienePreset, widthFoldBaseHygienePreset},
}

// Old and new bytes are ordered Canonical, DBSafe, Meaning.
var widthFoldGoldens = []struct {
	name, input string
	old, want   [3]string
}{
	{"voiced", "ｶﾞ", [3]string{"カ\u3099", "カ\u3099", "カ\u3099"}, [3]string{"カ", "カ", "ガ"}},
	{"semivoiced", "ﾊﾟ", [3]string{"ハ\u309a", "ハ\u309a", "ハ\u309a"}, [3]string{"ハ", "ハ", "パ"}},
	{"vu", "ｳﾞ", [3]string{"ウ\u3099", "ウ\u3099", "ウ\u3099"}, [3]string{"ウ", "ウ", "ヴ"}},
	{"va", "ﾜﾞ", [3]string{"ワ\u3099", "ワ\u3099", "ワ\u3099"}, [3]string{"ワ", "ワ", "ヷ"}},
	{"vo", "ｦﾞ", [3]string{"ヲ\u3099", "ヲ\u3099", "ヲ\u3099"}, [3]string{"ヲ", "ヲ", "ヺ"}},
	{"uncomposable", "ｱﾞ", [3]string{"ア\u3099", "ア\u3099", "ア\u3099"}, [3]string{"ア", "ア", "ア\u3099"}},
	{"standalone_voiced", "ﾞ", [3]string{"\u3099", "\u3099", "\u3099"}, [3]string{"", "", "\u3099"}},
	{"standalone_semivoiced", "ﾟ", [3]string{"\u309a", "\u309a", "\u309a"}, [3]string{"", "", "\u309a"}},
	{"multiple", "ｶﾞﾊﾟｳﾞ", [3]string{"カ\u3099ハ\u309aウ\u3099", "カ\u3099ハ\u309aウ\u3099", "カ\u3099ハ\u309aウ\u3099"}, [3]string{"カハウ", "カハウ", "ガパヴ"}},
	{"mixed_width", " Ｇｏ ｶﾞ １２３！ ", [3]string{"go カ\u3099 123!", "Go カ\u3099 123!", "go カ\u3099 123"}, [3]string{"go カ 123!", "Go カ 123!", "go ガ 123"}},
	{"latin", "Aﾞ", [3]string{"a\u3099", "A\u3099", "a\u3099"}, [3]string{"a", "A", "a"}},
	{"native_kana", "カﾞ", [3]string{"カ\u3099", "カ\u3099", "カ\u3099"}, [3]string{"カ", "カ", "ガ"}},
	{"native_mark", "ｶ\u3099", [3]string{"カ", "カ", "カ\u3099"}, [3]string{"カ", "カ", "ガ"}},
	{"mark_before", "ｶ\u0301ﾞ", [3]string{"カ\u3099", "カ\u3099", "カ\u0301\u3099"}, [3]string{"カ", "カ", "ガ\u0301"}},
	{"mark_after", "ｶﾞ\u0323", [3]string{"カ\u3099", "カ\u3099", "カ\u3099\u0323"}, [3]string{"カ", "カ", "ガ\u0323"}},
	{"reordering", "\u0300ﾞ", [3]string{"\u3099", "\u3099", "\u0300\u3099"}, [3]string{"", "", "\u3099\u0300"}},
	{"hangul", "ᄀﾞᅡ", [3]string{"ᄀ\u3099ᅡ", "ᄀ\u3099ᅡ", "ᄀ\u3099ᅡ"}, [3]string{"가", "가", "ᄀ\u3099ᅡ"}},
	{"normalize_before_and_after", "éﾞ\u1bf3", [3]string{"e\u3099\u1bf3", "e\u3099\u1bf3", "e\u3099\u1bf3"}, [3]string{"e\u1bf3", "e\u1bf3", "e\u1bf3"}},
	{"stable_key_changes", "＜\u0338", [3]string{"<", "<", "\u0338"}, [3]string{"<", "<", ""}},
	{"html_width", "&#xff76;&#xff9e;", [3]string{"&#xff76;&#xff9e;", "&#xff76;&#xff9e;", "カ\u3099"}, [3]string{"&#xff76;&#xff9e;", "&#xff76;&#xff9e;", "ガ"}},
}

func TestWidthFoldPresetGoldens(t *testing.T) {
	for i, preset := range widthFoldPresets[:3] {
		pipe, base := preset.head(WithWidthFold()), preset.base(WithWidthFold())
		for _, tt := range widthFoldGoldens {
			t.Run(preset.name+"/"+tt.name, func(t *testing.T) {
				if old, err := base.Run(tt.input); err != nil || old != tt.old[i] {
					t.Fatalf("frozen master(%+q) = %+q, %v; want %+q", tt.input, old, err, tt.old[i])
				}
				got, err := pipe.Run(tt.input)
				if err != nil || got != tt.want[i] {
					t.Errorf("Run(%+q) = %+q, %v; want %+q", tt.input, got, err, tt.want[i])
				}
				if again, err := pipe.Run(got); err != nil || again != got {
					t.Errorf("not idempotent: %+q -> %+q -> %+q, %v", tt.input, got, again, err)
				}
			})
		}
	}
}

func TestWidthFoldPresetBoundaries(t *testing.T) {
	for _, preset := range widthFoldPresets[:2] {
		pipe := preset.head(WithWidthFold())
		for _, n := range []int{29, 30, 31, 32} {
			for _, tt := range []struct{ name, input, want string }{
				{"marks", strings.Repeat("ﾞ", n), ""},
				{"latin", "a" + strings.Repeat("ﾞ", n), "a"},
				{"kana", "ｶ" + strings.Repeat("\u0301", n) + "ﾞ", "カ"},
				{"repeated_kana", strings.Repeat("ｶﾞ", n), strings.Repeat("カ", n)},
			} {
				t.Run(fmt.Sprintf("%s/%s/%d", preset.name, tt.name, n), func(t *testing.T) {
					got, err := pipe.Run(tt.input)
					if err != nil || got != tt.want {
						t.Errorf("Run(%+q) = %+q, %v; want %+q", tt.input, got, err, tt.want)
					}
					if again, err := pipe.Run(got); err != nil || again != got {
						t.Errorf("not idempotent: %+q -> %+q -> %+q, %v", tt.input, got, again, err)
					}
				})
			}
		}
	}
}

// This public-stage model preserves the frozen prefix/suffix, including HTML
// decoding before Meaning's normalization. The only correction is after Width.
func widthFoldPresetModel(name string, widthFold, renormalize bool) Pipeline {
	pipe := New().SanitizeUTF8()
	if name == "Meaning" {
		pipe = pipe.DecodeHTMLEntities().NormalizeUnicodeLatin()
	} else {
		pipe = pipe.NormalizeUnicode()
	}
	if widthFold {
		pipe = pipe.FoldWidth()
		if renormalize {
			if name == "Meaning" {
				pipe = pipe.NormalizeUnicodeLatin()
			} else {
				pipe = pipe.NormalizeUnicode()
			}
		}
	}
	if name != "DBSafe" {
		pipe = pipe.FoldCase()
	}
	if name == "Meaning" {
		pipe = pipe.PreserveMeaningPunct()
	}
	return pipe.TrimSpace().CollapseWhitespace()
}

var widthFoldPresetSeeds = []string{
	"", "Ｇｏ １２３！", "\xffｶﾞ\x00\xfe", "\ufffdｶﾞ",
	"0000000000\xf2Ă", "00000000000000000000000000000ꮒ0",
	"Ꮳꮳ Ᏸᏸ", "Café किताब مَكتَب", "ｶ\u034fﾞ", "ｶ\u200dﾞ",
	"++", "＋＋", "ŉﾞ\u0300\u1715", "&amp;amp;", "&amp;#xff76;&amp;#xff9e;",
}

func addWidthFoldPresetSeeds(f *testing.F) {
	for _, tt := range widthFoldGoldens {
		f.Add(tt.input)
	}
	for _, seed := range widthFoldPresetSeeds {
		f.Add(seed)
	}
	for _, n := range []int{29, 30, 31, 32} {
		f.Add("a" + strings.Repeat("ﾞ", n))
		f.Add("ｶ" + strings.Repeat("\u0301", n) + "ﾞ")
	}
}

func widthFoldCorpus(visit func(string)) {
	for _, tt := range widthFoldGoldens {
		visit(tt.input)
	}
	for _, s := range widthFoldPresetSeeds {
		visit(s)
	}
	widths := []rune{'\u3000'}
	for r := rune(0xff00); r <= 0xffef; r++ {
		widths = append(widths, r)
	}
	// Exhaustive width pairs and mark combinations run without race instrumentation
	// to keep this single-goroutine deterministic comparison inexpensive under -race.
	for _, r := range widths {
		visit(string(r))
		if !raceEnabled {
			for _, q := range widths {
				visit(string([]rune{r, q}))
			}
		}
	}
	if !raceEnabled {
		for r := rune(0); r <= utf8.MaxRune; r++ {
			if unicode.IsMark(r) {
				for _, mark := range []rune{'ﾞ', 'ﾟ'} {
					visit(string([]rune{r, mark}))
					visit(string([]rune{mark, r}))
					for _, base := range []rune{'A', 'é', 'ŉ', 'Ａ', 'カ', 'ｶ'} {
						visit(string([]rune{base, r, mark}))
						visit(string([]rune{base, mark, r}))
					}
				}
			}
		}
	}
	for _, n := range []int{29, 30, 31, 32, 63, 128} {
		for _, unit := range []string{"ｶﾞ", "Aﾞ", "\u0301", "ﾞ", "ﾟ", "ᄀﾞᅡ"} {
			visit(strings.Repeat(unit, n))
			visit("ｶ" + strings.Repeat(unit, n) + "ﾞ")
		}
	}
	rng := rand.New(rand.NewPCG(0xf2, 0xf94af78))
	alphabet := []rune("AéŉＡＧｶﾞﾊﾟカガ\u0300\u0301\u0323\u0338\u034f\u3099\u309a\u1bf3\u1715각ꮳᏣمَक़&;#+! １２\x00\ufffd")
	for class := range 3 {
		for range 1000 {
			var b strings.Builder
			for range rng.IntN(80) {
				switch class {
				case 0:
					b.WriteRune(alphabet[rng.IntN(len(alphabet))])
				case 1:
					b.WriteRune(rune(rng.IntN(utf8.MaxRune + 1)))
				case 2:
					b.WriteByte(byte(rng.IntN(256)))
				}
			}
			visit(b.String())
		}
	}
}

func TestWidthFoldPresetDifferential(t *testing.T) {
	var inputs []string
	widthFoldCorpus(func(input string) { inputs = append(inputs, input) })
	corpus := "full"
	if raceEnabled {
		corpus = "reduced race"
	}
	t.Logf("%s corpus: %d deterministic inputs", corpus, len(inputs))
	for i, preset := range widthFoldPresets {
		for _, widthFold := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/width=%t", preset.name, widthFold), func(t *testing.T) {
				var opts []PresetOption
				if widthFold {
					opts = append(opts, WithWidthFold())
				}
				head, base := preset.head(opts...), preset.base(opts...)
				model, uncorrected := base, base
				if i < 3 {
					model = widthFoldPresetModel(preset.name, widthFold, true)
					uncorrected = widthFoldPresetModel(preset.name, widthFold, false)
				}
				for _, input := range inputs {
					old, oldErr := base.Run(input)
					want, wantErr := old, oldErr
					if i < 3 {
						original, originalErr := uncorrected.Run(input)
						if original != old || originalErr != oldErr { //nolint:errorlint // Compare BASE error identity.
							t.Fatalf("uncorrected model(%+q) = %+q, %v; master = %+q, %v", input, original, originalErr, old, oldErr)
						}
						if widthFold {
							want, wantErr = model.Run(input)
						}
					}
					got, gotErr := head.Run(input)
					if got != want || gotErr != wantErr { //nolint:errorlint // Compare oracle error identity.
						t.Fatalf("Run(%+q) = %+q, %v; model = %+q, %v; master = %+q", input, got, gotErr, want, wantErr, old)
					}
				}
				t.Logf("compared %d deterministic inputs", len(inputs))
			})
		}
	}
}
