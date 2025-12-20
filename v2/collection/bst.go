package collection

import "cmp"

// BST is a generic Binary Search Tree.
// Elements must be ordered (satisfy cmp.Ordered constraint).
type BST[T cmp.Ordered] struct {
	root *bstNode[T]
	size int
}

type bstNode[T cmp.Ordered] struct {
	value T
	left  *bstNode[T]
	right *bstNode[T]
}

// NewBST creates a new empty Binary Search Tree.
//
// Example:
//
//	tree := NewBST[int]()
//	tree.Insert(5, 3, 7)
func NewBST[T cmp.Ordered]() *BST[T] {
	return &BST[T]{}
}

// NewBSTFrom creates a BST from a slice.
func NewBSTFrom[T cmp.Ordered](items []T) *BST[T] {
	tree := NewBST[T]()
	tree.Insert(items...)
	return tree
}

// Insert adds one or more values to the tree.
// Duplicate values are ignored.
//
// Example:
//
//	tree.Insert(5, 3, 7, 1, 4)
func (t *BST[T]) Insert(values ...T) {
	for _, v := range values {
		t.insert(v)
	}
}

func (t *BST[T]) insert(value T) {
	if t.root == nil {
		t.root = &bstNode[T]{value: value}
		t.size++
		return
	}

	current := t.root
	for {
		if value < current.value {
			if current.left == nil {
				current.left = &bstNode[T]{value: value}
				t.size++
				return
			}
			current = current.left
		} else if value > current.value {
			if current.right == nil {
				current.right = &bstNode[T]{value: value}
				t.size++
				return
			}
			current = current.right
		} else {
			// Duplicate, ignore
			return
		}
	}
}

// Contains returns true if the value is in the tree.
func (t *BST[T]) Contains(value T) bool {
	return t.find(value) != nil
}

func (t *BST[T]) find(value T) *bstNode[T] {
	current := t.root
	for current != nil {
		if value < current.value {
			current = current.left
		} else if value > current.value {
			current = current.right
		} else {
			return current
		}
	}
	return nil
}

// Remove deletes a value from the tree.
// Returns true if the value was found and removed.
func (t *BST[T]) Remove(value T) bool {
	var parent *bstNode[T]
	current := t.root
	isLeft := false

	// Find the node to remove
	for current != nil && current.value != value {
		parent = current
		if value < current.value {
			current = current.left
			isLeft = true
		} else {
			current = current.right
			isLeft = false
		}
	}

	if current == nil {
		return false
	}

	// Case 1: No children
	if current.left == nil && current.right == nil {
		if parent == nil {
			t.root = nil
		} else if isLeft {
			parent.left = nil
		} else {
			parent.right = nil
		}
	} else if current.left == nil {
		// Case 2: Only right child
		if parent == nil {
			t.root = current.right
		} else if isLeft {
			parent.left = current.right
		} else {
			parent.right = current.right
		}
	} else if current.right == nil {
		// Case 3: Only left child
		if parent == nil {
			t.root = current.left
		} else if isLeft {
			parent.left = current.left
		} else {
			parent.right = current.left
		}
	} else {
		// Case 4: Two children - find in-order successor
		successor := current.right
		successorParent := current

		for successor.left != nil {
			successorParent = successor
			successor = successor.left
		}

		current.value = successor.value

		if successorParent == current {
			successorParent.right = successor.right
		} else {
			successorParent.left = successor.right
		}
	}

	t.size--
	return true
}

// Len returns the number of elements in the tree.
func (t *BST[T]) Len() int {
	return t.size
}

// IsEmpty returns true if the tree has no elements.
func (t *BST[T]) IsEmpty() bool {
	return t.size == 0
}

// Clear removes all elements from the tree.
func (t *BST[T]) Clear() {
	t.root = nil
	t.size = 0
}

// Min returns the minimum value in the tree.
// Returns zero value and false if tree is empty.
func (t *BST[T]) Min() (T, bool) {
	if t.root == nil {
		var zero T
		return zero, false
	}

	current := t.root
	for current.left != nil {
		current = current.left
	}
	return current.value, true
}

// Max returns the maximum value in the tree.
// Returns zero value and false if tree is empty.
func (t *BST[T]) Max() (T, bool) {
	if t.root == nil {
		var zero T
		return zero, false
	}

	current := t.root
	for current.right != nil {
		current = current.right
	}
	return current.value, true
}

// InOrder returns values in sorted order (left, root, right).
//
// Example:
//
//	tree.Insert(5, 3, 7, 1, 4)
//	tree.InOrder() // [1, 3, 4, 5, 7]
func (t *BST[T]) InOrder() []T {
	var result []T
	t.inOrderTraverse(t.root, &result)
	return result
}

func (t *BST[T]) inOrderTraverse(node *bstNode[T], result *[]T) {
	if node == nil {
		return
	}
	t.inOrderTraverse(node.left, result)
	*result = append(*result, node.value)
	t.inOrderTraverse(node.right, result)
}

// PreOrder returns values in pre-order (root, left, right).
func (t *BST[T]) PreOrder() []T {
	var result []T
	t.preOrderTraverse(t.root, &result)
	return result
}

func (t *BST[T]) preOrderTraverse(node *bstNode[T], result *[]T) {
	if node == nil {
		return
	}
	*result = append(*result, node.value)
	t.preOrderTraverse(node.left, result)
	t.preOrderTraverse(node.right, result)
}

// PostOrder returns values in post-order (left, right, root).
func (t *BST[T]) PostOrder() []T {
	var result []T
	t.postOrderTraverse(t.root, &result)
	return result
}

func (t *BST[T]) postOrderTraverse(node *bstNode[T], result *[]T) {
	if node == nil {
		return
	}
	t.postOrderTraverse(node.left, result)
	t.postOrderTraverse(node.right, result)
	*result = append(*result, node.value)
}

// LevelOrder returns values level by level (breadth-first).
func (t *BST[T]) LevelOrder() []T {
	if t.root == nil {
		return nil
	}

	var result []T
	queue := []*bstNode[T]{t.root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node.value)

		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}

	return result
}

// Height returns the height of the tree.
// An empty tree has height 0; a single node has height 1.
func (t *BST[T]) Height() int {
	return t.height(t.root)
}

func (t *BST[T]) height(node *bstNode[T]) int {
	if node == nil {
		return 0
	}
	leftHeight := t.height(node.left)
	rightHeight := t.height(node.right)
	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// ForEach calls a function for each value in sorted order.
func (t *BST[T]) ForEach(fn func(T)) {
	t.forEach(t.root, fn)
}

func (t *BST[T]) forEach(node *bstNode[T], fn func(T)) {
	if node == nil {
		return
	}
	t.forEach(node.left, fn)
	fn(node.value)
	t.forEach(node.right, fn)
}

// Values returns all values in sorted order.
// Alias for InOrder().
func (t *BST[T]) Values() []T {
	return t.InOrder()
}
