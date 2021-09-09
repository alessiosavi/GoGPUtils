package goutils

import (
	"testing"
)

// const testFile string = `/tmp/GoGPUtils/array/arrayutils.go`
const testFile string = `../string/stringutils.go`

func TestCreateBenchmarkSignature(t *testing.T) {
	data, err := CreateBenchmarkSignature(testFile)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestCreateTestSignature(t *testing.T) {
	data, err := CreateTestSignature(testFile)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}
