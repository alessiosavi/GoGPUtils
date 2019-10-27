package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RandomGenerator is delegated to generate random without call seed every time
type RandomGenerator struct {
	randomizer *rand.Rand
}

// InitRandomizer initialize a new RandomGenerator
func InitRandomizer() RandomGenerator {
	var random RandomGenerator
	random.randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
	return random
}

// Check is a helper function which streamlines error checking
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// RandomInt initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt(min, max int) int {
	return rander.randomizer.Intn(max-min) + min
}

// RandomInt32 initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt32(min, max int32) int32 {
	return rander.randomizer.Int31n(max-min) + min
}

// RandomInt64 initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt64(min, max int64) int64 {
	return rander.randomizer.Int63n(max-min) + min
}

// RandomFloat32 initalizate a new seed using the UNIX Nano time and return a float32 between the 2 input value
func (rander RandomGenerator) RandomFloat32(min, max float32) float32 {
	return min + rander.randomizer.Float32()*(max-min)
}

// RandomFloat64 initalizate a new seed using the UNIX Nano time and return a float64 between the 2 input value
func (rander RandomGenerator) RandomFloat64(min, max float64) float64 {
	return min + rander.randomizer.Float64()*(max-min)
}

// RandomIntArray return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomIntArray(min, max, len int) []int {
	array := make([]int, len)
	for i := 0; i < len; i++ {
		array[i] = rander.RandomInt(min, max)
	}
	return array
}

// RandomInt32Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomInt32Array(min, max int32, len int) []int32 {
	array := make([]int32, len)
	for i := 0; i < len; i++ {
		array[i] = rander.RandomInt32(min, max)
	}
	return array
}

// RandomInt64Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomInt64Array(min, max int64, len int) []int64 {
	array := make([]int64, len)
	for i := 0; i < len; i++ {
		array[i] = rander.RandomInt64(min, max)
	}
	return array
}

// RandomFloat32Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomFloat32Array(min, max float32, len int) []float32 {
	array := make([]float32, len)
	for i := 0; i < len; i++ {
		array[i] = rander.RandomFloat32(min, max)
	}
	return array
}

// RandomFloat64Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomFloat64Array(min, max float64, len int) []float64 {
	array := make([]float64, len)
	for i := 0; i < len; i++ {
		array[i] = rander.RandomFloat64(min, max)
	}
	return array
}

// RandomInt initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandomInt32 initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt32(min, max int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(max-min) + min
}

// RandomInt64 initalizate a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// RandomFloat64 initalizate a new seed using the UNIX Nano time and return a float64 between the 2 input value
func RandomFloat64(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// RandomFloat32 initalizate a new seed using the UNIX Nano time and return a float32 between the 2 input value
func RandomFloat32(min, max float32) float32 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float32()*(max-min)
}

// GenerateSequentialIntArray is delegated to generate an array of sequential number
func GenerateSequentialIntArray(length int) []int {
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = i
	}
	return array
}

// GenerateSequentialFloat32Array is delegated to generate an array of sequential number
func GenerateSequentialFloat32Array(length int) []float32 {
	array := make([]float32, length)
	for i := 0; i < length; i++ {
		array[i] = float32(i)
	}
	return array
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
