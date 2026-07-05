---
layout: default
title: stringutil
parent: Packages
nav_order: 2
---

# stringutil

String manipulation utilities and similarity algorithms.

```go
import "github.com/alessiosavi/GoGPUtils/stringutil"
```

The `stringutil` package provides a comprehensive set of string manipulation utilities organized into several categories: search and indexing, transformation, validation, similarity algorithms, and string cleaning. All functions are designed to be nil-safe and handle edge cases gracefully, with full Unicode (UTF-8) support.

---

## Table of Contents

- [Search and Indexing](#search-and-indexing)
- [Transformation](#transformation)
- [Validation](#validation)
- [Similarity Algorithms](#similarity-algorithms)
- [String Cleaning](#string-cleaning)
- [Usage Examples](#usage-examples)

---

## Search and Indexing

Functions for finding substrings, prefixes, suffixes, and extracting content between markers.

### AllIndexes

```go
func AllIndexes(s, substr string) []int
```

Returns all starting positions of `substr` in `s`. Returns `nil` if `substr` is empty or `s` doesn't contain `substr`.

**Example:**

```go
indices := stringutil.AllIndexes("banana", "an")
// indices = [1, 3]
```

### HasAnyPrefix

```go
func HasAnyPrefix(s string, prefixes ...string) bool
```

Reports whether `s` starts with any of the given prefixes.

**Example:**

```go
if stringutil.HasAnyPrefix(url, "http://", "https://") {
    // URL uses HTTP or HTTPS
}
```

### HasAnySuffix

```go
func HasAnySuffix(s string, suffixes ...string) bool
```

Reports whether `s` ends with any of the given suffixes.

### ContainsAny

```go
func ContainsAny(s string, substrs ...string) bool
```

Reports whether `s` contains any of the given substrings.

**Example:**

```go
if stringutil.ContainsAny(text, "error", "fail", "warning") {
    // Log contains an error indicator
}
```

### ContainsAll

```go
func ContainsAll(s string, substrs ...string) bool
```

Reports whether `s` contains all of the given substrings.

### Between

```go
func Between(s, start, end string) (string, bool)
```

Extracts the substring between `start` and `end` markers. Returns empty string and `false` if markers are not found in proper order.

**Example:**

```go
content, ok := stringutil.Between("[hello]", "[", "]")
// content = "hello", ok = true
```

### BetweenAll

```go
func BetweenAll(s, start, end string) []string
```

Extracts all substrings between `start` and `end` markers.

**Example:**

```go
results := stringutil.BetweenAll("a[1]b[2]c[3]", "[", "]")
// results = ["1", "2", "3"]
```

### CommonPrefix

```go
func CommonPrefix(strs ...string) string
```

Returns the longest common prefix of the given strings. Returns empty string if no common prefix or fewer than 2 strings.

**Example:**

```go
prefix := stringutil.CommonPrefix("interstellar", "internet", "internal")
// prefix = "inter"
```

### CommonSuffix

```go
func CommonSuffix(strs ...string) string
```

Returns the longest common suffix of the given strings.

---

## Transformation

Functions for transforming strings: reversing, truncating, padding, case conversion, and more.

### Reverse

```go
func Reverse(s string) string
```

Returns `s` with its characters in reverse order. Correctly handles multi-byte UTF-8 characters.

**Example:**

```go
rev := stringutil.Reverse("hello")  // "olleh"
rev := stringutil.Reverse("日本語")   // "語本日"
```

### Truncate

```go
func Truncate(s string, maxLen int, suffix string) string
```

Shortens `s` to `maxLen` characters, appending `suffix` if truncated. The total length including suffix will not exceed `maxLen`.

**Example:**

```go
truncated := stringutil.Truncate("Hello World", 8, "...")
// truncated = "Hello..."
```

### TruncateWords

```go
func TruncateWords(s string, maxLen int, suffix string) string
```

Truncates `s` at a word boundary, appending `suffix` if truncated. Attempts to break at word boundaries rather than mid-word.

### PadLeft

```go
func PadLeft(s string, length int, padChar rune) string
```

Pads `s` on the left with `padChar` to reach the target length. If `s` is already >= length, returns `s` unchanged.

**Example:**

```go
padded := stringutil.PadLeft("42", 5, '0')
// padded = "00042"
```

### PadRight

```go
func PadRight(s string, length int, padChar rune) string
```

Pads `s` on the right with `padChar` to reach the target length.

**Example:**

```go
padded := stringutil.PadRight("42", 5, '0')
// padded = "42000"
```

### PadCenter

```go
func PadCenter(s string, length int, padChar rune) string
```

Centers `s` by adding `padChar` on both sides. If odd padding is needed, the extra character goes on the right.

**Example:**

```go
centered := stringutil.PadCenter("hello", 11, '*')
// centered = "***hello***"
```

### RemoveAll

```go
func RemoveAll(s string, substrs ...string) string
```

Removes all occurrences of the given substrings from `s`.

**Example:**

```go
clean := stringutil.RemoveAll("hello world", "l", "o")
// clean = "he wrd"
```

### Capitalize

```go
func Capitalize(s string) string
```

Returns `s` with the first character uppercased and the rest lowercased.

**Example:**

```go
result := stringutil.Capitalize("hELLO")
// result = "Hello"
```

### Title

```go
func Title(s string) string
```

Returns `s` with the first character of each word uppercased.

**Example:**

```go
result := stringutil.Title("hello world")
// result = "Hello World"
```

### SwapCase

```go
func SwapCase(s string) string
```

Swaps the case of each letter in `s`.

**Example:**

```go
result := stringutil.SwapCase("Hello World")
// result = "hELLO wORLD"
```

### SnakeCase

```go
func SnakeCase(s string) string
```

Converts `s` to `snake_case`.

**Example:**

```go
result := stringutil.SnakeCase("HelloWorld")
// result = "hello_world"
```

### CamelCase

```go
func CamelCase(s string) string
```

Converts `s` to `camelCase`.

**Example:**

```go
result := stringutil.CamelCase("hello_world")
// result = "helloWorld"
```

### PascalCase

```go
func PascalCase(s string) string
```

Converts `s` to `PascalCase`.

**Example:**

```go
result := stringutil.PascalCase("hello_world")
// result = "HelloWorld"
```

### KebabCase

```go
func KebabCase(s string) string
```

Converts `s` to `kebab-case`.

**Example:**

```go
result := stringutil.KebabCase("HelloWorld")
// result = "hello-world"
```

### Words

```go
func Words(s string) []string
```

Splits `s` into words, treating any non-alphanumeric character as a separator.

**Example:**

```go
words := stringutil.Words("hello, world!")
// words = ["hello", "world"]
```

### Lines

```go
func Lines(s string) []string
```

Splits `s` into lines. Unlike `strings.Split`, handles `\r\n` properly.

**Example:**

```go
lines := stringutil.Lines("a\nb\nc")
// lines = ["a", "b", "c"]
```

### CountLines

```go
func CountLines(s string) int
```

Returns the number of lines in `s`. An empty string returns 0; a string without newlines returns 1.

### Wrap

```go
func Wrap(s string, width int) string
```

Wraps text at the specified width, breaking at word boundaries. Preserves existing line breaks.

### Indent

```go
func Indent(s, prefix string) string
```

Adds `prefix` to the beginning of each line in `s`.

**Example:**

```go
result := stringutil.Indent("a\nb\nc", "  ")
// result = "  a\n  b\n  c"
```

### Dedent

```go
func Dedent(s string) string
```

Removes common leading whitespace from all lines.

**Example:**

```go
result := stringutil.Dedent("  a\n  b\n  c")
// result = "a\nb\nc"
```

### StripTags

```go
func StripTags(s string) string
```

Removes HTML/XML tags from `s`. This is a simple implementation that may not handle all edge cases.

**Example:**

```go
result := stringutil.StripTags("<p>Hello <b>World</b></p>")
// result = "Hello World"
```

### SplitN

```go
func SplitN(s, sep string, n int) []string
```

Splits `s` by `sep` into at most `n` parts. If `n <= 0`, returns all parts (same as `strings.Split`).

### SplitAfter

```go
func SplitAfter(s, sep string) []string
```

Splits `s` after each instance of `sep`.

### Join

```go
func Join(elems []string, sep string) string
```

Concatenates elements with `sep`.

### SplitAndTrim

```go
func SplitAndTrim(s, sep string) []string
```

Splits `s` by `sep`, trims whitespace from each token, and drops any tokens that are empty after trimming.

**Example:**

```go
result := stringutil.SplitAndTrim("  a , b ,  ,  c  ", ",")
// result = ["a", "b", "c"]
```

### Repeat

```go
func Repeat(s string, n int) string
```

Returns `s` repeated `n` times. If `n <= 0`, returns empty string.

**Example:**

```go
result := stringutil.Repeat("ab", 3)
// result = "ababab"
```

---

## Validation

Functions for checking string properties.

### IsEmpty

```go
func IsEmpty(s string) bool
```

Reports whether `s` is empty (zero length).

### IsBlank

```go
func IsBlank(s string) bool
```

Reports whether `s` contains only whitespace characters.

### IsAlpha

```go
func IsAlpha(s string) bool
```

Reports whether `s` contains only alphabetic characters.

### IsAlphanumeric

```go
func IsAlphanumeric(s string) bool
```

Reports whether `s` contains only letters and digits.

### IsNumeric

```go
func IsNumeric(s string) bool
```

Reports whether `s` contains only numeric digits.

### IsUpper

```go
func IsUpper(s string) bool
```

Reports whether all letters in `s` are uppercase. Returns `true` for strings with no letters.

### IsLower

```go
func IsLower(s string) bool
```

Reports whether all letters in `s` are lowercase. Returns `true` for strings with no letters.

### IsASCII

```go
func IsASCII(s string) bool
```

Reports whether `s` contains only ASCII characters.

### IsPrintable

```go
func IsPrintable(s string) bool
```

Reports whether `s` contains only printable characters.

### IsPalindrome

```go
func IsPalindrome(s string, normalize bool) bool
```

Reports whether `s` reads the same forwards and backwards. Case-sensitive; ignores whitespace and punctuation only if `normalize` is `true`.

**Example:**

```go
stringutil.IsPalindrome("racecar", false)                    // true
stringutil.IsPalindrome("A man a plan a canal Panama", true) // true (normalized)
```

---

## Unicode Utilities

### RuneCount

```go
func RuneCount(s string) int
```

Returns the number of runes (Unicode code points) in `s`. This differs from `len(s)`, which returns bytes.

**Example:**

```go
stringutil.RuneCount("日本語")  // 3
len("日本語")                  // 9 (bytes)
```

### SafeSlice

```go
func SafeSlice(s string, start, end int) string
```

Safely slices `s` by rune indices, returning an empty string for invalid ranges.

### NthRune

```go
func NthRune(s string, n int) (rune, bool)
```

Returns the rune at rune position `n` (0-indexed). Returns `(0, false)` if `n` is out of bounds.

---

## Similarity Algorithms

The `stringutil` package includes a comprehensive set of string similarity and distance algorithms. These are useful for fuzzy matching, spell checking, deduplication, and record linkage.

### Levenshtein Distance

```go
func LevenshteinDistance(s1, s2 string) int
```

Returns the minimum number of single-character edits (insertions, deletions, substitutions) required to change `s1` into `s2`.

| Aspect | Complexity               |
| ------ | ------------------------ |
| Time   | O(len(s1) × len(s2))     |
| Space  | O(min(len(s1), len(s2))) |

**Example:**

```go
distance := stringutil.LevenshteinDistance("kitten", "sitting")
// distance = 3
```

### Levenshtein Similarity

```go
func LevenshteinSimilarity(s1, s2 string) float64
```

Returns a similarity score between 0 and 1 based on Levenshtein distance. 1 means identical strings.

**Example:**

```go
score := stringutil.LevenshteinSimilarity("hello", "hallo")
// score ≈ 0.8
```

### Damerau-Levenshtein Distance

```go
func DamerauLevenshteinDistance(s1, s2 string) int
```

Extends Levenshtein to include transpositions (swapping two adjacent characters) as a single edit operation.

| Aspect | Complexity           |
| ------ | -------------------- |
| Time   | O(len(s1) × len(s2)) |
| Space  | O(len(s1) × len(s2)) |

**Example:**

```go
distance := stringutil.DamerauLevenshteinDistance("ca", "ac")
// distance = 1 (transposition)
// LevenshteinDistance("ca", "ac") would return 2
```

### Jaro Similarity

```go
func JaroSimilarity(s1, s2 string) float64
```

Returns the Jaro similarity between two strings. Returns a value between 0 (completely different) and 1 (identical). The algorithm considers the number of matching characters and transpositions.

| Aspect | Complexity           |
| ------ | -------------------- |
| Time   | O(len(s1) × len(s2)) |
| Space  | O(len(s1) + len(s2)) |

**Example:**

```go
score := stringutil.JaroSimilarity("martha", "marhta")
// score ≈ 0.944
```

### Jaro-Winkler Similarity

```go
func JaroWinklerSimilarity(s1, s2 string, prefixScale float64) float64
```

Returns the Jaro-Winkler similarity between two strings. This is an extension of Jaro that gives more weight to strings with a common prefix. The `prefixScale` parameter (0 to 0.25) determines how much weight to give to the common prefix. Standard value is 0.1.

| Aspect | Complexity           |
| ------ | -------------------- |
| Time   | O(len(s1) × len(s2)) |
| Space  | O(len(s1) + len(s2)) |

**Example:**

```go
score := stringutil.JaroWinklerSimilarity("martha", "marhta", 0.1)
// score ≈ 0.961
```

### Dice Coefficient

```go
func DiceCoefficient(s1, s2 string) float64
```

Returns the Sørensen–Dice coefficient comparing bigrams. Returns a value between 0 and 1, where 1 means identical sets of bigrams. This metric is useful for comparing short strings or when order matters less.

| Aspect | Complexity           |
| ------ | -------------------- |
| Time   | O(len(s1) + len(s2)) |
| Space  | O(len(s1) + len(s2)) |

**Example:**

```go
coefficient := stringutil.DiceCoefficient("night", "nacht")
// coefficient ≈ 0.25
```

### Hamming Distance

```go
func HammingDistance(s1, s2 string) int
```

Returns the number of positions where corresponding characters differ. Only defined for strings of equal length. Returns -1 if strings have different lengths.

| Aspect | Complexity               |
| ------ | ------------------------ |
| Time   | O(min(len(s1), len(s2))) |
| Space  | O(1)                     |

**Example:**

```go
distance := stringutil.HammingDistance("karolin", "kathrin")
// distance = 3
```

### Longest Common Subsequence

```go
func LongestCommonSubsequence(s1, s2 string) int
```

Returns the length of the longest common subsequence. A subsequence is a sequence that can be derived by deleting some elements without changing the order of remaining elements.

| Aspect | Complexity               |
| ------ | ------------------------ |
| Time   | O(len(s1) × len(s2))     |
| Space  | O(min(len(s1), len(s2))) |

**Example:**

```go
length := stringutil.LongestCommonSubsequence("ABCDGH", "AEDFHR")
// length = 3 (subsequence: "ADH")
```

### Longest Common Substring

```go
func LongestCommonSubstring(s1, s2 string) string
```

Returns the longest common contiguous substring.

| Aspect | Complexity               |
| ------ | ------------------------ |
| Time   | O(len(s1) × len(s2))     |
| Space  | O(min(len(s1), len(s2))) |

**Example:**

```go
substring := stringutil.LongestCommonSubstring("ABABC", "BABCA")
// substring = "BABC"
```

### Cosine Similarity

```go
func CosineSimilarity(s1, s2 string, n int) float64
```

Computes the cosine similarity of two strings based on their character n-gram vectors. Returns a value between 0 and 1. This is useful for comparing longer texts.

| Aspect | Complexity           |
| ------ | -------------------- |
| Time   | O(len(s1) + len(s2)) |
| Space  | O(len(s1) + len(s2)) |

**Example:**

```go
score := stringutil.CosineSimilarity("hello world", "hello there", 2)
// score ≈ 0.5
```

---

## String Cleaning

The `stringutil` package provides a composable API for cleaning and normalizing strings, particularly useful for preparing text for databases, search indexes, or URL slugs.

### CleanString

```go
func CleanString(input string, options ...CleanOption) (string, error)
```

Applies the specified cleaning options to the input string. Options are applied in a fixed order for consistency and correctness:

1. HTML stripping (remove tags, decode entities)
2. Unicode normalization (NFKD + diacritic removal)
3. Database sanitization (UTF-8 validation, NUL removal, truncation)

If no options are provided, the input string is returned unchanged.

**Example:**

```go
// Apply all three cleaning modes
result, err := stringutil.CleanString(
    "<p>Héllo &amp; Wörld</p>",
    stringutil.WithHTMLStrip(),
    stringutil.WithUnicodeNorm(),
    stringutil.WithDBSanitize(50),
)
// result = "Hello & World"

// Apply only Unicode normalization
result, _ := stringutil.CleanString("café", stringutil.WithUnicodeNorm())
// result = "cafe"
```

### WithHTMLStrip

```go
func WithHTMLStrip() CleanOption
```

Enables HTML tag removal and entity decoding. All HTML/XML tags are stripped, and HTML entities are decoded to their Unicode equivalents.

**Example:**

```go
result, _ := stringutil.CleanString("<p>Hello &amp; World</p>", stringutil.WithHTMLStrip())
// result = "Hello & World"
```

### WithUnicodeNorm

```go
func WithUnicodeNorm() CleanOption
```

Enables Unicode normalization: NFKD decomposition followed by removal of combining marks (diacritics). Converts characters like "é" → "e", "ñ" → "n", "ü" → "u".

**Example:**

```go
result, _ := stringutil.CleanString("café résumé", stringutil.WithUnicodeNorm())
// result = "cafe resume"
```

### WithDBSanitize

```go
func WithDBSanitize(maxLen int) CleanOption
```

Enables database sanitization:

- Replaces invalid UTF-8 sequences with U+FFFD (replacement character)
- Replaces NUL bytes (`\x00`) with empty string
- Optionally truncates to `maxLen` runes (0 = no truncation)

**Example:**

```go
result, _ := stringutil.CleanString("Hello\x00World", stringutil.WithDBSanitize(0))
// result = "HelloWorld"

result, _ := stringutil.CleanString("Hello World", stringutil.WithDBSanitize(5))
// result = "Hello"
```

### NormalizeUnicode

```go
func NormalizeUnicode(s string) (string, error)
```

Applies NFKD normalization and removes combining marks (diacritics) from the input string.

**Example:**

```go
result, _ := stringutil.NormalizeUnicode("café résumé")
// result = "cafe resume"

result, _ := stringutil.NormalizeUnicode("naïve")
// result = "naive"
```

### StripHTMLEntities

```go
func StripHTMLEntities(s string) string
```

Removes all HTML/XML tags and decodes HTML entities.

**Example:**

```go
result := stringutil.StripHTMLEntities("<p>Hello &amp; World</p>")
// result = "Hello & World"

result := stringutil.StripHTMLEntities("Price: &euro;10")
// result = "Price: €10"
```

### SanitizeUTF8

```go
func SanitizeUTF8(s string) string
```

Ensures the string contains only valid UTF-8 and removes NUL bytes. Invalid UTF-8 byte sequences are replaced with U+FFFD.

**Example:**

```go
result := stringutil.SanitizeUTF8("Hello\x00World")
// result = "HelloWorld"

result := stringutil.SanitizeUTF8("Hello\xffWorld")
// result = "Hello\uFFFDWorld"
```

### TruncateRunes

```go
func TruncateRunes(s string, maxLen int) string
```

Truncates `s` to at most `maxLen` runes. Unlike byte-level truncation, this is Unicode-safe and will never split a multi-byte character.

**Example:**

```go
result := stringutil.TruncateRunes("Hello, 世界!", 8)
// result = "Hello, 世界" (correct, not "Hello, \xe4")
```

### RemoveNonPrintable

```go
func RemoveNonPrintable(s string) string
```

Removes all non-printable characters from `s`, except for common whitespace (space, tab, newline, carriage return).

**Example:**

```go
result := stringutil.RemoveNonPrintable("Hello\x07World")
// result = "HelloWorld" (bell character removed)
```

### NormalizeWhitespace

```go
func NormalizeWhitespace(s string) string
```

Collapses all consecutive whitespace characters into a single space, and trims leading/trailing whitespace.

**Example:**

```go
result := stringutil.NormalizeWhitespace("  Hello   World  \n\t ")
// result = "Hello World"
```

### RemoveAccents

```go
func RemoveAccents(s string) (string, error)
```

Alias for `NormalizeUnicode` that removes diacritical marks from characters.

**Example:**

```go
result, _ := stringutil.RemoveAccents("café")
// result = "cafe"
```

### ToASCII

```go
func ToASCII(s string) (string, error)
```

Converts a Unicode string to its closest ASCII representation by removing diacritics, replacing non-letter/non-number characters with spaces, and collapsing whitespace.

**Example:**

```go
result, _ := stringutil.ToASCII("Héllo, Wörld!")
// result = "Hello  World"
```

### Slugify

```go
func Slugify(s string) (string, error)
```

Converts a string to a URL-friendly slug. Normalizes Unicode, lowercases, replaces non-alphanumeric characters with hyphens, collapses multiple hyphens, and trims leading/trailing hyphens.

**Example:**

```go
result, _ := stringutil.Slugify("Hello, World!")
// result = "hello-world"

result, _ := stringutil.Slugify("Café Résumé")
// result = "cafe-resume"
```

---

## Usage Examples

### Case Conversions

```go
stringutil.SnakeCase("HelloWorld")     // "hello_world"
stringutil.CamelCase("hello_world")    // "helloWorld"
stringutil.PascalCase("hello_world")   // "HelloWorld"
stringutil.KebabCase("HelloWorld")     // "hello-world"
```

### Padding and Truncation

```go
stringutil.PadLeft("42", 5, '0')       // "00042"
stringutil.PadRight("Go", 5, '-')      // "Go---"
stringutil.PadCenter("hi", 6, '*')     // "**hi**"
stringutil.Truncate("Hello World", 8, "...")  // "Hello..."
```

### Similarity Scoring

```go
// Edit distance
stringutil.LevenshteinDistance("kitten", "sitting")     // 3
stringutil.DamerauLevenshteinDistance("ca", "ac")       // 1

// Similarity scores (0.0 to 1.0)
stringutil.JaroSimilarity("martha", "marhta")           // ~0.944
stringutil.JaroWinklerSimilarity("martha", "marhta", 0.1) // ~0.961
stringutil.DiceCoefficient("night", "nacht")            // ~0.25
stringutil.LevenshteinSimilarity("hello", "hallo")      // ~0.8
```

### String Cleaning Pipeline

```go
// Clean HTML, normalize Unicode, and sanitize for database storage
result, err := stringutil.CleanString(
    "<div class='main'>Café &amp; Résumé</div>",
    stringutil.WithHTMLStrip(),
    stringutil.WithUnicodeNorm(),
    stringutil.WithDBSanitize(100),
)
// result = "Cafe & Resume"

// Generate a URL-friendly slug
slug, err := stringutil.Slugify("Héllo Wörld 2024!")
// slug = "hello-world-2024"

// Remove accents for search indexing
searchable, err := stringutil.RemoveAccents("café résumé naïve")
// searchable = "cafe resume naive"
```

### Working with Unicode

```go
// Count runes, not bytes
stringutil.RuneCount("日本語")  // 3
len("日本語")                  // 9

// Safe slicing by rune index
stringutil.SafeSlice("日本語", 0, 2)  // "日本"

// Get nth rune
r, ok := stringutil.NthRune("hello", 1)
// r = 'e', ok = true

// Reverse Unicode strings correctly
stringutil.Reverse("日本語")  // "語本日"
```

### Validation

```go
stringutil.IsEmpty("")              // true
stringutil.IsBlank("   ")           // true
stringutil.IsAlpha("hello")         // true
stringutil.IsNumeric("12345")       // true
stringutil.IsAlphanumeric("abc123") // true
stringutil.IsASCII("hello")         // true
stringutil.IsASCII("héllo")         // false
stringutil.IsPalindrome("racecar", false)  // true
```
