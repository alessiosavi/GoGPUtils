package dbutils

import (
	"database/sql"
	"log"
)

func ExecuteStatement(connection *sql.DB, queries ...string) error {
	for _, query := range queries {
		if _, err := connection.Exec(query); err != nil {
			log.Println("ERROR Executing the following query:\n" + query)
			return err
		}
	}
	return nil
}
