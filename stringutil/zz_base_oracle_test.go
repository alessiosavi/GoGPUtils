package stringutil

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

// BASE 4c05787953a944b57478180f4c27100869db12f6: bodies unchanged except names.
func baseDamerauLevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}

	if len2 == 0 {
		return len1
	}

	// Create distance matrix
	d := make([][]int, len1+1)
	for i := range d {
		d[i] = make([]int, len2+1)
		d[i][0] = i
	}

	for j := 0; j <= len2; j++ {
		d[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			d[i][j] = min(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)

			// Transposition
			if i > 1 && j > 1 && r1[i-1] == r2[j-2] && r1[i-2] == r2[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+cost)
			}
		}
	}

	return d[len1][len2]
}

func baseDiceCoefficient(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Generate bigrams
	bigrams1 := baseBigrams(strings.ToLower(s1))
	bigrams2 := baseBigrams(strings.ToLower(s2))

	if len(bigrams1) == 0 && len(bigrams2) == 0 {
		return 1.0
	}

	if len(bigrams1) == 0 || len(bigrams2) == 0 {
		return 0.0
	}

	// Count intersection
	intersection := 0

	counted := make(map[string]int)
	for bg := range bigrams1 {
		counted[bg] = bigrams1[bg]
	}

	for bg, count := range bigrams2 {
		if counted[bg] > 0 {
			commonCount := min(count, counted[bg])
			intersection += commonCount
		}
	}

	// Calculate total bigrams
	total1 := 0
	for _, count := range bigrams1 {
		total1 += count
	}

	total2 := 0
	for _, count := range bigrams2 {
		total2 += count
	}

	return 2.0 * float64(intersection) / float64(total1+total2)
}

func baseBigrams(s string) map[string]int {
	runes := []rune(s)
	if len(runes) < 2 {
		return nil
	}

	result := make(map[string]int)

	for i := range len(runes) - 1 {
		bg := string(runes[i : i+2])
		result[bg]++
	}

	return result
}

func baseCosineSimilarity(s1, s2 string, n int) float64 {
	if n <= 0 {
		n = 2
	}

	if s1 == s2 {
		return 1.0
	}

	ngrams1 := baseNgrams(strings.ToLower(s1), n)
	ngrams2 := baseNgrams(strings.ToLower(s2), n)

	if len(ngrams1) == 0 && len(ngrams2) == 0 {
		return 1.0
	}

	if len(ngrams1) == 0 || len(ngrams2) == 0 {
		return 0.0
	}

	// Calculate dot product and magnitudes
	dotProduct := 0.0
	mag1 := 0.0
	mag2 := 0.0

	// Build combined key set
	allKeys := make(map[string]bool)
	for k := range ngrams1 {
		allKeys[k] = true
	}

	for k := range ngrams2 {
		allKeys[k] = true
	}

	for k := range allKeys {
		v1 := float64(ngrams1[k])
		v2 := float64(ngrams2[k])
		dotProduct += v1 * v2
		mag1 += v1 * v1
		mag2 += v2 * v2
	}

	if mag1 == 0 || mag2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(mag1) * math.Sqrt(mag2))
}

func baseNgrams(s string, n int) map[string]int {
	runes := []rune(s)
	if len(runes) < n {
		return nil
	}

	result := make(map[string]int)

	for i := 0; i <= len(runes)-n; i++ {
		ng := string(runes[i : i+n])
		result[ng]++
	}

	return result
}

func pr3Pair(t *testing.T, a, b string) {
	t.Helper()
	if got, want := DamerauLevenshteinDistance(a, b), baseDamerauLevenshteinDistance(a, b); got != want {
		t.Fatalf("OSA(%q,%q)=%d want=%d", a, b, got, want)
	}
	if got, want := DiceCoefficient(a, b), baseDiceCoefficient(a, b); math.Float64bits(got) != math.Float64bits(want) {
		t.Fatalf("Dice(%q,%q)=%.17g want=%.17g", a, b, got, want)
	}
	for _, n := range []int{-1, 0, 1, 2, 3, 5} {
		if got, want := CosineSimilarity(a, b, n), baseCosineSimilarity(a, b, n); math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("Cosine(%q,%q,%d)=%.17g want=%.17g", a, b, n, got, want)
		}
	}
}

