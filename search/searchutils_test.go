package searchutils

import (
	"math"
	"os"
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const length int = math.MaxInt16

func TestLinearSearchOK(t *testing.T) {
	array := helper.GenerateSequentialArray[int](length)
	data := LinearSearch(array, length-1)
	if data == -1 {
		t.Error(data)
	}
}
func TestLinearSearchKO(t *testing.T) {
	array := helper.GenerateSequentialArray[int](length)
	data := LinearSearch(array, -1)
	if data != -1 {
		t.Error(data)
	}
}

func BenchmarkLinearSearch(t *testing.B) {
	array := helper.GenerateSequentialArray[int](length)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LinearSearch(array, length-1)
	}
}

func BenchmarkContainsStringByte(t *testing.B) {
	data, err := os.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	for i := 0; i < t.N; i++ {
		ContainsStringByte(data, "amor")
	}
}

func BenchmarkContainsStringsByte(t *testing.B) {
	data, err := os.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	target := []string{"amor", "Beatrice"}
	for i := 0; i < t.N; i++ {
		ContainsStringsByte(data, target)
	}
}

func BenchmarkContainsWhichStrings(t *testing.B) {
	data, err := os.ReadFile("../testdata/files/dante.txt")
	if err != nil {
		return
	}
	target := []string{"amor", "Beatrice"}
	for i := 0; i < t.N; i++ {
		ContainsWhichStringsByte(data, target)
	}
}
