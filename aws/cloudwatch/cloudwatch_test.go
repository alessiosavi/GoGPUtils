package cloudwatchutils

import (
	"fmt"
	"testing"
)

func TestGetLogGroups(t *testing.T) {

	// Chiama la funzione GetLogGroups
	groups, err := GetLogGroups()
	if err != nil {
		t.Errorf("Errore durante la chiamata GetLogGroups: %v", err)
	}

	if len(groups) == 0 {
		t.Errorf("Nessun gruppo di log trovato")
	}

	//Controlla alcuni gruppi per verificarne la correttezza
	for _, group := range groups {
		fmt.Println(group.LogGroupName)

	}

}
