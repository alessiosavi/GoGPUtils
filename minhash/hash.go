package minhash

const (
	fnvOffset uint64 = 14695981039346656037
	fnvPrime  uint64 = 1099511628211
	prime     uint64 = 1<<61 - 1
)

func mix64(z uint64) uint64 {
	z = (z ^ z>>30) * 0xbf58476d1ce4e5b9
	z = (z ^ z>>27) * 0x94d049bb133111eb
	return z ^ z>>31
}

func fnvBytes(h uint64, s string) uint64 {
	for i := range len(s) {
		h = (h ^ uint64(s[i])) * fnvPrime
	}
	return h
}

func fnvUint64(h, x uint64) uint64 {
	for range 8 {
		h = (h ^ (x & 255)) * fnvPrime
		x >>= 8
	}
	return h
}

// reduce folds a uint64 at bit 61. The sum is at most p+7.
func reduce(x uint64) uint64 {
	x = (x & prime) + (x >> 61)
	if x >= prime {
		x -= prime
	}
	return x
}
