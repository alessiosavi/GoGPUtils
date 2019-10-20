package stringutils

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
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

// RemoveFromByte Remove a given element from a string
// NOTE: Panic in case of index out of bound
func RemoveFromByte(s []byte, i int) []byte {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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

// CountLines return the number of lines in the given string
func CountLines(fileContet string) int {
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

// UpperizeString is delegated to upperize the case of a lower case character
func UpperizeString(str *string) string {
	for i, c := range *str {
		ascii := int(c)
		if !(ascii > 64 && ascii < 91) {
			*str = ReplaceAtIndex(*str, rune(ascii-32), i)
		}
	}
	return *str
}

// LowerizeString is delegated to lowerize the case of an upper case character
func LowerizeString(str *string) string {
	for i, c := range *str {
		ascii := int(c)
		if !(ascii > 96 && ascii < 123) {
			*str = ReplaceAtIndex(*str, rune(ascii+32), i)
		}
	}
	return *str
}

// RemoveNonAscii is delegated to clean the text from the NON ASCII character
func RemoveNonAscii(str string) string {
	var b bytes.Buffer
	b.Grow(len(str))
	for _, c := range str {
		if IsASCIIRune(c) {
			b.WriteRune(c)
		}
	}
	return RemoveWhiteSpaceString(b.String())
}
