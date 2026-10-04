package sliceutil

import (
	"cmp"
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

var benchSliceTotal int

// BenchmarkFilter measures Filter keeping about half of the elements; n is
// the number of input elements.
func BenchmarkFilter(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Filter(s, func(v int) bool { return v%2 == 0 })
		}
	})
}

// BenchmarkFilterInPlace measures FilterInPlace keeping about half of the
// elements; n is the number of elements. FilterInPlace overwrites its input,
// so each iteration first copies the pristine input (the copy is measured).
func BenchmarkFilterInPlace(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		src := benchkit.Ints(n)
		work := make([]int, n)
		b.ReportAllocs()
		for b.Loop() {
			copy(work, src)
			FilterInPlace(work, func(v int) bool { return v%2 == 0 })
		}
	})
}

// BenchmarkMap measures doubling n deterministic input elements into a new slice.
func BenchmarkMap(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Map(s, func(v int) int { return v * 2 })
		}
	})
}

// BenchmarkMapWithIndex measures adding each index to its value across n input
// elements, allocating the result.
func BenchmarkMapWithIndex(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			MapWithIndex(s, func(i, v int) int { return i + v })
		}
	})
}

// BenchmarkReduce measures summing n input elements from an initial zero.
func BenchmarkReduce(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Reduce(s, 0, func(acc, v int) int { return acc + v })
		}
	})
}

// BenchmarkContains measures a worst-case search for an absent value among n
// input elements. String targets have the same length as the eight-byte inputs.
func BenchmarkContains(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					Contains(s, -1)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				b.ReportAllocs()
				for b.Loop() {
					Contains(s, "zzzzzzzz")
				}
			})
		}},
	)
}

// BenchmarkContainsFunc measures a worst-case search over n nonnegative input
// elements for an absent negative value.
func BenchmarkContainsFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			ContainsFunc(s, func(v int) bool { return v < 0 })
		}
	})
}

// BenchmarkIndexOf measures a worst-case search for an absent value among n
// input elements. String targets have the same length as the eight-byte inputs.
func BenchmarkIndexOf(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					IndexOf(s, -1)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				b.ReportAllocs()
				for b.Loop() {
					IndexOf(s, "zzzzzzzz")
				}
			})
		}},
	)
}

// BenchmarkIndexOfFunc measures a worst-case search over n nonnegative input
// elements for an absent negative value.
func BenchmarkIndexOfFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			IndexOfFunc(s, func(v int) bool { return v < 0 })
		}
	})
}

// BenchmarkLastIndexOf measures a worst-case search for an absent value among n
// input elements. String targets have the same length as the eight-byte inputs.
func BenchmarkLastIndexOf(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					LastIndexOf(s, -1)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				b.ReportAllocs()
				for b.Loop() {
					LastIndexOf(s, "zzzzzzzz")
				}
			})
		}},
	)
}

// BenchmarkUnique measures deduplicating n input elements into a new slice.
// Setup repeats the first half in the second half, giving about 50% duplicates.
func BenchmarkUnique(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				copy(s[n/2:], s[:n/2])
				b.ReportAllocs()
				for b.Loop() {
					Unique(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				copy(s[n/2:], s[:n/2])
				b.ReportAllocs()
				for b.Loop() {
					Unique(s)
				}
			})
		}},
	)
}

// BenchmarkUniqueFunc measures deduplicating n input elements by identity key.
// Setup repeats the first half in the second half, giving about 50% duplicates.
func BenchmarkUniqueFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		copy(s[n/2:], s[:n/2])
		b.ReportAllocs()
		for b.Loop() {
			UniqueFunc(s, func(v int) int { return v })
		}
	})
}

// BenchmarkChunk measures splitting n input elements into chunks of 64.
// The result shares input storage; the input remains unchanged.
func BenchmarkChunk(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Chunk(s, 64)
		}
	})
}

// BenchmarkFlatten measures allocating a flat copy of n total elements.
// Setup partitions the input into groups of sqrt(n) elements.
func BenchmarkFlatten(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		groups := Chunk(s, int(math.Sqrt(float64(n))))
		b.ReportAllocs()
		for b.Loop() {
			Flatten(groups)
		}
	})
}

// BenchmarkReverse measures reversing n input elements into a new slice.
func BenchmarkReverse(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Reverse(s)
		}
	})
}

