package redshiftutils

import (
	"database/sql"
	"os"
	"testing"
)

var connection *sql.DB

func TestUnloadDB(t *testing.T) {
	UnloadDB(connection, "public", "prod-demand-planning-forecast-temp", "dump/REDSHIFT_PROD", os.Getenv("test_role_redshift"))

}

func TestLoadDB(t *testing.T) {
	LoadDB(connection, "public", "qa-demand-planning-forecast-temp", "dump/REDSHIFT_PARSED", os.Getenv("test_role_redshift"))
}

func TestPhysicalDelete(t *testing.T) {
	if err := PhysicalDelete(connection); err != nil {
		panic(err)
	}
}

func TestSetAutoOptimization(t *testing.T) {
	if err := SetAutoOptimization(connection); err != nil {
		panic(err)
	}
}
