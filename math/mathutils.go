package mathutils

import (
	"log"
	"strconv"

	"github.com/alessiosavi/GoGPUtils/helper"
)

// InitIntArray is delegated to initialize a new array of the given dimension, populated with the same input value
func InitIntArray(dimension, value int) []int {
	if dimension <= 0 {
		return nil
	}

	array := make([]int, dimension)
	for i := 0; i < dimension; i++ {
		array[i] = value
	}
	return array
}

// SumIntArray return the sum of every element contained in the array
func SumIntArray(integers []int) int {
	sum := 0
	for i := range integers {
		sum += integers[i]
	}
	return sum
}

// SubtractIntArray return the subtract of every element contained in the array
func SubtractIntArray(integers []int) int {
	subtract := 0
	for i := range integers {
		subtract -= integers[i]
	}
	return subtract
}

// SumIntArrays is delegated to sum the two given array
func SumIntArrays(a1, a2 []int) []int {
	if a1 == nil || a2 == nil || len(a1) != len(a2) {
		return nil
	}
	total := make([]int, len(a1))
	length := len(a1)
	for i := 0; i < length; i++ {
		total[i] = a1[i] + a2[i]
	}
	return total
}

// SubtractIntArrays is delegated to sum the two given array
func SubtractIntArrays(a1, a2 []int) []int {
	if a1 == nil || a2 == nil || len(a1) != len(a2) {
		return nil
	}
	total := make([]int, len(a1))
	length := len(a1)
	for i := 0; i < length; i++ {
		total[i] = a1[i] - a2[i]
	}
	return total
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

// CreateEmptyMatrix is delegated to initialize a new empty matrix
func CreateEmptyMatrix(r, c int) [][]int {
	if r <= 1 || c <= 1 {
		return nil
	}
	matrix := make([][]int, r)
	for rowsIndex := range matrix {
		matrix[rowsIndex] = make([]int, c)
	}
	return matrix
}

// DumpMatrix is delegated to print the given matrix
func DumpMatrix(m [][]int) {
	if m == nil {
		return
	}
	for i := range m {
		log.Println(m[i])
	}
	log.Println("Rows: " + strconv.Itoa(len(m)) + " Columns: " + strconv.Itoa(len(m[0])))
}

// InitRandomMatrix is delegated to initialize a random matrix with the given dimension
func InitRandomMatrix(r, c int) [][]int {
	m := CreateEmptyMatrix(r, c)
	randomizer := helper.InitRandomizer()
	for i := range m {
		m[i] = randomizer.RandomIntArray(0, 100, c)
	}
	return m
}

// InitStaticMatrix is delegated to initialize a matrix with the given dimension using the same value for each field
func InitStaticMatrix(r, c, value int) [][]int {
	m := CreateEmptyMatrix(r, c)
	for i := range m {
		m[i] = InitIntArray(c, value)
	}
	return m
}

// SumMatrix is delegated to sum the given matrix
func SumMatrix(m1, m2 [][]int) [][]int {
	if m1 == nil || m2 == nil || len(m1) != len(m2) {
		return nil
	}
	sum := make([][]int, len(m1))
	length := len(m1)
	for i := 0; i < length; i++ {
		sum[i] = SumIntArrays(m1[i], m2[i])
	}
	return sum
}

// SubtractMatrix is delegated to sum the given matrix
func SubtractMatrix(m1, m2 [][]int) [][]int {
	if m1 == nil || m2 == nil || len(m1) != len(m2) {
		return nil
	}
	total := make([][]int, len(m1))
	length := len(m1)
	for i := 0; i < length; i++ {
		total[i] = SubtractIntArrays(m1[i], m2[i])
	}
	return total
}

// MultiplyMatrix is delegated to execute the multiplication between the given matrix
func MultiplyMatrix(m1, m2 [][]int) [][]int {
	if m1 == nil || m2 == nil {
		return nil
	}

	if len(m1) == 0 || len(m2) == 0 {
		log.Println("Matrix empty")
		return nil
	}
	if len(m1[0]) != len(m2) {
		log.Println("Different size\nM1:")
		DumpMatrix(m1)
		log.Println("M2:")
		DumpMatrix(m2)
		return nil
	}
	total := InitStaticMatrix(len(m1), len(m2[0]), -1)
	for i := range m1 {
		arrayM1 := m1[i]
		for k := 0; k < len(m2); k++ {
			arrayM2 := make([]int, len(arrayM1))
			for j := range m2 {
				arrayM2[j] = m2[j][k]
			}
			data := MultiplySumArray(arrayM1, arrayM2)
			total[i][k] = data
		}
	}
	//DumpMatrix(total)
	return total
}

// MultiplySumArray is delegated to multiply the given array and sum every number of the result array
func MultiplySumArray(a, b []int) int {
	total := make([]int, len(a))
	for i := range a {
		total[i] = a[i] * b[i]
	}
	return SumIntArray(total)
}
