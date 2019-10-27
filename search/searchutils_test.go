package searchutils

import (
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const length int = 1000000

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
	data := LinearSearchParallelInt(array, length+22, 10)
	if data != length+22 {
		t.Error(data)
	}
}
func TestOddLinearSearchParallelIntKO(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length + 23)
	data := LinearSearchParallelInt(array, -1, 10)
	if data != -1 {
		t.Error(data)
	}
}

func TestLinearSearchParallelIntOK(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchParallelInt(array, length-1, 10)
	if data != length-1 {
		t.Error(data)
	}
}

func TestLinearSearchParallelIntKO(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchParallelInt(array, -1, 10)
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
		LinearSearchParallelInt(array, length-1, 10)
	}
}
