package cloudwatchutils

import (
	"fmt"
	"testing"
	"time"
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

	//Controlla alcuni gruppi per verificarne la correttezza
	for _, group := range groups {
		fmt.Println(group.Arn)
	}

}

/*
 * This function return the export of a given log group
 *
 */
func TestExportLog(t *testing.T) {
	// Chiamata alla funzione ExportLog per esportare i log
	output, err := ExportLog("friends-s3", "/aws/lambda/VisitCountFunction", "export_test", time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		fmt.Println("Errore durante l'esportazione dei log:", err)
		return
	}

	fmt.Println("Attività di esportazione dei log creata con successo:", output)

}
