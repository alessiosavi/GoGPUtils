package collection

import (
	"slices"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkStack_Push measures pushing one item onto a stack of n items and
// popping it again. The Pop and its success check are measured; a spare slot
// in the backing slice is reused, keeping the depth bounded at n or n+1.
func BenchmarkStack_Push(b *testing.B) {
	b.Run("case=push+pop", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewStackWithCapacity[int](n + 1)
			s.PushAll(items...)
			b.ReportAllocs()
			for b.Loop() {
				s.Push(items[n/2])
				if _, ok := s.Pop(); !ok {
					b.Fatal("stack unexpectedly empty")
				}
			}
		})
	})
}

// BenchmarkStack_Pop measures popping from a stack of n items and pushing the
// returned item back. The success check and restoring Push are measured; the
// backing slice's spare capacity is reused.
func BenchmarkStack_Pop(b *testing.B) {
	b.Run("case=pop+repush", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			s := NewStackWithCapacity[int](n + 1)
			s.PushAll(benchkit.Ints(n)...)
			b.ReportAllocs()
			for b.Loop() {
				item, ok := s.Pop()
				if !ok {
					b.Fatal("stack unexpectedly empty")
				}
				s.Push(item)
			}
		})
	})
}

// BenchmarkQueue_Enqueue measures adding one item to a queue of n items and
// dequeuing its front. The Dequeue and success check are measured. Dequeue
// advances the slice start, so this bounded cycle includes occasional growth
// of the backing array as its remaining capacity is consumed.
func BenchmarkQueue_Enqueue(b *testing.B) {
	b.Run("case=enqueue+dequeue", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			q := NewQueueWithCapacity[int](n + 1)
			q.EnqueueAll(items...)
			b.ReportAllocs()
			for b.Loop() {
				q.Enqueue(items[n/2])
				if _, ok := q.Dequeue(); !ok {
					b.Fatal("queue unexpectedly empty")
				}
			}
		})
	})
}

// BenchmarkQueue_Dequeue measures removing the front of a queue of n items
// and re-enqueuing that item. The success check and restoring Enqueue are
// measured. Values rotate and depth stays bounded; advancing the slice start
// includes occasional backing-array growth during re-enqueue.
func BenchmarkQueue_Dequeue(b *testing.B) {
	b.Run("case=dequeue+reenqueue", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			q := NewQueueWithCapacity[int](n + 1)
			q.EnqueueAll(benchkit.Ints(n)...)
			b.ReportAllocs()
			for b.Loop() {
				item, ok := q.Dequeue()
				if !ok {
					b.Fatal("queue unexpectedly empty")
				}
				q.Enqueue(item)
			}
		})
	})
}

// BenchmarkSet_Add measures adding one element to a set of n elements.
// case=existing re-adds a present element (no growth). case=new adds an
// absent element and removes it again, so the set stays at n elements; the
// Remove is measured.
func BenchmarkSet_Add(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		b.Run("case=existing", func(b *testing.B) {
			s := NewSetFrom(items)
			b.ReportAllocs()
			for b.Loop() {
				s.Add(items[n/2])
			}
		})
		b.Run("case=new", func(b *testing.B) {
			s := NewSetFrom(items)
			absent := -1 // benchkit.Ints never yields negatives
			b.ReportAllocs()
			for b.Loop() {
				s.Add(absent)
				s.Remove(absent)
			}
		})
	})
}

// BenchmarkSet_Contains measures a hit or miss in a prebuilt set of n
// deterministic integer elements; the expected-result check is measured.
func BenchmarkSet_Contains(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		s := NewSetFrom(items)
		for _, tc := range []struct {
			name  string
			value int
			want  bool
		}{
			{"hit", items[n/2], true},
			{"miss", -1, false},
		} {
			b.Run("case="+tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if s.Contains(tc.value) != tc.want {
						b.Fatal("unexpected membership result")
					}
				}
			})
		}
	})
}

