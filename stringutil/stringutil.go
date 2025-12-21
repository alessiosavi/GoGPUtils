package stringutil

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// AllIndexes returns all starting positions of substr in s.
// Returns nil if substr is empty or s doesn't contain substr.
//
// Example:
//
//	indices := AllIndexes("banana", "an")
//	// indices = [1, 3]
func AllIndexes(s, substr string) []int {
	if substr == "" || len(s) < len(substr) {
		return nil
	}

	var indices []int

	offset := 0

	for {
		i := strings.Index(s[offset:], substr)
		if i == -1 {
			break
		}

		indices = append(indices, offset+i)
		offset += i + 1
	}

	return indices
}

// HasAnyPrefix reports whether s starts with any of the given prefixes.
//
// Example:
//
//	if HasAnyPrefix(url, "http://", "https://") { ... }
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}

// HasAnySuffix reports whether s ends with any of the given suffixes.
func HasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}

// ContainsAny reports whether s contains any of the given substrings.
//
// Example:
//
//	if ContainsAny(text, "error", "fail", "warning") { ... }
func ContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}

// ContainsAll reports whether s contains all of the given substrings.
func ContainsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}

	return true
}

// Reverse returns s with its characters in reverse order.
// Correctly handles multi-byte UTF-8 characters.
//
// Example:
//
//	rev := Reverse("hello")  // "olleh"
//	rev := Reverse("日本語")   // "語本日"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// IsPalindrome reports whether s reads the same forwards and backwards.
// Case-sensitive and ignores whitespace/punctuation only if normalize is true.
//
// Example:
//
//	IsPalindrome("racecar", false)           // true
//	IsPalindrome("A man a plan a canal Panama", true)  // true (normalized)
func IsPalindrome(s string, normalize bool) bool {
	if normalize {
		s = strings.ToLower(s)

		var b strings.Builder

		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				b.WriteRune(r)
			}
		}

		s = b.String()
	}

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}

	return true
}

// Truncate shortens s to maxLen characters, appending suffix if truncated.
// The total length including suffix will not exceed maxLen.
//
// Example:
//
//	Truncate("Hello World", 8, "...")  // "Hello..."
func Truncate(s string, maxLen int, suffix string) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	suffixRunes := []rune(suffix)
	if len(suffixRunes) >= maxLen {
		return string(suffixRunes[:maxLen])
	}

	truncateAt := maxLen - len(suffixRunes)

	return string(runes[:truncateAt]) + suffix
}

// TruncateWords truncates s at a word boundary, appending suffix if truncated.
// Attempts to break at word boundaries rather than mid-word.
func TruncateWords(s string, maxLen int, suffix string) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	suffixRunes := []rune(suffix)
	if len(suffixRunes) >= maxLen {
		return string(suffixRunes[:maxLen])
	}

	truncateAt := maxLen - len(suffixRunes)

	// Find last space before truncate point
	lastSpace := -1

	for i := truncateAt - 1; i >= 0; i-- {
		if unicode.IsSpace(runes[i]) {
			lastSpace = i

			break
		}
	}

	if lastSpace > 0 {
		return string(runes[:lastSpace]) + suffix
	}

	return string(runes[:truncateAt]) + suffix
}

// PadLeft pads s on the left with padChar to reach the target length.
// If s is already >= length, returns s unchanged.
//
// Example:
//
//	PadLeft("42", 5, '0')  // "00042"
func PadLeft(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	padding := length - len(runes)

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))

	for range padding {
		b.WriteRune(padChar)
	}

	b.WriteString(s)

	return b.String()
}

// PadRight pads s on the right with padChar to reach the target length.
//
// Example:
//
//	PadRight("42", 5, '0')  // "42000"
func PadRight(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	padding := length - len(runes)

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))
	b.WriteString(s)

	for range padding {
		b.WriteRune(padChar)
	}

	return b.String()
}

// PadCenter centers s by adding padChar on both sides.
// If odd padding needed, extra character goes on the right.
//
// Example:
//
//	PadCenter("hello", 11, '*')  // "***hello***"
func PadCenter(s string, length int, padChar rune) string {
	runes := []rune(s)
	if len(runes) >= length {
		return s
	}

	totalPadding := length - len(runes)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	var b strings.Builder

	b.Grow(length * utf8.RuneLen(padChar))

	for range leftPadding {
		b.WriteRune(padChar)
	}

	b.WriteString(s)

	for range rightPadding {
		b.WriteRune(padChar)
	}

	return b.String()
}

// RemoveAll removes all occurrences of the given substrings from s.
//
// Example:
//
//	clean := RemoveAll("hello world", "l", "o")  // "he wrd"
func RemoveAll(s string, substrs ...string) string {
	for _, substr := range substrs {
		s = strings.ReplaceAll(s, substr, "")
	}

	return s
}

