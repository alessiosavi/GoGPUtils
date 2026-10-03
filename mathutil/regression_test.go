package mathutil

import (
	"math"
	"slices"
	"testing"
	"time"
)

func TestSqrt_HandlesExtremeValues(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"large", 1e100, 1e50},
		{"small", 1e-100, 1e-50},
		{"infinity", math.Inf(1), math.Inf(1)},
		{"nan", math.NaN(), math.NaN()},
		{"negative infinity", math.Inf(-1), 0},
		{"negative", -4, 0},
		{"negative zero", math.Copysign(0, -1), math.Copysign(0, -1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.x)
			if math.IsNaN(tt.want) {
				if !math.IsNaN(got) {
					t.Fatalf("Sqrt(%v) = %v, want NaN", tt.x, got)
				}

				return
			}

			if got != tt.want || math.Signbit(got) != math.Signbit(tt.want) {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func checkArangeTerminates[T Number](t *testing.T, start, stop, step T, want []T) {
	t.Helper()

	done := make(chan []T, 1)
	go func() {
		done <- Arange(start, stop, step)
	}()

	select {
	case got := <-done:
		if !slices.Equal(got, want) {
			t.Errorf("Arange(%v, %v, %v) = %v, want %v", start, stop, step, got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("Arange did not terminate within one second")
	}
}

func TestArange_StopsAtIntegerOverflow(t *testing.T) {
	t.Run("unsigned", func(t *testing.T) {
		checkArangeTerminates(t, uint8(250), uint8(255), uint8(4), []uint8{250, 254})
	})
	t.Run("signed ascending", func(t *testing.T) {
		checkArangeTerminates(t, int8(122), int8(127), int8(4), []int8{122, 126})
	})
	t.Run("signed descending", func(t *testing.T) {
		checkArangeTerminates(t, int8(-123), int8(-128), int8(-4), []int8{-123, -127})
	})
	t.Run("uint64", func(t *testing.T) {
		checkArangeTerminates(t, uint64(math.MaxUint64-5), uint64(math.MaxUint64), uint64(4), []uint64{math.MaxUint64 - 5, math.MaxUint64 - 1})
	})
	t.Run("int64 minimum step", func(t *testing.T) {
		checkArangeTerminates(t, int64(0), int64(math.MinInt64), int64(math.MinInt64), []int64{0})
	})
	t.Run("cross zero", func(t *testing.T) {
		checkArangeTerminates(t, int64(math.MinInt64), int64(math.MaxInt64), int64(math.MaxInt64), []int64{math.MinInt64, -1, math.MaxInt64 - 1})
	})
	t.Run("named integer", func(t *testing.T) {
		type counter uint8

		checkArangeTerminates(t, counter(250), counter(255), counter(4), []counter{250, 254})
	})
}

func TestArange_StopsWhenFloatCannotAdvance(t *testing.T) {
	t.Run("float64 ascending", func(t *testing.T) {
		checkArangeTerminates(t, 1e16, 1e16+10, 0.5, []float64{1e16})
	})
	t.Run("float64 descending", func(t *testing.T) {
		checkArangeTerminates(t, -1e16, -1e16-10, -0.5, []float64{-1e16})
	})
	t.Run("float32", func(t *testing.T) {
		checkArangeTerminates(t, float32(1e8), float32(1e8+16), float32(0.5), []float32{1e8})
	})
	t.Run("infinite start", func(t *testing.T) {
		checkArangeTerminates(t, math.Inf(-1), 0, 1, []float64{math.Inf(-1)})
	})
	t.Run("fractional steps", func(t *testing.T) {
		checkArangeTerminates(t, 0.0, 1.0, 0.25, []float64{0, 0.25, 0.5, 0.75})
	})
}
