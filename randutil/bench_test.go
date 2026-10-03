package randutil

import (
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkSecureBytes measures generating n bytes with the real crypto RNG.
// Every iteration allocates a fresh result; no restoration is needed.
func BenchmarkSecureBytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := SecureBytes(n); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSecureString measures generating n ASCII bytes from AlphaNumeric
// with the real crypto RNG. Each iteration generates a fresh string.
func BenchmarkSecureString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := SecureString(n, AlphaNumeric); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkGenerator_Int measures a draw in [0, 1000) from a generator seeded
// with 1 before the loop. Its random state advances each iteration.
func BenchmarkGenerator_Int(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Int(1000)
	}
}

// BenchmarkGenerator_String measures generating n ASCII bytes from AlphaNumeric.
// The generator is seeded with 1 before the loop and advances each iteration.
func BenchmarkGenerator_String(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		b.ReportAllocs()
		for b.Loop() {
			g.String(n, AlphaNumeric)
		}
	})
}

// BenchmarkGenerator_ShuffleInts measures shuffling n integers in place with a
// generator seeded with 1. Repeatedly shuffling the same slice is the same
// workload, so no restoration is needed; the random state advances each time.
func BenchmarkGenerator_ShuffleInts(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		items := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			g.ShuffleInts(items)
		}
	})
}

// BenchmarkSecureInt measures a draw in [0, 1000) using the real crypto RNG.
// There is no fixture to restore between draws.
func BenchmarkSecureInt(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := SecureInt(1000); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSecureInt64 measures a draw in [0, 1000) using the real crypto RNG.
// There is no fixture to restore between draws.
func BenchmarkSecureInt64(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := SecureInt64(1000); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSecureID measures generating 16 bytes with the real crypto RNG and
// encoding them as 32 hex characters. Each iteration generates a fresh ID.
func BenchmarkSecureID(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := SecureID(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSecureChoice measures selecting one of n unchanged input elements
// with the real crypto RNG. Inputs are deterministic integers or 8-byte ASCII
// alphanumeric strings, and every input slice is nonempty.
func BenchmarkSecureChoice(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				items := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := SecureChoice(items); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				items := benchkit.Strings(n, 8, AlphaNumeric)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := SecureChoice(items); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkNewGenerator measures constructing a fresh generator each iteration,
// including its automatic seeding from the real crypto RNG.
func BenchmarkNewGenerator(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewGenerator()
	}
}

// BenchmarkNewGeneratorWithSeed measures constructing a fresh generator with
// seed 1 each iteration. Construction is the entire measured operation.
func BenchmarkNewGeneratorWithSeed(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewGeneratorWithSeed(1)
	}
}

// BenchmarkGenerator_IntRange measures a draw in inclusive [10, 20] from a
// generator seeded with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_IntRange(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.IntRange(10, 20)
	}
}

// BenchmarkGenerator_Int64 measures a draw in [0, 1000) from a generator seeded
// with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Int64(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Int64(1000)
	}
}

// BenchmarkGenerator_Int64Range measures a draw in inclusive [10, 20] from a
// generator seeded with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Int64Range(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Int64Range(10, 20)
	}
}

// BenchmarkGenerator_Float64 measures a draw in [0, 1) from a generator seeded
// with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Float64(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Float64()
	}
}

// BenchmarkGenerator_Float64Range measures a draw in [10, 20) from a generator
// seeded with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Float64Range(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Float64Range(10, 20)
	}
}

// BenchmarkGenerator_Float32 measures a draw in [0, 1) from a generator seeded
// with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Float32(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Float32()
	}
}

// BenchmarkGenerator_Bool measures drawing one boolean from a generator seeded
// with 1 before the loop; its state advances each iteration.
func BenchmarkGenerator_Bool(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Bool()
	}
}

// BenchmarkGenerator_Bytes measures generating n bytes into a fresh slice.
// The generator is seeded with 1 before the loop and advances each iteration.
func BenchmarkGenerator_Bytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		b.ReportAllocs()
		for b.Loop() {
			g.Bytes(n)
		}
	})
}

// BenchmarkGenerator_AlphaNumericString measures generating n ASCII
// alphanumeric bytes. The generator is seeded with 1 before the loop and
// advances each iteration, producing a fresh string.
func BenchmarkGenerator_AlphaNumericString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		b.ReportAllocs()
		for b.Loop() {
			g.AlphaNumericString(n)
		}
	})
}

// BenchmarkGenerator_ChoiceInt measures selecting one of n unchanged,
// deterministic integers. The generator is seeded with 1 before the loop
// and advances each iteration; every input slice is nonempty.
func BenchmarkGenerator_ChoiceInt(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		items := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			g.ChoiceInt(items)
		}
	})
}

// BenchmarkGenerator_ChoiceString measures selecting one of n unchanged,
// deterministic 8-byte ASCII alphanumeric strings. The generator is seeded
// with 1 before the loop and advances each iteration; inputs are nonempty.
func BenchmarkGenerator_ChoiceString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		items := benchkit.Strings(n, 8, AlphaNumeric)
		b.ReportAllocs()
		for b.Loop() {
			g.ChoiceString(items)
		}
	})
}

