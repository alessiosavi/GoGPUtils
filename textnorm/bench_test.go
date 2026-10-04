package textnorm

import (
	"strings"
	"testing"
	"unicode"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkPipeline_Run executes a prebuilt NormalizeUnicode, FoldCase,
// CollapseWhitespace pipeline on n bytes of Unicode text. case=empty preserves
// the empty-pipeline workload on n bytes of ASCII text. Inputs are immutable;
// no restoration is measured.
func BenchmarkPipeline_Run(b *testing.B) {
	b.Run("case=unicode", func(b *testing.B) {
		p := New().NormalizeUnicode().FoldCase().CollapseWhitespace()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("case=empty", func(b *testing.B) {
		p := New()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.ASCII)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkSearchPreset measures preset construction (mode=build) and
// execution of a prebuilt preset on n bytes of Mixed text (mode=run).
// No options are supplied; input is immutable and no restoration is measured.
func BenchmarkSearchPreset(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			SearchPreset()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := SearchPreset()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkCanonicalPreset measures preset construction (mode=build) and
// execution of a prebuilt preset on n bytes of Mixed text (mode=run).
// No options are supplied; input is immutable and no restoration is measured.
func BenchmarkCanonicalPreset(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			CanonicalPreset()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := CanonicalPreset()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkDBSafePreset measures preset construction (mode=build) and
// execution of a prebuilt preset on n bytes of Mixed text (mode=run).
// No options are supplied; input is immutable and no restoration is measured.
// case=dirty prefixes Mixed text with NUL and an invalid byte, within n bytes,
// preserving the old benchmark's sanitization workload.
func BenchmarkDBSafePreset(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			DBSafePreset()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := DBSafePreset()
		b.Run("case=mixed", func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := benchkit.Text(n, benchkit.Mixed)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(text); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
		b.Run("case=dirty", func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := "\x00\xff" + benchkit.Text(n-2, benchkit.Mixed)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(text); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	})
}

// benchmarkRuneSet satisfies the FilterRunes set interface without an external
// test import. It keeps letters, numbers and whitespace.
type benchmarkRuneSet struct{}

func (benchmarkRuneSet) Contains(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)
}

// BenchmarkNew constructs a fresh empty pipeline. There is no text input,
// deferred stage or mutable state to restore.
func BenchmarkNew(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New()
		}
	})
}

// BenchmarkPipeline_Then appends a non-nil lowercasing Stage to an empty
// pipeline in mode=build; mode=run executes the prebuilt stage on n bytes of
// Unicode text. Each iteration uses immutable input; no restoration is needed.
func BenchmarkPipeline_Then(b *testing.B) {
	stage := func(s string) (string, error) { return strings.ToLower(s), nil }
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().Then(stage)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().Then(stage)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_FoldCase measures full Unicode case folding. mode=build
// appends to an empty pipeline; mode=run executes the prebuilt stage on n
// bytes of immutable Unicode text. No restoration is measured.
func BenchmarkPipeline_FoldCase(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().FoldCase()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().FoldCase()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_Lower measures Unicode lowercasing. mode=build appends
// to an empty pipeline; mode=run executes the prebuilt stage on n bytes of
// immutable Unicode text. No restoration is measured.
func BenchmarkPipeline_Lower(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().Lower()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().Lower()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_MapRunes maps runes with unicode.ToLower. mode=build
// appends to an empty pipeline; mode=run executes the prebuilt stage on n
// bytes of immutable Unicode text. No restoration is measured.
func BenchmarkPipeline_MapRunes(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().MapRunes(unicode.ToLower)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().MapRunes(unicode.ToLower)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_FilterRunes keeps letters, numbers and whitespace.
// mode=build appends to an empty pipeline; mode=run executes the prebuilt
// stage on n bytes of immutable Unicode text, rejecting marks, punctuation
// and emoji. No restoration is measured.
func BenchmarkPipeline_FilterRunes(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().FilterRunes(benchmarkRuneSet{})
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().FilterRunes(benchmarkRuneSet{})
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_DecodeHTMLEntities decodes HTML entities. mode=build
// appends to an empty pipeline; mode=run executes the prebuilt stage on n
// bytes of Mixed text prefixed with "&amp; &lt; " within that byte budget.
// Input is immutable; no restoration is measured.
func BenchmarkPipeline_DecodeHTMLEntities(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().DecodeHTMLEntities()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().DecodeHTMLEntities()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "&amp; &lt; " + benchkit.Text(n-len("&amp; &lt; "), benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_RemoveFormatChars removes format and control runes.
// mode=build appends to an empty pipeline; mode=run executes the prebuilt
// stage on n bytes of Unicode text prefixed with U+200B, U+202E and NUL
// within that byte budget. Input is immutable; no restoration is measured.
func BenchmarkPipeline_RemoveFormatChars(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().RemoveFormatChars()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().RemoveFormatChars()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "\u200b\u202e\x00" + benchkit.Text(n-len("\u200b\u202e\x00"), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_PreserveMeaningPunct preserves meaningful punctuation.
// mode=build appends to an empty pipeline; mode=run executes the prebuilt
// stage on n bytes of Unicode text prefixed with "4.5 c++ 100% ! " within
// that byte budget. Input is immutable; no restoration is measured.
func BenchmarkPipeline_PreserveMeaningPunct(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().PreserveMeaningPunct()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().PreserveMeaningPunct()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "4.5 c++ 100% ! " + benchkit.Text(n-len("4.5 c++ 100% ! "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_RemoveAccents removes accents. mode=build appends to an
// empty pipeline; mode=run executes the prebuilt stage on n bytes of
// immutable Unicode text with composed/decomposed accents. No restoration
// is measured.
func BenchmarkPipeline_RemoveAccents(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().RemoveAccents()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().RemoveAccents()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_NormalizeUnicodeLatin removes Latin accents while keeping
// non-Latin combining marks. mode=build appends to an empty pipeline;
// mode=run executes the prebuilt stage on n bytes of Unicode text prefixed
// with "é क़ " within that byte budget. Input is immutable; no restoration
// is measured.
func BenchmarkPipeline_NormalizeUnicodeLatin(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().NormalizeUnicodeLatin()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().NormalizeUnicodeLatin()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "é क़ " + benchkit.Text(n-len("é क़ "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_TrimSpace trims whitespace. mode=build appends to an
// empty pipeline; mode=run executes the prebuilt stage on n bytes of Mixed
// text with three-byte whitespace runs at both ends, included in n.
// Input is immutable; no restoration is measured.
func BenchmarkPipeline_TrimSpace(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().TrimSpace()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().TrimSpace()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := " \t\n" + benchkit.Text(n-6, benchkit.Mixed) + "\n\t "
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_CollapseWhitespace collapses whitespace. mode=build
// appends to an empty pipeline; mode=run executes the prebuilt stage on n
// bytes of Mixed text prefixed with x, a space, a tab, a newline and a space,
// all within that byte budget, so
// an internal whitespace run is collapsed. Input is immutable; no
// restoration is measured.
func BenchmarkPipeline_CollapseWhitespace(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().CollapseWhitespace()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().CollapseWhitespace()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "x \t\n " + benchkit.Text(n-len("x \t\n "), benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_FoldWidth folds fullwidth characters. mode=build
// appends to an empty pipeline; mode=run executes the prebuilt stage on n
// bytes of Unicode text prefixed with "Ｇｏ " within that byte budget.
// Input is immutable; no restoration is measured.
func BenchmarkPipeline_FoldWidth(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().FoldWidth()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().FoldWidth()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "Ｇｏ " + benchkit.Text(n-len("Ｇｏ "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_SplitTokens constructs an empty token pipeline in
// mode=build; mode=run executes it on n bytes of immutable Unicode text,
// including whitespace tokenization. The source string pipeline is empty;
// fresh tokens are returned each time and no restoration is measured.
func BenchmarkPipeline_SplitTokens(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().SplitTokens()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().SplitTokens()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_NormalizeUnicode measures the NormalizeUnicode stage.
// mode=build appends the stage to an empty pipeline. mode=run runs a
// pipeline holding only that stage on n bytes of Unicode text.
func BenchmarkPipeline_NormalizeUnicode(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().NormalizeUnicode()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().NormalizeUnicode()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkPipeline_SanitizeUTF8 measures sanitization. mode=build appends
// to an empty pipeline; mode=run executes the prebuilt stage on n bytes.
// case=valid uses Unicode text; case=dirty reserves two bytes for NUL and
// an invalid byte before the Unicode text. Inputs are immutable; no
// restoration is measured.
func BenchmarkPipeline_SanitizeUTF8(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			New().SanitizeUTF8()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := New().SanitizeUTF8()
		b.Run("case=valid", func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := benchkit.Text(n, benchkit.Unicode)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(text); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
		b.Run("case=dirty", func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := "\x00\xff" + benchkit.Text(n-2, benchkit.Unicode)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(text); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	})
}

// BenchmarkWithWidthFold measures the option through CanonicalPreset.
// mode=build constructs the preset with the option; mode=run executes that
// prebuilt preset on n bytes of Unicode text prefixed with "Ｇｏ " within
// that byte budget. Input is immutable; no restoration is measured.
func BenchmarkWithWidthFold(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			CanonicalPreset(WithWidthFold())
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := CanonicalPreset(WithWidthFold())
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "Ｇｏ " + benchkit.Text(n-len("Ｇｏ "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkMeaningPreset measures preset construction (mode=build) and
// execution of a prebuilt preset on n bytes of Mixed text (mode=run).
// No options are supplied; input is immutable and no restoration is measured.
func BenchmarkMeaningPreset(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			MeaningPreset()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := MeaningPreset()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkHygienePreset measures preset construction (mode=build) and
// execution of a prebuilt preset on n bytes of Mixed text (mode=run).
// No options are supplied; input is immutable and no restoration is measured.
func BenchmarkHygienePreset(b *testing.B) {
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			HygienePreset()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := HygienePreset()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Mixed)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_Then appends a non-nil stage to an empty token
// pipeline in mode=build; mode=run executes the prebuilt token pipeline on
// n bytes of Unicode text. The stage allocates a new slice and lowercases
// every token. The source string pipeline is empty and input is immutable;
// no restoration is measured.
func BenchmarkTokenPipeline_Then(b *testing.B) {
	base := New().SplitTokens()
	stage := func(tokens []string) ([]string, error) {
		out := make([]string, len(tokens))
		for i, token := range tokens {
			out[i] = strings.ToLower(token)
		}
		return out, nil
	}
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.Then(stage)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.Then(stage)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_MapTokens maps tokens using strings.ToLower.
// mode=build appends to a fixed empty token pipeline; mode=run executes the
// prebuilt token pipeline on n bytes of Unicode text, including tokenization.
// The source string pipeline is empty; input is immutable and no restoration
// is measured.
func BenchmarkTokenPipeline_MapTokens(b *testing.B) {
	base := New().SplitTokens()
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.MapTokens(strings.ToLower)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.MapTokens(strings.ToLower)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_FilterTokens keeps tokens longer than four bytes.
// mode=build appends to a fixed empty token pipeline; mode=run executes the
// prebuilt token pipeline on n bytes of Unicode text, including tokenization.
// The source string pipeline is empty; input is immutable and no restoration
// is measured.
func BenchmarkTokenPipeline_FilterTokens(b *testing.B) {
	base := New().SplitTokens()
	keep := func(token string) bool { return len(token) > 4 }
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.FilterTokens(keep)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.FilterTokens(keep)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_DedupTokens drops repeated tokens. mode=build
// appends to a fixed empty token pipeline; mode=run executes it on n bytes
// of Unicode text prefixed with "café café " within that byte budget, so
// even the smallest case has duplicates. The source string pipeline is
// empty; tokenization is measured, input is immutable and no restoration
// is measured.
func BenchmarkTokenPipeline_DedupTokens(b *testing.B) {
	base := New().SplitTokens()
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.DedupTokens()
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.DedupTokens()
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "café café " + benchkit.Text(n-len("café café "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_RemoveStopwords filters a fixed set containing café
// and 東京. mode=build appends to a fixed empty token pipeline; mode=run
// executes it on n bytes of Unicode text prefixed with "café " within that
// byte budget. The source string pipeline is empty; tokenization is measured,
// input and set are immutable and no restoration is measured.
func BenchmarkTokenPipeline_RemoveStopwords(b *testing.B) {
	base := New().SplitTokens()
	set := map[string]struct{}{"café": {}, "東京": {}}
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.RemoveStopwords(set)
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.RemoveStopwords(set)
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := "café " + benchkit.Text(n-len("café "), benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_JoinTokens joins lowercase, deduplicated tokens
// with a space. mode=build adds the join to a fixed token pipeline with
// MapTokens and DedupTokens; mode=run executes the prebuilt pipeline on n
// bytes of Unicode text, including tokenization and both token stages.
// The source string pipeline is empty; input is immutable and no restoration
// is measured.
func BenchmarkTokenPipeline_JoinTokens(b *testing.B) {
	base := New().SplitTokens().MapTokens(strings.ToLower).DedupTokens()
	b.Run("mode=build", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			base.JoinTokens(" ")
		}
	})
	b.Run("mode=run", func(b *testing.B) {
		p := base.JoinTokens(" ")
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.Unicode)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := p.Run(text); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkTokenPipeline_Run executes a prebuilt three-stage token pipeline:
// MapTokens(strings.ToLower), FilterTokens(keep tokens longer than four bytes),
// and DedupTokens. n is bytes of Unicode text; tokenization is measured and
// the source string pipeline is empty. Input is immutable and no restoration
// is measured.
func BenchmarkTokenPipeline_Run(b *testing.B) {
	p := New().SplitTokens().
		MapTokens(strings.ToLower).
		FilterTokens(func(token string) bool { return len(token) > 4 }).
		DedupTokens()
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.Unicode)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := p.Run(text); err != nil {
				b.Fatal(err)
			}
		}
	})
}
