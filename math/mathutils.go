package mathutils

import (
	"fmt"
	"github.com/alessiosavi/GoGPUtils/datastructure/types"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
)

// NewArray is delegated to initialize a new array of the given dimension, populated with the same input value
func NewArray[T types.Number](n int, value T) []T {
	if n <= 0 {
		return nil
	}
	var array = make([]T, n)
	for i := range array {
		array[i] = value
	}
	return array
}

// SumArray return the sum of every element contained in the array
func SumArray[T types.Number](array []T) T {
	var sum T = 0
	for i := range array {
		sum += array[i]
	}
	return sum
}

// SubtractArray subtract of every element contained in the array and return the result
func SubtractArray[T types.Number](integers []T) T {
	var subtract T = 0
	for i := range integers {
		subtract -= integers[i]
	}
	return subtract
}

// SumArrays is delegated to sum the two given array
func SumArrays[T types.Number](a1, a2 []T) []T {
	if a1 == nil || a2 == nil || len(a1) != len(a2) {
		return nil
	}
	total := make([]T, len(a1))
	length := len(a1)
	for i := 0; i < length; i++ {
		total[i] = a1[i] + a2[i]
	}
	return total
}

// SubtractArrays is delegated to sum the two given array
func SubtractArrays[T types.Number](a1, a2 []T) []T {
	if a1 == nil || a2 == nil || len(a1) != len(a2) {
		return nil
	}
	total := make([]T, len(a1))
	length := len(a1)
	for i := 0; i < length; i++ {
		total[i] = a1[i] - a2[i]
	}
	return total
}

// MaxIndex return the index that contains the max value for the given array
func MaxIndex[T types.Number](array []T) int {
	var index int
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MinIndex return the index that contains the min value for the given array
func MinIndex[T types.Number](array []T) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] < array[index] {
			index = i
		}
	}
	return index
}

// Mode is delegated to calculate the mode of the given array
func Mode[T types.Number](array []T) []T {
	// Save the number of occurrence for every number of the array
	var mode = make(map[T]int, len(array))
	for i := range array {
		mode[array[i]]++
	}
	var _max int
	var maxs []T

	// Avoid taking care about value that does not appear at least 2 time
	_max = 2
	for i := T(1); i < T(len(mode)); i++ {
		if mode[i] >= _max {
			_max = mode[i]
		}
	}
	for i := range mode {
		if mode[i] == _max {
			maxs = append(maxs, i)
		}
	}
	return maxs
}

// Median is delegated to calculate the median for the given array
func Median[T types.Number](arr []T) float64 {
	var array = make([]T, len(arr))
	// Avoid modifying the input array
	copy(array, arr)
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	if len(array)%2 != 0 {
		index := (len(array)) / 2
		return float64(array[index])
	}
	n1 := array[(len(array)-1)/2]
	n2 := array[(len(array)-1)/2+1]
	return float64(n1+n2) / 2.0
}

// Average is delegated to calculate the average of an int array
func Average[T types.Number](array []T) float64 {
	var total T
	// Same as len(array) == 0
	if array == nil {
		return 0
	} else if len(array) == 1 {
		return float64(array[0])
	}

	for i := range array {
		total += array[i]
	}
	return float64(total) / float64(len(array))
}

// StandardDeviation is delegated to calculate the STD for the given array
func StandardDeviation[T types.Number](array []T) float64 {
	// 1. Calculate average
	mean := Average(array)
	// 2. Subtract every term for the average and square the result. Sum every terms
	var sum float64
	for i := range array {
		sum += math.Pow(float64(array[i])-mean, 2)
	}
	// 3. Multiplying by 1/N (divide for N)
	sum /= float64(len(array))

	// 4.  Take the square root
	return math.Sqrt(sum)
}

func Variance[T types.Number](array []T) float64 {
	// 1. Work out the Mean (the simple average of the numbers)
	// 2. Then for each number: subtract the Mean and square the result (the squared difference).
	// 3. Then work out the average of those squared differences. (Why Square?)
	mean := Average(array)
	var sum float64

	for i := range array {
		sum += math.Pow(float64(array[i])-mean, 2)
	}
	sum /= float64(len(array))
	return math.Sqrt(sum)
}

