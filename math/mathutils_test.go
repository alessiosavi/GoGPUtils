package mathutils

import (
	"github.com/alessiosavi/GoGPUtils/datastructure/types"
	"math"
	"reflect"
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const total int = 1000

var prime = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199}

var randomizer = helper.InitRandomizer()

func BenchmarkSumArray(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		SumArray[int](array)
	}
}

func BenchmarkMaxIndex(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		MaxIndex[int](array)
	}
}

func BenchmarkAverage(t *testing.B) {
	array := randomizer.RandomIntArray(0, math.MaxInt8, total)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		Average[int](array)
	}
}

func TestCreateEmptyMatrix(t *testing.T) {
	m := InitMatrix[int](5, 10)
	DumpMatrix(m)
}

//func TestInitRandomMatrix(t *testing.T) {
//	/*m := */ InitRandomMatrix(5, 10)
//	//	DumpMatrix(m)
//}

func TestInitStaticMatrix(t *testing.T) {
	/*m :=*/ InitMatrixCustom[int](5, 10, 1)
	//DumpMatrix(m)
}

//func BenchmarkInitRandomMatrix(t *testing.B) {
//	for i := 0; i < t.N; i++ {
//		InitRandomMatrix(5, 10)
//	}
//}

func TestGenerateFibonacci(t *testing.T) {
	if len(GenerateFibonacci(100)) != 11 {
		t.Fail()
	}
}

//func TestGenerateFibonacciN(t *testing.T) {
//	GenerateFibonacciN(999009)
//}

func TestSumMatrix(t *testing.T) {
	m1 := InitMatrixCustom(5, 10, 1)
	m2 := InitMatrixCustom(5, 10, 1)
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
	data := randomizer.RandomIntArray(0, 100, 1000)
	t.Log(MultiplySumArray(data, data))
}
func BenchmarkMultiplySumArray1000(t *testing.B) {
	data := randomizer.RandomIntArray(0, 100, 1000)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		MultiplySumArray(data, data)
	}
}

//func BenchmarkMultiplyMatrixLegacy100x100(t *testing.B) {
//	m1 := InitRandomMatrix(100, 100)
//	m2 := InitRandomMatrix(100, 100)
//	t.ResetTimer()
//	for i := 0; i < t.N; i++ {
//		MultiplyMatrixLegacy(m1, m2)
//	}
//}

//func BenchmarkMultiplyMatrix100x100(t *testing.B) {
//	m1 := InitRandomMatrix(100, 100)
//	m2 := InitRandomMatrix(100, 100)
//	t.ResetTimer()
//	for i := 0; i < t.N; i++ {
//		MultiplyMatrix(m1, m2)
//	}
//}

// generateTestMatrix1 is delegated to generate a matrix for test purpose
func generateTestMatrix1() [][]int {
	matrix := make([][]int, 2)
	matrix[0] = []int{1, 2, 3}
	matrix[1] = []int{4, 5, 6}
	return matrix
}

// generateTestMatrix2 is delegated to generate a matrix for test purpose
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
	var array1 = []int{1, 2, 3, 4, 5, 6, 7}
	var result1 = []int{6, 5, 4, 3, 2, 1, 0}
	var array2 = []int{7, 6, 5, 4, 3, 2, 1}
	var result2 = []int{0, 1, 2, 3, 4, 5, 6}
	var array3 = []int{1, 9, 9, 2, 10, 3}
	var result3 = []int{4, 1, 2, 5, 3, 0}
	if !reflect.DeepEqual(SortMaxIndex(array1), result1) {
		t.Error("Error on [", array1, "]")
	}
	if !reflect.DeepEqual(SortMaxIndex(array2), result2) {
		t.Error("Error on [", array2, "]")
	}
	if !reflect.DeepEqual(SortMaxIndex(array3), result3) {
		t.Error("Error on [", array3, "]")
	}
}