// BenchmarkGenerator_ShuffleStrings measures shuffling n deterministic 8-byte
// ASCII alphanumeric strings in place with a generator seeded with 1.
// Repeatedly shuffling the same slice is the same workload, so no restoration
// is needed; the random state advances each iteration.
func BenchmarkGenerator_ShuffleStrings(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		items := benchkit.Strings(n, 8, AlphaNumeric)
		b.ReportAllocs()
		for b.Loop() {
			g.ShuffleStrings(items)
		}
	})
}

// BenchmarkGenerator_SampleInts measures selecting k=n/10 elements without
// replacement from n unchanged, deterministic integers. The API's full copy
// and shuffle are measured. The generator is seeded with 1 before the loop
// and advances each iteration; no fixture restoration is needed.
func BenchmarkGenerator_SampleInts(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		items := benchkit.Ints(n)
		k := n / 10
		b.ReportAllocs()
		for b.Loop() {
			g.SampleInts(items, k)
		}
	})
}

// BenchmarkChoice measures selecting one of n unchanged input elements using
// a generator seeded with 1 before the loop that advances each iteration.
// Inputs are deterministic integers or 8-byte ASCII alphanumeric strings;
// every input slice is nonempty.
func BenchmarkChoice(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					Choice(g, items)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Strings(n, 8, AlphaNumeric)
				b.ReportAllocs()
				for b.Loop() {
					Choice(g, items)
				}
			})
		}},
	)
}

// BenchmarkChoiceN measures selecting k=n/10 elements with replacement from n
// unchanged inputs: deterministic integers or 8-byte ASCII alphanumeric strings.
// The generator is seeded with 1 before the loop and advances each iteration;
// every call allocates its result and no fixture restoration is needed.
func BenchmarkChoiceN(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Ints(n)
				k := n / 10
				b.ReportAllocs()
				for b.Loop() {
					ChoiceN(g, items, k)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Strings(n, 8, AlphaNumeric)
				k := n / 10
				b.ReportAllocs()
				for b.Loop() {
					ChoiceN(g, items, k)
				}
			})
		}},
	)
}

// BenchmarkSample measures selecting k=n/10 elements without replacement from
// n unchanged inputs: deterministic integers or 8-byte ASCII alphanumeric
// strings. The API's full copy and shuffle are measured. The generator is
// seeded with 1 before the loop and advances each iteration; no fixture
// restoration is needed.
func BenchmarkSample(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Ints(n)
				k := n / 10
				b.ReportAllocs()
				for b.Loop() {
					Sample(g, items, k)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Strings(n, 8, AlphaNumeric)
				k := n / 10
				b.ReportAllocs()
				for b.Loop() {
					Sample(g, items, k)
				}
			})
		}},
	)
}

// BenchmarkShuffle measures shuffling n deterministic integers or 8-byte ASCII
// alphanumeric strings in place with a generator seeded with 1. Repeatedly
// shuffling the same slice is the same workload, so no restoration is needed;
// the random state advances each iteration.
func BenchmarkShuffle(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					Shuffle(g, items)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Strings(n, 8, AlphaNumeric)
				b.ReportAllocs()
				for b.Loop() {
					Shuffle(g, items)
				}
			})
		}},
	)
}

// BenchmarkShuffleCopy measures copying and shuffling n unchanged inputs:
// deterministic integers or 8-byte ASCII alphanumeric strings. The API's copy
// is measured. The generator is seeded with 1 before the loop and advances
// each iteration; no fixture restoration is needed.
func BenchmarkShuffleCopy(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					ShuffleCopy(g, items)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				g := NewGeneratorWithSeed(1)
				items := benchkit.Strings(n, 8, AlphaNumeric)
				b.ReportAllocs()
				for b.Loop() {
					ShuffleCopy(g, items)
				}
			})
		}},
	)
}

// BenchmarkGenerator_WeightedChoice measures selecting an index from n unchanged
// weights drawn deterministically from [1, 2), so their total is positive.
// The generator is seeded with 1 before the loop and advances each iteration.
func BenchmarkGenerator_WeightedChoice(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		g := NewGeneratorWithSeed(1)
		weights := benchkit.Floats(n)
		for i := range weights {
			weights[i]++
		}
		b.ReportAllocs()
		for b.Loop() {
			g.WeightedChoice(weights)
		}
	})
}

// BenchmarkGenerator_Probability measures a Bernoulli draw with probability
// 0.75 from a generator seeded with 1 before the loop; its state advances
// each iteration.
func BenchmarkGenerator_Probability(b *testing.B) {
	g := NewGeneratorWithSeed(1)
	b.ReportAllocs()
	for b.Loop() {
		g.Probability(0.75)
	}
}

// BenchmarkSequence measures allocating n sequential integers starting at 0.
// Every iteration creates a fresh slice; there is no fixture to restore.
func BenchmarkSequence(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			Sequence(0, n)
		}
	})
}

// BenchmarkRange measures allocating n integers from [0, n).
// Every iteration creates a fresh slice; there is no fixture to restore.
func BenchmarkRange(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			Range(0, n)
		}
	})
}

// BenchmarkRangeStep measures allocating n integers from [0, 2*n) with step 2.
// Every iteration creates a fresh slice; there is no fixture to restore.
func BenchmarkRangeStep(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			RangeStep(0, 2*n, 2)
		}
	})
}
