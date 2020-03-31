package binarytree

import "fmt"

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

type Tree struct {
	Root *Node
}

func initNode(val int) *Node {
	return &Node{Value: val, Left: nil, Right: nil}
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

func (n *Node) visitPreOrder() {
	if n.Left != nil {
		n.Left.visitPreOrder()
	}
	fmt.Printf(" %d", n.Value)
	if n.Right != nil {
		n.Right.visitPreOrder()
	}
}

func (t *Tree) VisitPreOrder() {
	if t.Root == nil {
		panic("Empty array")
	}
	fmt.Printf("[")
	t.Root.visitPreOrder()
	fmt.Printf("]\n")
}

func (t *Tree) InitTree(val int) {
	t.Root = initNode(val)
}

func (t *Tree) Insert(val int) {
	t.Root.insert(val)
}
