package stringutil

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// Fixtures are deterministic, prepared before b.Loop, and immutable during
// measurement. No benchmark includes restoration work.
func benchTextCases(b *testing.B, sizes []int, fn func(*testing.B, int, benchkit.TextKind)) {
	b.Helper()
	for _, c := range []struct {
		name string
		kind benchkit.TextKind
	}{
		{"ascii", benchkit.ASCII},
		{"unicode", benchkit.Unicode},
	} {
		b.Run("case="+c.name, func(b *testing.B) {
			benchkit.Run(b, sizes, func(b *testing.B, n int) {
				fn(b, n, c.kind)
			})
		})
	}
}

// benchFitText preserves UTF-8 and pads with a one-byte character so n always
// describes the actual input length after a fixture transformation.
func benchFitText(text string, n int, pad string) string {
	cut := min(n, len(text))
	for cut > 0 && cut < len(text) && !utf8.RuneStart(text[cut]) {
		cut--
	}

	return text[:cut] + strings.Repeat(pad, n-cut)
}

// benchHTMLText adds one complete tag or entity every eight words. The final
// plain-text fragment fills the byte budget without cutting through markup.
func benchHTMLText(n int, kind benchkit.TextKind, entities bool) string {
	words := strings.Fields(benchkit.Text(n, kind))
	var text strings.Builder
	text.Grow(n)
	for i, word := range words {
		if i%8 == 0 {
			marker := "<br>"
			if entities && i%16 == 0 {
				marker = "&amp;"
			}
			if text.Len()+len(marker) > n {
				break
			}
			text.WriteString(marker)
		}
		if text.Len()+len(word)+1 > n {
			break
		}
		text.WriteString(word)
		text.WriteByte(' ')
	}
	text.WriteString(benchkit.Text(n-text.Len(), kind))

	return text.String()
}

func benchCaseText(n int, kind benchkit.TextKind) string {
	upper := false
	text := strings.Map(func(r rune) rune {
		upper = !upper
		if upper {
			return unicode.ToUpper(r)
		}

		return unicode.ToLower(r)
	}, benchkit.Text(n, kind))

	return benchFitText(text, n, " ")
}

// benchPair returns equal rune lengths with exactly one in four positions
// changed; the replacement is absent from both source alphabets.
func benchPair(n int, kind benchkit.TextKind) (string, string) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	if kind == benchkit.Unicode {
		alphabet = "éàüλΩЖ東京韓😀"
	}
	first := benchkit.Strings(1, n, alphabet)[0]
	second := []rune(first)
	for i := 3; i < len(second); i += 4 {
		second[i] = '#'
	}

	return first, string(second)
}

func benchPredicateText(n int, kind benchkit.TextKind, ascii, nonASCII, pad string) string {
	alphabet := []rune(ascii)
	if kind == benchkit.Unicode {
		alphabet = []rune(nonASCII)
	}
	text := strings.Map(func(r rune) rune {
		return alphabet[int(r)%len(alphabet)]
	}, benchkit.Text(n, kind))

	return benchFitText(text, n, pad)
}

func benchWordText(n int, sep string) string {
	return strings.Join(benchkit.Words(n), sep)[:n]
}

// benchBetweenText brackets groups of sixteen words. A shorter final group
// keeps both delimiters inside the exact n-byte budget.
func benchBetweenText(n int) string {
	words := strings.Fields(benchkit.Text(n, benchkit.ASCII))
	var text strings.Builder
	text.Grow(n)
	for start := 0; start < len(words) && text.Len()+2 < n; start += 16 {
		content := strings.Join(words[start:min(start+16, len(words))], " ")
		content = content[:min(len(content), n-text.Len()-2)]
		text.WriteByte('[')
		text.WriteString(content)
		text.WriteByte(']')
		if text.Len() < n {
			text.WriteByte(' ')
		}
	}
	text.WriteString(strings.Repeat(" ", n-text.Len()))

	return text.String()
}

// BenchmarkReverse reverses immutable ASCII or Unicode text; n is input bytes.
func BenchmarkReverse(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			Reverse(text)
		}
	})
}

// BenchmarkLevenshteinDistance measures equal-length strings differing at exactly
// one quarter of positions; n is runes per string (including the old long case).
func BenchmarkLevenshteinDistance(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			LevenshteinDistance(first, second)
		}
	})
}