// BenchmarkReverseInPlace measures reversing n elements in place. Each iteration
// first copies the pristine input into the work buffer; the copy is measured.
func BenchmarkReverseInPlace(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		src := benchkit.Ints(n)
		work := make([]int, n)
		b.ReportAllocs()
		for b.Loop() {
			copy(work, src)
			ReverseInPlace(work)
		}
	})
}

// BenchmarkIntersect measures Intersect on two inputs of n elements each.
// Setup takes overlapping windows of 3*n/2 deterministic elements, giving
// about 50% overlap. Both inputs remain unchanged.
func BenchmarkIntersect(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Ints(3 * n / 2)
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Intersect(a, other)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Strings(3*n/2, 8, "abcdefgh")
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Intersect(a, other)
				}
			})
		}},
	)
}

// BenchmarkDifference measures Difference on two inputs of n elements each.
// Setup takes overlapping windows of 3*n/2 deterministic elements, giving
// about 50% overlap. Both inputs remain unchanged.
func BenchmarkDifference(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Ints(3 * n / 2)
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Difference(a, other)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Strings(3*n/2, 8, "abcdefgh")
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Difference(a, other)
				}
			})
		}},
	)
}

// BenchmarkUnion measures Union on two inputs of n elements each.
// Setup takes overlapping windows of 3*n/2 deterministic elements, giving
// about 50% overlap. Both inputs remain unchanged.
func BenchmarkUnion(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Ints(3 * n / 2)
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Union(a, other)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				src := benchkit.Strings(3*n/2, 8, "abcdefgh")
				a, other := src[:n], src[n/2:]
				b.ReportAllocs()
				for b.Loop() {
					Union(a, other)
				}
			})
		}},
	)
}

// BenchmarkGroupBy measures grouping n input elements into up to eight groups.
// Integer keys use modulo eight; string keys use the first character.
func BenchmarkGroupBy(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					GroupBy(s, func(v int) int { return v % 8 })
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				b.ReportAllocs()
				for b.Loop() {
					GroupBy(s, func(v string) string { return v[:1] })
				}
			})
		}},
	)
}

// BenchmarkPartition measures splitting n input elements by parity into two
// new slices, with about half in each.
func BenchmarkPartition(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Partition(s, func(v int) bool { return v%2 == 0 })
		}
	})
}

// BenchmarkTake measures copying the first half.
// n is the number of input elements; the input remains unchanged.
func BenchmarkTake(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Take(s, n/2)
		}
	})
}

// BenchmarkTakeLast measures copying the last half.
// n is the number of input elements; the input remains unchanged.
func BenchmarkTakeLast(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			TakeLast(s, n/2)
		}
	})
}

// BenchmarkDrop measures dropping the first half and copying the rest.
// n is the number of input elements; the input remains unchanged.
func BenchmarkDrop(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Drop(s, n/2)
		}
	})
}

// BenchmarkDropLast measures dropping the last half and copying the rest.
// n is the number of input elements; the input remains unchanged.
func BenchmarkDropLast(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			DropLast(s, n/2)
		}
	})
}

// BenchmarkTakeWhile measures scanning a nonnegative prefix of n/2 elements
// and copying the prefix. n is the total number of input elements;
// setup sets the second half to -1.
func BenchmarkTakeWhile(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		for i := n / 2; i < n; i++ {
			s[i] = -1
		}
		b.ReportAllocs()
		for b.Loop() {
			TakeWhile(s, func(v int) bool { return v >= 0 })
		}
	})
}

// BenchmarkDropWhile measures scanning a nonnegative prefix of n/2 elements
// and copying the remaining suffix. n is the total number of input elements;
// setup sets the second half to -1.
func BenchmarkDropWhile(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		for i := n / 2; i < n; i++ {
			s[i] = -1
		}
		b.ReportAllocs()
		for b.Loop() {
			DropWhile(s, func(v int) bool { return v >= 0 })
		}
	})
}

// BenchmarkAll measures a worst-case full scan of n nonnegative input
// elements with a predicate that always succeeds.
func BenchmarkAll(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			All(s, func(v int) bool { return v >= 0 })
		}
	})
}

// BenchmarkAny measures a worst-case full scan of n nonnegative input
// elements with a predicate that never matches.
func BenchmarkAny(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Any(s, func(v int) bool { return v < 0 })
		}
	})
}

// BenchmarkNone measures a worst-case full scan of n nonnegative input
// elements with a predicate that never matches.
func BenchmarkNone(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			None(s, func(v int) bool { return v < 0 })
		}
	})
}