// CountLines returns the number of lines in s.
// An empty string returns 0; a string without newlines returns 1.
func CountLines(s string) int {
	if s == "" {
		return 0
	}

	count := 1

	for _, r := range s {
		if r == '\n' {
			count++
		}
	}
	// Don't count trailing newline as extra line
	if strings.HasSuffix(s, "\n") {
		count--
	}

	return count
}

// Lines splits s into lines. Unlike strings.Split, handles \r\n properly.
//
// Example:
//
//	lines := Lines("a\nb\nc")  // ["a", "b", "c"]
func Lines(s string) []string {
	if s == "" {
		return nil
	}
	// Normalize line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")

	// Remove trailing empty line if string ended with newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// IsBlank reports whether s contains only whitespace characters.
func IsBlank(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// IsEmpty reports whether s is empty (zero length).
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsAlpha reports whether s contains only alphabetic characters.
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

// IsAlphanumeric reports whether s contains only letters and digits.
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// IsNumeric reports whether s contains only numeric digits.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// IsUpper reports whether all letters in s are uppercase.
// Returns true for strings with no letters.
func IsUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}

	return true
}

// IsLower reports whether all letters in s are lowercase.
// Returns true for strings with no letters.
func IsLower(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}

	return true
}

// IsASCII reports whether s contains only ASCII characters.
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}

	return true
}

// IsPrintable reports whether s contains only printable characters.
func IsPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}

// Capitalize returns s with the first character uppercased and the rest lowercased.
//
// Example:
//
//	Capitalize("hELLO")  // "Hello"
func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// Title returns s with the first character of each word uppercased.
//
// Example:
//
//	Title("hello world")  // "Hello World"
func Title(s string) string {
	return strings.Title(s) //nolint:staticcheck // strings.Title is deprecated but suitable here
}