// BenchmarkJaroWinklerSimilarity measures strings with a shared three-rune prefix
// and 25% changed positions, using prefixScale=0.1; n is runes per string.
func BenchmarkJaroWinklerSimilarity(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			JaroWinklerSimilarity(first, second, 0.1)
		}
	})
}

// BenchmarkSnakeCase converts mixed-case ASCII or Unicode text; n is input bytes.
func BenchmarkSnakeCase(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchCaseText(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			SnakeCase(text)
		}
	})
}

// BenchmarkAllIndexes finds all spaces in immutable ASCII text; n is input bytes.
func BenchmarkAllIndexes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.ASCII)
		b.ReportAllocs()
		for b.Loop() {
			AllIndexes(text, " ")
		}
	})
}

// BenchmarkIsPalindrome scans a true palindrome made by mirroring seeded text;
// n is input bytes. mode=literal compares runes directly; mode=normalized also
// lowercases and removes punctuation and whitespace. Both scan the full input.
func BenchmarkIsPalindrome(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		half := benchkit.Text(n/2, kind)
		text := half + Reverse(half)
		for _, c := range []struct {
			name      string
			normalize bool
		}{
			{"literal", false},
			{"normalized", true},
		} {
			b.Run("mode="+c.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if !IsPalindrome(text, c.normalize) {
						b.Fatal("expected a palindrome")
					}
				}
			})
		}
	})
}

// BenchmarkCleanString measures all options on immutable ASCII or Unicode HTML,
// with alternating entities/tags every eight words and a leading NUL; n is input
// bytes. Options are prepared outside the loop; the DB limit is half the input
// rune count. mode=all-options preserves the old composite benchmark.
func BenchmarkCleanString(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := "\x00" + benchHTMLText(n-1, kind, true)
		opts := []CleanOption{WithHTMLStrip(), WithUnicodeNorm(), WithDBSanitize(utf8.RuneCountInString(text) / 2)}
		b.Run("mode=all-options", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := CleanString(text, opts...); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkNormalizeUnicode normalizes immutable ASCII or Unicode text; n is input bytes.
func BenchmarkNormalizeUnicode(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := NormalizeUnicode(text); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSlugify slugifies immutable ASCII or Unicode text; n is input bytes.
func BenchmarkSlugify(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := Slugify(text); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkStripHTMLEntities strips tags and decodes entities from immutable text;
// n is input bytes, with a tag or entity every eight words.
func BenchmarkStripHTMLEntities(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchHTMLText(n, kind, true)
		b.ReportAllocs()
		for b.Loop() {
			StripHTMLEntities(text)
		}
	})
}

// BenchmarkSanitizeUTF8 validates or repairs immutable ASCII/Unicode text; n is
// input bytes. mode=valid preserves the old valid-input benchmark; mode=invalid
// replaces byte zero with NUL and the midpoint with 0xff before measurement.
func BenchmarkSanitizeUTF8(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		valid := benchkit.Text(n, kind)
		invalid := []byte(valid)
		invalid[0] = 0
		invalid[n/2] = 0xff
		for _, c := range []struct {
			name string
			text string
		}{
			{"valid", valid},
			{"invalid", string(invalid)},
		} {
			b.Run("mode="+c.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					SanitizeUTF8(c.text)
				}
			})
		}
	})
}

// BenchmarkNormalizeWhitespace collapses runs of spaces, tabs and newlines in
// immutable ASCII/Unicode text; n is the actual byte length after expanding
// spaces and fitting the fixture at a rune boundary.
func BenchmarkNormalizeWhitespace(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchFitText(strings.ReplaceAll(benchkit.Text(n, kind), " ", " \t\n "), n, " ")
		b.ReportAllocs()
		for b.Loop() {
			NormalizeWhitespace(text)
		}
	})
}

