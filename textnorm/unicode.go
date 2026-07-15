package textnorm

import (
	"strings"
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
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

// NormalizeUnicodeLatin appends a script-aware diacritic-removal stage:
// combining marks (Mn) are removed ONLY when their base character is Latin.
// Latin "café"→"cafe", but Devanagari matras and Arabic harakat — which are
// meaning-bearing vowels, not decorations — survive intact.
func (p Pipeline) NormalizeUnicodeLatin() Pipeline {
	return p.Then(normalizeUnicodeLatinStage)
}

func normalizeUnicodeLatinStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	decomposed := norm.NFD.String(s)
	var b strings.Builder
	b.Grow(len(decomposed))
	lastBaseLatin := false
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			if lastBaseLatin {
				continue
			}
			b.WriteRune(r)
			continue
		}
		lastBaseLatin = unicode.Is(unicode.Latin, r)
		b.WriteRune(r)
	}
	return norm.NFC.String(b.String()), nil
}