// Covariance is delegated to calculate the Covariance between the given arrays
func Covariance[T types.Number](arr1, arr2 []T) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}
	// 1. Calculate the mean
	avg1 := Average(arr1)
	avg2 := Average(arr2)

	var sum float64
	for i := range arr1 {
		sum += (float64(arr1[i]) - avg1) * (float64(arr2[i]) - avg2)
	}
	return sum / float64(len(arr1)-1)
}

// Correlation is delegated to calculate the correlation for the two given arrays
func Correlation[T types.Number](arr1, arr2 []T) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}

	// 1. Calculate the mean
	avg1 := Average(arr1)
	avg2 := Average(arr2)

	var sum float64
	var sum1, sum2 = make([]float64, len(arr1)), make([]float64, len(arr2))
	for i := range arr1 {
		sum1[i] = float64(arr1[i]) - avg1
		sum2[i] = float64(arr2[i]) - avg2
		sum += sum1[i] * sum2[i]
	}
	var pow1, pow2 float64

	for i := range arr1 {
		pow1 += math.Pow(sum1[i], 2)
	}
	for i := range arr1 {
		pow2 += math.Pow(sum2[i], 2)
	}
	return sum / math.Sqrt(pow1*pow2)
}

// InitMatrix is delegated to initialize a new empty matrix
func InitMatrix[T types.Number](r, c int) [][]T {
	if r <= 1 || c <= 1 {
		return nil
	}
	matrix := make([][]T, r)
	for rowsIndex := range matrix {
		matrix[rowsIndex] = make([]T, c)
	}
	return matrix
}

// DumpMatrix is delegated to print the given matrix
func DumpMatrix[T any](m [][]T) string {
	var sb strings.Builder
	if m == nil {
		return ""
	}
	for i := range m {
		sb.WriteString(fmt.Sprintf("%v", m[i]) + "\n")
	}
	sb.WriteString("\nRows: " + strconv.Itoa(len(m)) + " Columns: " + strconv.Itoa(len(m[0])))
	return sb.String()
}

// FIXME
// InitRandomMatrix is delegated to initialize a random matrix with the given dimension
//func InitRandomMatrix[T types.Number](r, c int) [][]T {
//	m := InitMatrix[T](r, c)
//	randomizer := helper.InitRandomizer()
//	for i := range m {
//		m[i] = randomizer.RandomIntArray(0, 100, c)
//	}
//	return m
//}

// InitMatrixCustom is delegated to initialize a matrix with the given dimension using the same value for each field
func InitMatrixCustom[T types.Number](r, c int, value T) [][]T {
	m := InitMatrix[T](r, c)
	for i := range m {
		m[i] = NewArray[T](c, value)
	}
	return m
}

// SumMatrix is delegated to sum the given matrix
func SumMatrix[T types.Number](m1, m2 [][]T) [][]T {
	if m1 == nil || m2 == nil || len(m1) != len(m2) {
		return nil
	}
	sum := make([][]T, len(m1))
	length := len(m1)
	for i := 0; i < length; i++ {
		sum[i] = SumArrays[T](m1[i], m2[i])
	}
	return sum
}

// SubtractMatrix is delegated to sum the given matrix
func SubtractMatrix[T types.Number](m1, m2 [][]T) [][]T {
	if m1 == nil || m2 == nil || len(m1) != len(m2) {
		return nil
	}
	total := make([][]T, len(m1))
	length := len(m1)
	for i := 0; i < length; i++ {
		total[i] = SubtractArrays[T](m1[i], m2[i])
	}
	return total
}

// MultiplyMatrix is delegated to execute the multiplication between the given matrix without extra allocation
func MultiplyMatrix[T types.Number](m1, m2 [][]T) [][]T {
	if m1 == nil || m2 == nil || len(m1) == 0 || len(m2) == 0 {
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

	n := len(m1)
	y := len(m1[0])
	m := len(m2[0])
	result := InitMatrix[T](n, m)

	for k := 0; k < y; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				result[i][j] += +m1[i][k] * m2[k][j]
			}
		}
	}

	return result
}

