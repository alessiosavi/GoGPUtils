package searchutils

import (
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const length int = 100000

func TestLinearSearchInt(t *testing.T) {
	array := helper.GenerateSequentialIntArray(length)
	result := make(chan int)
	go LinearSearchInt(array, length-1, result)
	data := <-result
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
	result := make(chan int)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		go LinearSearchInt(array, length-1, result)
		<-result
	}
}
func BenchmarkLinearSearchParallelInt(t *testing.B) {
	array := helper.GenerateSequentialIntArray(length)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LinearSearchParallelInt(array, length-1, 10)
	}
}