// BenchmarkWithUnicodeNorm measures the effect of
// the Unicode normalization option through CleanString on 1 KiB of immutable Mixed text.
// Option construction and application are measured together.
func BenchmarkWithUnicodeNorm(b *testing.B) {
	text := benchkit.Text(1024, benchkit.Mixed)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := CleanString(text, WithUnicodeNorm()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWithHTMLStrip measures the effect of
// the HTML option through CleanString on 1 KiB of immutable Mixed HTML, with a tag or entity every eight words.
// Option construction and application are measured together.
func BenchmarkWithHTMLStrip(b *testing.B) {
	text := benchHTMLText(1024, benchkit.Mixed, true)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := CleanString(text, WithHTMLStrip()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWithDBSanitize measures the effect of
// the DB option through CleanString on 1 KiB of immutable Mixed text with a leading NUL, truncating to 128 runes.
// Option construction and application are measured together.
func BenchmarkWithDBSanitize(b *testing.B) {
	text := "\x00" + benchkit.Text(1023, benchkit.Mixed)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := CleanString(text, WithDBSanitize(128)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRemoveAccents converts immutable ASCII or Unicode text; n is input bytes.
func BenchmarkRemoveAccents(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := RemoveAccents(text); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkToASCII converts immutable ASCII or Unicode text; n is input bytes.
func BenchmarkToASCII(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := ToASCII(text); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkTruncateRunes truncates immutable ASCII/Unicode text to half its
// rune count; n is input bytes.
func BenchmarkTruncateRunes(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		limit := utf8.RuneCountInString(text) / 2
		b.ReportAllocs()
		for b.Loop() {
			TruncateRunes(text, limit)
		}
	})
}

// BenchmarkRemoveNonPrintable removes bells from immutable ASCII/Unicode text;
// n is input bytes. The leading byte and all source spaces are bells.
func BenchmarkRemoveNonPrintable(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := "\x07" + strings.ReplaceAll(benchkit.Text(n-1, kind), " ", "\x07")
		b.ReportAllocs()
		for b.Loop() {
			RemoveNonPrintable(text)
		}
	})
}

// BenchmarkLevenshteinSimilarity measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkLevenshteinSimilarity(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			LevenshteinSimilarity(first, second)
		}
	})
}

// BenchmarkDamerauLevenshteinDistance measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkDamerauLevenshteinDistance(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			DamerauLevenshteinDistance(first, second)
		}
	})
}

// BenchmarkJaroSimilarity measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkJaroSimilarity(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			JaroSimilarity(first, second)
		}
	})
}

// BenchmarkDiceCoefficient measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkDiceCoefficient(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			DiceCoefficient(first, second)
		}
	})
}

// BenchmarkHammingDistance measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkHammingDistance(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if distance := HammingDistance(first, second); distance != n/4 {
				b.Fatalf("distance = %d, want %d", distance, n/4)
			}
		}
	})
}

// BenchmarkLongestCommonSubsequence measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkLongestCommonSubsequence(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			LongestCommonSubsequence(first, second)
		}
	})
}

// BenchmarkLongestCommonSubstring measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string.
func BenchmarkLongestCommonSubstring(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			LongestCommonSubstring(first, second)
		}
	})
}

// BenchmarkCosineSimilarity measures immutable equal-length strings with exactly
// one quarter of positions changed; n is runes per string. The n-gram width is 2.
func BenchmarkCosineSimilarity(b *testing.B) {
	benchTextCases(b, benchkit.SmallSizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		first, second := benchPair(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			CosineSimilarity(first, second, 2)
		}
	})
}

// BenchmarkHasAnyPrefix checks three affixes against immutable ASCII text;
// only the final candidate matches, so all candidates are visited. n is input bytes.
func BenchmarkHasAnyPrefix(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.ASCII)
		match := text[:3]
		b.ReportAllocs()
		for b.Loop() {
			if !HasAnyPrefix(text, "#", "@", match) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkHasAnySuffix checks three affixes against immutable ASCII text;
// only the final candidate matches, so all candidates are visited. n is input bytes.
func BenchmarkHasAnySuffix(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.ASCII)
		match := text[len(text)-3:]
		b.ReportAllocs()
		for b.Loop() {
			if !HasAnySuffix(text, "#", "@", match) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkContainsAny searches immutable ASCII text for three candidates,
// only the last present, at the final byte; n is input bytes. The true result
// requires scanning the input for every candidate.
func BenchmarkContainsAny(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n-1, benchkit.ASCII) + "#"
		b.ReportAllocs()
		for b.Loop() {
			if !ContainsAny(text, "!", "@", "#") {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkContainsAll searches immutable ASCII text for three markers, all
// present at the end; n is input bytes. The predicate holds after three scans.
func BenchmarkContainsAll(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n-3, benchkit.ASCII) + "#@!"
		b.ReportAllocs()
		for b.Loop() {
			if !ContainsAll(text, "#", "@", "!") {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkTruncate truncates immutable ASCII/Unicode text to half its
// rune count with a three-dot suffix; n is input bytes.
func BenchmarkTruncate(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		limit := utf8.RuneCountInString(text) / 2
		b.ReportAllocs()
		for b.Loop() {
			Truncate(text, limit, "...")
		}
	})
}

// BenchmarkTruncateWords truncates immutable ASCII/Unicode text to half its
// rune count with a three-dot suffix; n is input bytes.
func BenchmarkTruncateWords(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		limit := utf8.RuneCountInString(text) / 2
		b.ReportAllocs()
		for b.Loop() {
			TruncateWords(text, limit, "...")
		}
	})
}

// BenchmarkPadLeft pads immutable ASCII/Unicode text with spaces to twice
// its rune count; n is input bytes. All input sizes take the padding path.
func BenchmarkPadLeft(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		target := 2 * utf8.RuneCountInString(text)
		b.ReportAllocs()
		for b.Loop() {
			PadLeft(text, target, ' ')
		}
	})
}

