package textnorm

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/runes"
)

// baseSearchPreset freezes SearchPreset from master at 861ea36. Only the
// function name and final JoinTokens call are adapted to use the frozen helper
// below; keep the original stages, including the absence of post-filter NFC.
func baseSearchPreset(opts ...PresetOption) Pipeline {
	cfg := presetConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	pipe := New().SanitizeUTF8().NormalizeUnicode()
	if cfg.widthFold {
		pipe = pipe.FoldWidth()
	}

	keep := runes.Predicate(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)
	})

	return baseJoinTokens(pipe.
		FoldCase().
		FilterRunes(keep).
		TrimSpace().
		CollapseWhitespace().
		SplitTokens(), " ")
}

// baseJoinTokens has the verbatim body of JoinTokens at 861ea36, with the
// receiver changed to an argument. tp.source.Then is essential: it runs the
// source once before tp.Run runs that same source again.
func baseJoinTokens(tp TokenPipeline, sep string) Pipeline {
	return tp.source.Then(func(input string) (string, error) {
		tokens, err := tp.Run(input)
		if err != nil {
			return "", err
		}

		return strings.Join(tokens, sep), nil
	})
}

func checkSearchPresetAgainstBase(t *testing.T, input string) {
	t.Helper()
	for _, widthFold := range []bool{false, true} {
		var opts []PresetOption
		if widthFold {
			opts = append(opts, WithWidthFold())
		}
		want, err := baseSearchPreset(opts...).Run(input)
		if err != nil {
			t.Fatalf("base SearchPreset(%q, widthFold=%t): %v", input, widthFold, err)
		}
		got, err := SearchPreset(opts...).Run(input)
		if err != nil || got != want {
			t.Fatalf("SearchPreset(%q, widthFold=%t) = %q, %v; base = %q", input, widthFold, got, err, want)
		}
	}
}

func TestSearchPresetDifferential(t *testing.T) {
	for _, input := range searchPresetSeeds {
		t.Run(input, func(t *testing.T) {
			checkSearchPresetAgainstBase(t, input)
		})
	}
}