// BenchmarkCount measures counting even values among n input elements.
func BenchmarkCount(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Count(s, func(v int) bool { return v%2 == 0 })
		}
	})
}

// BenchmarkFind measures a worst-case successful forward search among n
// input elements. Setup places the only negative value last.
func BenchmarkFind(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		s[n-1] = -1
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := Find(s, func(v int) bool { return v < 0 }); !ok {
				b.Fatal("negative target not found")
			}
		}
	})
}

// BenchmarkFindLast measures a worst-case successful reverse search among n
// input elements. Setup places the only negative value first.
func BenchmarkFindLast(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		s[0] = -1
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := FindLast(s, func(v int) bool { return v < 0 }); !ok {
				b.Fatal("negative target not found")
			}
		}
	})
}

// BenchmarkMin measures selecting the minimum of n input elements. The input is nonempty and unchanged.
func BenchmarkMin(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := Min(s); !ok {
				b.Fatal("nonempty input reported empty")
			}
		}
	})
}

// BenchmarkMax measures selecting the maximum of n input elements. The input is nonempty and unchanged.
func BenchmarkMax(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := Max(s); !ok {
				b.Fatal("nonempty input reported empty")
			}
		}
	})
}

// BenchmarkMinFunc measures selecting the minimum of n input elements
// using cmp.Compare[int]. The input is nonempty and unchanged.
func BenchmarkMinFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := MinFunc(s, cmp.Compare[int]); !ok {
				b.Fatal("nonempty input reported empty")
			}
		}
	})
}

// BenchmarkMaxFunc measures selecting the maximum of n input elements
// using cmp.Compare[int]. The input is nonempty and unchanged.
func BenchmarkMaxFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := MaxFunc(s, cmp.Compare[int]); !ok {
				b.Fatal("nonempty input reported empty")
			}
		}
	})
}

// BenchmarkEqual measures a full comparison of two equal slices with n elements
// each. Setup copies the second slice and clones each string, so the string
// case compares equal content stored at different addresses.
func BenchmarkEqual(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				other := slices.Clone(s)
				b.ReportAllocs()
				for b.Loop() {
					Equal(s, other)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				other := slices.Clone(s)
				for i := range other {
					other[i] = strings.Clone(other[i])
				}
				b.ReportAllocs()
				for b.Loop() {
					Equal(s, other)
				}
			})
		}},
	)
}

// BenchmarkEqualFunc measures a full equality-callback comparison of two equal
// slices with n elements each. Setup copies the second slice.
func BenchmarkEqualFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		other := slices.Clone(s)
		b.ReportAllocs()
		for b.Loop() {
			EqualFunc(s, other, func(a, other int) bool { return a == other })
		}
	})
}

// BenchmarkForEach measures visiting n input elements and summing values in the callback. Resetting the total each iteration is measured;
// the final total is retained after the loop.
func BenchmarkForEach(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		total := 0
		b.ReportAllocs()
		for b.Loop() {
			total = 0
			ForEach(s, func(v int) { total += v })
		}
		benchSliceTotal = total
	})
}

// BenchmarkForEachWithIndex measures visiting n input elements and summing values
// plus indices in the callback. Resetting the total each iteration is measured;
// the final total is retained after the loop.
func BenchmarkForEachWithIndex(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		total := 0
		b.ReportAllocs()
		for b.Loop() {
			total = 0
			ForEachWithIndex(s, func(i, v int) { total += i + v })
		}
		benchSliceTotal = total
	})
}

// BenchmarkZip measures pairing two equally sized inputs of n elements each.
// The measured call includes the API's interface boxing and output allocation.
func BenchmarkZip(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		other := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Zip(s, other)
		}
	})
}

// BenchmarkZipWith measures adding paired values from two equally sized inputs
// of n elements each into a new slice.
func BenchmarkZipWith(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		other := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			ZipWith(s, other, func(a, other int) int { return a + other })
		}
	})
}

// BenchmarkPad measures appending n fill values of -1 to n input
// elements, producing a new slice of length 2*n.
func BenchmarkPad(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Pad(s, 2*n, -1)
		}
	})
}

// BenchmarkPadLeft measures prepending n fill values of -1 to n input
// elements, producing a new slice of length 2*n.
func BenchmarkPadLeft(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			PadLeft(s, 2*n, -1)
		}
	})
}

// BenchmarkRemoveAt measures removing the middle element (index n/2) from n
// input elements. The API allocates its result and leaves the input unchanged.
func BenchmarkRemoveAt(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			RemoveAt(s, n/2)
		}
	})
}

