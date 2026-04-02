package textnorm

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeUnicode appends a Unicode normalization stage.
func (p Pipeline) NormalizeUnicode() Pipeline {
	return p.Then(normalizeUnicodeStage)
}

// RemoveAccents appends a diacritic-removal stage.
func (p Pipeline) RemoveAccents() Pipeline {
	return p.Then(normalizeUnicodeStage)
}

func normalizeUnicodeStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}
