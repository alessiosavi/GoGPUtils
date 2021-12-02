package cloudwatchutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
)

func TestGetLogGroups(t *testing.T) {
	groups, err := GetLogGroups()
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(groups))
}
