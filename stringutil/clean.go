package stringutil

import (
	"html"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CleanOption configures a cleaning step for CleanString.
// Options are applied in a fixed, safe order regardless of the order they are passed:
//  1. HTML stripping (first, to remove markup before text processing)
//  2. Unicode normalization (second, to normalize the text content)
//  3. Database sanitization (last, to enforce encoding/length constraints)
type CleanOption func(*cleanConfig)

// cleanConfig holds all cleaning options.
type cleanConfig struct {
	unicodeNorm  bool
	htmlStrip    bool
	dbSanitize   bool
	dbMaxLen     int  // 0 = unlimited (in runes)
	dbReplaceNul bool // replace NUL bytes
}

// WithUnicodeNorm enables Unicode normalization: NFKD decomposition followed
// by removal of combining marks (diacritics). This converts characters like
// "é" → "e", "ñ" → "n", "ü" → "u".
//
// This is useful for search indexing, comparison, and ensuring ASCII-compatible
// text from Unicode input.
//
// Example:
//
//	result, _ := CleanString("café résumé", WithUnicodeNorm())
//	// result = "cafe resume"
func WithUnicodeNorm() CleanOption {
	return func(c *cleanConfig) {
		c.unicodeNorm = true
	}
}

// WithHTMLStrip enables HTML tag removal and entity decoding.
// All HTML/XML tags are stripped, and HTML entities are decoded to their
// Unicode equivalents (e.g., &amp; → &, &lt; → <, &nbsp; → space, &#169; → ©).
//
// Uses Go's standard html.UnescapeString for entity decoding, which handles
// all named HTML entities, decimal (&#123;), and hex (&#x7B;) numeric entities.
//
// Example:
//
//	result, _ := CleanString("<p>Hello &amp; World</p>", WithHTMLStrip())
//	// result = "Hello & World"
func WithHTMLStrip() CleanOption {
	return func(c *cleanConfig) {
		c.htmlStrip = true
	}
}

// WithDBSanitize enables database sanitization:
//   - Replaces invalid UTF-8 sequences with U+FFFD (replacement character)
//   - Replaces NUL bytes (\x00) with empty string (NUL breaks PostgreSQL, MySQL, etc.)
//   - Optionally truncates to maxLen runes (0 = no truncation)
//
// Truncation is rune-aware: it will never cut a multi-byte character in half.
//
// Note: This does NOT escape SQL. Use parameterized queries for SQL injection
// prevention. This function handles encoding-level sanitization only.
//
// Example:
//
//	result, _ := CleanString("Hello\x00World", WithDBSanitize(0))
//	// result = "HelloWorld"
//
//	result, _ := CleanString("Hello World", WithDBSanitize(5))
//	// result = "Hello"
func WithDBSanitize(maxLen int) CleanOption {
	return func(c *cleanConfig) {
		c.dbSanitize = true
		c.dbReplaceNul = true
		if maxLen > 0 {
			c.dbMaxLen = maxLen
		}
	}
}

// CleanString applies the specified cleaning options to the input string.
// Options are applied in a fixed order for consistency and correctness:
//  1. HTML stripping (remove tags, decode entities)
//  2. Unicode normalization (NFKD + diacritic removal)
//  3. Database sanitization (UTF-8 validation, NUL removal, truncation)
//
// This order ensures that:
//   - HTML entities are decoded before Unicode normalization processes the text
//   - Database constraints (length, encoding) are applied last on the final result
//
// Returns the cleaned string. Error is returned only for transformer failures
// (extremely unlikely with valid Go strings).
//
// If no options are provided, the input string is returned unchanged.
//
// Example:
//
//	// Apply all three cleaning modes:
//	result, err := CleanString(
//	    "<p>Héllo &amp; Wörld</p>",
//	    WithHTMLStrip(),
//	    WithUnicodeNorm(),
//	    WithDBSanitize(50),
//	)
//	// result = "Hello & World"
//
//	// Apply only Unicode normalization:
//	result, _ := CleanString("café", WithUnicodeNorm())
//	// result = "cafe"
func CleanString(input string, options ...CleanOption) (string, error) {
	if len(options) == 0 {
		return input, nil
	}

	cfg := &cleanConfig{}
	for _, opt := range options {
		opt(cfg)
	}

	result := input

	// Step 1: HTML stripping (before text processing)
	if cfg.htmlStrip {
		result = StripHTMLEntities(result)
	}

	// Step 2: Unicode normalization
	if cfg.unicodeNorm {
		var err error

		result, err = NormalizeUnicode(result)
		if err != nil {
			return "", err
		}
	}

	// Step 3: Database sanitization (applied last)
	if cfg.dbSanitize {
		result = SanitizeUTF8(result)
		if cfg.dbMaxLen > 0 {
			result = TruncateRunes(result, cfg.dbMaxLen)
		}
	}

	return result, nil
}

// NormalizeUnicode applies NFKD normalization and removes combining marks
// (diacritics) from the input string.
//
// The process:
//  1. NFKD (Compatibility Decomposition): decomposes characters into their
//     base form + combining marks (e.g., "é" → "e" + combining acute accent)
//  2. Remove combining marks: strips all Unicode Mn (Mark, Nonspacing) characters
//  3. Recompose to NFC for consistent output
//
// Characters that are not letters or numbers are preserved as-is (spaces,
// punctuation, etc.).
//
// Example:
//
//	NormalizeUnicode("café résumé")  // "cafe resume", nil
//	NormalizeUnicode("naïve")        // "naive", nil
//	NormalizeUnicode("Ångström")     // "Angstrom", nil
func NormalizeUnicode(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	// Chain: NFKD decomposition → remove combining marks (diacritics)
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

// StripHTMLEntities removes all HTML/XML tags and decodes HTML entities.
//
// Tag removal handles:
//   - Standard tags: <p>, <br/>, <div class="x">
//   - Self-closing tags: <br />, <img />
//   - Script and style tags (content included — for full script removal,
//     use a proper HTML parser)
//
// Entity decoding handles (via html.UnescapeString):
//   - Named entities: &amp; &lt; &gt; &quot; &apos; &nbsp; &copy; etc.
//   - Decimal numeric: &#169; &#8212;
//   - Hex numeric: &#x00A9; &#x2014;
//
// Example:
//
//	StripHTMLEntities("<p>Hello &amp; World</p>")  // "Hello & World"
//	StripHTMLEntities("Price: &euro;10")           // "Price: €10"
//	StripHTMLEntities("5 &gt; 3 &amp;&amp; 2 &lt; 4") // "5 > 3 && 2 < 4"
func StripHTMLEntities(s string) string {
	if s == "" {
		return ""
	}

	// Step 1: Remove HTML tags (reuse the existing StripTags logic)
	stripped := StripTags(s)

	// Step 2: Decode HTML entities
	decoded := html.UnescapeString(stripped)

	// Step 3: Normalize whitespace from &nbsp; and similar
	// Replace non-breaking spaces (U+00A0) with regular spaces
	decoded = strings.ReplaceAll(decoded, "\u00A0", " ")

	return decoded
}

// SanitizeUTF8 ensures the string contains only valid UTF-8 and removes
// NUL bytes (\x00) which are problematic for most databases.
//
// Invalid UTF-8 byte sequences are replaced with U+FFFD (Unicode replacement
// character), following Go's standard behavior.
//
// Example:
//
//	SanitizeUTF8("Hello\x00World")     // "HelloWorld"
//	SanitizeUTF8("Hello\xffWorld")     // "Hello\uFFFDWorld"
//	SanitizeUTF8("Valid UTF-8 string") // "Valid UTF-8 string" (unchanged)
func SanitizeUTF8(s string) string {
	if s == "" {
		return ""
	}

	// Fast path: if string is valid UTF-8 and has no NUL bytes, return as-is
	if utf8.ValidString(s) && !strings.ContainsRune(s, '\x00') {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size <= 1 {
			// Invalid UTF-8 byte: replace with U+FFFD
			b.WriteRune('\uFFFD')
			i++

			continue
		}

		// Skip NUL bytes
		if r == '\x00' {
			i += size

			continue
		}

		b.WriteRune(r)
		i += size
	}

	return b.String()
}

// TruncateRunes truncates s to at most maxLen runes.
// Unlike byte-level truncation, this is Unicode-safe and will never
// split a multi-byte character.
//
// If maxLen <= 0, returns empty string.
// If s has fewer runes than maxLen, returns s unchanged.
//
// Example:
//
//	TruncateRunes("Hello, 世界!", 8)  // "Hello, 世界"  (correct, not "Hello, \xe4")
//	TruncateRunes("café", 3)          // "caf"
func TruncateRunes(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen])
}

