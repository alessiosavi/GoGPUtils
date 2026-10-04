package collection

import (
	"math/rand"
	"testing"
)

func workspaceBalanced(lo, hi int, out *[]int) {
	if lo >= hi {
		return
	}
	m := (lo + hi) / 2
	*out = append(*out, m)
	workspaceBalanced(lo, m, out)
	workspaceBalanced(m+1, hi, out)
}

func TestIndexQueueDecisions(t *testing.T) {
	nodes := make([]bstNode[int], 100)
	q := bfsQueue[int]{buf: make([]*bstNode[int], 0, 64)}
	for i := range 64 {
		q.push(&nodes[i])
	}
	for range 32 {
		old := q.head
		q.pop()
		if q.buf[old] != nil {
			t.Fatal("pop did not clear")
		}
	}
	q.push(&nodes[64])
	same(t, cap(q.buf), 64)
	same(t, q.head, 0)
	same(t, len(q.buf), 33)
	for q.head < len(q.buf) {
		q.pop()
	}
	for _, p := range q.buf[:cap(q.buf)] {
		if p != nil {
			t.Fatal("stale node in workspace")
		}
	}
	q = bfsQueue[int]{buf: make([]*bstNode[int], 0, 64)}
	for i := range 64 {
		q.push(&nodes[i])
	}
	for range 31 {
		q.pop()
	}
	q.push(&nodes[64])
	same(t, cap(q.buf), 128)
}

func exerciseIndex(t *testing.T, root *bstNode[int], size int) {
	t.Helper()
	if root == nil {
		return
	}
	q := bfsQueue[int]{buf: make([]*bstNode[int], 0, min(size, 64))}
	q.push(root)
	frontier, maxcap := 1, cap(q.buf)
	var got []int
	for q.head < len(q.buf) {
		v := q.pop()
		got = append(got, v.value)
		if v.left != nil {
			q.push(v.left)
		}
		if v.right != nil {
			q.push(v.right)
		}
		frontier = max(frontier, len(q.buf)-q.head)
		maxcap = max(maxcap, cap(q.buf))
		if maxcap > max(64, 4*(frontier+1)) {
			t.Fatalf("cap=%d frontier=%d", maxcap, frontier)
		}
	}
	if len(got) != size {
		t.Fatalf("visited=%d want=%d", len(got), size)
	}
	t.Logf("nodes=%d maximum frontier=%d maximum capacity=%d", size, frontier, maxcap)
}

func TestIndexQueueShapes(t *testing.T) {
	for _, n := range []int{1, 8, 64, 97, 512, 1024, 65536} {
		var vals []int
		workspaceBalanced(0, n, &vals)
		tree := NewBSTFrom(vals)
		exerciseIndex(t, tree.root, tree.size)
	}
	for _, n := range []int{8, 64, 512} {
		tree := NewBST[int]()
		for i := range n {
			tree.Insert(i)
		}
		exerciseIndex(t, tree.root, tree.size)
	}
	r := rand.New(rand.NewSource(42))
	for range 10 {
		tree := NewBSTFrom(r.Perm(1024))
		exerciseIndex(t, tree.root, tree.size)
	}
	counts := append(append(append(make([]int, 32), make([]int, 30)...), 2), make([]int, 34)...)
	for i := range 32 {
		counts[i] = 2
	}
	for i := 32; i < 62; i++ {
		counts[i] = 1
	}
	nodes := make([]bstNode[int], 97)
	next := 1
	for i, c := range counts {
		nodes[i].value = i
		if c > 0 {
			nodes[i].left = &nodes[next]
			next++
		}
		if c > 1 {
			nodes[i].right = &nodes[next]
			next++
		}
	}
	exerciseIndex(t, &nodes[0], 97)
}

func TestIndexQueueCompactionTail(t *testing.T) {
	nodes := make([]bstNode[int], 65)
	q := bfsQueue[int]{buf: make([]*bstNode[int], 0, 64)}
	for i := range 64 {
		q.push(&nodes[i])
	}
	for range 32 {
		q.pop()
	}
	q.push(&nodes[64])
	if len(q.buf) != 33 || cap(q.buf) != 64 || q.head != 0 {
		t.Fatalf("did not compact: len=%d cap=%d head=%d", len(q.buf), cap(q.buf), q.head)
	}
	for i, p := range q.buf[len(q.buf):cap(q.buf)] {
		if p != nil {
			t.Fatalf("stale compaction tail at %d", len(q.buf)+i)
		}
	}
}

func TestIndexQueueBootstrap(t *testing.T) {
	var q bfsQueue[int]
	node := &bstNode[int]{value: 7}
	q.push(node)
	same(t, cap(q.buf), 8)
	same(t, q.pop(), node)
	same(t, q.buf[0], (*bstNode[int])(nil))
}

func TestBSTTraversalAllocations(t *testing.T) {
	for _, n := range []int{16, 1024, 65536} {
		var values []int
		workspaceBalanced(0, n, &values)
		tree := NewBSTFrom(values)
		for _, traversal := range []struct {
			name string
			run  func() []int
		}{
			{"InOrder", tree.InOrder}, {"PreOrder", tree.PreOrder},
			{"PostOrder", tree.PostOrder}, {"Values", tree.Values},
		} {
			allocs := testing.AllocsPerRun(10, func() { traversal.run() })
			t.Logf("%s n=%d AllocsPerRun=%g", traversal.name, n, allocs)
			if allocs != 1 {
				t.Fatalf("%s n=%d allocated %g times", traversal.name, n, allocs)
			}
		}
	}
}
