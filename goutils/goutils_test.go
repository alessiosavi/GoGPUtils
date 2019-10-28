package goutils

import (
	"testing"
	//arrayutils "github.com/alessiosavi/GoGPUtils/array"
)

func TestCreateBenchmarkSignature(t *testing.T) {
	data, err := CreateBenchmarkSignature("../math/mathutils.go")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestCreateBenchmarkSignature1(t *testing.T) {
	GenerateTestSignature("../", "/tmp/gotest")
}
