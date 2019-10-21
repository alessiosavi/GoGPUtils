package stringutils

import "testing"

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

func TestUpperizeString(t *testing.T) {
	dataOK := []string{`AAA`, `AbbbA`, `aBBBBa`}
	for _, item := range dataOK {
		t.Log(UpperizeString(&item))
	}
}

func TestLowerizeString(t *testing.T) {
	dataOK := []string{`AAA`, `AbbbA`, `aBBBBa`}
	for _, item := range dataOK {
		t.Log(LowerizeString(&item))
	}
}

func TestRemoveNonAscii(t *testing.T) {
	dataOK := []string{`AAÀ È Ì Ò Ù Ỳ Ǹ ẀBB`, `CCȨ Ç Ḑ Ģ Ḩ Ķ Ļ Ņ Ŗ Ş DD`, `$$♩ ♪ ♫ ♬ ♭ ♮ ♯##`}
	for _, item := range dataOK {
		t.Log(RemoveNonASCII(item))
	}
}