// BenchmarkPadRight pads immutable ASCII/Unicode text with spaces to twice
// its rune count; n is input bytes. All input sizes take the padding path.
func BenchmarkPadRight(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		target := 2 * utf8.RuneCountInString(text)
		b.ReportAllocs()
		for b.Loop() {
			PadRight(text, target, ' ')
		}
	})
}

// BenchmarkPadCenter pads immutable ASCII/Unicode text with spaces to twice
// its rune count; n is input bytes. All input sizes take the padding path.
func BenchmarkPadCenter(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		target := 2 * utf8.RuneCountInString(text)
		b.ReportAllocs()
		for b.Loop() {
			PadCenter(text, target, ' ')
		}
	})
}

// BenchmarkRemoveAll removes spaces, full stops and a from immutable ASCII text;
// n is input bytes.
func BenchmarkRemoveAll(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.ASCII)
		b.ReportAllocs()
		for b.Loop() {
			RemoveAll(text, " ", ".", "a")
		}
	})
}

// BenchmarkCountLines counts lines of seeded words separated by LF;
// n is input bytes and the fixture is immutable.
func BenchmarkCountLines(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchWordText(n, "\n")
		b.ReportAllocs()
		for b.Loop() {
			CountLines(text)
		}
	})
}

// BenchmarkLines splits immutable seeded words separated by LF or CRLF;
// n is input bytes, including separators.
func BenchmarkLines(b *testing.B) {
	for _, c := range []struct {
		name string
		sep  string
	}{
		{"lf", "\n"},
		{"crlf", "\r\n"},
	} {
		b.Run("case="+c.name, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := benchWordText(n, c.sep)
				b.ReportAllocs()
				for b.Loop() {
					Lines(text)
				}
			})
		})
	}
}

