package stringutils

import (
	"golang.org/x/exp/slices"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	arrayutils "github.com/alessiosavi/GoGPUtils/array"
)

const distance1 string = `../testdata/files/testDistance.txt`
const distance2 string = `../testdata/files/testDistance1.txt`
const danteDataset string = `../testdata/files/dante.txt`

func TestExtractUpperBlock(t *testing.T) {
	for _, test := range []string{"test_ID", "test_ID_test", "Test_ID_TEST", "test_ID_test_ID", "test ID", "test_ID test", "Test ID TEST", "test ID_test ID", "TestID", "TestId"} {

		//for _, test := range []string{"CompanyName"} {
		r := strings.NewReplacer(" ", "_")

		expected := r.Replace(strings.ToLower(test))
		if res := ExtractUpperBlock(test, r); res != expected {
			t.Errorf("Expected: %s | Found: %s", expected, res)
		}
	}
}
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

func TestContainsMultiple(t *testing.T) {

}

func TestRemoveNonAscii(t *testing.T) {
	testData := []string{`AÀ È Ì Ò Ù Ỳ Ǹ ẀA`, `AȨ Ç Ḑ Ģ Ḩ Ķ Ļ Ņ Ŗ ŞA`, `AA♩ ♪ ♫ ♬ ♭ ♮ ♯AA`}
	dataOK := []string{`A A`, `A A`, `AA AA`}
	var data []string
	for _, item := range testData {
		data = append(data, RemoveNonASCII(item))
	}
	if !reflect.DeepEqual(dataOK, data) {
		t.Errorf("Expected: %v | Found: %v", dataOK, data)
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
		if len(str) != 4 {
			t.Fail()
		}
	}
}

func TestRemoveDoubleWhiteSpace(t *testing.T) {
	data := []string{`  test`, `test  `, `te  st`}
	for _, item := range data {
		str := RemoveDoubleWhiteSpace(item)
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
	content, err := os.ReadFile(danteDataset)
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
		if len(str) != 4 {
			t.Fail()
		}
	}
}
func TestTrimNewLine(t *testing.T) {
	data := []string{"test\n", "\ntest", "\ntest\n", "\ntest \n"}
	for _, item := range data {
		str := Trim(item)
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

func TestJaroDistance(t *testing.T) {
	var str1, str2 = "MARTHA", "MARHTA"
	if JaroDistance(str1, str2) != 0.9444444444444445 {
		t.Fail()
	}
}

func TestDiceCoefficient(t *testing.T) {
	str1, str2 := "prova1", "prova2"
	diceCoeffiecnt := DiceCoefficient(str1, str2)
	if diceCoeffiecnt != 0.8 {
		t.Error(diceCoeffiecnt)
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

func BenchmarkLevenshteinDistanceLegacy(t *testing.B) {
	content, err := os.ReadFile(distance1)
	if err != nil {
		return
	}
	data1 := strings.ToUpper(string(content))
	content, err = os.ReadFile(distance2)
	if err != nil {
		return
	}
	data2 := strings.ToUpper(string(content))

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LevenshteinDistanceLegacy(data1, data2)
	}
}
func BenchmarkLevenshteinDistance(t *testing.B) {
	content, err := os.ReadFile(distance1)
	if err != nil {
		return
	}
	data1 := strings.ToUpper(string(content))
	content, err = os.ReadFile(distance2)
	if err != nil {
		return
	}
	data2 := strings.ToUpper(string(content))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		LevenshteinDistance(data1, data2)
	}
}

func BenchmarkDiceCoefficient(t *testing.B) {
	content, err := os.ReadFile(distance1)
	if err != nil {
		return
	}
	data1 := strings.ToUpper(string(content))
	content, err = os.ReadFile(distance2)
	if err != nil {
		return
	}
	data2 := strings.ToUpper(string(content))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		DiceCoefficient(data1, data2)
	}
}
func BenchmarkJaroDistance(t *testing.B) {
	content, err := os.ReadFile(distance1)
	if err != nil {
		return
	}
	data1 := strings.ToUpper(string(content))
	content, err = os.ReadFile(distance2)
	if err != nil {
		return
	}
	data2 := strings.ToUpper(string(content))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		JaroDistance(data1, data2)
	}
}

func BenchmarkContainsOnlyLetter(t *testing.B) {
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		RemoveFromString(data, len(data)-1)
	}
}

func BenchmarkExtractTextFromQuery(b *testing.B) {
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Split(strings.NewReader(data))
	}
}
func BenchmarkSplitBuiltin(b *testing.B) {
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
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
	content, err := os.ReadFile(danteDataset)
	if err != nil {
		return
	}
	data := string(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReverseString(data)
	}
}

func TestContainsMultiple1(t *testing.T) {
	type args struct {
		lower     bool
		s         string
		substring []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "testOK1", args: struct {
			lower     bool
			s         string
			substring []string
		}{lower: true, s: "this is a test", substring: []string{"this", "is"}}, want: true},
		{name: "testKO1", args: struct {
			lower     bool
			s         string
			substring []string
		}{lower: true, s: "this is a test", substring: []string{"not", "found"}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsMultiple(tt.args.lower, tt.args.s, tt.args.substring...); got != tt.want {
				t.Errorf("ContainsMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique2(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Ok",
			args: args{
				data: []string{"a", "b", "c", "c"},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "Ok",
			args: args{
				data: []string{"a", "b", "c"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Compact(tt.args.data)
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] < tt.want[j]
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexes(t *testing.T) {
	type args struct {
		s   string
		chs string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "OK",
			args: args{
				s:   "ABAABA",
				chs: "A",
			},
			want: []int{0, 2, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Indexes(tt.args.s, tt.args.chs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Indexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pad(t *testing.T) {

	expected := "000000012500"
	n := len(expected)
	type args struct {
		w string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				w: "12500",
			},
		},
		{
			name: "1",
			args: args{
				w: "012500",
			},
		},
		{
			name: "1",
			args: args{
				w: " 0012500",
			},
		},
		{
			name: "1",
			args: args{
				w: "000000012500",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pad(tt.args.w, "0", n); got != expected {
				t.Errorf("pad() = %v, want %v", got, expected)
			}
		})
	}
}

func TestTrimStrings(t *testing.T) {
	type args struct {
		vs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				vs: []string{"2023-01-01 "},
			},
			want: []string{"2023-01-01"},
		},
		{
			name: "2",
			args: args{
				vs: []string{"2023-01-01   "},
			},
			want: []string{"2023-01-01"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimStrings(tt.args.vs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
