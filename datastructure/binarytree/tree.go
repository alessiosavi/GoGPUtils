package binarytree

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

func (t *Tree) VisitPreOrder() []int {
	if t.Root == nil {
		panic("Empty array")
	}
	var result []int
	t.Root.visitPreOrder(&result)
	return result
}

func (t *Tree) VisitPostOrder() []int {
	if t.Root == nil {
		panic("Empty array")
	}
	var result []int
	t.Root.visitPostOrder(&result)
	return result
}

func (t *Tree) InitTree(val int) {
	t.Root = initNode(val)
}

func (t *Tree) Insert(val int) {
	t.Root.insert(val)
}
