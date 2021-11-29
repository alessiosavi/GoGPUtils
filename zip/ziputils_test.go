package ziputils

import (
	"archive/zip"
	"testing"
)

func TestReadZip01(t *testing.T) {
	//t.Log("----TestReadZip01---")
	var file = `../testdata/ziputils/test1.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 1 {
		t.Fail()
	}
}

func TestReadZip02(t *testing.T) {
	//t.Log("----TestReadZip02---")
	file := `../testdata/ziputils/test1.zip_error`
	data, err := ReadZip(file)
	if err == nil {
		t.Fail()
	}
	if data != nil {
		t.Fail()
	}
}

func TestReadZip03(t *testing.T) {
	//t.Log("----TestReadZip03---")
	file := `../testdata/ziputils/test2.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 2 {
		t.Fail()
	}
}

func TestReadZip04(t *testing.T) {
	//t.Log("----TestReadZip04---")
	file := `../testdata/ziputils/test3.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 2 {
		t.Fail()
	}
}

func TestReadZip05(t *testing.T) {
	//t.Log("----TestReadZip05---")
	file := `../testdata/ziputils/test4.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 3 {
		t.Fail()
	}
}

func TestReadZip06(t *testing.T) {
	//t.Log("----TestReadZip06---")
	file := `../testdata/ziputils/test5.zip`
	data, err := ReadZip(file)
	if err != nil {
		t.Fail()
	}
	if data == nil || len(data) != 5 {
		t.Fail()
	}
}

func TestReadZipFile(t *testing.T) {
	filename := `../testdata/ziputils/test1.zip`
	zf, err := zip.OpenReader(filename)
	if err != nil {
		t.Error("Error! ", err)
	}
	defer zf.Close()
	for _, file := range zf.File {
		if file.Mode().IsRegular() {
			_, err = ReadZipFile(file)
			if err != nil {
				t.Log("ReadZip | Unable to unzip file " + file.Name)
				t.Fail()
			}
		}
	}
}

func BenchmarkReadZipFile(t *testing.B) {
	filename := `../testdata/ziputils/test1.zip`
	for n := 0; n < t.N; n++ {
		zf, err := zip.OpenReader(filename)
		if err != nil {
			t.Error("Error! ", err)
		}

		for _, file := range zf.File {
			if file.Mode().IsRegular() {
				_, err := ReadZipFile(file)
				if err != nil {
					t.Error("ReadZip | Unable to unzip file " + file.Name)
				}
			}
		}
		zf.Close()
	}
}

func BenchmarkReadZip01(b *testing.B) {
	file := `../testdata/ziputils/test1.zip`
	for n := 0; n < b.N; n++ {
		data, err := ReadZip(file)
		if err != nil {
			b.Fail()
		}
		if data == nil || len(data) != 1 {
			b.Fail()
		}
	}
}
