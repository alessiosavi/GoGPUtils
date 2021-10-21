package httputils

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/helper"
)

var port string

func init() {
	for {

		_port := 		helper.RandomInt(1024, 65535)
		timeout := time.Second
		host := fmt.Sprintf("localhost:%d", _port)
		conn, err := net.DialTimeout("tcp", host, timeout)
		if conn != nil {
			fmt.Printf("Port %d is already used!", _port)
			conn.Close()
			continue
		}
		if err != nil {
			port = fmt.Sprintf("%d", _port)
			break
		}
	}
}

func TestCreateCookie(t *testing.T) {
	cookie, err := CreateCookie("alessio", "savi", "", "/", 60, false)
	if err != nil {
		t.Log(cookie)
		t.Log(err)
		t.Fail()
	}
}

func TestServeCookie(t *testing.T) {
	go ServeCookie("localhost", port, "", "alessio", "savi", "localhost", "/", 60, false)
	time.Sleep(time.Millisecond * 200)
	http.Get(`http://localhost:` + port + "/")
	time.Sleep(time.Millisecond * 200)
}
func TestDebugRequest(t *testing.T) {
	go DebugRequest("localhost", port, "")
	time.Sleep(time.Millisecond * 200)
	_, err := http.Get(`http://localhost:` + port + "/")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	time.Sleep(time.Millisecond * 200)
}

func TestServeHeaders(t *testing.T) {
	time.Sleep(time.Millisecond * 200)
	go ServeHeaders(nil, "localhost", port, "")
	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(`http://localhost:` + port + "/")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(resp)
	time.Sleep(time.Millisecond * 200)
}

func TestValidatePort(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test_ok1", args: args{port: 22}, want: true},
		{name: "test_ko1", args: args{port: -1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePort(tt.args.port); got != tt.want {
				t.Errorf("ValidatePort() = %v, want %v", got, tt.want)
			}
		})
	}
}
