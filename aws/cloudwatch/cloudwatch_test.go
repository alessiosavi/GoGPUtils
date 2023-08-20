package cloudwatchutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
	"time"
)

func TestGetLogGroups(t *testing.T) {
	groups, err := GetLogGroups()
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(groups))
}

func TestExportLog(t *testing.T) {
	exportLog, err := ExportLog("thb-batch-log", "/aws/lambda/prod-thb-sf-customer-merging-tool", "prod-thb-sf-customer-merging-tool",
		time.Date(2021, 12, 02, 0, 0, 0, 0, time.UTC), time.Now())
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(exportLog))
	for {
		task, err := DescribeExportTask(*exportLog.TaskId)
		if err != nil {
			panic(err)
		}
		t.Log(helper.MarshalIndent(task))
		time.Sleep(10 * time.Second)
	}
}

func TestDescribeExportTask(t *testing.T) {
	task, err := DescribeExportTask("4ff4e8a2-9995-49ba-921f-a0b4264d5a6a")
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(task))
}
