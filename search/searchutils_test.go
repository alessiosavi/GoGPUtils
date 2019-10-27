package searchutils

import (
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const length int = 1000000

func TestLinearSearchInt(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	data := LinearSearchInt(array, length-1)
	t.Log(data)
}

func TestOddLinearSearchParallelInt(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length + 23)
	t.Log(LinearSearchParallelInt(array, length+22, 10))
}
func TestLinearSearchParallelInt(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	t.Log(LinearSearchParallelInt(array, length-1, 10))
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
