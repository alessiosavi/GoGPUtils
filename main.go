package main

import (
	"encoding/json"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
)

func main() {

	var a []T
	data := `[
  {
    "channel": "CHINA",
    "sftp_folder": "upload/",
    "sftp_conf": {
      "host": "sftp.zegna.com",
      "user": "fabricalab",
      "pass": "jJsF)m.7J<af$nmK}F:x?.(t\\~jqy_+",
      "port": 22,
      "timeout": 60
    }
  },
  {
    "channel": "ALL_EXCEPT_CHINA",
    "sftp_folder": "put/",
    "sftp_conf": {
      "host": "sfile.cegid.com",
      "user": "thod3pcl",
      "pass": "33bNIQUG",
      "port": 22,
      "timeout": 60
    }
  }
]
`
	err := json.Unmarshal([]byte(data), &a)
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(a))
}

type T struct {
	Channel    string `json:"channel"`
	SftpFolder string `json:"sftp_folder"`
	SftpConf   struct {
		Host    string `json:"host"`
		User    string `json:"user"`
		Pass    string `json:"pass"`
		Port    int    `json:"port"`
		Timeout int    `json:"timeout"`
	} `json:"sftp_conf"`
}
