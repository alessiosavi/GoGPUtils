package stringutils

import (
	"testing"
)

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

func BenchmarkTestIsUpperKO(b *testing.B) {
	data := `AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa`
	for n := 0; n < b.N; n++ {
		IsUpper(data)
	}
}

func BenchmarkTestIsUpperOK(b *testing.B) {
	data := `AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`
	for n := 0; n < b.N; n++ {
		IsUpper(data)
	}
}

func BenchmarkTestIsLowerOK(b *testing.B) {
	data := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`
	for n := 0; n < b.N; n++ {
		IsLower(data)
	}
}
func BenchmarkTestIsLowerKO(b *testing.B) {
	data := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaA`
	for n := 0; n < b.N; n++ {
		IsLower(data)
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
	dataKO := []string{`....`, `,,,,,`, `,,, ,,,,`, `<<<`, `!!!`}
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
	dataKO := []string{`....`, `,,,,,`, `,,, ,,,,`, `<<<`, `!!!`, `2`}
	for i := range dataOK {
		if !ContainsOnlyLetter(dataOK[i]) {
			t.Fail()
		}
	}
	for i := range dataKO {
		if ContainsOnlyLetter(dataKO[i]) {
			t.Fail()
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
