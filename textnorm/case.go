package textnorm

import (
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// FoldCase appends full Unicode case folding.
func (p Pipeline) FoldCase() Pipeline {
	return p.Then(func(s string) (string, error) {
		return cases.Fold().String(s), nil
	})
}

// Lower appends Unicode-aware lowercasing.
func (p Pipeline) Lower() Pipeline {
	return p.Then(func(s string) (string, error) {
		return cases.Lower(language.Und).String(s), nil
	})
}

// MapRunes appends a rune-mapping stage.
func (p Pipeline) MapRunes(fn func(rune) rune) Pipeline {
	if fn == nil {
		return p
	}

	return p.Then(func(s string) (string, error) {
		result, _, err := transform.String(transform.Chain(runes.Map(fn)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}

// FilterRunes appends a rune-filtering stage.
func (p Pipeline) FilterRunes(keep runes.Set) Pipeline {
	if keep == nil {
		return p
	}

	return p.Then(func(s string) (string, error) {
		remove := runes.Predicate(func(r rune) bool {
			return !keep.Contains(r)
		})
		result, _, err := transform.String(transform.Chain(runes.Remove(remove)), s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}

var _ = unicode.MaxRune
