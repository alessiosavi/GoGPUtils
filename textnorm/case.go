package textnorm

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// FoldCase appends full Unicode case folding.
// Cherokee case pairs fold to uppercase, their stable Unicode representation.
// This differs from Lower, which maps Cherokee to lowercase.
func (p Pipeline) FoldCase() Pipeline {
	return p.Then(func(s string) (string, error) {
		if isASCII(s) {
			return strings.ToLower(s), nil
		}
		return foldCherokee(cases.Fold().String(s)), nil
	})
}

// foldCherokee corrects the Cherokee case toggle in cases.Fold:
// https://github.com/golang/go/issues/46101 (golang/go#46101).
// Copy untouched spans verbatim so this correction preserves malformed bytes.
func foldCherokee(s string) string {
	// Lowercase Cherokee uses the UTF-8 prefixes EA AD (U+AB70–U+AB7F),
	// EA AE (U+AB80–U+ABBF), or E1 8F (U+13F8–U+13FD). Skip rune
	// decoding when none is present, including Korean and Greek Extended text.
	first := strings.Index(s, "\xea\xad")
	if i := strings.Index(s, "\xea\xae"); first < 0 || (i >= 0 && i < first) {
		first = i
	}
	if i := strings.Index(s, "\xe1\x8f"); first < 0 || (i >= 0 && i < first) {
		first = i
	}
	if first < 0 {
		return s
	}
	// Neither lead byte can be a UTF-8 continuation byte, so starting here is
	// safe even when the candidate or the preceding bytes are malformed.
	var b strings.Builder
	last := 0
	for i, r := range s[first:] {
		var upper rune
		switch {
		case r >= '\uAB70' && r <= '\uABBF':
			upper = r - '\uAB70' + '\u13A0'
		case r >= '\u13F8' && r <= '\u13FD':
			upper = r - '\u13F8' + '\u13F0'
		default:
			continue
		}
		i += first
		if last == 0 {
			// Both sides of every pair have the same UTF-8 length.
			b.Grow(len(s))
		}
		b.WriteString(s[last:i])
		b.WriteRune(upper)
		last = i + utf8.RuneLen(r)
	}
	if last == 0 {
		return s
	}
	b.WriteString(s[last:])
	return b.String()
}

// Lower appends Unicode-aware lowercasing.
func (p Pipeline) Lower() Pipeline {
	return p.Then(func(s string) (string, error) {
		if isASCII(s) {
			return strings.ToLower(s), nil
		}
		return cases.Lower(language.Und).String(s), nil
	})
}

// MapRunes appends a rune-mapping stage.
func (p Pipeline) MapRunes(fn func(rune) rune) Pipeline {
	if fn == nil {
		return p
	}

	return p.Then(func(s string) (string, error) {
		result, _, err := transform.String(runes.Map(fn), s)
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
		result, _, err := transform.String(runes.Remove(remove), s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}

var _ = unicode.MaxRune
