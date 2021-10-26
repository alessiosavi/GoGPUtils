package redshift

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/alessiosavi/GoGPUtils/helper"
	sqlutils "github.com/alessiosavi/GoGPUtils/sql"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"strings"
)

type Conf struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Host     string      `json:"host"`
	Port     json.Number `json:"port"`
	DBName   string      `json:"dbname"`
}

func (c *Conf) Validate() error {
	if stringutils.IsBlank(c.Username) {
		return fmt.Errorf("username is empty:[%+v]", helper.MarshalIndent(*c))
	}
	if stringutils.IsBlank(c.Password) {
		return fmt.Errorf("password is empty:[%+v]", helper.MarshalIndent(*c))
	}
	if stringutils.IsBlank(c.Host) {
		return fmt.Errorf("host is empty:[%+v]", helper.MarshalIndent(*c))
	}
	if stringutils.IsBlank(c.Port.String()) {
		return fmt.Errorf("password is empty:[%+v]", helper.MarshalIndent(*c))
	}
	if stringutils.IsBlank(c.DBName) {
		return fmt.Errorf("DBName is empty:[%+v]", helper.MarshalIndent(*c))
	}
	return nil
}

func (c *Conf) Load(confFile string) error {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &c); err != nil {
		return err
	}
	if err = c.Validate(); err != nil {
		log.Println(helper.MarshalIndent(c))
		return err
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
	return db, db.Ping()
}

// CreateTableByType is delegated to create the `CREATE TABLE` query for the given table
// tableName: Name of the table
// headers: List of headers necessarya to preserve orders
// tableType: Map of headers:type for the given table
func CreateTableByType(tableName string, headers []string, tableType map[string]string) string {
	var sb strings.Builder
	translator := sqlutils.GetRedshiftTranslator()
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")
	for _, header := range headers {
		//for k, v := range tableType {
		sb.WriteString("\t" + header + " " + translator[tableType[header]] + ",\n")
	}
	data := sb.String()
	data = strings.TrimSuffix(data, ",\n")
	data = data + ");"
	return data
}
