package glueutils

import (
	"github.com/alessiosavi/GoGPUtils/helper"
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
