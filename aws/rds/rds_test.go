package rds

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"testing"
)

func TestListRDS(t *testing.T) {
	ListRDS()
}
func TestDescribeInstance(t *testing.T) {
	clusters, err := ListRDS()
	if err != nil {
		t.Error(err)
	}
	for _, cluster := range clusters {
		DescribeInstanceByID(cluster)
	}

}

func TestDescribeInstanceByID(t *testing.T) {
	id := DescribeInstanceByID("")
	t.Log(helper.MarshalIndent(*id))
}