// SwapCase swaps the case of each letter in s.
//
// Example:
//
//	SwapCase("Hello World")  // "hELLO wORLD"
func SwapCase(s string) string {
	var b strings.Builder

	b.Grow(len(s))

	for _, r := range s {
		if unicode.IsUpper(r) {
			b.WriteRune(unicode.ToLower(r))
		} else if unicode.IsLower(r) {
			b.WriteRune(unicode.ToUpper(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// SnakeCase converts s to snake_case.
//
// Example:
//
//	SnakeCase("HelloWorld")  // "hello_world"
//	SnakeCase("helloWorld")  // "hello_world"
func SnakeCase(s string) string {
	var b strings.Builder

	b.Grow(len(s) + 10) // Extra space for underscores

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}

			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// CamelCase converts s to camelCase.
//
// Example:
//
//	CamelCase("hello_world")  // "helloWorld"
//	CamelCase("hello-world")  // "helloWorld"
func CamelCase(s string) string {
	var b strings.Builder

	b.Grow(len(s))

	capitalizeNext := false

	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			capitalizeNext = true

			continue
		}

		if capitalizeNext {
			b.WriteRune(unicode.ToUpper(r))

			capitalizeNext = false
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	return b.String()
}

// PascalCase converts s to PascalCase.
//
// Example:
//
//	PascalCase("hello_world")  // "HelloWorld"
func PascalCase(s string) string {
	result := CamelCase(s)
	if result == "" {
		return ""
	}

	runes := []rune(result)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// KebabCase converts s to kebab-case.
//
// Example:
//
//	KebabCase("HelloWorld")  // "hello-world"
func KebabCase(s string) string {
	var b strings.Builder

	b.Grow(len(s) + 10)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('-')
			}

			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// Words splits s into words, treating any non-alphanumeric character as separator.
//
// Example:
//
//	Words("hello, world!")  // ["hello", "world"]
func Words(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// RuneCount returns the number of runes (Unicode code points) in s.
// This differs from len(s), which returns bytes.
//
// Example:
//
//	RuneCount("日本語")  // 3
//	len("日本語")      // 9 (bytes)
func RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}

// SafeSlice safely slices s by rune indices, returning an empty string for invalid ranges.
// Useful when working with user input where indices might be out of bounds.
func SafeSlice(s string, start, end int) string {
	runes := []rune(s)

	if start < 0 {
		start = 0
	}

	if end > len(runes) {
		end = len(runes)
	}

	if start >= end || start >= len(runes) {
		return ""
	}

	return string(runes[start:end])
}

// NthRune returns the rune at position n (0-indexed).
// Returns (0, false) if n is out of bounds.
func NthRune(s string, n int) (rune, bool) {
	if n < 0 {
		return 0, false
	}

	for i, r := range s {
		if i == n {
			return r, true
		}
	}

	return 0, false
}

// CommonPrefix returns the longest common prefix of the given strings.
// Returns empty string if no common prefix or fewer than 2 strings.
//
// Example:
//
//	CommonPrefix("interstellar", "internet", "internal")  // "inter"
func CommonPrefix(strs ...string) string {
	if len(strs) < 2 {
		if len(strs) == 1 {
			return strs[0]
		}

		return ""
	}

	// Find shortest string to bound our search
	minLen := len(strs[0])
	for _, s := range strs[1:] {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	var prefix strings.Builder

	for i := 0; i < minLen; i++ {
		char := strs[0][i]
		for _, s := range strs[1:] {
			if s[i] != char {
				return prefix.String()
			}
		}

		prefix.WriteByte(char)
	}

	return prefix.String()
}

// CommonSuffix returns the longest common suffix of the given strings.
func CommonSuffix(strs ...string) string {
	if len(strs) < 2 {
		if len(strs) == 1 {
			return strs[0]
		}

		return ""
	}

	// Reverse all strings, find prefix, then reverse result
	reversed := make([]string, len(strs))
	for i, s := range strs {
		reversed[i] = Reverse(s)
	}

	return Reverse(CommonPrefix(reversed...))
}

// Repeat returns s repeated n times.
// If n <= 0, returns empty string.
//
// Example:
//
//	Repeat("ab", 3)  // "ababab"
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}

	return strings.Repeat(s, n)
}

// Between extracts the substring between start and end markers.
// Returns empty string and false if markers not found in proper order.
//
// Example:
//
//	Between("[hello]", "[", "]")  // "hello", true
func Between(s, start, end string) (string, bool) {
	startIdx := strings.Index(s, start)
	if startIdx == -1 {
		return "", false
	}

	startIdx += len(start)

	endIdx := strings.Index(s[startIdx:], end)
	if endIdx == -1 {
		return "", false
	}

	return s[startIdx : startIdx+endIdx], true
}

// BetweenAll extracts all substrings between start and end markers.
//
// Example:
//
//	BetweenAll("a[1]b[2]c[3]", "[", "]")  // ["1", "2", "3"]
func BetweenAll(s, start, end string) []string {
	var results []string

	remaining := s

	for {
		result, ok := Between(remaining, start, end)
		if !ok {
			break
		}

		results = append(results, result)

		// Find the end marker and move past it
		idx := strings.Index(remaining, start)
		remaining = remaining[idx+len(start):]
		idx = strings.Index(remaining, end)
		remaining = remaining[idx+len(end):]
	}

	return results
}

// Wrap wraps text at the specified width, breaking at word boundaries.
// Preserves existing line breaks.
func Wrap(s string, width int) string {
	if width <= 0 {
		return s
	}

	var result strings.Builder
	for _, line := range Lines(s) {
		if result.Len() > 0 {
			result.WriteByte('\n')
		}

		result.WriteString(wrapLine(line, width))
	}

	return result.String()
}

func wrapLine(line string, width int) string {
	if len(line) <= width {
		return line
	}

	var result strings.Builder

	words := strings.Fields(line)
	lineLen := 0

	for i, word := range words {
		wordLen := utf8.RuneCountInString(word)

		if lineLen+wordLen > width && lineLen > 0 {
			result.WriteByte('\n')

			lineLen = 0
		} else if i > 0 {
			result.WriteByte(' ')

			lineLen++
		}

		result.WriteString(word)

		lineLen += wordLen
	}

	return result.String()
}

// Indent adds prefix to the beginning of each line in s.
//
// Example:
//
//	Indent("a\nb\nc", "  ")  // "  a\n  b\n  c"
func Indent(s, prefix string) string {
	lines := Lines(s)
	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}

// Dedent removes common leading whitespace from all lines.
//
// Example:
//
//	Dedent("  a\n  b\n  c")  // "a\nb\nc"
func Dedent(s string) string {
	lines := Lines(s)
	if len(lines) == 0 {
		return s
	}

	// Find minimum indentation (ignoring empty lines)
	minIndent := -1

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	// Remove minimum indent from each line
	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}

	return strings.Join(lines, "\n")
}

// StripTags removes HTML/XML tags from s.
// This is a simple implementation that may not handle all edge cases.
//
// Example:
//
//	StripTags("<p>Hello <b>World</b></p>")  // "Hello World"
func StripTags(s string) string {
	var result strings.Builder

	inTag := false

	for _, r := range s {
		if r == '<' {
			inTag = true

			continue
		}

		if r == '>' {
			inTag = false

			continue
		}

		if !inTag {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SplitN splits s by sep into at most n parts.
// If n <= 0, returns all parts (same as strings.Split).
// Wrapper around strings.SplitN for consistency.
func SplitN(s, sep string, n int) []string {
	if n <= 0 {
		return strings.Split(s, sep)
	}

	return strings.SplitN(s, sep, n)
}

// SplitAfter splits s after each instance of sep.
// Wrapper around strings.SplitAfter for consistency.
func SplitAfter(s, sep string) []string {
	return strings.SplitAfter(s, sep)
}

// Join concatenates elements with sep.
// Wrapper around strings.Join for API completeness.
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}
