package a

// Foo is benchmarked.
func Foo() int { return 1 }

// Bar has no benchmark.
func Bar() int { return 2 }

func helper() int { return 3 }

// T is an exported type.
type T struct{}

// M is benchmarked as BenchmarkT_M.
func (T) M() {}

// N is exempted.
func (*T) N() {}

type u struct{}

// X has an unexported receiver, so it is not counted.
func (u) X() {}

// G has two type parameters.
type G[K comparable, V any] struct{}

// Get is benchmarked as BenchmarkG_Get.
func (g *G[K, V]) Get() {}

// S has one type parameter.
type S[E any] struct{}

// Len has no benchmark.
func (s S[E]) Len() int { return 0 }
