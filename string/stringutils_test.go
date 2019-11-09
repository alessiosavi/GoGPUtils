package stringutils

import (
	"io/ioutil"
	"strings"
	"testing"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
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

func TestRemoveDoubleWhiteSpace(t *testing.T) {
	data := []string{`  test`, `test  `, `te  st`}
	for _, item := range data {
		str := RemoveDoubleWhiteSpace(item)
		t.Log("Data ->|"+item+"|Found: |"+str+"| Len: ", len(str))
		if len(str) != 5 {
			t.Fail()
		}
	}
}

func TestIsASCII(t *testing.T) {
	data := "This is a simple! Text< \n"
	if !IsASCII(data) {
		t.Fail()
	}
}

func TestCountLines(t *testing.T) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	lines := CountLines(data)
	if lines != 19567 {
		t.Error(lines)
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
func TestTrimNewLine(t *testing.T) {
	data := []string{"test\n", "\ntest", "\ntest\n", "\ntest \n"}
	for _, item := range data {
		str := Trim(item)
		t.Log("Data ->|"+item+"|Found: |"+str+"| Len: ", len(str))
		if len(str) != 4 {
			t.Fail()
		}
	}
}

func TestCheckPalindrome(t *testing.T) {
	data := []string{`aba`, `abba`, `abcba`}
	for _, item := range data {
		if !CheckPalindrome(item) {
			t.Fail()
		}
	}
}

func TestReverseString(t *testing.T) {
	data := "Golang is better than Java <3"
	test := `3< avaJ naht retteb si gnaloG`
	if ReverseString(data) != test {
		t.Fail()
	}
}
func TestExtractTextFromQuery(t *testing.T) {
	data := "Golang is better than java :D"
	ignore := []string{"java"}
	if len(ExtractTextFromQuery(data, nil)) != 6 {
		t.Fail()
	}
	if len(ExtractTextFromQuery(data, ignore)) != 5 {
		t.Fail()
	}

}
func TestCheckPresence(t *testing.T) {
	data := []string{"test1", "test2", "test3", "test4"}
	if !CheckPresence(data, "test4") {
		t.Fail()
	}
	if CheckPresence(data, "test5") {
		t.Fail()
	}
}

func TestLevenshteinDistanceLegacy(t *testing.T) {
	str1 := `kitten kitten kitten kitten kitten kitten`
	str2 := `sitting sitting sitting sitting sitting`
	distance := LevenshteinDistanceLegacy(str1, str2)
	if distance != 21 {
		t.Error(distance)
	}
}

func BenchmarkLevenshteinDistanceLegacy(t *testing.B) {
	str1 := `kitten kitten kitten kitten kitten kitten`
	str2 := `sitting sitting sitting sitting sitting`
	for i := 0; i < t.N; i++ {
		LevenshteinDistanceLegacy(str1, str2)
	}
}
func BenchmarkLevenshteinDistance(t *testing.B) {
	str1 := `kitten kitten kitten kitten kitten kitten`
	str2 := `sitting sitting sitting sitting sitting`
	for i := 0; i < t.N; i++ {
		LevenshteinDistance(str1, str2)
	}
}

func TestLevenshteinDistance(t *testing.T) {
	str1 := `kitten kitten kitten kitten kitten kitten`
	str2 := `sitting sitting sitting sitting sitting`
	distance := LevenshteinDistance(str1, str2)
	if distance != 21 {
		t.Error(distance)
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

func BenchmarkRemoveFromString(t *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		RemoveFromString(data, len(data)-1)
	}
}

func BenchmarkRandomString(t *testing.B) {
	for n := 0; n < t.N; n++ {
		RandomString(5000)
	}
}

func BenchmarkExtractTextFromQuery(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractTextFromQuery(data, nil)
	}
}
func BenchmarkCheckPresence(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPresence([]string{"amor, Beatrice"}, data)
	}
}
func BenchmarkIsUpper(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.ToUpper(string(content))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsUpper(data)
	}
}

func BenchmarkIsLower(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := strings.ToLower(string(content))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsLower(data)
	}

}

func BenchmarkRemoveWhiteSpace(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveWhiteSpace(data)
	}
}
func BenchmarkIsASCII(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsASCII(data)
	}
}

func BenchmarkSplit(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Split(data)
	}
}
func BenchmarkSplitBuiltin(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.Split(data, "\n")
	}
}

func BenchmarkExtractString(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	initial := "LA DIVINA COMMEDIA"
	final := "altre stelle."
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractString(&data, initial, final)
	}
}

func BenchmarkRemoveNonASCII(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveNonASCII(data)
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
		arrayutils.JoinStrings(data, " ")
	}
}

func BenchmarkTrim(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Trim(data)
	}
}

func BenchmarkRemoveDoubleWhiteSpace(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveDoubleWhiteSpace(data)
	}
}
func BenchmarkCountLines(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountLines(data)
	}
}

func BenchmarkReverseString(b *testing.B) {
	content, err := ioutil.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReverseString(data)
	}
}
