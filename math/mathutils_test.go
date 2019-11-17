package mathutils

import (
	"math"
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const total int = 1000

var prime []int = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199}

var randomizer helper.RandomGenerator = helper.InitRandomizer()

func TestInitInitArray(t *testing.T) {
	array := InitIntArray(10, 1)
	if SumIntArray(array) != 10 {
		t.Fail()
	}
}

func BenchmarkSumIntArray(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumIntArray(array)
	}
}

func BenchmarkSumInt32Array(t *testing.B) {
	array := randomizer.RandomInt32Array(0, math.MaxInt32, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumInt32Array(array)
	}
}

func BenchmarkSumInt64Array(t *testing.B) {
	array := randomizer.RandomInt64Array(0, math.MaxInt64, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumInt64Array(array)
	}
}

func BenchmarkSumFloat32Array(t *testing.B) {
	array := randomizer.RandomFloat32Array(0, math.MaxInt64, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumFloat32Array(array)
	}
}

func BenchmarkSumFloat64Array(t *testing.B) {
	array := randomizer.RandomFloat64Array(0, math.MaxInt64, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumFloat64Array(array)
	}
}

func BenchmarkMaxIntIndex(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxIntIndex(array)
	}
}

func BenchmarkMaxInt32Index(t *testing.B) {
	array := randomizer.RandomInt32Array(0, math.MaxInt32, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxInt32Index(array)
	}
}

func BenchmarkMaxInt64Index(t *testing.B) {
	array := randomizer.RandomInt64Array(0, math.MaxInt64, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxInt64Index(array)
	}
}

func BenchmarkMaxFloat32Index(t *testing.B) {
	array := randomizer.RandomFloat32Array(0, math.MaxFloat32, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxFloat32Index(array)
	}
}

func BenchmarkMaxFloat64Index(t *testing.B) {
	array := randomizer.RandomFloat64Array(0, math.MaxFloat64, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxFloat64Index(array)
	}
}

func BenchmarkAverageInt(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		AverageInt(array)
	}
}

func BenchmarkAverageInt32(t *testing.B) {
	array := randomizer.RandomInt32Array(0, math.MaxInt32, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		AverageInt32(array)
	}
}

func BenchmarkAverageInt64(t *testing.B) {
	array := randomizer.RandomInt64Array(0, math.MaxInt64, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		AverageInt64(array)
	}
}

func BenchmarkAverageFloat32(t *testing.B) {
	array := randomizer.RandomFloat32Array(0, math.MaxFloat32, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		AverageFloat32(array)
	}
}
func BenchmarkAverageFloat64(t *testing.B) {
	array := randomizer.RandomFloat64Array(0, math.MaxFloat64, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		AverageFloat64(array)
	}
}

func TestCreateEmptyMatrix(t *testing.T) {
	/*m := */ CreateEmptyMatrix(5, 10)
	//DumpMatrix(m)
}

func TestInitRandomMatrix(t *testing.T) {
	/*m := */ InitRandomMatrix(5, 10)
	//	DumpMatrix(m)
}

func TestInitStaticMatrix(t *testing.T) {
	/*m :=*/ InitStaticMatrix(5, 10, 1)
	//DumpMatrix(m)
}
func BenchmarkInitRandomMatrix(t *testing.B) {
	for i := 0; i < t.N; i++ {
		InitRandomMatrix(5, 10)
	}
}

func TestSumMatrix(t *testing.T) {
	m1 := InitStaticMatrix(5, 10, 1)
	m2 := InitStaticMatrix(5, 10, 1)
	m3 := SumMatrix(m1, m2)
	t.Log(DumpMatrix(m3))
}

func TestMultiplyMatrixLegacy(t *testing.T) {
	m1 := generateTestMatrix1()
	m2 := generateTestMatrix2()
	m3 := MultiplyMatrixLegacy(m1, m2)
	t.Log(DumpMatrix(m3))
}

func TestMultiplyMatrix(t *testing.T) {
	m1 := generateTestMatrix1()
	m2 := generateTestMatrix2()
	m3 := MultiplyMatrix(m1, m2)
	t.Log(DumpMatrix(m3))
}

func TestMultiplySumArray1000(t *testing.T) {
	randomizer := helper.InitRandomizer()
	data := randomizer.RandomIntArray(0, 100, 1000)
	t.Log(MultiplySumArray(data, data))
}
func BenchmarkMultiplySumArray1000(t *testing.B) {
	randomizer := helper.InitRandomizer()
	data := randomizer.RandomIntArray(0, 100, 1000)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		MultiplySumArray(data, data)
	}
}

func BenchmarkMultiplyMatrixLegacy100x100(t *testing.B) {
	m1 := InitRandomMatrix(100, 100)
	m2 := InitRandomMatrix(100, 100)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		MultiplyMatrixLegacy(m1, m2)
	}
}

func BenchmarkMultiplyMatrix100x100(t *testing.B) {
	m1 := InitRandomMatrix(100, 100)
	m2 := InitRandomMatrix(100, 100)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		MultiplyMatrix(m1, m2)
	}
}

// generateTestMatrix1 is delegated to generate a matrix for test purpouse
func generateTestMatrix1() [][]int {
	matrix := make([][]int, 2)
	matrix[0] = []int{1, 2, 3}
	matrix[1] = []int{4, 5, 6}
	return matrix
}

