package mathutils

// SumIntArray return the of every element contained in the array
func SumIntArray(integers []int) int {
	sum := 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// MaxInt return the index that contains the max value for the given array
func MaxIntIndex(array []int) int {
	index := 0
	lenght := len(array)
	for i := 1; i < lenght; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}