func TestPR3Adversarial(t *testing.T) {
	cases := [][2]string{
		{"ca", "ac"}, {"CA", "ABC"}, {"aab", "aba"}, {"ab", "ba"},
		{"aaaa", "aa"}, {"abababab", "abab"}, {"éééé", "éé"},
		{"", ""}, {"", "é"}, {"é", ""}, {"a", "b"}, {"", "ab"},
		{"A", "a"}, {"AB", "ab"}, {"ẞİᎠ", "ßiꭰ"}, {"東京😀", "東😀京"},
		{"\xff", "\xfe"}, {"\xff\xfe", "��"}, {"\xffa", "�a"},
		{"\xc0\xaf", "\xe0\x80\xaf"}, {"\x00\xff\xc3", "\x00��"},
		{"abc", "a"}, {"a", "abc"}, {"ab", "a"}, {"abcde", "abcdf"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) { pr3Pair(t, c[0], c[1]) })
	}
	for _, n := range []int{31, 32, 33, 64, 512, int(^uint(0) >> 1)} {
		a, b := "abé", "ac東"
		if got, want := CosineSimilarity(a, b, n), baseCosineSimilarity(a, b, n); math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("large width n=%d got=%g want=%g", n, got, want)
		}
	}
}

func TestPR3Random(t *testing.T) {
	for class := range 7 {
		t.Run(fmt.Sprint(class), func(t *testing.T) {
			rng := rand.New(rand.NewPCG(0x707233, uint64(class)))
			makeText := func(alphabet string) string {
				runes := []rune(alphabet)
				out := make([]rune, rng.IntN(25))
				for i := range out {
					out[i] = runes[rng.IntN(len(runes))]
				}
				return string(out)
			}
			for i := range 20000 {
				var a, b string
				switch class {
				case 0:
					a, b = makeText("abABcdXY012 !"), makeText("abABcdXY012 !")
				case 1:
					alphabet := "ab"
					if i%2 == 1 {
						alphabet = "abc"
					}
					a, b = makeText(alphabet), makeText(alphabet)
					if i%3 == 0 && len(a) > 1 {
						x := []byte(a)
						j := rng.IntN(len(x) - 1)
						x[j], x[j+1] = x[j+1], x[j]
						b = string(x)
					}
				case 2:
					a, b = makeText("ẞßİiᎠꭰÉé東😀Σσς"), makeText("ẞßİiᎠꭰÉé東😀Σσς")
				case 3:
					bytesText := func() string {
						v := make([]byte, 1+rng.IntN(24))
						for j := range v {
							v[j] = byte(rng.IntN(256))
						}
						v[0] = 0xff
						return string(v)
					}
					a, b = bytesText(), bytesText()
				case 4:
					a, b = "�"+makeText("a�é\x00"), "�"+makeText("a�é\x00")
				case 5:
					a, b = "", makeText("ẞaé�")
					if i%2 == 1 {
						a, b = b, a
					}
				case 6:
					a = makeText("ẞabé�\x00")
					if i%2 == 1 {
						a += "\xff"
					}
					b = a
				}
				pr3Pair(t, a, b)
			}
		})
	}
}

func TestPR3Exhaustive(t *testing.T) {
	words := []string{""}
	frontier := []string{""}
	for length := 1; length <= 5; length++ {
		var next []string
		for _, s := range frontier {
			for _, r := range "ab" {
				next = append(next, s+string(r))
			}
		}
		words = append(words, next...)
		frontier = next
	}
	for _, a := range words {
		for _, b := range words {
			pr3Pair(t, a, b)
		}
	}
	t.Logf("%d exhaustive binary pairs", len(words)*len(words))
}

func TestPR3LoweredKeys(t *testing.T) {
	for value := range 65536 {
		raw := string([]byte{byte(value >> 8), byte(value)})
		s := strings.ToLower(raw)
		pr3Pair(t, raw, "a�é")
		if !utf8.ValidString(s) {
			t.Fatalf("invalid lower %04x: %q", value, s)
		}
		if got, want := bigrams(s), baseBigrams(s); !reflect.DeepEqual(got, want) {
			t.Fatalf("bigrams lower(%04x) got=%v want=%v", value, got, want)
		}
		for n := 1; n <= 3; n++ {
			if got, want := ngrams(s, n), baseNgrams(s, n); !reflect.DeepEqual(got, want) {
				t.Fatalf("ngrams lower(%04x),n=%d got=%v want=%v", value, n, got, want)
			}
		}
	}
	t.Log("65,536 two-byte inputs: valid lowered UTF-8, helper parity and public Dice/Cosine parity")
}

func TestPR3Allocs(t *testing.T) {
	a, b := strings.Repeat("abé", 24), strings.Repeat("baé", 24)
	t.Logf("diagnostic only: OSA=%g Dice=%g Cosine=%g",
		testing.AllocsPerRun(5, func() { DamerauLevenshteinDistance(a, b) }),
		testing.AllocsPerRun(5, func() { DiceCoefficient(a, b) }),
		testing.AllocsPerRun(5, func() { CosineSimilarity(a, b, 2) }))
}