func TestFindIndexValue(t *testing.T) {
	var array = []int{1, 2, 3, 4, 5, 6, 6, 7}
	t.Log(FindIndexValue(array, 6))
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
	var x = []float64{1, 2, 3}
	var y = []float64{2, 4, 6}
	taxiNorm := ManhattanDistance(x, y)
	if taxiNorm != 6 {
		t.Error(taxiNorm)
	}
}

func TestEuclideanDistance(t *testing.T) {
	var x = []float64{1, 2, 3}
	var y = []float64{2, 4, 6}
	distance := EuclideanDistance(x, y)
	if !(distance > 3.741 && distance < 3.742) {
		t.Error(distance)
	}
}

func TestMode(t *testing.T) {
	noMode := []int{0, 1, 2, 3, 4, 5, 6}
	oneMode := []int{0, 1, 2, 3, 4, 5, 6, 6}
	twoMode := []int{0, 1, 2, 3, 4, 5, 5, 6, 6}
	m1 := Mode[int](noMode)
	if len(m1) > 0 {
		t.Error("Err", m1)
	}
	m2 := Mode[int](oneMode)
	if len(m2) > 1 {
		t.Error("Err", m2)
	}
	m3 := Mode[int](twoMode)
	if len(m3) > 2 {
		t.Error("Err", m3)
	}
}

func TestMedian(t *testing.T) {
	median := []int{6, 5, 4, 3, 2, 1, 0}
	m1 := Median[int](median)
	if m1 != 3 {
		t.Error(m1)
	}
	median = append(median, 7)
	m1 = Median[int](median)
	if m1 != 3.5 {
		t.Error(m1)
	}
}

func TestStandardDeviation(t *testing.T) {
	median := []int{9, 2, 5, 4, 12, 7, 8, 11, 9, 3, 7, 4, 12, 5, 4, 10, 9, 6, 9, 4}
	m1 := StandardDeviation[int](median)
	if !(m1 > 2.9832 && m1 < 2.9833) {
		t.Error(m1)
	}
}

func TestVariance(t *testing.T) {
	median := []int{600, 470, 170, 430, 300}
	m1 := Variance[int](median)
	if !(m1 > 147.322 && m1 < 147.323) {
		t.Error(m1)
	}
}

func TestCovariance(t *testing.T) {
	arr1 := []int{1692, 1978, 1884, 2151, 2519}
	arr2 := []int{68, 102, 110, 112, 154}
	cv := Covariance[int](arr1, arr2)
	if cv != 9107.3 {
		t.Error(cv)
	}
}

func TestCorrelationInt(t *testing.T) {
	arr1 := []int{1692, 1978, 1884, 2151, 2519}
	arr2 := []int{68, 102, 110, 112, 154}
	cv := Correlation[int](arr1, arr2)
	if !(cv > 0.949 && cv < 0.950) {
		t.Error(cv)
	}
}

func TestCorrelationFloat64(t *testing.T) {
	arr1 := []float64{1691.75, 1977.80, 1884.09, 2151.13, 2519.36}
	arr2 := []float64{68.96, 100.11, 109.06, 112.18, 154.12}
	cv := Correlation[float64](arr1, arr2)
	if !(cv > 0.954 && cv < 0.955) {
		t.Error(cv)
	}
}

func TestSumArrays1(t *testing.T) {
	type args[T types.Number] struct {
		a1 []T
		a2 []T
	}
	type testCase[T types.Number] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args[int]{a1: []int{1, 2}, a2: []int{2, 1}},
			want: []int{3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumArrays(tt.args.a1, tt.args.a2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumArrays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransponse1D(t *testing.T) {
	type args[T any] struct {
		v []T
	}

	type testCase[T any] struct {
		name string
		args args[T]
		want [][1]T
	}

	tests := []testCase[int]{
		{
			name: "ok1",
			args: args[int]{[]int{1, 2, 3, 4}},
			want: [][1]int{{1}, {2}, {3}, {4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Transponse1D(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transponse1D() = %v, want %v", got, tt.want)
			}
		})
	}

}