// MultiplyMatrixLegacy is delegated to execute the multiplication between the given matrix
func MultiplyMatrixLegacy[T types.Number](m1, m2 [][]T) [][]T {
	if m1 == nil || m2 == nil || len(m1) == 0 || len(m2) == 0 {
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

	total := InitMatrix[T](len(m1), len(m2[0]))
	for i := range m1 {
		arrayM1 := m1[i]
		for k := 0; k < len(m2); k++ {
			arrayM2 := make([]T, len(arrayM1))
			for j := range m2 {
				arrayM2[j] = m2[j][k]
			}
			data := MultiplySumArray(arrayM1, arrayM2)
			total[i][k] = data
		}
	}
	return total
}

// MultiplySumArray is delegated to multiply the given array and sum every number of the result array
func MultiplySumArray[T types.Number](a, b []T) T {
	if len(a) != len(b) {
		panic("Different length ...")
	}
	total := make([]T, len(a))
	for i := range a {
		total[i] = a[i] * b[i]
	}
	return SumArray[T](total)
}

// SumArraysPadded is delegated to sum 2 array of different length.
func SumArraysPadded[T types.Number](n1, n2 []T) []T {
	var (
		result []T
		odd    int
		length int
		sum    T
	)
	if len(n1) > len(n2) {
		length = len(n1)
	} else {
		length = len(n2)
	}
	n1 = PadArray(n1, length)
	n2 = PadArray(n2, length)

	for i := length - 1; i >= 0; i-- {
		sum = n1[i] + n2[i] + T(odd)
		if sum > 9 {
			odd = 1
			sum -= 10
		} else {
			odd = 0
		}
		result = append(result, sum)
	}
	if odd != 0 {
		result = append(result, T(odd))
	}
	reversed := arrayutils.ReverseArray[T](result)
	return reversed
}

// CalculateMaxPrimeFactor is delegated to calculate the max prime factor for the given input
func CalculateMaxPrimeFactor(n int64) int64 {
	var maxPrime int64 = -1
	var i int64
	for n%2 == 0 {
		n /= 2
	}

	for i = 3; float64(i) <= math.Sqrt(float64(n)); i += 2 {
		for n%i == 0 {
			n /= i
		}
	}
	if n > 2 {
		maxPrime = n
	}
	return maxPrime
}

// IsPrime is delegated to verify if the given number is Prime
func IsPrime(n int) bool {
	if n <= 3 {
		return n > 1
	} else if n%2 == 0 || n%3 == 0 {
		return false
	}
	i := 5
	mult := float64(2)
	for int(math.Pow(float64(i), mult)) <= n {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}

// GenerateFibonacci is delegated to generate the Fibonacci sequence
func GenerateFibonacci[T types.Number](max T) []T {
	var array []T
	// Hardcoded for enhance for performance
	array = append(array, 1, 1, 2)
	i := 3
	var value = array[i-1] + array[i-2]
	for value < max {
		array = append(array, value)
		i++
		value = array[i-1] + array[i-2]
	}
	return array
}

//// GenerateFibonacciN is delegated to generate N Fibonacci number
//func GenerateFibonacciN[T types.Number](max T) []T {
//	var array []T
//	// Hardcoded for enhance for performance
//	array = append(array, 1, 1, 2)
//	i := 3
//	for T(len(array)) < max && array[len(array)-1] <= T(math.MaxFloat64) {
//		array = append(array, array[i-1]+array[i-2])
//		i++
//	}
//	return array
//}

// ExtractEvenValuedNumber Is delegated to extract only the even number from the input array
func ExtractEvenValuedNumber(array []int64) []int64 {
	var result []int64
	for i := range array {
		if array[i]%2 == 0 {
			result = append(result, array[i])
		}
	}
	return result
}

// FindDivisor is delegated to find every divisor for the inpuT types.Number
func FindDivisor(n int) []int {
	var count int
	var divisor []int
	_max := int(math.Sqrt(float64(n)))
	for i := 1; i <= _max; i++ {
		if n%i == 0 {
			div := n / i
			divisor = append(divisor, div)
			if div != i {
				count += 2
			} else {
				count++
			}
			divisor = append(divisor, i)
		}
	}
	sort.Ints(divisor)
	return divisor
}

// PadArray is delegated to return a new padded array with length n
func PadArray[T types.Number](array []T, n int) []T {
	var result []T
	var length = len(array)
	if n != length {
		result = make([]T, n-length)
		for i := 0; i < n-length; i++ {
			result[i] = 0
		}
		result = append(result, array...)
	} else {
		return array
	}
	//log.Println("Input: ", result, " Output: ", array)
	return result
}

// FindIndexValue is delegated to retrieve the index of the given value into the input array.
func FindIndexValue[T types.Number](array []T, value T) []int {
	var indexes []int
	for i := range array {
		if array[i] == value {
			indexes = append(indexes, i)
		}
	}
	//log.Println("Found value [", value, "] at index [", indexes, "]")
	return indexes
}

// SortMaxIndex is delegated to return an array that contains the position of the order value (from max to min) of the given array
// {1, 9, 2, 10, 3} -> [3 1 4 2 0] || {7, 6, 5, 4, 3, 2, 1} -> [0 1 2 3 4 5 6] || {1, 2, 3, 4, 5, 6, 7} -> [6 5 4 3 2 1 0]
func SortMaxIndex[T types.Number](array []T) []int {
	var (
		result, additional []int
		index              int
		value              T
		arrayCopy          = make([]T, len(array))
	)
	copy(arrayCopy, array)
	for len(array) > 0 {
		// Retrieve the max index
		index = MaxIndex[T](array)
		value = array[index]
		// Find the value(s)
		additional = FindIndexValue[T](arrayCopy, value)
		result = append(result, additional...)
		array = arrayutils.RemoveByValue[T](array, value)
	}
	return result
}

// SimilarityPreCheck is delegated to verify that the given array have the correct size
func SimilarityPreCheck[T types.Number](a, b []T) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	return true
}

// CosineSimilarity is delegated to calculate the Cosine Similarity for the given array
func CosineSimilarity(a, b []float64) float64 {
	if !SimilarityPreCheck[float64](a, b) {
		return -1
	}

	// Calculate numerator
	var numerator float64
	for i := range a {
		numerator += a[i] * b[i]
	}

	// Calculate first term of denominator
	var den1 float64
	for i := range a {
		den1 += math.Pow(a[i], 2)
	}
	den1 = math.Sqrt(den1)

	// Calculate second term of denominator
	var den2 float64
	for i := range b {
		den2 += math.Pow(b[i], 2)
	}
	den2 = math.Sqrt(den2)

	result := numerator / (den1 * den2)
	return result
}

// EuclideanDistance is delegated to calculate the Euclidean distance for the given array
func EuclideanDistance(v1, v2 []float64) float64 {
	if !SimilarityPreCheck[float64](v1, v2) {
		return -1
	}
	var euclidean float64
	for i := range v1 {
		euclidean += math.Pow(v1[i]-v2[i], 2)
	}
	return math.Sqrt(euclidean)
}

// ManhattanDistance is delegated to calculate the Manhattan norm for the given array
func ManhattanDistance(v1, v2 []float64) float64 {
	if !SimilarityPreCheck[float64](v1, v2) {
		return -1
	}
	var taxicab float64
	for i := range v1 {
		taxicab += math.Abs(v2[i] - v1[i])
	}
	return taxicab
}

// Max is delegated to return the max from the two given numbers
func Max[T types.Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min is delegated to return the min from the two given numbers
func Min[T types.Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Maxs is delegated to return the max value with a variable number of input int
func Maxs[T types.Number](a ...T) T {
	return a[MaxIndex(a)]
}

// Mins is delegated to return the min value with a variable number of input int
func Mins[T types.Number](a ...T) T {
	return a[MinIndex(a)]
}
