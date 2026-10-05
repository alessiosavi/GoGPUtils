package textnorm

import (
	"unicode"

	"golang.org/x/text/runes"
)

// Frozen preset constructors from master at f94af78; only their names change.
// Keep these bodies verbatim so the width-fold model is checked against BASE.
func widthFoldBaseCanonicalPreset(opts ...PresetOption) Pipeline {
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

func widthFoldBaseDBSafePreset(opts ...PresetOption) Pipeline {
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

func widthFoldBaseMeaningPreset(opts ...PresetOption) Pipeline {
	cfg := presetConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	pipe := New().SanitizeUTF8().DecodeHTMLEntities().NormalizeUnicodeLatin()
	if cfg.widthFold {
		pipe = pipe.FoldWidth()
	}

	return pipe.
		FoldCase().
		PreserveMeaningPunct().
		TrimSpace().
		CollapseWhitespace()
}

func widthFoldBaseSearchPreset(opts ...PresetOption) Pipeline {
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
		NormalizeUnicode().
		TrimSpace().
		CollapseWhitespace().
		SplitTokens().
		JoinTokens(" ")
}

func widthFoldBaseHygienePreset(opts ...PresetOption) Pipeline {
	cfg := presetConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	pipe := New().SanitizeUTF8().DecodeHTMLEntities().RemoveFormatChars()
	if cfg.widthFold {
		pipe = pipe.FoldWidth()
	}

	return pipe.
		TrimSpace().
		CollapseWhitespace()
}
