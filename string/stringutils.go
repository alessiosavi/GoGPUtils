package stringutils

import (
	"bufio"
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
		if ascii < 65 || ascii > 90 {
			return false
		}
	}
	return true
}

// IsLower verify that a string does contains only lower char
func IsLower(str string) bool {
	for i := range str {
		ascii := int(str[i])
		if ascii < 97 || ascii > 122 {
			return false
		}
	}
	return true
}

// ContainsLetter verity that the given string contains, at least, an ASCII character
func ContainsLetter(str string) bool {
	for _, charVariable := range str {
		if (charVariable >= 'a' && charVariable <= 'z') || (charVariable >= 'A' && charVariable <= 'Z') {
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
	for _, str := range strs {
		sb.WriteString(str)
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
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

// RemoveFromString Remove a given element from a string
func RemoveFromString(s []byte, i int) []byte {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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
