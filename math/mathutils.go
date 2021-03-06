package mathutils

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
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
	var index int
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] > array[index] {
			index = i
		}
	}
	return index
}

// MinIntIndex return the index that contains the min value for the given array
func MinIntIndex(array []int) int {
	index := 0
	length := len(array)
	for i := 1; i < length; i++ {
		if array[i] < array[index] {
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

// ModeInt is delegated to calculate the mode of the given array
func ModeInt(array []int) []int {
	// Save the number of occurrence for every number of the array
	var mode = make(map[int]int)
	for i := range array {
		mode[array[i]]++
	}
	var max int
	var maxs []int

	// Avoid to take care about value that does not appear at least 2 time
	max = 2
	for i := 1; i < len(mode); i++ {
		if mode[i] >= max {
			max = mode[i]
		}
	}
	for i := range mode {
		if mode[i] == max {
			maxs = append(maxs, i)
		}
	}
	return maxs
}

// MedianInt is delegated to calculate the median for the given INT array
func MedianInt(arr []int) float64 {
	var array = make([]int, len(arr))
	// Avoid to modify the input array
	copy(array, arr)
	sort.Ints(array)
	if len(array)%2 != 0 {
		index := (len(array)) / 2
		return float64(array[index])
	}
	n1 := array[(len(array)-1)/2]
	n2 := array[(len(array)-1)/2+1]
	return float64(n1+n2) / 2.0
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
	return float64(total) / float64(len(array))
}

// StandardDeviationInt is delegated to calculate the STD for the given array
func StandardDeviationInt(array []int) float64 {
	// 1. Calculate average
	mean := AverageInt(array)
	// 2. Subtract every terms for the average and square the result. Sum every terms
	var sum float64
	for i := range array {
		sum += math.Pow(float64(array[i])-mean, 2)
	}
	// 3. Multiplying by 1/N (divide for N)
	sum /= float64(len(array))

	// 4.  Take the square root
	return math.Sqrt(sum)
}

// StandardDeviationInt is delegated to calculate the STD for the given array
func StandardDeviationFloat64(array []float64) float64 {
	// 1. Calculate average
	mean := AverageFloat64(array)
	// 2. Subtract every terms for the average and square the result. Sum every terms
	var sum float64
	for i := range array {
		sum += math.Pow(array[i]-mean, 2)
	}
	// 3. Multiplying by 1/N (divide for N)
	sum /= float64(len(array))
	// 4.  Take the square root
	return math.Sqrt(sum)
}

func VarianceInt(array []int) float64 {
	// 1. Work out the Mean (the simple average of the numbers)
	// 2. Then for each number: subtract the Mean and square the result (the squared difference).
	// 3. Then work out the average of those squared differences. (Why Square?)
	mean := AverageInt(array)
	var sum float64

	for i := range array {
		sum += math.Pow(float64(array[i])-mean, 2)
	}
	sum /= float64(len(array))
	return math.Sqrt(sum)
}

// CovarianceInt is delegated to calculate the Covariance between the given arrays
func CovarianceInt(arr1, arr2 []int) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}
	// 1. Calculate the mean
	avg1 := AverageInt(arr1)
	avg2 := AverageInt(arr2)

	var sum float64
	for i := range arr1 {
		sum += (float64(arr1[i]) - avg1) * (float64(arr2[i]) - avg2)
	}
	return sum / float64(len(arr1)-1)
}

// CovarianceInt is delegated to calculate the Covariance between the given arrays
func CovarianceFloat64(arr1, arr2 []float64) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}
	// 1. Calculate the mean
	avg1 := AverageFloat64(arr1)
	avg2 := AverageFloat64(arr2)

	var sum float64
	for i := range arr1 {
		sum += (arr1[i] - avg1) * (arr2[i] - avg2)
	}
	return sum / float64(len(arr1)-1)
}

// CorrelationInt is delegated to calculate the correlation for the two given arrays
func CorrelationInt(arr1, arr2 []int) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}

	// 1. Calculate the mean
	avg1 := AverageInt(arr1)
	avg2 := AverageInt(arr2)

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

// CorrelationFloat64 is delegated to calculate the correlation for the two given arrays
func CorrelationFloat64(arr1, arr2 []float64) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		log.Fatal("CovarianceInt | Input array have a different shape: Array1 [", arr1, "], Array2: [", arr2, "]")
	}

	// 1. Calculate the mean
	avg1 := AverageFloat64(arr1)
	avg2 := AverageFloat64(arr2)

	var sum float64
	for i := range arr1 {
		// TODO: Save calculation instead of recalculate at the step below
		sum += (arr1[i] - avg1) * (arr2[i] - avg2)
	}
	var pow1, pow2 float64

	for i := range arr1 {
		pow1 += math.Pow(arr1[i]-avg1, 2)
	}
	for i := range arr1 {
		pow2 += math.Pow(arr2[i]-avg2, 2)
	}
	return sum / math.Sqrt(pow1*pow2)
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

