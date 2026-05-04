package textnorm

import (
	"unicode"

	"golang.org/x/text/runes"
)

type presetConfig struct {
	widthFold bool
}

// PresetOption customizes preset builders.
type PresetOption func(*presetConfig)

// WithWidthFold enables explicit width folding in a preset pipeline.
func WithWidthFold() PresetOption {
	return func(cfg *presetConfig) {
		cfg.widthFold = true
	}
}

// SearchPreset builds a search-key pipeline (UTF-8 sanitize, Unicode
// normalize, fold case, keep letters/numbers/spaces, trim, collapse).
//
// To dedup repeated tokens or strip stopwords, compose on top:
//
//	textnorm.SearchPreset().
//	    SplitTokens().
//	    DedupTokens().
//	    RemoveStopwords(stopwords.English()).
//	    JoinTokens(" ")
func SearchPreset(opts ...PresetOption) Pipeline {
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

	return pipe.
		FoldCase().
		FilterRunes(keep).
		TrimSpace().
		CollapseWhitespace().
		SplitTokens().
		JoinTokens(" ")
}

// CanonicalPreset builds a general-purpose canonicalization pipeline.
func CanonicalPreset(opts ...PresetOption) Pipeline {
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

	return pipe.
		FoldCase().
		TrimSpace().
		CollapseWhitespace()
}

// DBSafePreset builds a persistence-safe normalization pipeline.
func DBSafePreset(opts ...PresetOption) Pipeline {
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

	return pipe.
		TrimSpace().
		CollapseWhitespace()
}
