package redshiftutils

import (
	"database/sql"
	secretsutils "github.com/alessiosavi/GoGPUtils/aws/secrets"
	"log"
	"testing"
)

func Test_PhysicalDelete(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	err = PhysicalDelete(connection)
	if err != nil {
		panic(err)
	}
}

func Test_VACUUM(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	err = VACUUM(connection)
	if err != nil {
		panic(err)
	}
}

func InitRedshiftConnection() (*sql.DB, error) {
	log.Println("Initializing connection for Redshift")
	var c Conf
	if err := secretsutils.UnmarshalSecret("prod/redshift", &c); err != nil {
		return nil, err
	}
	c.Host = "localhost"
	c.Port = "5439"
	connection, err := MakeRedshfitConnection(c)
	if err != nil {
		return nil, err
	}
	return connection, err
}

func TestManualSnapshot(t *testing.T) {
	ManualSnapshot()
}

func TestSetAutoOptimization(t *testing.T) {
	connection, err := InitRedshiftConnection()
	if err != nil {
		panic(err)
	}
	if err = SetAutoOptimization(connection); err != nil {
		panic(err)
	}
}
