package httputils

import (
	"testing"
)

func TestCreateCookie(t *testing.T) {
	cookie, err := CreateCookie("alessio", "savi", "", "/", 60, false)
	if err != nil {
		t.Log(cookie)
		t.Log(err)
		t.Fail()
	}
}

func TestServeCookie(t *testing.T) {
	ServeCookie("localhost", "8080", "", "alessio", "savi", "localhost", "/", 60, false)
}
func TestDebugRequest(t *testing.T) {
	DebugRequest("localhost", "8080", "")
}
