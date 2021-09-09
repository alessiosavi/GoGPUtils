package byteutils

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const danteDataset string = `../testdata/files/dante.txt`

func BenchmarkTestIsUpperByteOK(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	content = bytes.ToUpper(content)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IsUpperByte(content)
	}
}

func BenchmarkTestIsLowerByteKO(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	content = bytes.ToLower(content)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IsLowerByte(content)
	}
}
