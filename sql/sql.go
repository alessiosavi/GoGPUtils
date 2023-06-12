package sqlutils

import (
	"database/sql"
	"fmt"
)

func ExecuteStatement(connection *sql.DB, queries ...string) error {
	for _, query := range queries {
		if _, err := connection.Exec(query); err != nil {
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
	replacer["time"] = "TIMESTAMPTZ"
	return replacer
}

// GetOracleTranslator is delegated to create a translator for csvutils.GetCSVDataType in order to create table
func GetOracleTranslator(charLength int) map[string]string {
	var replacer = make(map[string]string)
	replacer["string"] = fmt.Sprintf("VARCHAR2(%d)", charLength)
	replacer["int"] = "INTEGER"
	replacer["float"] = "FLOAT"
	replacer["bool"] = "VARCHAR(5)" // TRUE,FALSE -> 5 Character
	return replacer
}
