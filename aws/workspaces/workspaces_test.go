package workspaces

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
)

func TestGetWorkspaces(t *testing.T) {
	workspaces, err := GetWorkspaces("")
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(workspaces))
}
