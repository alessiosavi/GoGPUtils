package mathutils

// SumIntArray return the of every element contained in the array
func SumIntArray(integers []int) int {
	sum := 0
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

// AverageFloat64 is delegated to calculate the average of an float64 array
func AverageFloat64(array []float64) float64 {
	var total float64
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}

	for i := range array {
		total += array[i]
	}
	return total / float64(len(array))
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