// BenchmarkBST_Insert measures inserting one value into a tree of n nodes.
// Shuffled fixtures use Sizes; sorted fixtures use SmallSizes because their
// construction is quadratic. case=new inserts above the maximum (the deep end
// of a sorted tree), checks size, and removes that leaf; checks and removal
// are measured. case=existing reinserts a resident value without growing.
func BenchmarkBST_Insert(b *testing.B) {
	for _, order := range []struct {
		name   string
		sizes  []int
		values func(int) []int
	}{
		{"shuffled", benchkit.Sizes, benchkit.Ints},
		{"sorted", benchkit.SmallSizes, benchkit.SortedInts},
	} {
		b.Run("order="+order.name, func(b *testing.B) {
			benchkit.Run(b, order.sizes, func(b *testing.B, n int) {
				items := order.values(n)
				b.Run("case=new", func(b *testing.B) {
					tree := NewBSTFrom(items)
					value := slices.Max(items) + 1
					b.ReportAllocs()
					for b.Loop() {
						tree.Insert(value)
						if tree.Len() != n+1 {
							b.Fatal("insertion did not add one node")
						}
						if !tree.Remove(value) {
							b.Fatal("inserted leaf unexpectedly missing")
						}
					}
				})
				b.Run("case=existing", func(b *testing.B) {
					tree := NewBSTFrom(items)
					b.ReportAllocs()
					for b.Loop() {
						tree.Insert(items[n/2])
					}
				})
			})
		})
	}
}

// BenchmarkBST_Contains measures a hit or miss in a prebuilt shuffled tree
// of n nodes; the expected-result check is measured.
func BenchmarkBST_Contains(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		tree := NewBSTFrom(items)
		for _, tc := range []struct {
			name  string
			value int
			want  bool
		}{
			{"hit", items[n/2], true},
			{"miss", -1, false},
		} {
			b.Run("case="+tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if tree.Contains(tc.value) != tc.want {
						b.Fatal("unexpected membership result")
					}
				}
			})
		}
	})
}

// BenchmarkBST_InOrder measures allocating the sorted traversal of a
// prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_InOrder(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.InOrder()
		}
	})
}

// Constructor results escape so the benchmark includes the returned object's
// allocation, as well as any backing storage. Each sink retains only one result.
var (
	benchmarkBSTResult   *BST[int]
	benchmarkQueueResult *Queue[int]
	benchmarkSetResult   *Set[int]
	benchmarkStackResult *Stack[int]
)

// BenchmarkNewBST measures allocating an empty int BST whose pointer
// escapes to a sink. This scalar constructor has no input-size parameter.
func BenchmarkNewBST(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkBSTResult = NewBST[int]()
	}
}

// BenchmarkNewQueue measures allocating an empty int Queue whose pointer
// escapes to a sink. This scalar constructor has no input-size parameter.
func BenchmarkNewQueue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkQueueResult = NewQueue[int]()
	}
}

// BenchmarkNewSet measures allocating an empty int Set whose pointer
// escapes to a sink. This scalar constructor has no input-size parameter.
func BenchmarkNewSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkSetResult = NewSet[int]()
	}
}

// BenchmarkNewStack measures allocating an empty int Stack whose pointer
// escapes to a sink. This scalar constructor has no input-size parameter.
func BenchmarkNewStack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkStackResult = NewStack[int]()
	}
}

// BenchmarkNewQueueWithCapacity measures allocating an empty int Queue
// and its reserved storage; n is the capacity hint. The returned pointer escapes.
func BenchmarkNewQueueWithCapacity(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkQueueResult = NewQueueWithCapacity[int](n)
		}
	})
}

// BenchmarkNewSetWithCapacity measures allocating an empty int Set
// and its reserved storage; n is the capacity hint. The returned pointer escapes.
func BenchmarkNewSetWithCapacity(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSetResult = NewSetWithCapacity[int](n)
		}
	})
}

// BenchmarkNewStackWithCapacity measures allocating an empty int Stack
// and its reserved storage; n is the capacity hint. The returned pointer escapes.
func BenchmarkNewStackWithCapacity(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkStackResult = NewStackWithCapacity[int](n)
		}
	})
}

