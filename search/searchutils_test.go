package searchutils

import (
	"io/ioutil"
	"math"
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const length int = math.MaxInt16
const thread int = 5

func TestLinearSearchIntOK(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchInt(array, length-1)
	if data == -1 {
		t.Error(data)
	}
}
func TestLinearSearchIntKO(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchInt(array, -1)
	if data != -1 {
		t.Error(data)
	}
}

func TestOddLinearSearchParallelIntOK(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length + 23)
	data := LinearSearchParallelInt(array, length+22, thread)
	if data != length+22 {
		t.Error(data)
	}
}
func TestOddLinearSearchParallelIntKO(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length + 23)
	data := LinearSearchParallelInt(array, -1, thread)
	if data != -1 {
		t.Error(data)
	}
}

func TestLinearSearchParallelIntOK(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchParallelInt(array, length-1, thread)
	if data != length-1 {
		t.Error(data)
	}
}

func TestLinearSearchParallelIntKO(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchParallelInt(array, -1, thread)
	if data != -1 {
		t.Error(data)
	}
}

func BenchmarkLinearSearchInt(t *testing.B) {
	array := helper.GenerateSequentialIntArray(length)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LinearSearchInt(array, length-1)
	}
}
func BenchmarkLinearSearchParallelInt(t *testing.B) {
	array := helper.GenerateSequentialIntArray(length)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LinearSearchParallelInt(array, length-1, thread)
	}
}

func BenchmarkContainsStringByte(t *testing.B) {
	data, err := ioutil.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	for i := 0; i < t.N; i++ {
		ContainsStringByte(data, "amor")
	}
}

func BenchmarkContainsStringsByte(t *testing.B) {
	data, err := ioutil.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	target := []string{"amor", "Beatrice"}
	for i := 0; i < t.N; i++ {
		ContainsStringsByte(data, target)
	}
}

func BenchmarkContainsWhichStrings(t *testing.B) {
	data, err := ioutil.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	target := []string{"amor", "Beatrice"}
	for i := 0; i < t.N; i++ {
		ContainsWhichStringsByte(data, target)
	}
}

func TestLinearSearchParallelInt(t *testing.T) {}

func TestLinearSearchParallelIntHelper(t *testing.T) {}
func BenchmarkLinearSearchParallelIntHelper(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestLinearSearchInt(t *testing.T) {}

func TestContainsString(t *testing.T) {}
func BenchmarkContainsString(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestContainsStringByte(t *testing.T) {}

func TestContainsStrings(t *testing.T) {}
func BenchmarkContainsStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestContainsStringsByte(t *testing.T) {}

func TestContainsWhichStrings(t *testing.T)     {}
func TestContainsWhichStringsByte(t *testing.T) {}
func BenchmarkContainsWhichStringsByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
