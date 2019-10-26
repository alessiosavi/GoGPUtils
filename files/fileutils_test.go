package fileutils

import (
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/helper"
)

const gogputils string = `/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/GoGPUtils`
const codeFolder string = `/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils`

func TestCountLinesFile(t *testing.T) {
	file := `../tests/files/test1.txt`
	lines, err := CountLinesFile(file, "", -1)
	if err != nil || lines != 112 {
		t.Log(err)
		t.Fail()
	}

	_, err = CountLinesFile(file+"test", "", -1)
	if err == nil {
		t.Log(err)
		t.Fail()
	}
}

func BenchmarkCountLinesFile(b *testing.B) {
	file := `../tests/files/test1.txt`
	for i := 0; i < b.N; i++ {
		_, err := CountLinesFile(file, "", -1)
		if err != nil {
			b.Fail()
		}
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

func TestListFile(t *testing.T) {
	t.Log(ListFile(codeFolder))
}

func BenchmarkListFile(t *testing.B) {
	for n := 0; n < t.N; n++ {
		ListFile(codeFolder)
	}
}

func TestFindFilesSensitive(t *testing.T) {
	if len(FindFiles(codeFolder, `FindMe`, false)) != 2 {
		t.Fail()
	}
}

func TestFindFilesInsensitive(t *testing.T) {
	if len(FindFiles(codeFolder, `findme`, false)) != 2 {
		t.Fail()
	}
}

func BenchmarkFindFilesSensitive(t *testing.B) {
	for n := 0; n < t.N; n++ {
		FindFiles(codeFolder, `FindMe`, true)
	}
}
func BenchmarkFindFilesInsensitive(t *testing.B) {
	for n := 0; n < t.N; n++ {
		FindFiles(codeFolder, `findme`, true)
	}
}

func TestGetFileSize(t *testing.T) {
	size, err := GetFileSize(gogputils)
	if err != nil {
		t.Fail()
	}

	kbSize := size / 1024
	t.Log(helper.ByteCountIEC(size))
	t.Log(helper.ByteCountSI(size))
	t.Log(size)
	t.Log(kbSize, "K")
}

func TestGetFileSize2(t *testing.T) {
	size, err := GetFileSize2(gogputils)
	if err != nil {
		t.Fail()
	}
	t.Log(size)
}

func BenchmarkGetFileSize(t *testing.B) {
	for n := 0; n < t.N; n++ {
		_, err := GetFileSize(gogputils)
		if err != nil {
			t.Fail()
		}
	}
}

func BenchmarkGetFileSize2(t *testing.B) {
	for n := 0; n < t.N; n++ {
		_, err := GetFileSize2(gogputils)
		if err != nil {
			t.Fail()
		}
	}
}

func TestFileExists(t *testing.T) {
	if !FileExists(gogputils) {
		t.Fail()
	}
}