// BenchmarkNewBSTFrom measures allocating a tree from n input values. The
// shuffled case uses Sizes; the sorted case uses SmallSizes because building
// an unbalanced sorted tree is quadratic. The returned pointer escapes.
func BenchmarkNewBSTFrom(b *testing.B) {
	for _, order := range []struct {
		name   string
		sizes  []int
		values func(int) []int
	}{
		{"shuffled", benchkit.Sizes, benchkit.Ints},
		{"sorted", benchkit.SmallSizes, benchkit.SortedInts},
	} {
		b.Run("order="+order.name, func(b *testing.B) {
			benchkit.Run(b, order.sizes, func(b *testing.B, n int) {
				items := order.values(n)
				b.ReportAllocs()
				for b.Loop() {
					benchmarkBSTResult = NewBSTFrom(items)
				}
			})
		})
	}
}

// BenchmarkNewSetFrom measures allocating a set from n deterministic integer
// elements; these fixtures are distinct. The returned pointer escapes.
func BenchmarkNewSetFrom(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSetResult = NewSetFrom(items)
		}
	})
}

// BenchmarkBST_Remove measures removing a known leaf from a shuffled tree of
// n nodes and inserting it again. The checked Remove and restoring Insert are
// measured. The last distinct inserted value is a leaf, verified in setup, so
// reinsertion preserves the tree shape.
func BenchmarkBST_Remove(b *testing.B) {
	b.Run("case=leaf+reinsert", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			tree := NewBSTFrom(items)
			value := items[n-1]
			leaf := tree.find(value)
			if leaf == nil || leaf.left != nil || leaf.right != nil {
				b.Fatal("removal fixture is not a leaf")
			}
			b.ReportAllocs()
			for b.Loop() {
				if !tree.Remove(value) {
					b.Fatal("leaf unexpectedly missing")
				}
				tree.Insert(value)
			}
		})
	})
}

// BenchmarkBST_Clear measures refilling an empty tree with n shuffled values
// and clearing it. Node allocations during Insert and the size check before
// Clear are measured, so every iteration clears n nodes.
func BenchmarkBST_Clear(b *testing.B) {
	b.Run("case=refill+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			tree := NewBST[int]()
			b.ReportAllocs()
			for b.Loop() {
				tree.Insert(items...)
				if tree.Len() != n {
					b.Fatal("tree refill did not produce n nodes")
				}
				tree.Clear()
			}
		})
	})
}

// BenchmarkBST_Len measures reading the stored node count
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_Len(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.Len()
		}
	})
}

// BenchmarkBST_IsEmpty measures checking a prebuilt nonempty shuffled tree
// for emptiness; n is the number of resident nodes.
func BenchmarkBST_IsEmpty(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.IsEmpty()
		}
	})
}

// BenchmarkBST_PreOrder measures allocating a root-left-right traversal
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_PreOrder(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.PreOrder()
		}
	})
}

// BenchmarkBST_PostOrder measures allocating a left-right-root traversal
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_PostOrder(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.PostOrder()
		}
	})
}

// BenchmarkBST_LevelOrder measures allocating a breadth-first traversal
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_LevelOrder(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.LevelOrder()
		}
	})
}

// BenchmarkBST_Height measures recursively computing the height
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_Height(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.Height()
		}
	})
}

// BenchmarkBST_Values measures allocating the sorted values
// of a prebuilt shuffled tree; n is the number of resident nodes.
func BenchmarkBST_Values(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			tree.Values()
		}
	})
}

// BenchmarkBST_Min measures finding the minimum in a prebuilt
// shuffled tree of n nodes; the success check is included.
func BenchmarkBST_Min(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := tree.Min(); !ok {
				b.Fatal("tree unexpectedly empty")
			}
		}
	})
}

// BenchmarkBST_Max measures finding the maximum in a prebuilt
// shuffled tree of n nodes; the success check is included.
func BenchmarkBST_Max(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		tree := NewBSTFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := tree.Max(); !ok {
				b.Fatal("tree unexpectedly empty")
			}
		}
	})
}

