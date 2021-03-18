package sftp

import (
	"bytes"
	"errors"
	"fmt"
	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	httputils "github.com/alessiosavi/GoGPUtils/http"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

var DEFAULT_KEY_EXCHANGE_ALGO = []string{"diffie-hellman-group-exchange-sha256"}

type SFTPConf struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"pass"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
}
type SFTPClient struct {
	Client *sftp.Client
}


func (c *SFTPConf) validate() error {
	if stringutils.IsBlank(c.Host) {
		return errors.New("SFTP host not provided")
	}
	if stringutils.IsBlank(c.User) {
		return errors.New("SFTP user not provided")
	}
	if stringutils.IsBlank(c.Password) {
		return errors.New("SFTP password not provided")
	}
	if !httputils.ValidatePort(c.Port) {
		return errors.New("SFTP port not provided")
	}
	return nil
}

// Create a new SFTP connection by given parameters
func (c *SFTPConf) NewConn(keyExchanges []string) (*SFTPClient, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	// Add default key exchange algorithm
	for _, algo := range DEFAULT_KEY_EXCHANGE_ALGO {
		if !arrayutils.InStrings(keyExchanges, algo) {
			keyExchanges = append(keyExchanges, algo)
		}
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(c.Timeout) * time.Second,
	}

	config.Config.KeyExchanges = append(config.Config.KeyExchanges, keyExchanges...)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	log.Println("Connecting to: " + addr)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(conn) // create sftp client
	if err != nil {
		return nil, err
	}
	return &SFTPClient{Client: client}, nil
}

func (c *SFTPClient) Get(remoteFile string) (*bytes.Buffer, error) {
	srcFile, err := c.Client.Open(remoteFile)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()
	var buf *bytes.Buffer
	_, err = io.Copy(buf, srcFile)
	return buf, err
}
func (c *SFTPClient) Put(data []byte, path string) error {
	f, err := c.Client.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (c *SFTPClient) CreateDirectory(path string) error {
	return c.Client.MkdirAll(path)

}

func (c *SFTPClient) DeleteDirectory(path string) error {
	return c.Client.RemoveDirectory(path)
}

func (c *SFTPClient) Exist(path string) (bool, error) {
	_, err := c.Client.Lstat(path)
	return err != nil, err
}

func (c *SFTPClient) IsDir(path string) (bool, error) {
	lstat, err := c.Client.Lstat(path)
	if err != nil {
		return false, err
	}
	return lstat.IsDir(), nil
}
func (c *SFTPClient) IsFile(path string) (bool, error) {
	lstat, err := c.Client.Lstat(path)
	if err != nil {
		return false, err
	}
	return !lstat.IsDir(), nil
}
