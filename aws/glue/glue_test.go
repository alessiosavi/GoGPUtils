package glueutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"strings"
	"testing"
)

func TestDeleteWorkflow(t *testing.T) {
	workflow, err := DeleteWorkflow("manella", true)
	if err != nil {
		t.Log(helper.MarshalIndent(err))
	}
	t.Log(helper.MarshalIndent(workflow))
}

func TestDeleteTable(t *testing.T) {

	databases, err := ListDatabases()
	if err != nil {
		panic(err)
	}
	for _, database := range databases {

		if _, err = DeleteDatabase(database); err != nil {
			panic(err)
		}
	}
}

func TestListJob(t *testing.T) {
	jobs, err := ListJobs()
	if err != nil {
		return
	}
	t.Log(helper.MarshalIndent(jobs))
}

func TestListJobs(t *testing.T) {
	jobs, err := ListJobs()
	if err != nil {
		t.Fatal(err)
	}
	for _, j := range jobs {
		if strings.HasPrefix(j, "qa-insert-sales-export") {
			log.Println(j)
			if err = PushRepo(j); err != nil {
				t.Fatal(err)
			}
		}
	}
}
