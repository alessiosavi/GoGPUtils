package stringutils

import (
	"bufio"
	"bytes"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unsafe"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	mathutils "github.com/alessiosavi/GoGPUtils/math"
)

// ExtractTextFromQuery is delegated to retrieve the list of word involved in the query.
// It can be viewed as a tokenzier that use whitespace for delimit the word
func ExtractTextFromQuery(target string, ignore []string) []string {
	var queries []string
	rgxp := regexp.MustCompile(`(\w+)`)
	// Extract the list of word
	for _, item := range rgxp.FindAllString(target, -1) {
		if !CheckPresence(ignore, item) {
			queries = append(queries, item)
		}
	}
	return queries
}

// CheckPresence verify that the given array contains the target string
func CheckPresence(array []string, target string) bool {
	for i := range array {
		if array[i] == target {
			return true
		}
	}
	return false
}

// IsUpper verify that a string does contains only upper character
func IsUpper(str string) bool {
	for i := range str {
		ascii := int(str[i])
		if !(ascii > 64 && ascii < 91) {
			return false
		}
	}
	return true
}

// IsLower verify that a string does contains only lower character
func IsLower(str string) bool {
	for i := range str {
		ascii := int(str[i])
		if !(ascii > 96 && ascii < 123) {
			return false
		}
	}
	return true
}

// ContainsLetter verity that the given string contains, at least, an ASCII alphabet characters
// Note: whitespace is allowed
func ContainsLetter(str string) bool {
	for i := range str {
		if (str[i] >= 'a' && str[i] <= 'z') || (str[i] >= 'A' && str[i] <= 'Z') || str[i] == ' ' {
			return true
		}
	}
	return false
}

// ContainsOnlyLetter verity that the given string contains, only, ASCII alphabet characters
// Note, whitespace is allowed
func ContainsOnlyLetter(str string) bool {
	for i := range str {
		if !((str[i] >= 'a' && str[i] <= 'z') || (str[i] >= 'A' && str[i] <= 'Z') || str[i] == ' ') {
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
		json = arrayutils.JoinStrings([]string{json, `"`, values[i], `":"`, values[i+1], `",`}, "")
	}
	json = strings.TrimSuffix(json, `,`)
	json += `}`
	return json
}

// RemoveDoubleWhiteSpace is delegated to remove the whitespace from the given string
// FIXME: memory unefficient, use 2n size, use RemoveFromString method instead
func RemoveDoubleWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for i := range str {
		if !(str[i] == 32 && (i+1 < len(str) && str[i+1] == 32)) {
			b.WriteRune(rune(str[i]))
		}
	}
	return b.String()
}

// RemoveWhiteSpace is delegated to remove the whitespace from the given string
// FIXME: memory unefficient, use 2n size, use RemoveFromString method instead
func RemoveWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for i := range str {
		if str[i] != 32 {
			b.WriteRune(rune(str[i]))
		}
	}
	return b.String()
}

// IsASCII is delegated to verify if a given string is ASCII compliant
func IsASCII(s string) bool {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

// IsASCIIRune is delegated to verify if the given character is ASCII compliant
func IsASCIIRune(r rune) bool {
	return r < 128
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
func Split(data string) []string {
	var linesList []string
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		linesList = append(linesList, scanner.Text())
	}
	return linesList
}

// CountLines return the number of lines in the given string
func CountLines(str string) int {
	scanner := bufio.NewScanner(strings.NewReader(str)) // Create a scanner for iterate the string
	counter := 0
	for scanner.Scan() {
		counter++
	}
	return counter
}

//ExtractString is delegated to filter the content of the given data delimited by 'first' and 'last' string
func ExtractString(data *string, first, last string) string {
	// Find the first instance of 'start' in the give string data
	startHeder := strings.Index(*data, first)
	if startHeder != -1 { // Found !
		startHeder += len(first) // Remove the first word
		// Check the first occurrence of 'last' that delimit the string to return
		endHeader := strings.Index((*data)[startHeder:], last)
		// Ok, seems good, return the content of the string delimited by 'first' and 'last'
		if endHeader != -1 {
			return (*data)[startHeder : startHeder+endHeader]
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
	for _, c := range str {
		if IsASCIIRune(c) {
			b.WriteRune(c)
		}
	}
	return RemoveDoubleWhiteSpace(b.String())
}

// IsBlank is delegated to verify that the does not contains only empty char
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
// It also convert every new line in a space
func Trim(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	length := len(str)
	for i := 0; i < length; i++ {
		// Write character
		if str[i] > 32 {
			b.WriteByte(str[i])
			// Convert new line in space
		} else if str[i] == 10 {
			b.WriteByte(32)
			// Print the space only if followed by an ASCII character
		} else if i+1 < length && (str[i] < 33 && str[i+1] > 32) {
			b.WriteByte(str[i])
		}
	}
	var data string = b.String()
	length = len(data)
	if data[0] == 32 {
		data = data[1:]
		length--
	}
	if data[length-1] == 32 {
		data = data[:length-1]
	}
	return data
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString is delegated to create a random string with whitespace included as fast as possible
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i > -1; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
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
	var builder strings.Builder
	for i := length; i >= 0; i-- {
		builder.WriteByte(str[i])
	}
	return builder.String()
}

// LevenshteinDistanceLegacy is delegated to calculate the Levenshtein distance for the given string
func LevenshteinDistanceLegacy(str1, str2 string) int {
	d := make([][]int, len(str1)+1)
	for i := range d {
		d[i] = make([]int, len(str2)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(str2); j++ {
		for i := 1; i <= len(str1); i++ {
			if str1[i-1] == str2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}
	}
	return d[len(str1)][len(str2)]
}

// LevenshteinDistance is an optimized version for calculate the levenstein distance
func LevenshteinDistance(str1, str2 string) int {
	var n, m int = len(str1), len(str2)
	if n == 0 {
		return m
	} else if m == 0 {
		return n
	}

	var p []int = make([]int, n+1)
	var d []int = make([]int, n+1)
	var _d []int

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
			d[i] = mathutils.MinInt(mathutils.MinInt(d[i-1]+1, p[i]+1), p[i-1]+cost)
		}

		// copy current distance counts to 'previous row' distance counts
		_d = p
		p = d
		d = _d
	}
	return p[n]
}

// JaroDistance is delegated to calculate the Jaro distance from the two given string
func JaroDistance(str1, str2 string) float64 {
	if str1 == str2 {
		return 1
	}

	if len(str1) == 0 || len(str2) == 0 {
		return 0
	}

	match_distance := len(str1)
	if len(str2) > match_distance {
		match_distance = len(str2)
	}
	match_distance = match_distance/2 - 1
	str1_matches := make([]bool, len(str1))
	str2_matches := make([]bool, len(str2))
	matches := 0.
	transpositions := 0.
	for i := range str1 {
		start := i - match_distance
		if start < 0 {
			start = 0
		}
		end := i + match_distance + 1
		if end > len(str2) {
			end = len(str2)
		}
		for k := start; k < end; k++ {
			if str2_matches[k] {
				continue
			}
			if str1[i] != str2[k] {
				continue
			}
			str1_matches[i] = true
			str2_matches[k] = true
			matches++
			break
		}
	}
	if matches == 0 {
		return 0
	}
	k := 0
	for i := range str1 {
		if !str1_matches[i] {
			continue
		}
		for !str2_matches[k] {
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
	if len(string1) == 0 && len(string2) == 0 {
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
	var i, j int = 0, 0

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
