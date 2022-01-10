package ssmutils

import (
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
