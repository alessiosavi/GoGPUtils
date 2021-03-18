package processing

import (
	"bytes"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"github.com/alessiosavi/GoGPUtils/helper"
	"os"
	"testing"
)

func TestDetectCarriageReturn(t *testing.T) {
	files, err := fileutils.ListFile(".")
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
	fakeCrs := generateFakeDataLineTerminator(CR)
	terminator, err := ReplaceLineTerminator(fakeCrs, []byte(LF))
	if err != nil {
		t.Error(err)
	}
	if bytes.Contains(terminator, []byte(CR)) {
		t.Fail()
	}
}

func generateFakeDataLineTerminator(terminator LineTerminatorType) []byte {
	var fakeDataCR []byte
	for i := 0; i < 100; i++ {
		fakeDataCR = append(fakeDataCR, helper.RandomByte(255)...)
		fakeDataCR = append(fakeDataCR, []byte(terminator)...)
	}
	return fakeDataCR
}
