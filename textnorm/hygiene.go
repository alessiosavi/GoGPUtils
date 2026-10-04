// Byte-hygiene stages for text that will be sent to an external model or
// stored verbatim. These stages never touch case, punctuation, or diacritics.

package textnorm

import (
	"html"
	"strings"
	"unicode"
	"unicode/utf8"
)

// DecodeHTMLEntities appends an HTML-entity decoding stage ("&amp;" → "&").
func (p Pipeline) DecodeHTMLEntities() Pipeline {
	return p.Then(func(s string) (string, error) {
		return html.UnescapeString(s), nil
	})
}

// RemoveFormatChars appends a stage that drops Unicode format (Cf) and
// control (Cc) runes — zero-width spaces, BiDi marks, NULs — except '\n' and
// '\t', which downstream whitespace collapsing normalizes.
func (p Pipeline) RemoveFormatChars() Pipeline {
	return p.Then(func(s string) (string, error) {
		first := -1
		for i, r := range s {
			if (r != '\n' && r != '\t' && (unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Cc, r))) || (r == utf8.RuneError && func() bool { _, width := utf8.DecodeRuneInString(s[i:]); return width == 1 }()) {
				first = i
				break
			}
		}
		if first < 0 {
			return strings.Clone(s), nil
		}
		var b strings.Builder
		b.Grow(len(s))
		b.WriteString(s[:first])
		for _, r := range s[first:] {
			if r != '\n' && r != '\t' && (unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Cc, r)) {
				continue
			}
			b.WriteRune(r)
		}
		return b.String(), nil
	})
}
