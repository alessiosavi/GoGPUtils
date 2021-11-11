package helper

import (
	"math"
	"testing"
)

const total int = 1000

func TestRandomIntn(t *testing.T) {
	for n := 0; n < 1000; n++ {
		if RandomInt(0, 1000) > 1000 {
			t.Fail()
		}
	}
}

func TestRandomInt32(t *testing.T) {
	for n := 0; n < 1000; n++ {
		if RandomInt32(0, 1000) > 1000 {
			t.Fail()
		}
	}
}

func TestRandomInt64(t *testing.T) {
	for n := 0; n < 1000; n++ {
		if RandomInt64(0, 1000) > 1000 {
			t.Fail()
		}
	}
}

func TestRandomFloat32(t *testing.T) {
	for n := 0; n < 1000; n++ {
		if RandomFloat32(0, 1000) > 1000 {
			t.Fail()
		}
	}
}

func TestRandomFloat64(t *testing.T) {
	for n := 0; n < 1000; n++ {
		if RandomFloat64(0, 1000) > 1000 {
			t.Fail()
		}
	}
}

func BenchmarkRandomIntn(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomInt(0, math.MaxInt8)
	}
}

func BenchmarkRandomInt32(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomInt32(0, math.MaxInt32)
	}
}

func BenchmarkRandomInt64(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomInt64(0, math.MaxInt64)
	}
}

func BenchmarkRandomFloat32(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomFloat32(0, math.MaxFloat32)
	}
}

func BenchmarkRandomFloat64(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomFloat64(0, math.MaxFloat64)
	}
}

func BenchmarkRandomIntnR(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomInt(0, math.MaxInt8)
	}
}

func BenchmarkRandomInt32R(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomInt32(0, math.MaxInt32)
	}
}

func BenchmarkRandomInt64R(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomInt64(0, math.MaxInt64)
	}
}

func BenchmarkRandomFloat32R(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomFloat32(0, math.MaxFloat32)
	}
}

func BenchmarkRandomFloat64R(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomFloat64(0, math.MaxFloat64)
	}
}

func BenchmarkRandomIntnRArray(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomIntArray(0, math.MaxInt8, total)
	}
}

func BenchmarkRandomInt32RArray(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomInt32Array(0, math.MaxInt32, total)
	}
}

func BenchmarkRandomInt64RArray(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomInt64Array(0, math.MaxInt64, total)
	}
}

func BenchmarkRandomFloat32Array(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomFloat32Array(0, math.MaxFloat32, total)
	}
}

func BenchmarkRandomFloat64RArray(t *testing.B) {
	randomer := InitRandomizer()
	for n := 0; n < t.N; n++ {
		randomer.RandomFloat64Array(0, math.MaxFloat64, total)
	}
}

//func TestConvertSize(t *testing.T) {
//	t.Log(ConvertSize(1024, "KB"))
//	t.Log(ConvertSize(1000000, "MB"))
//	t.Log(ConvertSize(1024, "GB"))
//}

func TestRandomGenerator_RandomString(t *testing.T) {
	var r = InitRandomizer()
	t.Log(r.RandomString(10))
}

func BenchmarkRandomString(t *testing.B) {
	var r = InitRandomizer()
	for n := 0; n < t.N; n++ {
		r.RandomString(5000)
	}
}
