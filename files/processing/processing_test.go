package processing

import (
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"os"
	"testing"
)

func TestDetectCarriageReturn(t *testing.T) {
	files := fileutils.ListFile(".")
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
