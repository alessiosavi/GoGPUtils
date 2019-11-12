package httputils

import (
	"net/http"
	"testing"
	"time"
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
	go ServeCookie("localhost", "9999", "", "alessio", "savi", "localhost", "/", 60, false)
	time.Sleep(time.Millisecond * 200)
	http.Get(`http://localhost:9999/`)
	time.Sleep(time.Millisecond * 200)
}
func TestDebugRequest(t *testing.T) {
	go DebugRequest("localhost", "9999", "")
	time.Sleep(time.Millisecond * 200)
	_, err := http.Get(`http://localhost:9999/`)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	time.Sleep(time.Millisecond * 200)
}

func TestServeHeaders(t *testing.T) {
	go ServeHeaders(nil, "localhost", "9999", "")
	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(`http://localhost:9999/`)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(resp)
	time.Sleep(time.Millisecond * 200)
}

func BenchmarkCreateCookie(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
func TestSetHeaders(t *testing.T) {}
func BenchmarkSetHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkServeHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkServeCookie(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkDebugRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
