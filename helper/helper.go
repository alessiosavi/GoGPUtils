package helper

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

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

// RandomInt initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt(min, max int) int {
	return rander.randomizer.Intn(max-min) + min
}

// RandomInt32 initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt32(min, max int32) int32 {
	return rander.randomizer.Int31n(max-min) + min
}

// RandomInt64 initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func (rander RandomGenerator) RandomInt64(min, max int64) int64 {
	return rander.randomizer.Int63n(max-min) + min
}

// RandomFloat32 initialize a new seed using the UNIX Nano time and return a float32 between the 2 input value
func (rander RandomGenerator) RandomFloat32(min, max float32) float32 {
	return min + rander.randomizer.Float32()*(max-min)
}

// RandomFloat64 initialize a new seed using the UNIX Nano time and return a float64 between the 2 input value
func (rander RandomGenerator) RandomFloat64(min, max float64) float64 {
	return min + rander.randomizer.Float64()*(max-min)
}

// RandomIntArray return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomIntArray(min, max, length int) []int {
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = rander.RandomInt(min, max)
	}
	return array
}

// RandomInt32Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomInt32Array(min, max int32, length int) []int32 {
	array := make([]int32, length)
	for i := 0; i < length; i++ {
		array[i] = rander.RandomInt32(min, max)
	}
	return array
}

// RandomInt64Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomInt64Array(min, max int64, length int) []int64 {
	array := make([]int64, length)
	for i := 0; i < length; i++ {
		array[i] = rander.RandomInt64(min, max)
	}
	return array
}

// RandomFloat32Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomFloat32Array(min, max float32, length int) []float32 {
	array := make([]float32, length)
	for i := 0; i < length; i++ {
		array[i] = rander.RandomFloat32(min, max)
	}
	return array
}

// RandomFloat64Array return a new array with random number from min to max of length len
func (rander RandomGenerator) RandomFloat64Array(min, max float64, length int) []float64 {
	array := make([]float64, length)
	for i := 0; i < length; i++ {
		array[i] = rander.RandomFloat64(min, max)
	}
	return array
}

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	LETTER_IDX_BYTES = 6                       // 6 bits to represent a letter index
	LETTER_IDX_MASK  = 1<<LETTER_IDX_BYTES - 1 // All 1-bits, as many as LETTER_IDX_BYTES
	LETTER_IDX_MAX   = 63 / LETTER_IDX_BYTES   // # of letter indices fitting in 63 bits
)

func (rander RandomGenerator) RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rander.randomizer.Int63(), LETTER_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = rander.randomizer.Int63(), LETTER_IDX_MAX
		}
		if idx := int(cache & LETTER_IDX_MASK); idx < len(LETTER_BYTES) {
			b[i] = LETTER_BYTES[idx]
			i--
		}
		cache >>= LETTER_IDX_BYTES
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// RandomByte is delegated to generate a byte array with the given input length
func RandomByte(length int) ([]byte, error) {
	data := make([]byte, length)
	_, err := crand.Read(data)
	return data, err
}

// RandomInt initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomInt32 initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt32(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}

// RandomInt64 initialize a new seed using the UNIX Nano time and return an integer between the 2 input value
func RandomInt64(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// RandomFloat64 initialize a new seed using the UNIX Nano time and return a float64 between the 2 input value
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomFloat32 initialize a new seed using the UNIX Nano time and return a float32 between the 2 input value
func RandomFloat32(min, max float32) float32 {
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

// ConvertSize is delegated to return the dimension related to the input byte size
func ConvertSize(bytes float64, dimension string) float64 {
	var value float64
	dimension = strings.ToUpper(dimension)
	switch dimension {
	case "KB", "KILOBYTE":
		value = bytes / 1000
	case "MB", "MEGABYTE":
		value = bytes / math.Pow(1000, 2)
	case "GB", "GIGABYTE":
		value = bytes / math.Pow(1000, 3)
	case "TB", "TERABYTE":
		value = bytes / math.Pow(1000, 4)
	case "PB", "PETABYTE":
		value = bytes / math.Pow(1000, 5)
	case "XB", "EXABYTE":
		value = bytes / math.Pow(1000, 6)
	case "ZB", "ZETTABYTE":
		value = bytes / math.Pow(1000, 7)
	}
	return value
}

func Marshal(data interface{}) string {
	if indent, err := json.Marshal(data); err == nil {
		return string(indent)
	}
	return fmt.Sprintf("%+v\n", data)

}

func MarshalIndent(data interface{}) string {
	if indent, err := json.MarshalIndent(data, " ", "  "); err == nil {
		return string(indent)
	}
	return fmt.Sprintf("%+v\n", data)

}
