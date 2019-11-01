package byteutils

// IsUpperByte verify that a string does contains only upper character (byte version)
func IsUpperByte(str []byte) bool {
	for i := range str {
		ascii := str[i]
		if !(ascii > 64 && ascii < 91) {
			return false
		}
	}
	return true
}

// IsLowerByte verify that a string does contains only lower character (byte version)
func IsLowerByte(str []byte) bool {
	for i := range str {
		ascii := str[i]
		if !(ascii > 96 && ascii < 123) {
			return false
		}
	}
	return true
}
