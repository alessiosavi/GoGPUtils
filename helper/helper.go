package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Check is a helper function which streamlines error checking
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Random initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandomFloat initalizate a new seed using the UNIX Nano time and return a float64 between the 2 input value
func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// ByteCountSI convert the byte in input to MB/KB/TB ecc
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountIEC convert the byte in input to MB/KB/TB ecc
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

// RecognizeFormat is delegated to valutate the extension and return the properly Mimetype by a given format type
// reurn: (Mimetype http compliant,Content-Disposition header value)
func RecognizeFormat(input string) (string, string) {
	// Find the last occurrence of the dot
	// extract only the extension of the file by slicing the string
	var contentDisposition string
	var mimeType string
	contentDisposition = `inline; filename="` + input + `"`
	switch input[strings.LastIndex(input, ".")+1:] {
	case "doc":
		mimeType = "application/msword"
	case "docx":
		mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "pdf":
		mimeType = "application/pdf"
	default:
		mimeType = "application/octet-stream"
		contentDisposition = `inline; filename="` + input + `"`
	}
	return mimeType, contentDisposition
}