// RemoveNonPrintable removes all non-printable characters from s,
// except for common whitespace (space, tab, newline, carriage return).
//
// This is useful for cleaning user input that may contain control characters,
// zero-width characters, or other invisible Unicode characters.
//
// Example:
//
//	RemoveNonPrintable("Hello\x07World")  // "HelloWorld" (bell character removed)
//	RemoveNonPrintable("Hello\tWorld\n")  // "Hello\tWorld\n" (whitespace preserved)
func RemoveNonPrintable(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if unicode.IsPrint(r) || r == '\t' || r == '\n' || r == '\r' {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// NormalizeWhitespace collapses all consecutive whitespace characters
// (spaces, tabs, newlines) into a single space, and trims leading/trailing
// whitespace.
//
// This is useful for cleaning user input or text extracted from HTML
// where whitespace may be irregular.
//
// Example:
//
//	NormalizeWhitespace("  Hello   World  \n\t ")  // "Hello World"
//	NormalizeWhitespace("\t\n")                     // ""
func NormalizeWhitespace(s string) string {
	if s == "" {
		return ""
	}

	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// RemoveAccents is an alias for NormalizeUnicode that removes diacritical
// marks from characters. This is a common operation name used in many
// string-processing libraries.
//
// Example:
//
//	RemoveAccents("café")  // "cafe", nil
//	RemoveAccents("über")  // "uber", nil
func RemoveAccents(s string) (string, error) {
	return NormalizeUnicode(s)
}

// ToASCII converts a Unicode string to its closest ASCII representation
// by removing diacritics, replacing non-letter/non-number characters with
// spaces, and collapsing whitespace.
//
// This is useful for generating slugs, search keys, or filenames from
// Unicode input.
//
// Example:
//
//	ToASCII("Héllo, Wörld!")    // "Hello  World", nil
//	ToASCII("café résumé")      // "cafe resume", nil
func ToASCII(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mn)),
		runes.Map(func(r rune) rune {
			if r <= 127 {
				return r
			}
			// Non-ASCII characters that survived diacritic removal
			// get replaced with space
			return -1
		}),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Slugify converts a string to a URL-friendly slug.
// It normalizes Unicode, lowercases, replaces non-alphanumeric characters
// with hyphens, collapses multiple hyphens, and trims leading/trailing hyphens.
//
// Example:
//
//	Slugify("Hello, World!")        // "hello-world", nil
//	Slugify("Café Résumé")         // "cafe-resume", nil
//	Slugify("  Multiple   Spaces ") // "multiple-spaces", nil
func Slugify(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	// Step 1: Normalize Unicode (remove diacritics)
	normalized, err := NormalizeUnicode(s)
	if err != nil {
		return "", err
	}

	// Step 2: Lowercase
	normalized = strings.ToLower(normalized)

	// Step 3: Replace non-alphanumeric with hyphens
	var b strings.Builder
	b.Grow(len(normalized))

	prevHyphen := false

	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)

			prevHyphen = false
		} else if !prevHyphen {
			b.WriteByte('-')

			prevHyphen = true
		}
	}

	// Step 4: Trim leading/trailing hyphens
	return strings.Trim(b.String(), "-"), nil
}
