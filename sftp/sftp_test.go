package sftputils

import (
	"io"
	"testing"
)

func TestCopyFile(t *testing.T) {
	sftpConf := SFTPConf{
		Host:     "test.rebex.net",
		User:     "demo",
		Password: "password",
		Port:     22,
		Timeout:  5,
	}
	conn, err := sftpConf.NewConn()
	if err != nil {
		panic(err)
	}
	defer conn.Client.Close()
	list, err := conn.List(".")
	if err != nil {
		panic(err)
	}
	if len(list) == 0 {
		t.Fail()
	}

	get, err := conn.Get("readme.txt")
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(get)
	if err != nil {
		panic(err)
	}
	if len(all) == 0 {
		t.Fail()
	}

}
