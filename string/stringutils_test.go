package stringutils

import (
	"golang.org/x/exp/slices"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const distance1 string = `../testdata/files/testDistance.txt`
const distance2 string = `../testdata/files/testDistance1.txt`
const danteDataset string = `../testdata/files/dante.txt`

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
		ExtractString(data, initial, final)
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
			name: "2",
			args: args{
				w: "012500",
			},
		},
		{
			name: "3",
			args: args{
				w: " 0012500",
			},
		},
		{
			name: "4",
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

func TestExtractString(t *testing.T) {
	type args struct {
		data  string
		first string
		last  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok_1",
			args: args{
				data:  "This is a test for extract string",
				first: "This is a",
				last:  "for extract string",
			},
			want: "test",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractString(tt.args.data, tt.args.first, tt.args.last); got != tt.want {
				t.Errorf("ExtractString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUpper(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok_1",
			args: args{
				str: "AAA",
			},
			want: true,
		},
		{
			name: "ko_1",
			args: args{
				str: "AaA",
			},
			want: false,
		},
		{
			name: "ko_2",
			args: args{
				str: "A!",
			},
			want: true,
		},
		{
			name: "ko_3",
			args: args{
				str: "a!",
			},
			want: false,
		},
		{
			name: "ko_4",
			args: args{
				str: "!",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUpper(tt.args.str); got != tt.want {
				t.Errorf("IsUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLower(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok_1",
			args: args{
				str: "aaa",
			},
			want: true,
		},
		{
			name: "ko_1",
			args: args{
				str: "AaA",
			},
			want: false,
		},
		{
			name: "ko_2",
			args: args{
				str: "a!",
			},
			want: true,
		},
		{
			name: "ko_3",
			args: args{
				str: "A!",
			},
			want: false,
		},
		{
			name: "ko_4",
			args: args{
				str: "!",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLower(tt.args.str); got != tt.want {
				t.Errorf("IsUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveNonASCII(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				str: "AÀ È Ì Ò Ù Ỳ Ǹ ẀA",
			},
			want: "A A",
		},
		{
			name: "2",
			args: args{
				str: `AȨ Ç Ḑ Ģ Ḩ Ķ Ļ Ņ Ŗ ŞA`,
			},
			want: "A A",
		},
		{
			name: "3",
			args: args{
				str: `AA♩ ♪ ♫ ♬ ♭ ♮ ♯AA`,
			},
			want: `AA AA`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveNonASCII(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNonASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				str: "a",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				str: "a a",
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				str: "A",
			},
			want: true,
		},
		{
			name: "4",
			args: args{
				str: "A A",
			},
			want: true,
		},
		{
			name: "5",
			args: args{
				str: "a!",
			},
			want: false,
		},
		{
			name: "5",
			args: args{
				str: "!!",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.args.str); got != tt.want {
				t.Errorf("IsAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}
