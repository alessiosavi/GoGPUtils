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

	exportLog, err := ExportLog("thb-batch-log", "/aws/lambda/prod-go-centric-parser", "prod-go-centric-parser",
		time.Date(2021, 12, 02, 0, 0, 0, 0, time.UTC), time.Date(2021, 12, 8, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	}
	t.Log(helper.MarshalIndent(exportLog))
}