// InitMatrix is delegated to initialize a new empty matrix
func InitMatrix(r, c int) [][]int {
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
func DumpMatrix(m [][]int) string {
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

// InitRandomMatrix is delegated to initialize a random matrix with the given dimension
func InitRandomMatrix(r, c int) [][]int {
	m := InitMatrix(r, c)
	randomizer := helper.InitRandomizer()
	for i := range m {
		m[i] = randomizer.RandomIntArray(0, 100, c)
	}
	return m
}

// InitMatrixCustom is delegated to initialize a matrix with the given dimension using the same value for each field
func InitMatrixCustom(r, c, value int) [][]int {
	m := InitMatrix(r, c)
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

// MultiplyMatrix is delegated to execute the multiplication between the given matrix without extra allocation
func MultiplyMatrix(m1, m2 [][]int) [][]int {
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
	result := InitMatrix(n, m)

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
func MultiplyMatrixLegacy(m1, m2 [][]int) [][]int {
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

	total := InitMatrix(len(m1), len(m2[0]))
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
	return total
}

// MultiplySumArray is delegated to multiply the given array and sum every number of the result array
func MultiplySumArray(a, b []int) int {
	if len(a) != len(b) {
		log.Println("Different length ...")
		return -1
	}
	total := make([]int, len(a))
	for i := range a {
		total[i] = a[i] * b[i]
	}
	return SumIntArray(total)
}

// SumArrays is delegated to sum 2 array of different length.
func SumArrays(n1, n2 []int) []int {
	var (
		result []int
		odd    int
		length int
		sum    int
	)
	if len(n1) > len(n2) {
		length = len(n1)
	} else {
		length = len(n2)
	}
	n1 = PadArray(n1, length)
	n2 = PadArray(n2, length)

	for i := length - 1; i >= 0; i-- {
		sum = n1[i] + n2[i] + odd
		if sum > 9 {
			odd = 1
			sum -= 10
		} else {
			odd = 0
		}
		result = append(result, sum)
	}
	if odd != 0 {
		result = append(result, odd)
	}
	reversed := arrayutils.ReverseArrayInt(result)
	return reversed
}

// CalculateMaxPrimeFactor is delegated to calculate the max prime factor for the input number
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
func GenerateFibonacci(max int64) []int64 {
	var array []int64
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

// GenerateFibonacciN is delegated to generate N Fibonacci number
func GenerateFibonacciN(max int) []float64 {
	var array []float64
	// Hardcoded for enhance for performance
	array = append(array, 1, 1, 2)
	i := 3
	for len(array) < max && array[len(array)-1] <= math.MaxFloat64 {
		array = append(array, array[i-1]+array[i-2])
		i++
	}
	return array
}

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

// FindDivisor is delegated to find every divisor for the input number
func FindDivisor(n int) []int {
	var count int
	var divisor []int
	max := int(math.Sqrt(float64(n)))
	for i := 1; i <= max; i++ {
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
func PadArray(array []int, n int) []int {
	var result []int
	var length = len(array)
	if n != length {
		result = make([]int, n-length)
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
func FindIndexValue(array []int, value int) []int {
	var indexs []int
	for i := range array {
		if array[i] == value {
			indexs = append(indexs, i)
		}
	}
	//log.Println("Found value [", value, "] at index [", indexs, "]")
	return indexs
}

// SortMaxIndex is delegated to return an array that contains the position of the order value (from max to min) of the given array
// {1, 9, 2, 10, 3} -> [3 1 4 2 0] || {7, 6, 5, 4, 3, 2, 1} -> [0 1 2 3 4 5 6] || {1, 2, 3, 4, 5, 6, 7} -> [6 5 4 3 2 1 0]
func SortMaxIndex(array []int) []int {
	var (
		result, additional []int
		index, value       int
		arrayCopy          = make([]int, len(array))
	)
	copy(arrayCopy, array)
	for len(array) > 0 {
		// Retrieve the max index
		index = MaxIntIndex(array)
		value = array[index]
		// Find the value(s)
		additional = FindIndexValue(arrayCopy, value)
		result = append(result, additional...)
		array = arrayutils.RemoveIntByValue(array, value)
	}
	return result
}

// SimilarityPreCheck is delegated to verify that the given array have the correct size
func SimilarityPreCheck(a, b []float64) bool {
	if len(a) == 0 || len(b) == 0 {
		log.Println("CosineSimilarity | Nil input data")
		return false
	}

	if len(a) != len(b) {
		log.Printf("CosineSimilarity | Input vectors have different size")
		return false
	}

	return true
}

// CosineSimilarity is delegated to calculate the Cosine Similarity for the given array
func CosineSimilarity(a, b []float64) float64 {
	if !SimilarityPreCheck(a, b) {
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

// EuclideanDistance is delegated to calculate the euclidean distance for the given array
func EuclideanDistance(v1, v2 []float64) float64 {
	if !SimilarityPreCheck(v1, v2) {
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
	if !SimilarityPreCheck(v1, v2) {
		return -1
	}
	var taxicab float64
	for i := range v1 {
		taxicab += math.Abs(v2[i] - v1[i])
	}
	return taxicab
}

// MaxInt is delegated to return the max int from the two given int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt is delegated to return the min int from the two given int
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxFloat64 is delegated to return the max int from the two given int
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 is delegated to return the min int from the two given int
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// MaxIntMultiple is delegated to return the max value with a variable number of input int
func MaxIntMultiple(a ...int) int {
	return a[MaxIntIndex(a)]
}

// MinIntMultiple is delegated to return the min value with a variable number of input int
func MinIntMultiple(a ...int) int {
	return a[MinIntIndex(a)]
}