// BenchmarkBST_RangeSearch measures collecting the numeric middle 10% of
// the value range in a prebuilt shuffled tree; n is the number of resident
// nodes. Both range bounds are inclusive and prepared outside the loop.
func BenchmarkBST_RangeSearch(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		tree := NewBSTFrom(items)
		low, high := slices.Min(items), slices.Max(items)
		span := high - low
		mid := low + span/2
		low, high = mid-span/20, mid+span/20
		b.ReportAllocs()
		for b.Loop() {
			tree.RangeSearch(low, high)
		}
	})
}

// BenchmarkBST_ForEach measures visiting n shuffled tree nodes and
// summing their integer values. The accumulator reset and callback additions
// are measured; the final checksum is checked outside the loop.
func BenchmarkBST_ForEach(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		tree := NewBSTFrom(items)
		want := 0
		for _, v := range items {
			want += v
		}
		sum := 0
		visit := func(v int) { sum += v }
		b.ReportAllocs()
		for b.Loop() {
			sum = 0
			tree.ForEach(visit)
		}
		if sum != want {
			b.Fatal("callback checksum mismatch")
		}
	})
}

// BenchmarkSet_ForEach measures visiting n set elements and
// summing their integer values. The accumulator reset and callback additions
// are measured; the final checksum is checked outside the loop.
func BenchmarkSet_ForEach(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		s := NewSetFrom(items)
		want := 0
		for _, v := range items {
			want += v
		}
		sum := 0
		visit := func(v int) { sum += v }
		b.ReportAllocs()
		for b.Loop() {
			sum = 0
			s.ForEach(visit)
		}
		if sum != want {
			b.Fatal("callback checksum mismatch")
		}
	})
}

// BenchmarkQueue_Peek measures reading the front of a prefilled
// queue of n items; the success check is measured.
func BenchmarkQueue_Peek(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		q := NewQueueWithCapacity[int](n)
		q.EnqueueAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := q.Peek(); !ok {
				b.Fatal("queue unexpectedly empty")
			}
		}
	})
}

// BenchmarkQueue_Len measures reading the stored length;
// n is the number of integer items in the prefilled queue.
func BenchmarkQueue_Len(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		q := NewQueueWithCapacity[int](n)
		q.EnqueueAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			q.Len()
		}
	})
}

// BenchmarkQueue_IsEmpty measures checking a nonempty structure for emptiness;
// n is the number of integer items in the prefilled queue.
func BenchmarkQueue_IsEmpty(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		q := NewQueueWithCapacity[int](n)
		q.EnqueueAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			q.IsEmpty()
		}
	})
}

// BenchmarkQueue_Values measures allocating a FIFO copy of a
// prefilled queue; n is the number of resident integer items.
func BenchmarkQueue_Values(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		q := NewQueueWithCapacity[int](n)
		q.EnqueueAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			q.Values()
		}
	})
}

// BenchmarkQueue_Clear measures appending n integer items to an empty
// queue with capacity n and clearing it again. The refill, size check,
// and Clear are measured; the backing storage is reused and state is bounded.
func BenchmarkQueue_Clear(b *testing.B) {
	b.Run("case=refill+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			q := NewQueueWithCapacity[int](n)
			b.ReportAllocs()
			for b.Loop() {
				q.EnqueueAll(items...)
				if q.Len() != n {
					b.Fatal("queue refill did not produce n items")
				}
				q.Clear()
			}
		})
	})
}

// BenchmarkQueue_EnqueueAll measures appending n integer items to an empty
// queue with capacity n and clearing it again. The refill, size check,
// and Clear are measured; the backing storage is reused and state is bounded.
func BenchmarkQueue_EnqueueAll(b *testing.B) {
	b.Run("case=append+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			q := NewQueueWithCapacity[int](n)
			b.ReportAllocs()
			for b.Loop() {
				q.EnqueueAll(items...)
				if q.Len() != n {
					b.Fatal("queue refill did not produce n items")
				}
				q.Clear()
			}
		})
	})
}

// BenchmarkStack_Peek measures reading the top of a prefilled
// stack of n items; the success check is measured.
func BenchmarkStack_Peek(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewStackWithCapacity[int](n)
		s.PushAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			if _, ok := s.Peek(); !ok {
				b.Fatal("stack unexpectedly empty")
			}
		}
	})
}

