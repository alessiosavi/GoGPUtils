package textnorm

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeUnicode appends NFD decomposition, combining-mark (Mn) removal,
// and NFC recomposition.
//
// Decomposition guarantees hold for valid UTF-8. A malformed byte can prevent
// decomposition of the following character, leaving its accent in the output.
// Pipelines over untrusted bytes should start with SanitizeUTF8(), as every
// preset does.
func (p Pipeline) NormalizeUnicode() Pipeline {
	return p.Then(normalizeUnicodeStage)
}

// RemoveAccents appends a diacritic-removal stage.
// It performs the same transformation as NormalizeUnicode.
//
// Decomposition guarantees hold for valid UTF-8. A malformed byte can prevent
// decomposition of the following character, leaving its accent in the output.
// Pipelines over untrusted bytes should start with SanitizeUTF8(), as every
// preset does.
func (p Pipeline) RemoveAccents() Pipeline {
	return p.Then(normalizeUnicodeStage)
}

func normalizeUnicodeStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if isASCII(s) {
		// Preserve x/text v0.42.0 transform.String ownership: its 128-byte
		// initial buffer aliases unchanged input; larger input is detached.
		if len(s) > 128 {
			return strings.Clone(s), nil
		}
		return s, nil
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
//
// Decomposition guarantees hold for valid UTF-8. A malformed byte can prevent
// decomposition of the following character, leaving its accent in the output.
// Pipelines over untrusted bytes should start with SanitizeUTF8(), as every
// preset does.
func (p Pipeline) NormalizeUnicodeLatin() Pipeline {
	return p.Then(normalizeUnicodeLatinStage)
}

func normalizeUnicodeLatinStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if isASCII(s) {
		// Preserve the original Builder's detached result.
		return strings.Clone(s), nil
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

// isASCII reports whether every input byte is below 0x80.
func isASCII(s string) bool {
	for len(s) >= 8 {
		// The compiler combines these byte accesses into a single word load.
		word := uint64(s[0]) | uint64(s[1])<<8 | uint64(s[2])<<16 | uint64(s[3])<<24 |
			uint64(s[4])<<32 | uint64(s[5])<<40 | uint64(s[6])<<48 | uint64(s[7])<<56
		if word&0x8080808080808080 != 0 {
			return false
		}
		s = s[8:]
	}
	for i := range len(s) {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}
