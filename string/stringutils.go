package stringutils

import (
	"bufio"
	"bytes"
	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	mathutils "github.com/alessiosavi/GoGPUtils/math"
	"io"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

var BOMS = [][]byte{
	{0xef, 0xbb, 0xbf},       // UTF-8
	{0xfe, 0xff},             // UTF-16 BE
	{0xff, 0xfe},             // UTF-16 LE
	{0x00, 0x00, 0xfe, 0xff}, // UTF-32 BE
	{0xff, 0xfe, 0x00, 0x00}, // UTF-32 LE

}
var CUTSET = " \n\r\t"

func Indexes(s string, chs string) []int {
	var ret []int
	for i := 0; i <= len(s)-len(chs); i++ {
		if s[i:i+len(chs)] == chs {
			ret = append(ret, i)
		}
	}
	return ret
}

func HasPrefixes(prefixs []string, target string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return false
}

// IsUpper verify that a string contains only upper character
func IsUpper(s string) bool {
	for _, r := range s {
		switch {
		case !unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter}, r):
			continue
		case !unicode.IsUpper(r):
			return false
		}
	}
	return true
}

// IsLower verify that a string contains only upper character
func IsLower(s string) bool {
	for _, r := range s {
		switch {
		case !unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter}, r):
			continue

		case !unicode.IsLower(r):
			return false
		}
	}
	return true
}

func IsAlpha(str string) bool {
	for _, r := range str {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)) {
			return false
		}
	}
	return true
}

// ContainsMultiple is delegated to verify if the given string 's' contains all the 'substring' present
func ContainsMultiple(lower bool, s string, substring ...string) bool {
	if lower {
		s = strings.ToLower(s)
		for i := range substring {
			substring[i] = strings.ToLower(substring[i])
		}
	}
	for _, toFind := range substring {
		if !strings.Contains(s, toFind) {
			return false
		}
	}
	return true
}

// CreateJSON is delegated to create a simple json object for the key pair in input
func CreateJSON(values []string) string {
	json := `{`
	length := len(values)

	// Not a key-value list
	if length%2 != 0 {
		log.Fatal("Not a key-value list")
		return ""
	}
	for i := 0; i < length; i += 2 {
		json = strings.Join([]string{json, `"`, values[i], `":"`, values[i+1], `",`}, "")
	}
	json = strings.TrimSuffix(json, `,`)
	json += `}`
	return json
}

// IsASCII is delegated to verify if a given string is ASCII compliant
func IsASCII(s string) bool {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// IsASCIIRune is delegated to verify if the given character is ASCII compliant
func IsASCIIRune(r rune) bool {
	return r < utf8.RuneSelf
}

// RemoveFromString Remove a given character in position i from the input string
func RemoveFromString(data string, i int) string {
	if i >= len(data) {
		return data
	}
	s := []byte(data)
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return string(s[:len(s)-1])
}

// Split is delegated to split the string by the new line
func Split(data io.Reader) []string {
	var linesList []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		linesList = append(linesList, scanner.Text())
	}
	return linesList
}

// CountLines return the number of lines in the given string
func CountLines(str string) int {
	scanner := bufio.NewScanner(strings.NewReader(str)) // Create a scanner for iterate the string
	n := 0
	for scanner.Scan() {
		n++
	}
	return n
}

// ExtractString is delegated to filter the content of the given data delimited by 'first' and 'last' string
func ExtractString(data string, first, last string) string {
	// Find the first instance of 'start' in the give string data
	startHeder := strings.Index(data, first)
	if startHeder != -1 { // Found !
		startHeder += len(first) // Remove the first word
		// Check the first occurrence of 'last' that delimit the string to return
		endHeader := strings.Index(data[startHeder:], last)
		// Ok, seems good, return the content of the string delimited by 'first' and 'last'
		if endHeader != -1 {
			return strings.Trim(data[startHeder:startHeder+endHeader], CUTSET)
		}
	}
	return ""
}

// ReplaceAtIndex is delegated to replace the character related to the index with the input rune
func ReplaceAtIndex(str, replacement string, index int) string {
	return str[:index] + replacement + str[index+1:]
}

// RemoveNonASCII is delegated to clean the text from the NON ASCII character
func RemoveNonASCII(str string) string {
	var b bytes.Buffer
	b.Grow(len(str))
	for _, r := range str {
		if IsASCIIRune(r) {
			b.WriteRune(r)
		}
	}

	scanner := bufio.NewScanner(bytes.NewReader(b.Bytes()))
	scanner.Split(bufio.ScanWords)
	var sb strings.Builder
	v := scanner.Text()
	sb.WriteString(v)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteRune(' ')
	}

	return sb.String()[:sb.Len()-1]
}

// IsBlank is delegated to verify that the string does not contain only empty char
func IsBlank(str string) bool {
	// Check length
	if len(str) > 0 {
		// Iterate string
		for i := range str {
			// Check about char different from whitespace
			if str[i] > 32 {
				return false
			}
		}
	}
	return true
}

