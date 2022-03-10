package redshiftutils

import (
	"database/sql"
	secretutils "github.com/alessiosavi/GoGPUtils/aws/secrets"
	"log"
	"os"
	"testing"
)

func InitRedshiftConnection() (*sql.DB, error) {
	var c Conf
	if err := secretutils.UnmarshalSecret(os.Getenv("secret_redshift"), &c); err != nil {
		return nil, err
	}

	c.Host = "localhost"
	c.Port = "5439"
	log.Println("Initializing connection for Redshift @" + c.Host)
	connection, err := MakeRedshfitConnection(c)
	if err != nil {
		return nil, err
	}
	return connection, err
}

func TestUnloadDB(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	UnloadDB(connection, "public", "prod-demand-planning-forecast-temp", "dump/REDSHIFT_PROD", os.Getenv("test_role_redshift"))

}

func TestLoadDB(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	LoadDB(connection, "public", "qa-demand-planning-forecast-temp", "dump/REDSHIFT_PARSED", os.Getenv("test_role_redshift"))
}

func TestPhysicalDelete(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	if err = PhysicalDelete(connection); err != nil {
		panic(err)
	}
}
