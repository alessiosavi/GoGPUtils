package sftp

import (
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {
	sftpConf := SFTPConf{
		Host:     "localhost",
		User:     os.Getenv("ssh_user_test"),
		Password: os.Getenv("ssh_user_pass"),
		Port:     22,
		Timeout:  5,
	}
	conn, err := sftpConf.NewConn()
	if err != nil {
		panic(err)
	}
	defer conn.Client.Close()
	if err = conn.Put([]byte("this is a test!"), "/tmp/test/file.txt"); err != nil {
		panic(err)
	}
}