// BenchmarkStack_Len measures reading the stored length;
// n is the number of integer items in the prefilled stack.
func BenchmarkStack_Len(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewStackWithCapacity[int](n)
		s.PushAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			s.Len()
		}
	})
}

// BenchmarkStack_IsEmpty measures checking a nonempty structure for emptiness;
// n is the number of integer items in the prefilled stack.
func BenchmarkStack_IsEmpty(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewStackWithCapacity[int](n)
		s.PushAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			s.IsEmpty()
		}
	})
}

// BenchmarkStack_Values measures allocating a bottom-to-top copy of a
// prefilled stack; n is the number of resident integer items.
func BenchmarkStack_Values(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewStackWithCapacity[int](n)
		s.PushAll(benchkit.Ints(n)...)
		b.ReportAllocs()
		for b.Loop() {
			s.Values()
		}
	})
}

// BenchmarkStack_Clear measures appending n integer items to an empty
// stack with capacity n and clearing it again. The refill, size check,
// and Clear are measured; the backing storage is reused and state is bounded.
func BenchmarkStack_Clear(b *testing.B) {
	b.Run("case=refill+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewStackWithCapacity[int](n)
			b.ReportAllocs()
			for b.Loop() {
				s.PushAll(items...)
				if s.Len() != n {
					b.Fatal("stack refill did not produce n items")
				}
				s.Clear()
			}
		})
	})
}

// BenchmarkStack_PushAll measures appending n integer items to an empty
// stack with capacity n and clearing it again. The refill, size check,
// and Clear are measured; the backing storage is reused and state is bounded.
func BenchmarkStack_PushAll(b *testing.B) {
	b.Run("case=append+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewStackWithCapacity[int](n)
			b.ReportAllocs()
			for b.Loop() {
				s.PushAll(items...)
				if s.Len() != n {
					b.Fatal("stack refill did not produce n items")
				}
				s.Clear()
			}
		})
	})
}

// BenchmarkQueue_DequeueN measures draining n items from a prefilled queue
// and refilling it with the original items. The returned-length check and
// EnqueueAll are measured. DequeueN consumes the slice prefix, so the refill
// includes backing-array allocation in addition to the returned slice.
func BenchmarkQueue_DequeueN(b *testing.B) {
	b.Run("case=drain+refill", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			q := NewQueueWithCapacity[int](n)
			q.EnqueueAll(items...)
			b.ReportAllocs()
			for b.Loop() {
				if got := q.DequeueN(n); len(got) != n {
					b.Fatal("queue did not contain n items")
				}
				q.EnqueueAll(items...)
			}
		})
	})
}

// BenchmarkStack_PopAll measures draining n items from a prefilled stack
// in LIFO order and refilling it with the original input order. The length
// check and PushAll are measured. The result allocates; refill reuses capacity.
func BenchmarkStack_PopAll(b *testing.B) {
	b.Run("case=drain+refill", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewStackWithCapacity[int](n)
			s.PushAll(items...)
			b.ReportAllocs()
			for b.Loop() {
				if got := s.PopAll(); len(got) != n {
					b.Fatal("stack did not contain n items")
				}
				s.PushAll(items...)
			}
		})
	})
}

// BenchmarkStack_PopN measures draining n items from a prefilled stack
// in LIFO order and refilling it with the original input order. The length
// check and PushAll are measured. The result allocates; refill reuses capacity.
func BenchmarkStack_PopN(b *testing.B) {
	b.Run("case=drain+refill", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewStackWithCapacity[int](n)
			s.PushAll(items...)
			b.ReportAllocs()
			for b.Loop() {
				if got := s.PopN(n); len(got) != n {
					b.Fatal("stack did not contain n items")
				}
				s.PushAll(items...)
			}
		})
	})
}