// Trim is delegated to remove the initial, final whitespace and the double whitespace present in the data
func Trim(str string) string {
	if str == "" {
		return str
	}
	var b strings.Builder
	replacer := strings.NewReplacer("  ", " ", "\\n", "", "\\r", "", "\\t", "")
	str = replacer.Replace(str)
	b.Grow(len(str))

	for i := 0; i < len(str); i++ {
		// Write character
		if str[i] > 32 {
			b.WriteByte(str[i])
			// Print the space only if followed by an ASCII character
		} else if i+1 < len(str) && (str[i] < 33 && str[i+1] > 32) {
			b.WriteByte(str[i])
		}
	}

	return strings.Trim(b.String(), CUTSET)
}

// CheckPalindrome is delegated to verify if the given string is palindrome
func CheckPalindrome(str string) bool {
	length := len(str) - 1
	for i := range str {
		if str[i] != str[length-i] {
			return false
		}
	}
	return true
}

// ReverseString is delegated to return the reverse of the input string
func ReverseString(str string) string {
	length := len(str) - 1
	var sb strings.Builder
	for i := length; i >= 0; i-- {
		sb.WriteByte(str[i])
	}
	return sb.String()
}

// LevenshteinDistance is an optimized version for calculate the Levenshtein distance
func LevenshteinDistance(str1, str2 string) int {
	var n, m = len(str1), len(str2)
	if n == 0 {
		return m
	} else if m == 0 {
		return n
	}

	var p = make([]int, n+1)
	var d = make([]int, n+1)

	var i, j, cost int
	var t_j byte

	for i = 0; i <= n; i++ {
		p[i] = i
	}

	for j = 1; j <= m; j++ {
		t_j = str2[j-1]
		d[0] = j

		for i = 1; i <= n; i++ {
			if str1[i-1] == t_j {
				cost = 0
			} else {
				cost = 1
			}

			// minimum of cell to the left+1, to the top+1, diagonally left and up +cost
			d[i] = mathutils.Min[int](mathutils.Min[int](d[i-1]+1, p[i]+1), p[i-1]+cost)
		}

		// copy current distance counts to 'previous row' distance counts
		p, d = d, p
	}
	return p[n]
}

// JaroDistance is delegated to calculate the Jaro distance from the two given string
func JaroDistance(str1, str2 string) float64 {
	if str1 == str2 {
		return 1
	}

	if str1 == "" || str2 == "" {
		return 0
	}

	matchDistance := len(str1)
	if len(str2) > matchDistance {
		matchDistance = len(str2)
	}
	matchDistance = matchDistance/2 - 1
	str1Matches := make([]bool, len(str1))
	str2Matches := make([]bool, len(str2))
	matches := 0.
	transpositions := 0.
	for i := range str1 {
		start := i - matchDistance
		if start < 0 {
			start = 0
		}
		end := i + matchDistance + 1
		if end > len(str2) {
			end = len(str2)
		}
		for k := start; k < end; k++ {
			if str2Matches[k] {
				continue
			}
			if str1[i] != str2[k] {
				continue
			}
			str1Matches[i] = true
			str2Matches[k] = true
			matches++
			break
		}
	}
	if matches == 0 {
		return 0
	}
	k := 0
	for i := range str1 {
		if !str1Matches[i] {
			continue
		}
		for !str2Matches[k] {
			k++
		}
		if str1[i] != str2[k] {
			transpositions++
		}
		k++
	}
	transpositions /= 2
	return (matches/float64(len(str1)) +
		matches/float64(len(str2)) +
		(matches-transpositions)/matches) / 3
}

// DiceCoefficient is bigram position dependent implementation of the Dice coefficient
func DiceCoefficient(string1, string2 string) float64 {
	// Check for nil or empty string
	if string1 == "" && string2 == "" {
		return 0
	}
	if string1 == string2 {
		return 1
	}

	strlen1 := len(string1) - 1
	strlen2 := len(string2) - 1
	if strlen1 < 1 || strlen2 < 1 {
		return 0
	}

	var matches float64 = 0
	var i, j = 0, 0

	//get bigrams and compare
	for i < strlen1 && j < strlen2 {
		a := string(string1[i] + string1[i+1])
		b := string(string2[j] + string2[j+1])
		if strings.EqualFold(a, b) {
			matches += 2
		}
		i++
		j++
	}
	return matches / float64(strlen1+strlen2)
}

// Join is a quite efficient string "concatenator"
func Join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func GetFirstRune(data string) rune {
	for _, c := range data {
		return c
	}
	return 0
}

func ArrayToMap(slice []string) map[string]struct{} {
	result := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		result[s] = struct{}{}
	}
	return result
}

func TrimStrings(vs []string) []string {
	arrayutils.Apply(&vs, func(i int, s string) string {
		return strings.Trim(s, CUTSET)
	}, true)
	return vs
}

func Pad(word, v string, n int) string {
	var sb strings.Builder
	word = Trim(word)
	if len(word) > n {
		return word[:n]
	}
	for i := 0; i < n-len(word); i++ {
		sb.WriteString(v)
	}
	sb.WriteString(word)
	return sb.String()
}

func ToByte(slice []string, separator string) []byte {
	var sb bytes.Buffer
	for i := range slice {
		sb.WriteString(slice[i])
		sb.WriteString(separator)
	}
	return bytes.TrimSuffix(sb.Bytes(), []byte(separator))
}
