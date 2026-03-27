package stringutil

import (
	"strings"
	"testing"
)

// ============================================================================
// CleanString Composable API Tests
// ============================================================================

func TestCleanString_NoOptions(t *testing.T) {
	input := "Hello World"

	got, err := CleanString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != input {
		t.Errorf("CleanString(%q) with no options = %q, want %q", input, got, input)
	}
}

func TestCleanString_AllOptions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "html + unicode + db sanitize",
			input: "<p>Héllo &amp; Wörld</p>",
			want:  "Hello & World",
		},
		{
			name:  "complex html with entities and diacritics",
			input: "<div class='main'>Café &amp; Résumé</div>",
			want:  "Cafe & Resume",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "plain ascii no-op",
			input: "Hello World",
			want:  "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanString(tt.input,
				WithHTMLStrip(),
				WithUnicodeNorm(),
				WithDBSanitize(0),
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("CleanString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCleanString_OrderIndependence(t *testing.T) {
	input := "<p>Héllo &amp; Wörld</p>"

	// Order 1: HTML, Unicode, DB
	r1, err := CleanString(input, WithHTMLStrip(), WithUnicodeNorm(), WithDBSanitize(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Order 2: DB, Unicode, HTML
	r2, err := CleanString(input, WithDBSanitize(0), WithUnicodeNorm(), WithHTMLStrip())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Order 3: Unicode, HTML, DB
	r3, err := CleanString(input, WithUnicodeNorm(), WithHTMLStrip(), WithDBSanitize(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r1 != r2 || r2 != r3 {
		t.Errorf("Options order should not matter:\n  Order1: %q\n  Order2: %q\n  Order3: %q", r1, r2, r3)
	}
}

func TestCleanString_Idempotency(t *testing.T) {
	input := "<p>Héllo &amp; Wörld</p>"

	opts := []CleanOption{WithHTMLStrip(), WithUnicodeNorm(), WithDBSanitize(0)}

	r1, err := CleanString(input, opts...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Apply cleaning again to already-clean output
	r2, err := CleanString(r1, opts...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r1 != r2 {
		t.Errorf("CleanString should be idempotent:\n  First:  %q\n  Second: %q", r1, r2)
	}
}

func TestCleanString_SingleOption_HTMLOnly(t *testing.T) {
	input := "<b>Hello</b> &amp; <i>World</i>"

	got, err := CleanString(input, WithHTMLStrip())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello & World"
	if got != want {
		t.Errorf("CleanString with HTML only = %q, want %q", got, want)
	}
}

func TestCleanString_SingleOption_UnicodeOnly(t *testing.T) {
	input := "café résumé"

	got, err := CleanString(input, WithUnicodeNorm())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "cafe resume"
	if got != want {
		t.Errorf("CleanString with Unicode only = %q, want %q", got, want)
	}
}

func TestCleanString_SingleOption_DBOnly(t *testing.T) {
	input := "Hello\x00World"

	got, err := CleanString(input, WithDBSanitize(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HelloWorld"
	if got != want {
		t.Errorf("CleanString with DB only = %q, want %q", got, want)
	}
}

func TestCleanString_DBSanitize_WithTruncation(t *testing.T) {
	input := "Hello, 世界!"

	got, err := CleanString(input, WithDBSanitize(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello"
	if got != want {
		t.Errorf("CleanString with DB truncation = %q, want %q", got, want)
	}
}

// ============================================================================
// NormalizeUnicode Tests
// ============================================================================

func TestNormalizeUnicode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"ascii only", "Hello World", "Hello World"},
		{"french accents", "café résumé", "cafe resume"},
		{"german umlauts", "über Straße", "uber Straße"},
		{"spanish tilde", "señor niño", "senor nino"},
		{"nordic", "Ångström", "Angstrom"},
		{"naive", "naïve", "naive"},
		{"mixed", "Héllo Wörld", "Hello World"},
		{"combining acute", "e\u0301", "e"}, // e + combining acute accent → e
		{"precomposed e-acute", "é", "e"},
		{"numbers preserved", "café123", "cafe123"},
		{"punctuation preserved", "café, résumé!", "cafe, resume!"},
		{"spaces preserved", "  café  ", "  cafe  "},
		{"emoji preserved", "Hello 🌍", "Hello 🌍"},
		{"japanese", "日本語", "日本語"},             // CJK characters have no diacritics to remove
		{"korean", "한국어", "한국어"},               // Hangul stays as-is (decomposed and recomposed)
		{"cedilla", "façade", "facade"},
		{"circumflex", "crêpe", "crepe"},
		{"tilde n", "cañón", "canon"},
		{"double acute", "ő", "o"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeUnicode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("NormalizeUnicode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// StripHTMLEntities Tests
// ============================================================================

func TestStripHTMLEntities(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"no html", "Hello World", "Hello World"},
		{"simple tags", "<p>Hello</p>", "Hello"},
		{"nested tags", "<div><p>Hello <b>World</b></p></div>", "Hello World"},
		{"self-closing tags", "Hello<br/>World", "HelloWorld"},
		{"tags with attributes", `<a href="url">Link</a>`, "Link"},

		// Entity decoding
		{"amp entity", "Hello &amp; World", "Hello & World"},
		{"lt gt entities", "5 &lt; 10 &gt; 3", "5 < 10 > 3"},
		{"quot entity", "He said &quot;hello&quot;", `He said "hello"`},
		{"nbsp entity", "Hello&nbsp;World", "Hello World"},
		{"copy entity", "&copy; 2024", "© 2024"},
		{"euro entity", "Price: &euro;10", "Price: €10"},
		{"reg entity", "Brand&reg;", "Brand®"},

		// Numeric entities
		{"decimal entity", "&#169;", "©"},
		{"hex entity", "&#x00A9;", "©"},
		{"decimal mdash", "Hello&#8212;World", "Hello—World"},

		// Combined tags + entities
		{"tags and entities", "<p>Hello &amp; World</p>", "Hello & World"},
		{"complex mixed", "<div>&lt;script&gt;alert('xss')&lt;/script&gt;</div>", "<script>alert('xss')</script>"},
		{"multiple entities", "&lt;&gt;&amp;&quot;", `<>&"`},

		// Edge cases
		{"unclosed tag", "<p>Hello", "Hello"},
		// Note: bare angle brackets are treated as tags by StripTags. This is a
		// known limitation of the simple tag stripper. For proper HTML parsing,
		// use html/template or a full HTML parser. The "< 10" is treated as an
		// unclosed tag, so "10" is consumed as tag content.
		{"bare angle brackets", "5 > 3 < 10", "5  3 "},
		{"empty tags", "<><>Hello<>", "Hello"},
		{"only tags", "<p><br/></p>", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripHTMLEntities(tt.input)
			if got != tt.want {
				t.Errorf("StripHTMLEntities(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// SanitizeUTF8 Tests
// ============================================================================

func TestSanitizeUTF8(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"valid ascii", "Hello World", "Hello World"},
		{"valid unicode", "日本語", "日本語"},
		{"valid emoji", "Hello 🌍🎉", "Hello 🌍🎉"},
		{"nul byte removal", "Hello\x00World", "HelloWorld"},
		{"multiple nul bytes", "H\x00e\x00l\x00l\x00o", "Hello"},
		{"invalid utf8 byte", "Hello\xffWorld", "Hello\uFFFDWorld"},
		{"invalid utf8 sequence", "Hello\xc0\xafWorld", "Hello\uFFFD\uFFFDWorld"},
		{"nul and invalid mixed", "H\x00\xffW", "H\uFFFDW"},
		{"only nul bytes", "\x00\x00\x00", ""},
		{"valid with multibyte", "café ☕", "café ☕"},
		{"trailing invalid", "Hello\xff", "Hello\uFFFD"},
		{"leading invalid", "\xffHello", "\uFFFDHello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeUTF8(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeUTF8(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// TruncateRunes Tests
// ============================================================================

func TestTruncateRunes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"empty string", "", 5, ""},
		{"no truncation needed", "Hello", 10, "Hello"},
		{"exact length", "Hello", 5, "Hello"},
		{"truncate ascii", "Hello World", 5, "Hello"},
		{"truncate unicode", "日本語です", 3, "日本語"},
		{"zero length", "Hello", 0, ""},
		{"negative length", "Hello", -1, ""},
		{"single char", "Hello", 1, "H"},
		{"emoji truncation", "Hello 🌍🎉", 7, "Hello 🌍"},
		{"mixed unicode", "café", 3, "caf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateRunes(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateRunes(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

// ============================================================================
// RemoveNonPrintable Tests
// ============================================================================

func TestRemoveNonPrintable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"printable only", "Hello World", "Hello World"},
		{"bell character", "Hello\x07World", "HelloWorld"},
		{"null character", "Hello\x00World", "HelloWorld"},
		{"tab preserved", "Hello\tWorld", "Hello\tWorld"},
		{"newline preserved", "Hello\nWorld", "Hello\nWorld"},
		{"cr preserved", "Hello\rWorld", "Hello\rWorld"},
		{"escape character", "Hello\x1bWorld", "HelloWorld"},
		{"backspace", "Hello\x08World", "HelloWorld"},
		{"mixed control chars", "\x00\x07Hello\x1b\x08World\x7f", "HelloWorld"},
		{"unicode printable", "日本語 café", "日本語 café"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveNonPrintable(tt.input)
			if got != tt.want {
				t.Errorf("RemoveNonPrintable(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// NormalizeWhitespace Tests
// ============================================================================

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"single space", "Hello World", "Hello World"},
		{"multiple spaces", "Hello   World", "Hello World"},
		{"tabs", "Hello\tWorld", "Hello World"},
		{"newlines", "Hello\nWorld", "Hello World"},
		{"mixed whitespace", "  Hello \t World \n ! ", "Hello World !"},
		{"only whitespace", "   \t\n  ", ""},
		{"leading trailing", "  Hello  ", "Hello"},
		{"no change needed", "Hello", "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// RemoveAccents Tests (alias for NormalizeUnicode)
// ============================================================================

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"cafe", "café", "cafe"},
		{"uber", "über", "uber"},
		{"resume", "résumé", "resume"},
		{"empty", "", ""},
		{"ascii", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveAccents(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("RemoveAccents(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// ToASCII Tests
// ============================================================================

func TestToASCII(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"ascii only", "Hello World", "Hello World"},
		{"french accents", "café résumé", "cafe resume"},
		{"numbers", "café123", "cafe123"},
		{"punctuation", "Hello, World!", "Hello, World!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToASCII(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("ToASCII(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Slugify Tests
// ============================================================================

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"simple", "Hello World", "hello-world"},
		{"accented", "Café Résumé", "cafe-resume"},
		{"multiple spaces", "Hello   World", "hello-world"},
		{"special chars", "Hello, World!", "hello-world"},
		{"leading trailing spaces", "  Hello World  ", "hello-world"},
		{"numbers", "Product 123", "product-123"},
		{"mixed", "Héllo Wörld 2024!", "hello-world-2024"},
		{"already slug", "hello-world", "hello-world"},
		{"underscores", "hello_world", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Slugify(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestCleanString_VeryLongString(t *testing.T) {
	// Test with a very long string to ensure no performance issues
	input := strings.Repeat("Héllo Wörld ", 10000)

	got, err := CleanString(input, WithUnicodeNorm())
	if err != nil {
		t.Fatalf("unexpected error on long string: %v", err)
	}

	want := strings.Repeat("Hello World ", 10000)
	if got != want {
		t.Errorf("CleanString on long string produced unexpected output (length: %d vs %d)", len(got), len(want))
	}
}

func TestCleanString_OnlyWhitespace(t *testing.T) {
	input := "   \t\n   "

	got, err := CleanString(input, WithHTMLStrip(), WithUnicodeNorm(), WithDBSanitize(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Whitespace should be preserved (cleaning doesn't trim)
	if got != input {
		t.Errorf("CleanString on whitespace = %q, want %q", got, input)
	}
}

func TestCleanString_EmojiPreserved(t *testing.T) {
	input := "Hello 🌍🎉 World"

	got, err := CleanString(input, WithUnicodeNorm(), WithDBSanitize(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello 🌍🎉 World"
	if got != want {
		t.Errorf("CleanString should preserve emoji: %q, want %q", got, want)
	}
}

func TestCleanString_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"newlines", "Hello\nWorld", "Hello\nWorld"},
		{"tabs", "Hello\tWorld", "Hello\tWorld"},
		// Note: NFKD compatibility decomposition converts em-space (U+2003)
		// to regular space (U+0020). This is correct NFKD behavior.
		{"unicode em-space normalized", "Hello\u2003World", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanString(tt.input, WithUnicodeNorm(), WithDBSanitize(0))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("CleanString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCleanString_HTMLWithUnicode(t *testing.T) {
	input := "<p>Caf&eacute; R&eacute;sum&eacute;</p>"

	got, err := CleanString(input, WithHTMLStrip(), WithUnicodeNorm())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Cafe Resume"
	if got != want {
		t.Errorf("CleanString HTML+Unicode = %q, want %q", got, want)
	}
}

func TestCleanString_DBTruncateUnicode(t *testing.T) {
	// Ensure truncation is rune-aware
	input := "日本語です"

	got, err := CleanString(input, WithDBSanitize(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "日本語"
	if got != want {
		t.Errorf("CleanString DB truncate = %q, want %q", got, want)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkCleanString_AllOptions(b *testing.B) {
	input := "<p>Héllo &amp; Wörld café résumé</p>"
	opts := []CleanOption{WithHTMLStrip(), WithUnicodeNorm(), WithDBSanitize(100)}

	for b.Loop() {
		_, _ = CleanString(input, opts...)
	}
}

func BenchmarkNormalizeUnicode(b *testing.B) {
	input := "café résumé naïve Ångström über señor"
	for b.Loop() {
		_, _ = NormalizeUnicode(input)
	}
}

func BenchmarkStripHTMLEntities(b *testing.B) {
	input := "<div class='main'><p>Hello &amp; World</p><br/><p>Price: &euro;10 &copy; 2024</p></div>"
	for b.Loop() {
		StripHTMLEntities(input)
	}
}

func BenchmarkSanitizeUTF8(b *testing.B) {
	input := "Hello\x00World\xffTest\x00Valid"
	for b.Loop() {
		SanitizeUTF8(input)
	}
}

func BenchmarkSanitizeUTF8_Valid(b *testing.B) {
	input := "Hello World this is a perfectly valid UTF-8 string"
	for b.Loop() {
		SanitizeUTF8(input)
	}
}

func BenchmarkSlugify(b *testing.B) {
	input := "Héllo Wörld 2024! Café & Résumé"
	for b.Loop() {
		_, _ = Slugify(input)
	}
}

func BenchmarkNormalizeWhitespace(b *testing.B) {
	input := "  Hello   World  \t\n  How   are  you  "
	for b.Loop() {
		NormalizeWhitespace(input)
	}
}
