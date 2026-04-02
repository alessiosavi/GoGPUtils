package textnorm

import (
	"strings"
	"unicode/utf8"
)

// TrimSpace appends a trimming stage.
func (p Pipeline) TrimSpace() Pipeline {
	return p.Then(trimSpaceStage)
}

// CollapseWhitespace appends a whitespace-collapsing stage.
func (p Pipeline) CollapseWhitespace() Pipeline {
	return p.Then(collapseWhitespaceStage)
}

// SanitizeUTF8 appends a UTF-8 and NUL sanitization stage.
func (p Pipeline) SanitizeUTF8() Pipeline {
	return p.Then(sanitizeUTF8Stage)
}

func trimSpaceStage(s string) (string, error) {
	return strings.TrimSpace(s), nil
}

func collapseWhitespaceStage(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	return strings.Join(strings.Fields(s), " "), nil
}

func sanitizeUTF8Stage(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	if utf8.ValidString(s) && !strings.ContainsRune(s, '\x00') {
		return s, nil
	}

	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			b.WriteRune('\uFFFD')
			i++

			continue
		}

		if r == '\x00' {
			i += size

			continue
		}

		b.WriteRune(r)
		i += size
	}

	return b.String(), nil
}
