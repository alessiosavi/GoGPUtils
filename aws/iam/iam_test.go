package iamutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	list, err := ListRoles(nil)
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(list))

	var notEU []string
	var notThisYear []string
	var neverUsed []string

	for _, r := range list {
		if r.RoleLastUsed == nil {
			neverUsed = append(neverUsed, *r.RoleName)
		} else {
			if *r.RoleLastUsed.Region != "eu-west-1" {
				notEU = append(notEU, *r.RoleName)
			}
			if !r.RoleLastUsed.LastUsedDate.After(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) {
				notThisYear = append(notThisYear, *r.RoleName)
			}
		}

	}
	t.Log("NOT EU:")
	t.Log(helper.MarshalIndent(notEU))

	t.Log("NOT THIS YEAR:")
	t.Log(helper.MarshalIndent(notThisYear))

	t.Log("Never Used:")
	t.Log(helper.MarshalIndent(neverUsed))

}

func TestInfo(t *testing.T) {
	role, err := GetRole("thom-browne-remove-tables-role-eumrh440")
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(*role))
}