// BenchmarkSet_Remove measures removing a resident element from a set of n
// elements and adding it back. Contains verifies the destructive precondition
// each iteration; both that check and the restoring Add are measured.
func BenchmarkSet_Remove(b *testing.B) {
	b.Run("case=remove+readd", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewSetFrom(items)
			value := items[n/2]
			b.ReportAllocs()
			for b.Loop() {
				if !s.Contains(value) {
					b.Fatal("element to remove unexpectedly missing")
				}
				s.Remove(value)
				s.Add(value)
			}
		})
	})
}

// BenchmarkSet_Clear measures refilling an empty set with n integer elements
// and clearing it. Add, the size check, and Clear are measured, including map
// allocation and growth on refill because Clear replaces the backing map.
func BenchmarkSet_Clear(b *testing.B) {
	b.Run("case=refill+clear", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			s := NewSet[int]()
			b.ReportAllocs()
			for b.Loop() {
				s.Add(items...)
				if s.Len() != n {
					b.Fatal("set refill did not produce n elements")
				}
				s.Clear()
			}
		})
	})
}

// BenchmarkSet_Len measures reading the stored cardinality;
// n is the number of deterministic integer elements in the prebuilt set.
func BenchmarkSet_Len(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			s.Len()
		}
	})
}

// BenchmarkSet_IsEmpty measures checking a nonempty set for emptiness;
// n is the number of deterministic integer elements in the prebuilt set.
func BenchmarkSet_IsEmpty(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			s.IsEmpty()
		}
	})
}

// BenchmarkSet_Values measures allocating an unordered slice of the elements;
// n is the number of deterministic integer elements in the prebuilt set.
func BenchmarkSet_Values(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			s.Values()
		}
	})
}

// BenchmarkSet_Clone measures allocating an independent copy of the set;
// n is the number of deterministic integer elements in the prebuilt set.
func BenchmarkSet_Clone(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			s.Clone()
		}
	})
}

// BenchmarkSet_ContainsAll measures testing n queries against a set of n
// integer elements. case=all uses every resident element; case=missing-last
// replaces only the last query with an absent value, so both scan n queries.
// The expected-result check is measured.
func BenchmarkSet_ContainsAll(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		s := NewSetFrom(items)
		missing := slices.Clone(items)
		missing[n-1] = -1
		for _, tc := range []struct {
			name    string
			queries []int
			want    bool
		}{
			{"all", items, true},
			{"missing-last", missing, false},
		} {
			b.Run("case="+tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if s.ContainsAll(tc.queries...) != tc.want {
						b.Fatal("unexpected ContainsAll result")
					}
				}
			})
		}
	})
}

// BenchmarkSet_ContainsAny measures testing n queries against a set of n
// integer elements. Negative queries all miss; case=last-hit places a resident
// value last, so both cases scan all n queries. The result check is measured.
func BenchmarkSet_ContainsAny(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		s := NewSetFrom(items)
		missing := make([]int, n)
		for i, value := range items {
			missing[i] = -value - 1
		}
		lastHit := slices.Clone(missing)
		lastHit[n-1] = items[n/2]
		for _, tc := range []struct {
			name    string
			queries []int
			want    bool
		}{
			{"last-hit", lastHit, true},
			{"miss", missing, false},
		} {
			b.Run("case="+tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if s.ContainsAny(tc.queries...) != tc.want {
						b.Fatal("unexpected ContainsAny result")
					}
				}
			})
		}
	})
}

// BenchmarkSet_Union measures Union on two prebuilt integer
// sets with exactly 50% overlap; n is the number of elements in each operand.
// Building the result is measured; both input sets remain unchanged.
func BenchmarkSet_Union(b *testing.B) {
	b.Run("case=half-overlap", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			other := slices.Clone(items)
			for i := n / 2; i < n; i++ {
				other[i] = -other[i] - 1
			}
			a, c := NewSetFrom(items), NewSetFrom(other)
			b.ReportAllocs()
			for b.Loop() {
				a.Union(c)
			}
		})
	})
}

