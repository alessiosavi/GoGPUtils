package fileutils

import "testing"

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