// BenchmarkRemoveValue measures removing a value present exactly once among n
// input elements. Setup places -1 last, requiring a full scan;
// the API allocates its result and leaves the input unchanged.
func BenchmarkRemoveValue(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		s[n-1] = -1
		b.ReportAllocs()
		for b.Loop() {
			RemoveValue(s, -1)
		}
	})
}

// BenchmarkRemoveFirst measures removing a value present exactly once among n
// input elements. Setup places -1 last, requiring a full successful search;
// the API allocates its result and leaves the input unchanged.
func BenchmarkRemoveFirst(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		s[n-1] = -1
		b.ReportAllocs()
		for b.Loop() {
			RemoveFirst(s, -1)
		}
	})
}

// BenchmarkInsert measures inserting -1 at the middle index n/2 of n input
// elements. The API allocates its result and leaves the input unchanged.
func BenchmarkInsert(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			Insert(s, n/2, -1)
		}
	})
}

// BenchmarkShuffle measures copying and shuffling n input elements. Setup seeds
// the PRNG once with 42; it advances across iterations. The API copies the input,
// so no restoration is needed.
func BenchmarkShuffle(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		SeedShuffle(42)
		b.ReportAllocs()
		for b.Loop() {
			Shuffle(s)
		}
	})
}

// BenchmarkShuffleInPlace measures shuffling n elements in place. Each iteration
// first copies the pristine input into the work buffer; the copy is measured.
// Setup seeds the PRNG once with 42; it advances across iterations.
func BenchmarkShuffleInPlace(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		src := benchkit.Ints(n)
		work := make([]int, n)
		SeedShuffle(42)
		b.ReportAllocs()
		for b.Loop() {
			copy(work, src)
			ShuffleInPlace(work)
		}
	})
}

// BenchmarkSeedShuffle measures setting the fixed scalar seed 42, including
// the API's mutex. Every iteration overwrites the PRNG state; there is no n.
func BenchmarkSeedShuffle(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		SeedShuffle(42)
	}
}

// BenchmarkCompact measures compacting n input elements arranged in adjacent
// runs of four during setup. The API allocates its output without modifying
// the input, so no restoration copy is needed.
func BenchmarkCompact(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] = s[i-i%4]
				}
				b.ReportAllocs()
				for b.Loop() {
					Compact(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				for i := range s {
					s[i] = s[i-i%4]
				}
				b.ReportAllocs()
				for b.Loop() {
					Compact(s)
				}
			})
		}},
	)
}

// BenchmarkCompactFunc measures compacting n input elements arranged in adjacent
// runs of four during setup, using an equality callback. The API allocates its
// output without modifying the input, so no restoration copy is needed.
func BenchmarkCompactFunc(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		for i := range s {
			s[i] = s[i-i%4]
		}
		b.ReportAllocs()
		for b.Loop() {
			CompactFunc(s, func(a, other int) bool { return a == other })
		}
	})
}

// BenchmarkFlatMap measures flattening n total elements with an identity group
// transform. Setup partitions the input into groups of sqrt(n) elements;
// the measured call appends each unchanged group to a new output slice.
func BenchmarkFlatMap(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		groups := Chunk(s, int(math.Sqrt(float64(n))))
		b.ReportAllocs()
		for b.Loop() {
			FlatMap(groups, func(group []int) []int { return group })
		}
	})
}

// BenchmarkAssociate measures mapping n input elements to themselves by identity
// key; key and element types match. The input remains unchanged.
func BenchmarkAssociate(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				b.ReportAllocs()
				for b.Loop() {
					Associate(s, func(v int) int { return v })
				}
			})
		}},
		benchkit.TypeCase{Name: "string", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Strings(n, 8, "abcdefgh")
				b.ReportAllocs()
				for b.Loop() {
					Associate(s, func(v string) string { return v })
				}
			})
		}},
	)
}

// BenchmarkAssociateWith measures building a map from n input keys to doubled
// values. The input remains unchanged.
func BenchmarkAssociateWith(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			AssociateWith(s, func(v int) int { return v * 2 })
		}
	})
}

// BenchmarkMapErr measures the successful path over n input elements, doubling
// each value into a new slice and checking the returned error.
func BenchmarkMapErr(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			if _, err := MapErr(s, func(v int) (int, error) { return v * 2, nil }); err != nil {
				b.Fatal(err)
			}
		}
	})
}
