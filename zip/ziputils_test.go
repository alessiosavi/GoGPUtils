package ziputils

import (
	"testing"
)

func TestReadZip01(t *testing.T) {
	t.Log("----TestReadZip01---")
	file := `../tests/ziputils/test1.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	t.Log(data)
	if data == nil || len(data) != 1 {
		t.Fail()
	}
}

func TestReadZip02(t *testing.T) {
	t.Log("----TestReadZip02---")
	file := `../tests/ziputils/test1.zip_error`
	data, err := ReadZip(file)
	if err == nil {
		t.Fail()
	}
	if data != nil {
		t.Fail()
	}
}

func TestReadZip03(t *testing.T) {
	t.Log("----TestReadZip03---")
	file := `../tests/ziputils/test2.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 2 {
		t.Fail()
	}
}

func TestReadZip04(t *testing.T) {
	t.Log("----TestReadZip04---")
	file := `../tests/ziputils/test3.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 2 {
		t.Fail()
	}
}

func TestReadZip05(t *testing.T) {
	t.Log("----TestReadZip05---")
	file := `../tests/ziputils/test4.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 3 {
		t.Fail()
	}
}

func TestReadZip06(t *testing.T) {
	t.Log("----TestReadZip06---")
	file := `../tests/ziputils/test5.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 5 {
		t.Fail()
	}
}
