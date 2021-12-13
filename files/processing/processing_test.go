package processingutils

import (
	"bytes"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"os"
	"testing"
)

func TestDetectCarriageReturn(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	files, err := fileutils.ListFiles(".")
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, f := range files {
		fd, err := os.Open(f)
		if err != nil {
			panic(err)
		}
		carriageReturn, err := DetectLineTerminator(fd)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		t.Logf("File %s: %x\n", f, carriageReturn)
	}
}
func TestCR(t *testing.T) {
	var terminators = []LineTerminatorType{CR, LF, CRLF, LFCR, RS}
	for _, terminator := range terminators {
		fakeCrs := generateFakeDataLineTerminator(terminator)
		tt, err := DetectLineTerminator(bytes.NewReader(fakeCrs))
		if err != nil {
			t.Error(err)
		}
		if tt != terminator {
			t.Errorf("Line terminator [%x] not found! Found: [%x]\n", terminator, tt)
		}
	}

}

func generateFakeDataLineTerminator(terminator LineTerminatorType) []byte {
	var fakeDataCR []byte
	var terminators = []LineTerminatorType{CR, LF, CRLF, LFCR, RS}
	var index = 0
	for i := range terminators {
		if terminators[i] == terminator {
			index = i
		}
	}
	terminators = append(terminators[:index], terminators[index+1:]...)
	for i := 0; i < 1000; i++ {
		data := helper.RandomByte(255)
		for _, t := range terminators {
			data = bytes.ReplaceAll(data, []byte(t), []byte("x"))
		}
		fakeDataCR = append(fakeDataCR, data...)
		fakeDataCR = append(fakeDataCR, []byte(terminator)...)
	}
	for _, t := range terminators {
		fakeDataCR = append(fakeDataCR, []byte(t)...)
	}
	return fakeDataCR
}

//func TestToUTF8(t *testing.T) {
//	files, err := fileutils.ListFiles("/tmp/sap/data")
//	if err != nil {
//		panic(err)
//	}
//	for _, f := range files {
//
//		file, err := ioutil.ReadFile(f)
//		if err != nil {
//			panic(err)
//		}
//		file, err = ToUTF8(file)
//		if err != nil {
//			return
//		}
//
//		err = ioutil.WriteFile(f, file, 0755)
//		if err != nil {
//			panic(err)
//		}
//	}
//}
