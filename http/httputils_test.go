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
		_port := helper.InitRandomizer().RandomInt(81, 65535)
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
