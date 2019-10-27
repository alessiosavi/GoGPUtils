package mathutils

import (
	"math"
	"strings"
)

// SumIntArray return the of every element contained in the array
func SumIntArray(integers []int) int {
	sum := 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// SumInt32Array return the of every element contained in the array
func SumInt32Array(integers []int32) int32 {
	var sum int32
	sum = 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// SumInt64Array return the of every element contained in the array
func SumInt64Array(integers []int64) int64 {
	var sum int64
	sum = 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// SumFloat32Array return the of every element contained in the array
func SumFloat32Array(integers []float32) float32 {
	var sum float32
	sum = 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// SumFloat64Array return the of every element contained in the array
func SumFloat64Array(integers []float64) float64 {
	var sum float64
	sum = 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// MaxIntIndex return the index that contains the max value for the given array
func MaxIntIndex(array []int) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MaxInt32Index return the index that contains the max value for the given array
func MaxInt32Index(array []int32) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MaxInt64Index return the index that contains the max value for the given array
func MaxInt64Index(array []int64) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MaxFloat32Index return the index that contains the max value for the given array
func MaxFloat32Index(array []float32) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MaxFloat64Index return the index that contains the max value for the given array
func MaxFloat64Index(array []float64) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// AverageInt is delegated to calculate the average of an int array
func AverageInt(array []int) float64 {
	var total int
	// Same as len(array) == 0
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}

	for i := range array {
		total += array[i]
	}
	return float64(total / len(array))
}

// AverageInt32 is delegated to calculate the average of an int array
func AverageInt32(array []int32) float64 {
	var total int32
	// Same as len(array) == 0
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}

	for i := range array {
		total += array[i]
	}
	return float64(total / int32(len(array)))
}

// AverageInt64 is delegated to calculate the average of an int array
func AverageInt64(array []int64) float64 {
	var total int64
	// Same as len(array) == 0
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}

	for i := range array {
		total += array[i]
	}
	return float64(total / int64(len(array)))
}

// AverageFloat32 is delegated to calculate the average of an float32 array
func AverageFloat32(array []float32) float64 {
	var total float32
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}
	for i := range array {
		total += array[i]
	}
	return float64(total / float32(len(array)))
}

// AverageFloat64 is delegated to calculate the average of an float64 array
func AverageFloat64(array []float64) float64 {
	var total float64
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return array[0]
	}

	for i := range array {
		total += array[i]
	}
	return total / float64(len(array))
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
