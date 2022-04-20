package redshiftutils

import (
	"database/sql"
	secretutils "github.com/alessiosavi/GoGPUtils/aws/secrets"
	"log"
	"os"
	"testing"
)

var connection *sql.DB

func init() {
	var err error
	if connection, err = InitRedshiftConnection(); err != nil {
		panic(err)
	}

}
func InitRedshiftConnection() (*sql.DB, error) {
	var c Conf
	if err := secretutils.UnmarshalSecret(os.Getenv("secret_redshift"), &c); err != nil {
		return nil, err
	}

	c.Host = "localhost"
	c.Port = "5439"
	log.Println("Initializing connection for Redshift @" + c.Host)
	return MakeRedshfitConnection(c)
}

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
