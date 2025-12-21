package stringutil

import (
	"strings"
)

// LevenshteinDistance returns the minimum number of single-character edits
// (insertions, deletions, substitutions) required to change s1 into s2.
//
// Time complexity: O(len(s1) * len(s2))
// Space complexity: O(min(len(s1), len(s2)))
//
// Example:
//
//	LevenshteinDistance("kitten", "sitting")  // 3
func LevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	// Convert to runes for proper Unicode handling
	r1 := []rune(s1)
	r2 := []rune(s2)

	// Ensure r1 is the shorter one for space optimization
	if len(r1) > len(r2) {
		r1, r2 = r2, r1
	}

	// Use single row for space optimization
	prev := make([]int, len(r1)+1)
	curr := make([]int, len(r1)+1)

	// Initialize first row
	for i := range prev {
		prev[i] = i
	}

	for j := 1; j <= len(r2); j++ {
		curr[0] = j

		for i := 1; i <= len(r1); i++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			curr[i] = min(
				prev[i]+1,      // deletion
				curr[i-1]+1,    // insertion
				prev[i-1]+cost, // substitution
			)
		}

		prev, curr = curr, prev
	}

	return prev[len(r1)]
}

// LevenshteinSimilarity returns a similarity score between 0 and 1
// based on Levenshtein distance. 1 means identical strings.
//
// Example:
//
//	LevenshteinSimilarity("hello", "hallo")  // ~0.8
func LevenshteinSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	maxLen := max(len([]rune(s1)), len([]rune(s2)))
	if maxLen == 0 {
		return 1.0
	}

	distance := LevenshteinDistance(s1, s2)

	return 1.0 - float64(distance)/float64(maxLen)
}

