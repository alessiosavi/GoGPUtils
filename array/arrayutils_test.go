package arrayutils

import (
	"reflect"
	"strconv"
	"testing"
)

func TestRemoveElementsFromMatrixByIndex(t *testing.T) {

	var data [][]string = [][]string{
		[]string{"0", "0", "0", "0", "0", "0", "0", "0", "0"},
		[]string{"1", "1", "1", "1", "1", "1", "1", "1", "1"},
		[]string{"2", "2", "2", "2", "2", "2", "2", "2", "2"},
		[]string{"3", "3", "3", "3", "3", "3", "3", "3", "3"},
		[]string{"4", "4", "4", "4", "4", "4", "4", "4", "4"},
		[]string{"5", "5", "5", "5", "5", "5", "5", "5", "5"},
		[]string{"6", "6", "6", "6", "6", "6", "6", "6", "6"},
	}
	var target [][]string = [][]string{
		[]string{"0", "0", "0", "0", "0", "0", "0", "0", "0"},
		[]string{"1", "1", "1", "1", "1", "1", "1", "1", "1"},
		[]string{"4", "4", "4", "4", "4", "4", "4", "4", "4"},
		[]string{"5", "5", "5", "5", "5", "5", "5", "5", "5"},
		[]string{"6", "6", "6", "6", "6", "6", "6", "6", "6"},
	}

	data = RemoveElementsFromMatrixByIndex(data, []int{2, 3})

	if !reflect.DeepEqual(data, target) {
		t.Errorf("Got: %v\nExpected:%+v\n", data, target)
	}

}

func TestRemoveElementsFromString1(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	t.Log("Before remove:", data, "Len:", len(data))
	data = RemoveElementsFromStringByIndex(data, []int{0, 1, 2})
	t.Log("After remove:", data, "Len:", len(data))
	for i := range data {
		if data[i] != strconv.Itoa(i+4) {
			t.Fail()
		}
	}
}

func TestRemoveElementsFromString2(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	t.Log("Before remove:", data, "Len:", len(data))
	data = RemoveElementsFromStringByIndex(data, []int{5, 4, 3})
	t.Log("After remove:", data, "Len:", len(data))
	for i := range data {
		if data[i] != strconv.Itoa(i+1) {
			t.Fail()
		}
	}
}

func TestRemoveElementsFromString3(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	t.Log("Before remove:", data, "Len:", len(data))
	data = RemoveElementsFromStringByIndex(data, []int{1, 2, 3, 5, 4, 3, 0})
	t.Log("After remove:", data, "Len:", len(data))
	for i := range data {
		if data[i] != strconv.Itoa(i+1) {
			t.Fail()
		}
	}
}

func TestRemoveElement(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	lenBefore := len(data)
	t.Log("Before remove:", data, "Len:", len(data))
	data = RemoveElement(data, 0)
	t.Log("After remove:", data, "Len:", len(data))
	if lenBefore-1 != len(data) {
		t.Fail()
	}
}

func TestJoinStrings(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g"}
	joined := JoinStrings(data, "")
	if joined != "abcdefg" {
		t.Fail()
	}
	joined = JoinStrings(data, " ")
	if joined != "a b c d e f g" {
		t.Fail()
	}
}
func TestJoinInts(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	joined := JoinInts(data, "")
	if joined != "0123456" {
		t.Fail()
	}
	joined = JoinInts(data, " ")
	if joined != "0 1 2 3 4 5 6" {
		t.Fail()
	}
}
func TestReverseArrayInt(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	reversed := ReverseArrayInt(data)
	for i, j := len(reversed)-1, 0; i > 0; i, j = i-1, j+1 {
		if reversed[i] != data[j] {
			t.Fail()
		}
	}
}
func TestReverseArrayString(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	reversed := ReverseArrayString(data)
	for i, j := len(reversed)-1, 0; i > 0; i, j = i-1, j+1 {
		if reversed[i] != data[j] {
			t.Fail()
		}
	}
}
func TestRemoveIntByIndex(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	test := []int{1, 2, 3, 4, 5, 6}
	deleted := RemoveIntByIndex(data, 0)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveIntByIndex2(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	test := []int{0, 1, 2, 3, 4, 5}
	deleted := RemoveIntByIndex(data, 6)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveIntByIndex3(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	deleted := RemoveIntByIndex(data, 7)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveIntByIndex4(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	deleted := RemoveIntByIndex(data, -1)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveStringByIndex(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, 0)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveStringByIndex2(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"1", "2", "3", "4", "5"}
	deleted := RemoveStringByIndex(data, 5)
	if !reflect.DeepEqual(deleted, test) {
		t.Log(deleted)
		t.Log(test)
	}
}

func TestRemoveStringByIndex3(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, 7)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveStringByIndex4(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, -1)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}
func TestRemoveStrings(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStrings(data, test)
	if len(deleted) != 0 {
		t.Error(deleted)
	}
}

func TestRemoveStrings1(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"7", "8", "9", "10", "11", "12"}
	deleted := RemoveStrings(data, test)
	if !reflect.DeepEqual(data, deleted) {
		t.Log(data)
		t.Log(deleted)
	}
}