// BenchmarkIsBlank scans immutable ASCII/Unicode whitespace; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsBlank(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, " \t\n", "\u2003\u00a0\u3000", " ")
		b.ReportAllocs()
		for b.Loop() {
			if !IsBlank(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsAlpha scans immutable ASCII/Unicode letters; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsAlpha(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, "abcdefghijklmnopqrstuvwxyz", "αβγδεζηθéàü", "a")
		b.ReportAllocs()
		for b.Loop() {
			if !IsAlpha(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsAlphanumeric scans immutable ASCII/Unicode letters and digits; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsAlphanumeric(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, "abcxyz0123456789", "αβéü٠١٢٣٤٥٦٧٨٩", "a")
		b.ReportAllocs()
		for b.Loop() {
			if !IsAlphanumeric(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsNumeric scans immutable ASCII/Unicode digits; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsNumeric(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, "0123456789", "٠١٢٣٤٥٦٧٨٩", "0")
		b.ReportAllocs()
		for b.Loop() {
			if !IsNumeric(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsUpper scans immutable ASCII/Unicode uppercase letters; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsUpper(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞ", "A")
		b.ReportAllocs()
		for b.Loop() {
			if !IsUpper(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsLower scans immutable ASCII/Unicode lowercase letters; n is input
// bytes. Seeded text is mapped to valid characters, then fitted at a rune
// boundary. The predicate is true, forcing a full scan.
func BenchmarkIsLower(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchPredicateText(n, kind, "abcdefghijklmnopqrstuvwxyz", "αβγδεζηθικλμνξ", "a")
		b.ReportAllocs()
		for b.Loop() {
			if !IsLower(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkIsEmpty checks the empty-string domain, the only input satisfying
// the predicate. This is a constant-time length check with no size sweep.
func BenchmarkIsEmpty(b *testing.B) {
	b.Run("case=empty", func(b *testing.B) {
		text := benchkit.Text(0, benchkit.ASCII)
		b.ReportAllocs()
		for b.Loop() {
			if !IsEmpty(text) {
				b.Fatal("expected empty text")
			}
		}
	})
}

// BenchmarkIsASCII scans immutable ASCII text; n is input bytes. Every byte
// satisfies the predicate, forcing a full scan; Unicode is not a true domain.
func BenchmarkIsASCII(b *testing.B) {
	b.Run("case=ascii", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			text := benchkit.Text(n, benchkit.ASCII)
			b.ReportAllocs()
			for b.Loop() {
				if !IsASCII(text) {
					b.Fatal("expected ASCII text")
				}
			}
		})
	})
}

// BenchmarkIsPrintable scans immutable printable ASCII/Unicode text; n is
// input bytes. Every rune satisfies the predicate, forcing a full scan.
func BenchmarkIsPrintable(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			if !IsPrintable(text) {
				b.Fatal("expected predicate to hold")
			}
		}
	})
}

// BenchmarkCapitalize converts immutable mixed-case ASCII/Unicode text;
// n is input bytes after fitting the case-mapped fixture at a rune boundary.
func BenchmarkCapitalize(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchCaseText(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			Capitalize(text)
		}
	})
}

// BenchmarkSwapCase converts immutable mixed-case ASCII/Unicode text;
// n is input bytes after fitting the case-mapped fixture at a rune boundary.
func BenchmarkSwapCase(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchCaseText(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			SwapCase(text)
		}
	})
}

// BenchmarkKebabCase converts immutable mixed-case ASCII/Unicode text;
// n is input bytes after fitting the case-mapped fixture at a rune boundary.
func BenchmarkKebabCase(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchCaseText(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			KebabCase(text)
		}
	})
}

// BenchmarkTitle title-cases immutable ASCII/Unicode words; n is input bytes.
func BenchmarkTitle(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			Title(text)
		}
	})
}

// BenchmarkCamelCase converts immutable ASCII/Unicode text with underscore
// and hyphen separators; n is input bytes.
func BenchmarkCamelCase(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := strings.NewReplacer(" ", "_", ".", "-").Replace(benchkit.Text(n, kind))
		b.ReportAllocs()
		for b.Loop() {
			CamelCase(text)
		}
	})
}

// BenchmarkPascalCase converts immutable ASCII/Unicode text with underscore
// and hyphen separators; n is input bytes.
func BenchmarkPascalCase(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := strings.NewReplacer(" ", "_", ".", "-").Replace(benchkit.Text(n, kind))
		b.ReportAllocs()
		for b.Loop() {
			PascalCase(text)
		}
	})
}

// BenchmarkWords extracts words from immutable ASCII/Unicode text;
// n is input bytes.
func BenchmarkWords(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			Words(text)
		}
	})
}

// BenchmarkRuneCount counts runes in immutable ASCII/Unicode text;
// n is input bytes.
func BenchmarkRuneCount(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			RuneCount(text)
		}
	})
}

// BenchmarkSafeSlice extracts the middle half of immutable ASCII/Unicode
// text using valid rune indices; n is input bytes.
func BenchmarkSafeSlice(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		count := utf8.RuneCountInString(text)
		start, end := count/4, 3*count/4
		b.ReportAllocs()
		for b.Loop() {
			SafeSlice(text, start, end)
		}
	})
}

// BenchmarkNthRune retrieves the final rune from immutable ASCII/Unicode
// text, forcing a full scan; n is input bytes and the index is valid.
func BenchmarkNthRune(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		index := utf8.RuneCountInString(text) - 1
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := NthRune(text, index); !ok {
				b.Fatal("expected a rune at the valid index")
			}
		}
	})
}

// BenchmarkCommonPrefix finds the shared prefix of three immutable strings;
// n is bytes per string, of which n-1 are a shared ASCII/Unicode prefix.
func BenchmarkCommonPrefix(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		shared := benchkit.Text(n-1, kind)
		first, second, third := shared+"a", shared+"b", shared+"c"
		b.ReportAllocs()
		for b.Loop() {
			CommonPrefix(first, second, third)
		}
	})
}