// BenchmarkSet_Intersection measures Intersection on two prebuilt integer
// sets with exactly 50% overlap; n is the number of elements in each operand.
// Building the result is measured; both input sets remain unchanged.
func BenchmarkSet_Intersection(b *testing.B) {
	b.Run("case=half-overlap", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			other := slices.Clone(items)
			for i := n / 2; i < n; i++ {
				other[i] = -other[i] - 1
			}
			a, c := NewSetFrom(items), NewSetFrom(other)
			b.ReportAllocs()
			for b.Loop() {
				a.Intersection(c)
			}
		})
	})
}

// BenchmarkSet_Difference measures Difference on two prebuilt integer
// sets with exactly 50% overlap; n is the number of elements in each operand.
// Building the result is measured; both input sets remain unchanged.
func BenchmarkSet_Difference(b *testing.B) {
	b.Run("case=half-overlap", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			other := slices.Clone(items)
			for i := n / 2; i < n; i++ {
				other[i] = -other[i] - 1
			}
			a, c := NewSetFrom(items), NewSetFrom(other)
			b.ReportAllocs()
			for b.Loop() {
				a.Difference(c)
			}
		})
	})
}

// BenchmarkSet_SymmetricDifference measures SymmetricDifference on two prebuilt integer
// sets with exactly 50% overlap; n is the number of elements in each operand.
// Building the result is measured; both input sets remain unchanged.
func BenchmarkSet_SymmetricDifference(b *testing.B) {
	b.Run("case=half-overlap", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			other := slices.Clone(items)
			for i := n / 2; i < n; i++ {
				other[i] = -other[i] - 1
			}
			a, c := NewSetFrom(items), NewSetFrom(other)
			b.ReportAllocs()
			for b.Loop() {
				a.SymmetricDifference(c)
			}
		})
	})
}

// BenchmarkSet_IsSubset measures a true subset relation between a
// prebuilt set of n integer elements and its first n/2 fixture elements.
// n names the larger set; all n/2 subset entries are tested, including the
// expected-true result check.
func BenchmarkSet_IsSubset(b *testing.B) {
	b.Run("case=subset", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			full, subset := NewSetFrom(items), NewSetFrom(items[:n/2])
			b.ReportAllocs()
			for b.Loop() {
				if !subset.IsSubset(full) {
					b.Fatal("unexpected set relation")
				}
			}
		})
	})
}

// BenchmarkSet_IsSuperset measures a true superset relation between a
// prebuilt set of n integer elements and its first n/2 fixture elements.
// n names the larger set; all n/2 subset entries are tested, including the
// expected-true result check.
func BenchmarkSet_IsSuperset(b *testing.B) {
	b.Run("case=superset", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			items := benchkit.Ints(n)
			full, subset := NewSetFrom(items), NewSetFrom(items[:n/2])
			b.ReportAllocs()
			for b.Loop() {
				if !full.IsSuperset(subset) {
					b.Fatal("unexpected set relation")
				}
			}
		})
	})
}

// BenchmarkSet_Equal measures comparing two sets of n integer elements.
// case=equal scans all entries of independently built equal sets; case=different
// uses equally sized sets with 50% overlap, avoiding the size-only shortcut.
// The expected-result check is measured.
func BenchmarkSet_Equal(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		a := NewSetFrom(items)
		other := slices.Clone(items)
		for i := n / 2; i < n; i++ {
			other[i] = -other[i] - 1
		}
		for _, tc := range []struct {
			name string
			set  *Set[int]
			want bool
		}{
			{"equal", NewSetFrom(items), true},
			{"different", NewSetFrom(other), false},
		} {
			b.Run("case="+tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if a.Equal(tc.set) != tc.want {
						b.Fatal("unexpected equality result")
					}
				}
			})
		}
	})
}

// BenchmarkSet_Filter measures allocating a set containing the even-valued
// elements (about half) of a prebuilt set of n deterministic integer elements.
func BenchmarkSet_Filter(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			s.Filter(func(v int) bool { return v%2 == 0 })
		}
	})
}

// BenchmarkToSliceSorted measures copying and sorting a prebuilt set of n
// deterministic integer elements. The set itself remains unchanged.
func BenchmarkToSliceSorted(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		s := NewSetFrom(benchkit.Ints(n))
		b.ReportAllocs()
		for b.Loop() {
			ToSliceSorted(s)
		}
	})
}
