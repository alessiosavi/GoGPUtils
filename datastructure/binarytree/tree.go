package binarytree

import (
	"fmt"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
)

// Node is the atomic struct of the Tree
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// Tree wrap the node structure
type Tree struct {
	Root *Node
}

func initNode(val int) *Node {
	return &Node{Value: val, Left: nil, Right: nil}
}

// Print is delegated to print the current node
func (n *Node) Print() string {
	return fmt.Sprintf("%+v\n", *n)
}

func (n *Node) insert(val int) {
	if val <= n.Value {
		if n.Left == nil {
			n.Left = initNode(val)
		} else {
			n.Left.insert(val)
		}
	} else {
		if n.Right == nil {
			n.Right = initNode(val)
		} else {
			n.Right.insert(val)
		}
	}
}

func (n *Node) visitPreOrder(array *[]int) []int {
	if n.Left != nil {
		n.Left.visitPreOrder(array)
	}
	*array = append(*array, n.Value)
	if n.Right != nil {
		n.Right.visitPreOrder(array)
	}
	return *array
}

func (n *Node) visitPostOrder(array *[]int) []int {
	if n.Right != nil {
		n.Right.visitPostOrder(array)
	}
	*array = append(*array, n.Value)
	if n.Left != nil {
		n.Left.visitPostOrder(array)
	}
	return *array
}

func (n *Node) visitInOrder(array *[]int) {
	if n.Left != nil {
		n.Left.visitInOrder(array)
	}
	// fmt.Printf("%d ", n.Value)
	*array = append(*array, n.Value)
	if n.Right != nil {
		n.Right.visitInOrder(array)
	}
}

// VisitInOrder is delegated to traverse the Tree in order
func (t *Tree) VisitInOrder() []int {
	if t.Root == nil {
		panic("Empty array")
	}
	var result []int
	t.Root.visitInOrder(&result)
	return result
}

// VisitPreOrder is delegated to traverse the Tree in pre-order
func (t *Tree) VisitPreOrder() []int {
	if t.Root == nil {
		panic("Empty array")
	}
	var result []int
	t.Root.visitPreOrder(&result)
	return result
}

// VisitPostOrder is delegated to traverse the Tree in post order
func (t *Tree) VisitPostOrder() []int {
	if t.Root == nil {
		panic("Empty array")
	}
	var result []int
	t.Root.visitPostOrder(&result)
	return result
}

// InitTree is delegated to initialize a new Tree
func (t *Tree) InitTree(val int) {
	t.Root = initNode(val)
}

// Insert is delegated to insert a new node into the Tree
func (t *Tree) Insert(val int) {
	t.Root.insert(val)
}

// Remove is delegated to remove a node with the given value from the Tree
func (t *Tree) Remove(val int) {
	t.Root.remove(val)
}

func (n *Node) remove(val int) {
	if n.Value == val {
		var tmp *Node
		if n.Left != nil {
			tmp = n.Right
			*n = *n.Left
			n.Right = tmp
		} else {
			tmp = n.Left
			*n = *n.Right
			n.Left = tmp
		}
	} else if val <= n.Value && n.Left != nil {
		n.Left.remove(val)
	} else if val > n.Value && n.Right != nil {
		n.Right.remove(val)
	}
}

// Height is delegated to compute the length of the tree
func (t *Tree) Height() int {
	if t == nil || t.Root == nil {
		return 0
	}
	return t.Root.height()
}

func (n *Node) height() int {
	if n == nil {
		return 0
	}
	lheight := n.Left.height()
	rheight := n.Right.height()

	if lheight > rheight {
		return lheight + 1
	}
	return rheight + 1
}

/*
	Function to print level

order traversal a tree
*/
func (n *Node) print() string {

	h := n.height()
	printMap = make(map[int][]int, h)
	for i := 1; i <= h; i++ {
		n.printByLevel(i)
	}
	for key := h; key > 0; key-- {
		for j := h; j > key; j-- {
			for _, k := range printMap[j] {
				if arrayutils.InInt(printMap[key], k) {
					printMap[key] = arrayutils.RemoveByValue[int](printMap[key], k)
				}
			}
		}
	}

	s := fmt.Sprintf("Tree: %+v", printMap)
	printMap = nil

	return s
}

// printMap is delegated to save the value for the given level
var printMap map[int][]int

/* Print nodes at a given level */
func (n *Node) printByLevel(level int) {
	if n == nil {
		return
	}
	printMap[level] = append(printMap[level], n.Value)

	if level > 1 {
		n.Left.printByLevel(level - 1)
		n.Right.printByLevel(level - 1)
	}
}

func (t *Tree) Print() string {
	return t.Root.print()
}
