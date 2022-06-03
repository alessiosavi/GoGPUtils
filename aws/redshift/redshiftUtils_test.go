package redshiftutils

import (
	"database/sql"
	secretsutils "github.com/alessiosavi/GoGPUtils/aws/secrets"
	"log"
	"os"
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

func InitRedshiftConnection() (*sql.DB, error) {
	log.Println("Initializing connection for Redshift")
	var c Conf
	if err := secretsutils.UnmarshalSecret(os.Getenv("secret_redshift"), &c); err != nil {
		return nil, err
	}
	c.Host = "localhost"
	c.Port = "5440"
	connection, err := MakeRedshfitConnection(c)
	if err != nil {
		return nil, err
	}
	return connection, err
}
