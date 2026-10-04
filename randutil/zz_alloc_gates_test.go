package randutil

import "testing"

var allocIntSink int
var allocInt64Sink int64

// Keep these tests sequential and use the real crypto reader.
func TestSecureIntAllocations(t *testing.T) {
	t.Run("SecureInt", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() {
			var err error
			allocIntSink, err = SecureInt(1000)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Logf("SecureInt(1000) allocations=%g", got)
		if got > 1 {
			t.Fatalf("want at most 1 allocation, got %g", got)
		}
	})
	t.Run("SecureInt64", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() {
			var err error
			allocInt64Sink, err = SecureInt64(1000)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Logf("SecureInt64(1000) allocations=%g", got)
		if got > 1 {
			t.Fatalf("want at most 1 allocation, got %g", got)
		}
	})
}

var allocStringSink string

func TestSecureStringAllocations(t *testing.T) {
	for _, length := range []int{16, 1024} {
		got := testing.AllocsPerRun(100, func() {
			var err error
			allocStringSink, err = SecureString(length, AlphaNumeric)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Logf("SecureString(%d, AlphaNumeric) allocations=%g", length, got)
		if got > 2 {
			t.Errorf("SecureString(%d, AlphaNumeric): want at most 2 allocations, got %g", length, got)
		}
	}
}
