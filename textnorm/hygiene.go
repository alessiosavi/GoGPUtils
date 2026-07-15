// Package textnorm — hygiene.go: byte-hygiene stages for text that will be
// sent to an external model or stored verbatim. These stages never touch
// case, punctuation, or diacritics.
package textnorm

import (
	"html"
	"strings"
	"unicode"
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
		var b strings.Builder
		b.Grow(len(s))
		for _, r := range s {
			if r != '\n' && r != '\t' && (unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Cc, r)) {
				continue
			}
			b.WriteRune(r)
		}
		return b.String(), nil
	})
}
