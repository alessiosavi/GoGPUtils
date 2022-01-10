package ssmutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
)

func TestList(t *testing.T) {
	res, err := List()
	if err != nil {
		panic(err)
	}
	log.Println(res)
}

func TestDescribe(t *testing.T) {
	describe, err := Describe("/qa/autoexecution-secret")
	if err != nil {
		return
	}
	log.Println(helper.MarshalIndent(describe))
}
