package k

// T has a method M.
type T struct{}

// M expects BenchmarkT_M.
func (T) M() {}

// T_M also expects BenchmarkT_M.
func T_M() {}

// Foo is benchmarked in both test packages.
func Foo() {}
