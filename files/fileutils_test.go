package fileutils

import (
	"strings"
	"testing"
)

func TestCountLinesFile(t *testing.T) {

	file := `../tests/files/test1.txt`
	lines, err := CountLinesFile(file, -1)
	if err != nil || lines != 112 {
		t.Log(err)
		t.Fail()
	}

	_, err = CountLinesFile(file+"test", -1)
	if err == nil {
		t.Log(err)
		t.Fail()
	}
}

func BenchmarkCountLinesFile(b *testing.B) {
	file := `../tests/files/test1.txt`
	for i := 0; i < b.N; i++ {
		CountLinesFile(file, -1)
	}
}

func TestGetFileContentTypeKO(t *testing.T) {
	file := `../tests/files/test.txt`
	_, err := GetFileContentType(file)

	if err == nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	t.Log(err)
}

func TestGetFileContentTypeTXT(t *testing.T) {
	file := `../tests/files/test1.txt`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if !strings.Contains(fileType, "text/plain") {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypePDF(t *testing.T) {
	file := `../tests/files/test2.pdf`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	t.Log(fileType)
}

func TestGetFileContentTypeZIP(t *testing.T) {
	file := `../tests/ziputils/test1.zip`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "application/zip" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypeODT(t *testing.T) {
	file := `../tests/files/test3.odt`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "application/odt" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypeDOCX(t *testing.T) {
	file := `../tests/files/test4.docx`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "application/docx" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypeDOC(t *testing.T) {
	file := `../tests/files/test5.doc`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "application/doc" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypePickle(t *testing.T) {
	file := `../tests/files/test6.pkl`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "application/pickle" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypeMP4(t *testing.T) {
	file := `../tests/files/test7.mp4`
	fileType, err := GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "video/mp4" {
		t.Log(fileType)
		t.Fail()
	}

	file = `../tests/files/test8.mp4`
	fileType, err = GetFileContentType(file)

	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "video/mp4" {
		t.Log(fileType)
		t.Fail()
	}
}

func TestGetFileContentTypeBIN(t *testing.T) {
	file := `../tests/files/test9`
	fileType, err := GetFileContentType(file)
	if err != nil {
		t.Log("Error -> ", err)
		t.Fail()
	}
	if fileType != "elf/binary" {
		t.Log(fileType)
		t.Fail()
	}
}
