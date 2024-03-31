package cloudwatchutils

import (
	"fmt"
	"testing"
)

/*
 * This function controlls that there is a least 1 log group
 * and if there is return it
 */
func TestGetLogGroups(t *testing.T) {

	// Chiama la funzione GetLogGroups
	groups, err := GetLogGroups()
	if err != nil {
		t.Errorf("Errore durante la chiamata GetLogGroups: %v", err)
	}

	if len(groups) == 0 {
		t.Errorf("Nessun gruppo di log trovato")
	}
	for group := range groups {
		fmt.Println("log nome ", group)
	}

}
