package redshift

import (
	"database/sql"
	"fmt"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

type Conf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

func (c *Conf) Validate() error {

	if stringutils.IsBlank(c.Username) {
		return fmt.Errorf("username is empty:[%+v]", *c)
	}
	if stringutils.IsBlank(c.Password) {
		return fmt.Errorf("password is empty:[%+v]", *c)
	}
	if stringutils.IsBlank(c.Host) {
		return fmt.Errorf("host is empty:[%+v]", *c)
	}
	if stringutils.IsBlank(c.Port) {
		return fmt.Errorf("port is empty:[%+v]", *c)
	}
	if stringutils.IsBlank(c.DBName) {
		return fmt.Errorf("DBName is empty:[%+v]", *c)
	}
	return nil
}
func MakeRedshfitConnection(conf Conf) (*sql.DB, error) {
	var err error
	var db *sql.DB = nil
	if err := conf.Validate(); err != nil {
		return db, err
	}
	url := fmt.Sprintf("sslmode=require user=%v password=%v host=%v port=%v dbname=%v",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName)

	if db, err = sql.Open("postgres", url); err != nil {
		return nil, fmt.Errorf("redshift connect error : (%s)", err.Error())
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("redshift ping error : (%s)", err.Error())
	}
	return db, nil
}
