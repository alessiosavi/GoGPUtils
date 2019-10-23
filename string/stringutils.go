package stringutils

import (
	"bufio"
	"bytes"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

// ExtractTextFromQuery is delegated to retrieve the list of string involved in the query
func ExtractTextFromQuery(target string, ignore []string) []string {
	var queries []string
	rgxp := regexp.MustCompile(`(\w+)`)
	// Extract the list of word
	for _, item := range rgxp.FindAllString(target, -1) {
		if !CheckPresence(item, ignore) {
			queries = append(queries, item)
		}
	}
	return queries
}

// CheckPresence verify that the given array contains the target string
func CheckPresence(target string, array []string) bool {
	for i := range array {
		if array[i] == target {
			return true
		}
	}
	return false
}

// IsUpper verify that a string does contains only upper char
func IsUpper(str string) bool {
	for i := range str {
		ascii := int(str[i])
		if !(ascii > 64 && ascii < 91) {
			return false
		}
	}
	return true
}

// IsLower verify that a string does contains only lower char
func IsLower(str string) bool {
	for i := range str {
		ascii := int(str[i])
		if !(ascii > 96 && ascii < 123) {
			return false
		}
	}
	return true
}

// ContainsLetter verity that the given string contains, at least, an ASCII character
func ContainsLetter(str string) bool {
	for i := range str {
		if (str[i] >= 'a' && str[i] <= 'z') || (str[i] >= 'A' && str[i] <= 'Z') {
			return true
		}
	}
	return false
}

// CreateJSON is delegated to create a json object for the key pair in input
func CreateJSON(values ...string) string {
	json := `{`
	lenght := len(values)
	if lenght%2 != 0 {
		return ""
	}
	for i := 0; i < lenght; i += 2 {
		json = Join(json, `"`, values[i], `":"`, values[i+1], `",`)
	}
	json = strings.TrimSuffix(json, `,`)
	json += `}`
	return json
}

// Join is a quite efficient string concatenator
func Join(strs ...string) string {
	var sb strings.Builder
	for i := range strs {
		sb.WriteString(strs[i])
	}
	return sb.String()
}

// RemoveWhiteSpaceString is delegated to remove the whitespace from the given string
// FIXME: unefficient, use 2n size, use RemoveFromString method instead
func RemoveWhiteSpaceString(str string) string {
	var b strings.Builder
	defer b.Reset()
	b.Grow(len(str))
	for i := range str {
		if !(str[i] == 32 && (i+1 < len(str) && str[i+1] == 32)) {
			b.WriteRune(rune(str[i]))
		}
	}
	return b.String()
}

// IsASCII is delegated to verify if a given string is ASCII compliant
func IsASCII(s string) bool {
	for i := range s {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

// IsASCIIRune is delegated to verify if the given character is ASCII compliant
func IsASCIIRune(r rune) bool {
	return !(r > 127)
}

// RemoveFromString Remove a given element from a string
func RemoveFromString(data string, i int) string {
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

// CountLinesString return the number of lines in the given string
func CountLinesString(fileContet string) int {
	scanner := bufio.NewScanner(strings.NewReader(fileContet)) // Create a scanner for iterate the string
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
func ReplaceAtIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
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
	return RemoveWhiteSpaceString(b.String())
}

func IsBlank(str string) bool {
	// Check lenght
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

// Trim is delegated to remove the initial and final whitespace and the double whitespace
func Trim(str string) string {
	var b strings.Builder
	lenght := len(str)
	for i := 0; i < lenght; i++ {
		// Convert double space as a single space
		if str[i] > 32 {
			b.WriteByte(str[i])
		} else if i+1 < lenght && (str[i] < 33 && str[i+1] > 32) {
			b.WriteByte(str[i])
		}
	}
	var data string
	data = b.String()
	lenght = len(data)
	if str[0] == 32 {
		data = data[1:]
		lenght--
	}
	if data[lenght-1] == 32 {
		data = data[:lenght-2]
	}
	return data
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString is delegated to create a random string with whitespace included as fast as possible
func RandomString(n int) string {
	b := make([]byte, n)
	var src = rand.NewSource(time.Now().UnixNano())
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
