package sqlutils

import (
	"database/sql"
	"fmt"
	redshiftutils "github.com/alessiosavi/GoGPUtils/aws/redshift"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"strings"
)

func ExecuteStatement(connection *sql.DB, queries ...string) error {
	for _, query := range queries {
		if _, err := connection.Exec(query); err != nil {
			if strings.Contains(err.Error(), "stl_load_errors") {
				log.Println(helper.MarshalIndent(redshiftutils.GetCOPYErrors(connection)[0]))
			}
			return fmt.Errorf("ERROR Executing the following query:\n%s\n | Error: %s", query, err.Error())
		}
	}
	return nil
}

// GetRedshiftTranslator is delegated to create a translator for csvutils.GetCSVDataType in order to create table
func GetRedshiftTranslator() map[string]string {
	var replacer = make(map[string]string)
	replacer["string"] = "TEXT"
	replacer["int"] = "INTEGER"
	replacer["float"] = "FLOAT"
	replacer["bool"] = "BOOLEAN"
	return replacer
}

// GetOracleTranslator is delegated to create a translator for csvutils.GetCSVDataType in order to create table
func GetOracleTranslator(charLenght int) map[string]string {
	var replacer = make(map[string]string)
	replacer["string"] = fmt.Sprintf("VARCHAR2(%d)", charLenght)
	replacer["int"] = "INTEGER"
	replacer["float"] = "FLOAT"
	replacer["bool"] = "VARCHAR(5)" // TRUE,FALSE -> 5 Character
	return replacer
}