// BenchmarkCommonSuffix finds the shared suffix of three immutable strings;
// n is bytes per string, of which n-1 are a shared ASCII/Unicode suffix.
func BenchmarkCommonSuffix(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		shared := benchkit.Text(n-1, kind)
		first, second, third := "a"+shared, "b"+shared, "c"+shared
		b.ReportAllocs()
		for b.Loop() {
			CommonSuffix(first, second, third)
		}
	})
}

// BenchmarkRepeat repeats an immutable 16-byte ASCII seed; n is repetitions,
// so each independent call produces 16*n bytes.
func BenchmarkRepeat(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(16, benchkit.ASCII)
		b.ReportAllocs()
		for b.Loop() {
			Repeat(text, n)
		}
	})
}

// BenchmarkBetween finds the first bracketed group in immutable ASCII text
// with delimiters every sixteen words; n is input bytes, including delimiters.
// Both delimiters are present even in the smallest fixture.
func BenchmarkBetween(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchBetweenText(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := Between(text, "[", "]"); !ok {
				b.Fatal("expected a delimited group")
			}
		}
	})
}

// BenchmarkBetweenAll finds all bracketed groups in immutable ASCII text
// with delimiters every sixteen words; n is input bytes, including delimiters.
func BenchmarkBetweenAll(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchBetweenText(n)
		b.ReportAllocs()
		for b.Loop() {
			if groups := BetweenAll(text, "[", "]"); len(groups) == 0 {
				b.Fatal("expected delimited groups")
			}
		}
	})
}

// BenchmarkWrap wraps immutable ASCII/Unicode text at width 80; n is input bytes.
func BenchmarkWrap(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchkit.Text(n, kind)
		b.ReportAllocs()
		for b.Loop() {
			Wrap(text, 80)
		}
	})
}

// BenchmarkIndent adds four spaces per line of immutable seeded words;
// n is input bytes before indentation.
func BenchmarkIndent(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchWordText(n, "\n")
		b.ReportAllocs()
		for b.Loop() {
			Indent(text, "    ")
		}
	})
}

// BenchmarkDedent removes four leading spaces per line of immutable seeded
// words; n is input bytes including the indentation.
func BenchmarkDedent(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := "    " + benchWordText(n-4, "\n    ")
		b.ReportAllocs()
		for b.Loop() {
			Dedent(text)
		}
	})
}

// BenchmarkStripTags removes a complete tag every eight words from immutable
// ASCII/Unicode text; n is input bytes including tags.
func BenchmarkStripTags(b *testing.B) {
	benchTextCases(b, benchkit.Sizes, func(b *testing.B, n int, kind benchkit.TextKind) {
		text := benchHTMLText(n, kind, false)
		b.ReportAllocs()
		for b.Loop() {
			StripTags(text)
		}
	})
}

// BenchmarkSplitN splits immutable comma-separated seeded words; n is input
// bytes. mode=limited returns at most eight parts; mode=all uses limit zero.
func BenchmarkSplitN(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchWordText(n, ",")
		for _, c := range []struct {
			name  string
			limit int
		}{
			{"limited", 8},
			{"all", 0},
		} {
			b.Run("mode="+c.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					SplitN(text, ",", c.limit)
				}
			})
		}
	})
}

// BenchmarkSplitAfter splits immutable comma-separated seeded words while
// retaining separators; n is input bytes.
func BenchmarkSplitAfter(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchWordText(n, ",")
		b.ReportAllocs()
		for b.Loop() {
			SplitAfter(text, ",")
		}
	})
}

// BenchmarkJoin joins immutable seeded words with commas; n is input elements.
func BenchmarkJoin(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		words := benchkit.Words(n)
		b.ReportAllocs()
		for b.Loop() {
			Join(words, ",")
		}
	})
}

// BenchmarkSplitAndTrim splits immutable comma-separated seeded words with
// whitespace around tokens and every eighth token blank; n is input bytes.
// The prepared fixture is cropped to n bytes before measurement.
func BenchmarkSplitAndTrim(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		words := benchkit.Words(n)
		for i := 7; i < len(words); i += 8 {
			words[i] = ""
		}
		text := (" " + strings.Join(words, " , ") + " ")[:n]
		b.ReportAllocs()
		for b.Loop() {
			SplitAndTrim(text, ",")
		}
	})
}
