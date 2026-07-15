// Package textnorm — meaning.go implements meaning-preserving punctuation
// filtering: instead of deleting non-alphanumeric runes (which merges tokens,
// e.g. "2.5"→"25"), every rejected rune becomes a single space, and the small
// set of punctuation that carries meaning in product/search text is kept:
//
//   - '.' and ',' when BOTH neighbours are digits ("4.5", "1,000").
//     Separators are NOT unified: "1,000" and "1.000" stay distinct because
//     the same string means different numbers in different locales.
//   - '+' when it terminates an alphanumeric token ("s22+", "c++").
//   - '%' immediately after a digit ("100%").
package textnorm

import (
	"strings"
	"unicode"
)

// PreserveMeaningPunct appends the meaning-preserving punctuation stage.
func (p Pipeline) PreserveMeaningPunct() Pipeline {
	return p.Then(preserveMeaningPunct)
}

func preserveMeaningPunct(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s))
	for i, r := range rs {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r):
			b.WriteRune(r)
		case r == '.' || r == ',':
			if i > 0 && i < len(rs)-1 && unicode.IsDigit(rs[i-1]) && unicode.IsDigit(rs[i+1]) {
				b.WriteRune(r)
			} else {
				b.WriteByte(' ')
			}
		case r == '+':
			prevOK := i > 0 && (unicode.IsLetter(rs[i-1]) || unicode.IsDigit(rs[i-1]) || rs[i-1] == '+')
			nextOK := i == len(rs)-1 || unicode.IsSpace(rs[i+1]) || rs[i+1] == '+'
			if prevOK && nextOK {
				b.WriteRune(r)
			} else {
				b.WriteByte(' ')
			}
		case r == '%':
			if i > 0 && unicode.IsDigit(rs[i-1]) {
				b.WriteRune(r)
			} else {
				b.WriteByte(' ')
			}
		default:
			b.WriteByte(' ')
		}
	}
	return b.String(), nil
}