// generateTestMatrix2 is delegated to generate a matrix for test purpouse
func generateTestMatrix2() [][]int {
	matrix := make([][]int, 3)
	matrix[0] = []int{1, 4, 7}
	matrix[1] = []int{2, 5, 8}
	matrix[2] = []int{3, 6, 9}
	return matrix
}

func TestIsPrime(t *testing.T) {
	for _, item := range prime {
		if !IsPrime(item) {
			t.Fail()
		}
	}
}

func BenchmarkIsPrime(t *testing.B) {
	for i := 0; i < t.N; i++ {
		for _, item := range prime {
			if !IsPrime(item) {
				t.Fail()
			}
		}
	}
}

func TestPadArray(t *testing.T) {
	array := []int{1, 2, 3, 4}
	t.Log(PadArray(array, 5))
}

func TestSumArrays(t *testing.T) {
	array1 := []int{1, 1, 2, 3, 4}
	array2 := []int{9, 3, 3, 3}
	// 10567
	t.Log(SumArrays(array1, array2))
}
func TestSortMaxIndex(t *testing.T) {
	var array []int = []int{1, 2, 3, 4, 5, 6, 7}
	// var array []int = []int{7, 6, 5, 4, 3, 2, 1}
	//var array []int = []int{1, 9, 2, 10, 3}
	t.Log(SortMaxIndex(array))
}

func TestCosineSimilarity(t *testing.T) {
	a := []float64{2, 0, 1, 1, 0, 2, 1, 1}
	b := []float64{2, 1, 1, 0, 1, 1, 1, 1}
	similarity := CosineSimilarity(a, b)
	if !(similarity < 0.822 && similarity > 0.821) {
		t.Fail()
	}
}

func BenchmarkCosineSimilarity(t *testing.B) {
	a := []float64{2, 0, 1, 1, 0, 2, 1, 1}
	b := []float64{2, 1, 1, 0, 1, 1, 1, 1}
	for i := 0; i < t.N; i++ {
		CosineSimilarity(a, b)
	}
}

func TestManhattanDistance(t *testing.T) {
	var x []float64 = []float64{1, 2, 3}
	var y []float64 = []float64{2, 4, 6}
	taxiNorm := ManhattanDistance(x, y)
	if taxiNorm != 6 {
		t.Error(taxiNorm)
	}
}

func TestEuclideanDistance(t *testing.T) {
	var x []float64 = []float64{1, 2, 3}
	var y []float64 = []float64{2, 4, 6}
	distance := EuclideanDistance(x, y)
	if !(distance > 3.741 && distance < 3.742) {
		t.Error(distance)
	}
}

func TestModeInt(t *testing.T) {
	noMode := []int{0, 1, 2, 3, 4, 5, 6}
	oneMode := []int{0, 1, 2, 3, 4, 5, 6, 6}
	twoMode := []int{0, 1, 2, 3, 4, 5, 5, 6, 6}
	m1 := ModeInt(noMode)
	t.Log(m1)
	if len(m1) > 0 {
		t.Error("Err", m1)
	}
	m2 := ModeInt(oneMode)
	t.Log(m2)
	if len(m2) > 1 {
		t.Error("Err", m2)
	}
	m3 := ModeInt(twoMode)
	t.Log(m3)
	if len(m3) > 2 {
		t.Error("Err", m3)
	}
}

func TestMedianInt(t *testing.T) {
	median := []int{6, 5, 4, 3, 2, 1, 0}
	m1 := MedianInt(median)
	if m1 != 3 {
		t.Error(m1)
	}
	median = append(median, 7)
	m1 = MedianInt(median)
	if m1 != 3.5 {
		t.Error(m1)
	}
}

func TestStandardDeviationInt(t *testing.T) {
	median := []int{9, 2, 5, 4, 12, 7, 8, 11, 9, 3, 7, 4, 12, 5, 4, 10, 9, 6, 9, 4}
	m1 := StandardDeviationInt(median)
	if !(m1 > 2.9832 && m1 < 2.9833) {
		t.Error(m1)
	}
}

func TestVarianceInt(t *testing.T) {
	median := []int{600, 470, 170, 430, 300}
	m1 := VarianceInt(median)
	if !(m1 > 147.322 && m1 < 147.323) {
		t.Error(m1)
	}
}

func TestCovarianceInt(t *testing.T) {
	arr1 := []int{1692, 1978, 1884, 2151, 2519}
	arr2 := []int{68, 102, 110, 112, 154}
	cv := CovarianceInt(arr1, arr2)
	if cv != 9107.3 {
		t.Error(cv)
	}
}

func TestCorrelationInt(t *testing.T) {
	arr1 := []int{1692, 1978, 1884, 2151, 2519}
	arr2 := []int{68, 102, 110, 112, 154}
	cv := CorrelationInt(arr1, arr2)
	if !(cv > 0.949 && cv < 0.950) {
		t.Error(cv)
	}
}

func TestCorrelationFloat64(t *testing.T) {
	arr1 := []float64{1691.75, 1977.80, 1884.09, 2151.13, 2519.36}
	arr2 := []float64{68.96, 100.11, 109.06, 112.18, 154.12}
	cv := CorrelationFloat64(arr1, arr2)
	if !(cv > 0.954 && cv < 0.955) {
		t.Error(cv)
	}
}
