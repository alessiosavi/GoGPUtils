package stringutils

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

const danteDataset string = `../testdata/files/dante.txt`

func TestIsUpper(t *testing.T) {
	dataOK := []string{`AAA`, `BBB`, `ZZZ`}
	dataKO := []string{`aaa`, `bbb`, `zzz`, `<<<`, `!!!`}
	for i := range dataOK {
		if !IsUpper(dataOK[i]) {
			t.Fail()
		}
	}
	for i := range dataKO {
		if IsUpper(dataKO[i]) {
			t.Fail()
		}
	}
}

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

func BenchmarkTestIsUpperOK(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.ToUpper(string(content))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IsUpper(data)
	}
}

func BenchmarkTestIsLowerOK(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.ToLower(string(content))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IsLower(data)
	}
}
func BenchmarkTestIsLowerByteKO(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IsLowerByte(content)
	}
}

func BenchmarkCreateJSON(t *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.Split(string(content), "\n")
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		CreateJSON(data)
	}
}

func BenchmarkJoin(t *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.Split(string(content), "\n")
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Join(data)
	}
}

func TestIsLower(t *testing.T) {
	dataOK := []string{`AAA`, `BBB`, `ZZZ`}
	dataKO := []string{`aaa`, `bbb`, `zzz`}
	for i := range dataOK {
		if IsLower(dataOK[i]) {
			t.Fail()
		}
	}
	for i := range dataKO {
		if !IsLower(dataKO[i]) {
			t.Fail()
		}
	}
}
func TestContainsLetter(t *testing.T) {
	dataOK := []string{`baaaa`, `baa!aaa`, `baa aaa`, `!!!!a!!!`}
	dataKO := []string{`....`, `,,,,,`, `,,,,,,,`, `<<<`, `!!!`}
	for i := range dataOK {
		if !ContainsLetter(dataOK[i]) {
			t.Fail()
		}
	}
	for i := range dataKO {
		if ContainsLetter(dataKO[i]) {
			t.Fail()
		}
	}
}

func TestContainsOnlyLetter(t *testing.T) {
	dataOK := []string{`baaaa`, `baaaaa`, `baaa`, `a`}
	dataKO := []string{`....`, `,,,,,`, `<<<`, `!!!`, `2`}
	for i := range dataOK {
		if !ContainsOnlyLetter(dataOK[i]) {
			t.Error(dataOK[i])
		}
	}
	for i := range dataKO {
		if ContainsOnlyLetter(dataKO[i]) {
			t.Error(dataKO[i])
		}
	}
}

func BenchmarkContainsOnlyLetter(t *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		ContainsOnlyLetter(data)
	}
}

func TestRemoveNonAscii(t *testing.T) {
	dataOK := []string{`AÀ È Ì Ò Ù Ỳ Ǹ ẀA`, `AȨ Ç Ḑ Ģ Ḩ Ķ Ļ Ņ Ŗ ŞA`, `AA♩ ♪ ♫ ♬ ♭ ♮ ♯AA`}
	for _, item := range dataOK {
		t.Log(RemoveNonASCII(item))
	}
}

func TestRemoveFromString(t *testing.T) {
	data := []string{`test1`, `another test`, `another another test`}
	for _, item := range data {
		res := RemoveFromString(item, len(item)-1)
		if res != item[:len(item)-1] {
			t.Log(item)
			t.Error(res)
		}
	}
}

func BenchmarkRemoveFromString(t *testing.B) {
	data := `another another test`
	for n := 0; n < t.N; n++ {
		RemoveFromString(data, len(data)-1)
	}
}

func TestIsBlank(t *testing.T) {
	data := []string{``, ` `, `    `, `	`}
	for _, item := range data {
		if !IsBlank(item) {
			t.Log(item)
			t.Fail()
		}
	}
}

func TestTrimDoubleSpace(t *testing.T) {
	data := []string{`  test`, `test  `, `t  st`}
	for _, item := range data {
		str := Trim(item)
		t.Log("Data ->|"+item+"|Found: |"+str+"| Len: ", len(str))
		if len(str) != 4 {
			t.Fail()
		}
	}
}

func TestTrim(t *testing.T) {
	data := []string{` test`, `test `, `te s`}
	for _, item := range data {
		str := Trim(item)
		t.Log("Data ->|"+item+"|Found: |"+str+"| Len: ", len(str))
		if len(str) != 4 {
			t.Fail()
		}
	}
}

func BenchmarkRandomString(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomString(5000)
	}
}

func TestExtractTextFromQuery(t *testing.T) {}
func BenchmarkExtractTextFromQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestCheckPresence(t *testing.T) {}
func BenchmarkCheckPresence(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkIsUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestIsUpperByte(t *testing.T) {}
func BenchmarkIsUpperByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkIsLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestIsLowerByte(t *testing.T) {}
func BenchmarkIsLowerByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkContainsLetter(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func TestCreateJSON(t *testing.T) {}

func TestJoin(t *testing.T) {}

func TestRemoveWhiteSpaceString(t *testing.T) {}
func BenchmarkRemoveWhiteSpaceString(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestIsASCII(t *testing.T) {}
func BenchmarkIsASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestIsASCIIRune(t *testing.T) {}
func BenchmarkIsASCIIRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func TestSplit(t *testing.T) {}
func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestCountLinesString(t *testing.T) {}
func BenchmarkCountLinesString(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestExtractString(t *testing.T) {}
func BenchmarkExtractString(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestReplaceAtIndex(t *testing.T) {}
func BenchmarkReplaceAtIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestRemoveNonASCII(t *testing.T) {}
func BenchmarkRemoveNonASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkIsBlank(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func BenchmarkTrim(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestRandomString(t *testing.T) {}
