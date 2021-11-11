package byteutils

// IsUpperByte verify that a string contains only upper character (byte version)
func IsUpperByte(str []byte) bool {
	for i := range str {
		ascii := str[i]
		if !(ascii > 64 && ascii < 91) {
			return false
		}
	}
	return true
}

// IsLowerByte verify that a string contains only lower character (byte version)
func IsLowerByte(str []byte) bool {
	for i := range str {
		ascii := str[i]
		if !(ascii > 96 && ascii < 123) {
			return false
		}
	}
	return true
}

// RemoveFromByte Remove a given element from a string
func RemoveFromByte(s []byte, i int) []byte {
	if i > len(s) {
		return s
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
