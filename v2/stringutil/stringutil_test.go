package stringutil

import (
	"slices"
	"testing"
)

// ============================================================================
// AllIndexes Tests
// ============================================================================

func TestAllIndexes(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   []int
	}{
		{"multiple matches", "banana", "an", []int{1, 3}},
		{"single match", "hello", "ll", []int{2}},
		{"no match", "hello", "xyz", nil},
		{"empty substr", "hello", "", nil},
		{"substr longer than s", "hi", "hello", nil},
		{"overlapping matches", "aaa", "aa", []int{0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AllIndexes(tt.s, tt.substr)
			if !slices.Equal(got, tt.want) {
				t.Errorf("AllIndexes(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

// ============================================================================
// HasAnyPrefix/Suffix Tests
// ============================================================================

func TestHasAnyPrefix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefixes []string
		want     bool
	}{
		{"match first", "https://example.com", []string{"http://", "https://"}, true},
		{"match second", "http://example.com", []string{"ftp://", "http://"}, true},
		{"no match", "file://path", []string{"http://", "https://"}, false},
		{"empty prefixes", "hello", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyPrefix(tt.s, tt.prefixes...); got != tt.want {
				t.Errorf("HasAnyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAnySuffix(t *testing.T) {
	if !HasAnySuffix("file.txt", ".txt", ".doc") {
		t.Error("HasAnySuffix should match .txt")
	}
	if HasAnySuffix("file.pdf", ".txt", ".doc") {
		t.Error("HasAnySuffix should not match .pdf")
	}
}

// ============================================================================
// ContainsAny/All Tests
// ============================================================================

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		substrs []string
		want    bool
	}{
		{"match first", "error occurred", []string{"error", "warning"}, true},
		{"match second", "warning issued", []string{"error", "warning"}, true},
		{"no match", "info message", []string{"error", "warning"}, false},
		{"empty substrs", "hello", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAny(tt.s, tt.substrs...); got != tt.want {
				t.Errorf("ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAll(t *testing.T) {
	if !ContainsAll("foo bar baz", "foo", "bar") {
		t.Error("ContainsAll should find both substrings")
	}
	if ContainsAll("foo bar", "foo", "baz") {
		t.Error("ContainsAll should return false when one is missing")
	}
}

// ============================================================================
// Reverse Tests
// ============================================================================

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"ascii", "hello", "olleh"},
		{"unicode", "日本語", "語本日"},
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"with spaces", "hello world", "dlrow olleh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.input); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// IsPalindrome Tests
// ============================================================================

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		normalize bool
		want      bool
	}{
		{"simple palindrome", "racecar", false, true},
		{"not palindrome", "hello", false, false},
		{"normalized", "A man a plan a canal Panama", true, true},
		{"case sensitive", "Racecar", false, false},
		{"empty", "", false, true},
		{"single char", "a", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.s, tt.normalize); got != tt.want {
				t.Errorf("IsPalindrome(%q, %v) = %v, want %v", tt.s, tt.normalize, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Truncate Tests
// ============================================================================

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		suffix string
		want   string
	}{
		{"truncate needed", "Hello World", 8, "...", "Hello..."},
		{"no truncation", "Hello", 10, "...", "Hello"},
		{"exact length", "Hello", 5, "...", "Hello"},
		{"zero length", "Hello", 0, "...", ""},
		{"suffix longer than max", "Hello", 2, "...", ".."},
		{"unicode", "日本語です", 4, "...", "日..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.s, tt.maxLen, tt.suffix); got != tt.want {
				t.Errorf("Truncate(%q, %d, %q) = %q, want %q", tt.s, tt.maxLen, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestTruncateWords(t *testing.T) {
	got := TruncateWords("Hello wonderful world", 15, "...")
	want := "Hello..."
	if got != want {
		t.Errorf("TruncateWords() = %q, want %q", got, want)
	}
}

// ============================================================================
// Pad Tests
// ============================================================================

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		length  int
		padChar rune
		want    string
	}{
		{"pad zeros", "42", 5, '0', "00042"},
		{"no padding needed", "hello", 3, '*', "hello"},
		{"unicode pad", "hi", 5, '日', "日日日hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadLeft(tt.s, tt.length, tt.padChar); got != tt.want {
				t.Errorf("PadLeft(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.padChar, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	got := PadRight("42", 5, '0')
	want := "42000"
	if got != want {
		t.Errorf("PadRight() = %q, want %q", got, want)
	}
}

func TestPadCenter(t *testing.T) {
	got := PadCenter("hello", 11, '*')
	want := "***hello***"
	if got != want {
		t.Errorf("PadCenter() = %q, want %q", got, want)
	}
}

// ============================================================================
// CountLines Tests
// ============================================================================

func TestCountLines(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"single line", "hello", 1},
		{"multiple lines", "a\nb\nc", 3},
		{"trailing newline", "a\nb\n", 2},
		{"empty", "", 0},
		{"only newlines", "\n\n", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountLines(tt.s); got != tt.want {
				t.Errorf("CountLines(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"unix", "a\nb\nc", []string{"a", "b", "c"}},
		{"windows", "a\r\nb\r\nc", []string{"a", "b", "c"}},
		{"trailing newline", "a\nb\n", []string{"a", "b"}},
		{"empty", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lines(tt.s)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Lines(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Character Classification Tests
// ============================================================================

func TestIsBlank(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{" a ", false},
	}

	for _, tt := range tests {
		if got := IsBlank(tt.s); got != tt.want {
			t.Errorf("IsBlank(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"hello", true},
		{"Hello", true},
		{"hello123", false},
		{"", false},
		{"hello world", false},
	}

	for _, tt := range tests {
		if got := IsAlpha(tt.s); got != tt.want {
			t.Errorf("IsAlpha(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"123", true},
		{"123.45", false},
		{"", false},
		{"12a", false},
	}

	for _, tt := range tests {
		if got := IsNumeric(tt.s); got != tt.want {
			t.Errorf("IsNumeric(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsAlphanumeric(t *testing.T) {
	if !IsAlphanumeric("hello123") {
		t.Error("IsAlphanumeric should return true for 'hello123'")
	}
	if IsAlphanumeric("hello 123") {
		t.Error("IsAlphanumeric should return false for 'hello 123'")
	}
}

func TestIsUpper(t *testing.T) {
	if !IsUpper("HELLO") {
		t.Error("IsUpper should return true for 'HELLO'")
	}
	if IsUpper("Hello") {
		t.Error("IsUpper should return false for 'Hello'")
	}
	if !IsUpper("123") { // No letters
		t.Error("IsUpper should return true for string with no letters")
	}
}

func TestIsLower(t *testing.T) {
	if !IsLower("hello") {
		t.Error("IsLower should return true for 'hello'")
	}
	if IsLower("Hello") {
		t.Error("IsLower should return false for 'Hello'")
	}
}

func TestIsASCII(t *testing.T) {
	if !IsASCII("hello") {
		t.Error("IsASCII should return true for 'hello'")
	}
	if IsASCII("héllo") {
		t.Error("IsASCII should return false for 'héllo'")
	}
}

func TestIsPrintable(t *testing.T) {
	if !IsPrintable("hello world") {
		t.Error("IsPrintable should return true for 'hello world'")
	}
	if IsPrintable("hello\x00world") {
		t.Error("IsPrintable should return false for string with null char")
	}
}

// ============================================================================
// Case Conversion Tests
// ============================================================================

func TestCapitalize(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"hello", "Hello"},
		{"hELLO", "Hello"},
		{"", ""},
		{"a", "A"},
	}

	for _, tt := range tests {
		if got := Capitalize(tt.s); got != tt.want {
			t.Errorf("Capitalize(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestSwapCase(t *testing.T) {
	got := SwapCase("Hello World")
	want := "hELLO wORLD"
	if got != want {
		t.Errorf("SwapCase() = %q, want %q", got, want)
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"Hello", "hello"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := SnakeCase(tt.s); got != tt.want {
			t.Errorf("SnakeCase(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := CamelCase(tt.s); got != tt.want {
			t.Errorf("CamelCase(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestPascalCase(t *testing.T) {
	got := PascalCase("hello_world")
	want := "HelloWorld"
	if got != want {
		t.Errorf("PascalCase() = %q, want %q", got, want)
	}
}

func TestKebabCase(t *testing.T) {
	got := KebabCase("HelloWorld")
	want := "hello-world"
	if got != want {
		t.Errorf("KebabCase() = %q, want %q", got, want)
	}
}

// ============================================================================
// Words Tests
// ============================================================================

func TestWords(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{"hello, world!", []string{"hello", "world"}},
		{"one-two three_four", []string{"one", "two", "three", "four"}},
		{"", nil},
	}

	for _, tt := range tests {
		got := Words(tt.s)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Words(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

// ============================================================================
// RuneCount Tests
// ============================================================================

func TestRuneCount(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"hello", 5},
		{"日本語", 3},
		{"", 0},
	}

	for _, tt := range tests {
		if got := RuneCount(tt.s); got != tt.want {
			t.Errorf("RuneCount(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestSafeSlice(t *testing.T) {
	tests := []struct {
		s     string
		start int
		end   int
		want  string
	}{
		{"hello", 0, 3, "hel"},
		{"hello", -1, 3, "hel"},
		{"hello", 0, 100, "hello"},
		{"hello", 10, 20, ""},
		{"日本語", 0, 2, "日本"},
	}

	for _, tt := range tests {
		got := SafeSlice(tt.s, tt.start, tt.end)
		if got != tt.want {
			t.Errorf("SafeSlice(%q, %d, %d) = %q, want %q", tt.s, tt.start, tt.end, got, tt.want)
		}
	}
}

func TestNthRune(t *testing.T) {
	r, ok := NthRune("hello", 1)
	if !ok || r != 'e' {
		t.Errorf("NthRune('hello', 1) = %c, %v; want 'e', true", r, ok)
	}

	_, ok = NthRune("hello", 10)
	if ok {
		t.Error("NthRune should return false for out of bounds")
	}
}

// ============================================================================
// CommonPrefix/Suffix Tests
// ============================================================================

func TestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"normal", []string{"interstellar", "internet", "internal"}, "inter"},
		{"no common", []string{"abc", "xyz"}, ""},
		{"single string", []string{"hello"}, "hello"},
		{"empty slice", []string{}, ""},
		{"identical", []string{"abc", "abc", "abc"}, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommonPrefix(tt.strs...); got != tt.want {
				t.Errorf("CommonPrefix() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommonSuffix(t *testing.T) {
	got := CommonSuffix("testing", "running", "jumping")
	want := "ing"
	if got != want {
		t.Errorf("CommonSuffix() = %q, want %q", got, want)
	}
}

// ============================================================================
// Between Tests
// ============================================================================

func TestBetween(t *testing.T) {
	tests := []struct {
		s, start, end string
		want          string
		ok            bool
	}{
		{"[hello]", "[", "]", "hello", true},
		{"<tag>content</tag>", "<tag>", "</tag>", "content", true},
		{"no markers", "[", "]", "", false},
		{"[hello", "[", "]", "", false}, // only start marker
	}

	for _, tt := range tests {
		got, ok := Between(tt.s, tt.start, tt.end)
		if ok != tt.ok || got != tt.want {
			t.Errorf("Between(%q, %q, %q) = %q, %v; want %q, %v",
				tt.s, tt.start, tt.end, got, ok, tt.want, tt.ok)
		}
	}
}

func TestBetweenAll(t *testing.T) {
	got := BetweenAll("a[1]b[2]c[3]", "[", "]")
	want := []string{"1", "2", "3"}
	if !slices.Equal(got, want) {
		t.Errorf("BetweenAll() = %v, want %v", got, want)
	}
}

// ============================================================================
// Wrap Tests
// ============================================================================

func TestWrap(t *testing.T) {
	input := "This is a long line that should be wrapped at word boundaries"
	got := Wrap(input, 20)

	lines := Lines(got)
	for _, line := range lines {
		if RuneCount(line) > 20 {
			t.Errorf("Wrap produced line longer than width: %q", line)
		}
	}
}

// ============================================================================
// Indent/Dedent Tests
// ============================================================================

func TestIndent(t *testing.T) {
	got := Indent("a\nb\nc", "  ")
	want := "  a\n  b\n  c"
	if got != want {
		t.Errorf("Indent() = %q, want %q", got, want)
	}
}

func TestDedent(t *testing.T) {
	got := Dedent("  a\n  b\n  c")
	want := "a\nb\nc"
	if got != want {
		t.Errorf("Dedent() = %q, want %q", got, want)
	}
}

// ============================================================================
// StripTags Tests
// ============================================================================

func TestStripTags(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"<p>Hello</p>", "Hello"},
		{"<p>Hello <b>World</b></p>", "Hello World"},
		{"No tags here", "No tags here"},
		{"<script>alert('xss')</script>", "alert('xss')"},
	}

	for _, tt := range tests {
		if got := StripTags(tt.s); got != tt.want {
			t.Errorf("StripTags(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

// ============================================================================
// RemoveAll Tests
// ============================================================================

func TestRemoveAll(t *testing.T) {
	got := RemoveAll("hello world", "l", "o")
	want := "he wrd"
	if got != want {
		t.Errorf("RemoveAll() = %q, want %q", got, want)
	}
}

// ============================================================================
// Similarity Algorithm Tests
// ============================================================================

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"kitten", "sitting", 3},
		{"", "", 0},
		{"hello", "", 5},
		{"", "world", 5},
		{"hello", "hello", 0},
		{"ab", "ba", 2},
	}

	for _, tt := range tests {
		if got := LevenshteinDistance(tt.s1, tt.s2); got != tt.want {
			t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestLevenshteinSimilarity(t *testing.T) {
	sim := LevenshteinSimilarity("hello", "hallo")
	if sim < 0.7 || sim > 0.9 {
		t.Errorf("LevenshteinSimilarity('hello', 'hallo') = %f, expected ~0.8", sim)
	}

	if LevenshteinSimilarity("hello", "hello") != 1.0 {
		t.Error("Identical strings should have similarity 1.0")
	}
}

func TestDamerauLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"ca", "ac", 1},  // Single transposition
		{"abc", "acb", 1}, // Single transposition
		{"kitten", "sitting", 3},
	}

	for _, tt := range tests {
		if got := DamerauLevenshteinDistance(tt.s1, tt.s2); got != tt.want {
			t.Errorf("DamerauLevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestJaroSimilarity(t *testing.T) {
	sim := JaroSimilarity("martha", "marhta")
	if sim < 0.94 || sim > 0.96 {
		t.Errorf("JaroSimilarity('martha', 'marhta') = %f, expected ~0.944", sim)
	}

	if JaroSimilarity("hello", "hello") != 1.0 {
		t.Error("Identical strings should have Jaro similarity 1.0")
	}

	if JaroSimilarity("", "hello") != 0.0 {
		t.Error("Empty string should have Jaro similarity 0.0")
	}
}

func TestJaroWinklerSimilarity(t *testing.T) {
	sim := JaroWinklerSimilarity("martha", "marhta", 0.1)
	if sim < 0.96 || sim > 0.98 {
		t.Errorf("JaroWinklerSimilarity('martha', 'marhta') = %f, expected ~0.961", sim)
	}

	// Winkler should be >= Jaro due to common prefix boost
	jaro := JaroSimilarity("martha", "marhta")
	winkler := JaroWinklerSimilarity("martha", "marhta", 0.1)
	if winkler < jaro {
		t.Error("JaroWinkler should be >= Jaro for strings with common prefix")
	}
}

func TestDiceCoefficient(t *testing.T) {
	// Identical strings
	if DiceCoefficient("hello", "hello") != 1.0 {
		t.Error("Identical strings should have Dice coefficient 1.0")
	}

	// Completely different
	dice := DiceCoefficient("abc", "xyz")
	if dice != 0.0 {
		t.Errorf("Completely different strings should have Dice coefficient 0, got %f", dice)
	}

	// Partially similar
	dice = DiceCoefficient("night", "nacht")
	if dice < 0.2 || dice > 0.4 {
		t.Errorf("DiceCoefficient('night', 'nacht') = %f, expected ~0.25", dice)
	}
}

func TestHammingDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"karolin", "kathrin", 3},
		{"hello", "hello", 0},
		{"abc", "abcd", -1}, // Different lengths
	}

	for _, tt := range tests {
		if got := HammingDistance(tt.s1, tt.s2); got != tt.want {
			t.Errorf("HammingDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"ABCDGH", "AEDFHR", 3}, // ADH
		{"AGGTAB", "GXTXAYB", 4}, // GTAB
		{"", "abc", 0},
	}

	for _, tt := range tests {
		if got := LongestCommonSubsequence(tt.s1, tt.s2); got != tt.want {
			t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestLongestCommonSubstring(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   string
	}{
		{"ABABC", "BABCA", "BABC"},
		{"abcdef", "xyzabc", "abc"},
		{"", "abc", ""},
	}

	for _, tt := range tests {
		if got := LongestCommonSubstring(tt.s1, tt.s2); got != tt.want {
			t.Errorf("LongestCommonSubstring(%q, %q) = %q, want %q", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestCosineSimilarity(t *testing.T) {
	// Identical strings
	if CosineSimilarity("hello world", "hello world", 2) != 1.0 {
		t.Error("Identical strings should have cosine similarity 1.0")
	}

	// Similar strings should have high similarity
	sim := CosineSimilarity("hello world", "hello there", 2)
	if sim < 0.3 || sim > 0.7 {
		t.Errorf("CosineSimilarity('hello world', 'hello there') = %f, expected ~0.5", sim)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkReverse(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		Reverse(s)
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	s1 := "kitten"
	s2 := "sitting"
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkLevenshteinDistanceLong(b *testing.B) {
	s1 := "The quick brown fox jumps over the lazy dog"
	s2 := "The slow brown cat walks under the busy cat"
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkJaroWinklerSimilarity(b *testing.B) {
	s1 := "martha"
	s2 := "marhta"
	for i := 0; i < b.N; i++ {
		JaroWinklerSimilarity(s1, s2, 0.1)
	}
}

func BenchmarkSnakeCase(b *testing.B) {
	s := "ThisIsACamelCaseString"
	for i := 0; i < b.N; i++ {
		SnakeCase(s)
	}
}

func BenchmarkAllIndexes(b *testing.B) {
	s := "banana banana banana banana banana"
	for i := 0; i < b.N; i++ {
		AllIndexes(s, "an")
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	s := "A man a plan a canal Panama"
	for i := 0; i < b.N; i++ {
		IsPalindrome(s, true)
	}
}