// DamerauLevenshteinDistance extends Levenshtein to include transpositions
// (swapping two adjacent characters) as a single edit operation.
//
// Example:
//
//	DamerauLevenshteinDistance("ca", "ac")  // 1 (transposition)
//	LevenshteinDistance("ca", "ac")         // 2 (delete + insert)
func DamerauLevenshteinDistance(s1, s2 string) int {
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

// JaroSimilarity returns the Jaro similarity between two strings.
// Returns a value between 0 (completely different) and 1 (identical).
//
// The algorithm considers:
// - Number of matching characters
// - Number of transpositions
//
// Example:
//
//	JaroSimilarity("martha", "marhta")  // ~0.944
func JaroSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	// Calculate match window
	matchWindow := max(max(len1, len2)/2-1, 0)

	s1Matches := make([]bool, len1)
	s2Matches := make([]bool, len2)

	matches := 0
	transpositions := 0

	// Find matches
	for i := range len1 {
		start := max(0, i-matchWindow)
		end := min(len2, i+matchWindow+1)

		for j := start; j < end; j++ {
			if s2Matches[j] || r1[i] != r2[j] {
				continue
			}

			s1Matches[i] = true
			s2Matches[j] = true
			matches++

			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	j := 0

	for i := range len1 {
		if !s1Matches[i] {
			continue
		}

		for !s2Matches[j] {
			j++
		}

		if r1[i] != r2[j] {
			transpositions++
		}

		j++
	}

	m := float64(matches)

	return (m/float64(len1) + m/float64(len2) + (m-float64(transpositions)/2)/m) / 3
}

// JaroWinklerSimilarity returns the Jaro-Winkler similarity between two strings.
// This is an extension of Jaro that gives more weight to strings with a common prefix.
//
// The prefixScale parameter (0 to 0.25) determines how much weight to give
// to the common prefix. Standard value is 0.1.
//
// Example:
//
//	JaroWinklerSimilarity("martha", "marhta", 0.1)  // ~0.961
func JaroWinklerSimilarity(s1, s2 string, prefixScale float64) float64 {
	jaroSim := JaroSimilarity(s1, s2)

	// Clamp prefixScale to valid range
	if prefixScale < 0 {
		prefixScale = 0
	}

	if prefixScale > 0.25 {
		prefixScale = 0.25
	}

	// Find common prefix length (max 4 characters)
	r1 := []rune(s1)
	r2 := []rune(s2)
	prefixLen := 0
	maxPrefix := min(4, min(len(r1), len(r2)))

	for i := range maxPrefix {
		if r1[i] == r2[i] {
			prefixLen++
		} else {
			break
		}
	}

	return jaroSim + float64(prefixLen)*prefixScale*(1-jaroSim)
}

// DiceCoefficient returns the Sørensen–Dice coefficient comparing bigrams.
// Returns a value between 0 and 1, where 1 means identical sets of bigrams.
//
// This metric is useful for comparing short strings or when order matters less.
//
// Example:
//
//	DiceCoefficient("night", "nacht")  // ~0.25
func DiceCoefficient(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Generate bigrams
	bigrams1 := bigrams(strings.ToLower(s1))
	bigrams2 := bigrams(strings.ToLower(s2))

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

// bigrams generates a map of bigrams and their counts.
func bigrams(s string) map[string]int {
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

// HammingDistance returns the number of positions where corresponding
// characters differ. Only defined for strings of equal length.
// Returns -1 if strings have different lengths.
//
// Example:
//
//	HammingDistance("karolin", "kathrin")  // 3
func HammingDistance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)

	if len(r1) != len(r2) {
		return -1
	}

	distance := 0

	for i := range r1 {
		if r1[i] != r2[i] {
			distance++
		}
	}

	return distance
}

// LongestCommonSubsequence returns the length of the longest common subsequence.
// A subsequence is a sequence that can be derived by deleting some elements
// without changing the order of remaining elements.
//
// Example:
//
//	LongestCommonSubsequence("ABCDGH", "AEDFHR")  // 3 ("ADH")
func LongestCommonSubsequence(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 || len2 == 0 {
		return 0
	}

	// Space-optimized: use two rows
	prev := make([]int, len2+1)
	curr := make([]int, len2+1)

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			if r1[i-1] == r2[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}

		prev, curr = curr, prev
	}

	return prev[len2]
}

// LongestCommonSubstring returns the longest common contiguous substring.
//
// Example:
//
//	LongestCommonSubstring("ABABC", "BABCA")  // "BABC"
func LongestCommonSubstring(s1, s2 string) string {
	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 || len2 == 0 {
		return ""
	}

	maxLen := 0
	endIndex := 0

	// Use single row for space optimization
	prev := make([]int, len2+1)
	curr := make([]int, len2+1)

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			if r1[i-1] == r2[j-1] {
				curr[j] = prev[j-1] + 1
				if curr[j] > maxLen {
					maxLen = curr[j]
					endIndex = i
				}
			} else {
				curr[j] = 0
			}
		}

		prev, curr = curr, prev
	}

	if maxLen == 0 {
		return ""
	}

	return string(r1[endIndex-maxLen : endIndex])
}

// CosineSimilarity computes the cosine similarity of two strings based on
// their character n-gram vectors. Returns a value between 0 and 1.
//
// This is useful for comparing longer texts.
func CosineSimilarity(s1, s2 string, n int) float64 {
	if n <= 0 {
		n = 2
	}

	if s1 == s2 {
		return 1.0
	}

	ngrams1 := ngrams(strings.ToLower(s1), n)
	ngrams2 := ngrams(strings.ToLower(s2), n)

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

	return dotProduct / (sqrt(mag1) * sqrt(mag2))
}

// ngrams generates a map of n-grams and their counts.
func ngrams(s string, n int) map[string]int {
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

// Simple square root using Newton's method (to avoid math package import).
func sqrt(x float64) float64 {
	if x == 0 || x == 1 {
		return x
	}

	guess := x / 2
	for range 50 { // 50 iterations for precision
		newGuess := (guess + x/guess) / 2
		if abs(newGuess-guess) < 1e-15 {
			break
		}

		guess = newGuess
	}

	return guess
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}

	return x
}
