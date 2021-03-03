package processing

import (
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"testing"
)

func TestDetectCarriageReturn(t *testing.T) {
	files := fileutils.ListFile(".")
	for _, f := range files {
		carriageReturn, err := DetectLineTerminator(f)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		t.Logf("File %s: %s\n", f, carriageReturn)
	}
}
