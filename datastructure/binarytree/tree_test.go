package binarytree

import (
	"sort"
	"testing"
)

func Test_InitTree(t *testing.T) {
	var tree Tree
	tree.InitTree(0)
	if tree.Root == nil || tree.Root.Value != 0 {
		t.Error("Error, root value mismatch")
	}
	if tree.Root.Left != nil {
		t.Error("Expected nil left child")
	}

	if tree.Root.Right != nil {
		t.Error("Expected nil right child")
	}
	t.Log(tree.Root.Print())
}

func Test_InsertLeft(t *testing.T) {
	var tree Tree
	tree.InitTree(0)
	tree.Insert(-1)

	if tree.Root.Left == nil || tree.Root.Left.Value != -1 {
		t.Error("Expected a leaf")
	}
	if tree.Root.Right != nil {
		t.Error("Expected empty leaf")
	}
	t.Log(tree.Root.Print())
}

func Test_InsertRight(t *testing.T) {
	var tree Tree
	tree.InitTree(0)
	tree.Insert(1)

	if tree.Root.Right == nil || tree.Root.Right.Value != 1 {
		t.Error("Expected a leaf")
	}
	if tree.Root.Left != nil {
		t.Error("Expected empty leaf")
	}
	t.Log(tree.Root.Print())
}

func Test_MultipleInsertRight(t *testing.T) {
	var tree Tree
	var node Node
	tree.InitTree(0)
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	if tree.Root.Right == nil || tree.Root.Right.Value != 1 {
		t.Error("Expected a leaf")
	}
	if tree.Root.Left != nil {
		t.Error("Expected empty leaf")
	}
	node = *tree.Root.Right

	if node.Right == nil || node.Right.Value != 2 {
		t.Error("Expected a leaf")
	}
	if node.Left != nil {
		t.Error("Expected empty leaf")
	}
	node = *node.Right

	if node.Right == nil || node.Right.Value != 3 {
		t.Error("Expected a leaf")
	}
	if node.Left != nil {
		t.Error("Expected empty leaf")
	}
	t.Log(tree.Root.Print())
}

func Test_MultipleInsertLeft(t *testing.T) {
	var tree Tree
	var node Node
	tree.InitTree(0)
	tree.Insert(-1)
	tree.Insert(-2)
	tree.Insert(-3)

	if tree.Root.Left == nil || tree.Root.Left.Value != -1 {
		t.Error("Expected a leaf")
	}
	if tree.Root.Right != nil {
		t.Error("Expected empty leaf")
	}
	node = *tree.Root.Left

	if node.Left == nil || node.Left.Value != -2 {
		t.Error("Expected a leaf")
	}
	if node.Right != nil {
		t.Error("Expected empty leaf")
	}
	node = *node.Left

	if node.Left == nil || node.Left.Value != -3 {
		t.Error("Expected a leaf")
	}
	if node.Right != nil {
		t.Error("Expected empty leaf")
	}
	t.Log(tree.Root.Print())
}

func Test_VisitPreOrder(t *testing.T) {

	var tree Tree
	tree.InitTree(0)
	tree.Insert(10)
	tree.Insert(5)
	tree.Insert(20)
	tree.Insert(1)
	tree.Insert(30)
	tree.Insert(15)
	tree.Insert(6)
	tree.Insert(-1)
	res := tree.VisitPreOrder()
	if !sort.SliceIsSorted(res, func(i, j int) bool { return res[i] < res[j] }) {
		t.Error("Slice is not sorted!", res)
	}
	t.Log(tree.Root.Print())
}

func Test_VisitPostOrder(t *testing.T) {

	var tree Tree
	tree.InitTree(0)
	tree.Insert(10)
	tree.Insert(5)
	tree.Insert(20)
	tree.Insert(1)
	tree.Insert(30)
	tree.Insert(15)
	tree.Insert(6)
	tree.Insert(-1)
	res := tree.VisitPostOrder()
	if !sort.SliceIsSorted(res, func(i, j int) bool { return res[i] > res[j] }) {
		t.Error("Slice is not sorted!", res)
	}
	t.Log(tree.Root.Print())
}

func Test_Remove(t *testing.T) {
	var tree Tree
	value := 25

	tree.InitTree(10)
	tree.Insert(5)
	tree.Insert(4)
	tree.Insert(6)
	tree.Insert(20)
	tree.Insert(15)
	tree.Insert(value)
	tree.Insert(21)
	tree.Insert(26)

	res := tree.VisitPreOrder()
	res_n := len(res)
	tree.Remove(value)
	res = tree.VisitPreOrder()
	if res_n-1 != len(res) {
		t.Error("Expected 1 less length of the result")
	}
	for _, v := range res {
		if v == value {
			t.Errorf("Expected no more [%d]\n", value)
		}
	}
	t.Log(tree.Root.Print())
}
