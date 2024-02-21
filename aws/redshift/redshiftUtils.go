package redshiftutils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	sqlutils "github.com/alessiosavi/GoGPUtils/sql"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	_ "github.com/lib/pq"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

var redshiftClient *redshift.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		redshiftClient = redshift.New(redshift.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}
func ManualSnapshot() {

	//clusters, err := redshiftClient.DescribeClusters(context.Background(), &redshift.DescribeClustersInput{
	//	ClusterIdentifier: nil,
	//	Marker:            nil,
	//	MaxRecords:        nil,
	//	TagKeys:           nil,
	//	TagValues:         nil,
	//})
	//if err != nil {
	//	return
	//}
	//log.Println(helper.MarshalIndent(clusters))

	t := time.Now().Format(time.RFC3339)
	snapshot, err := redshiftClient.CreateClusterSnapshot(context.Background(), &redshift.CreateClusterSnapshotInput{
		ClusterIdentifier:  aws.String("qa-data-warehouse"),
		SnapshotIdentifier: aws.String(fmt.Sprintf("%s-%s", "qa-data-warehouse", t)),
	})
	if err != nil {
		return
	}
	log.Println(helper.MarshalIndent(snapshot))
}

func (c *Conf) Load(confFile string) error {
	data, err := os.ReadFile(confFile)
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
// headers: List of headers necessary to preserve orders
// tableType: Map of headers:type for the given table

func CreateTableByType(tableName string, headers []string, tableType map[string]string) string {
	var sb strings.Builder
	translator := sqlutils.GetRedshiftTranslator()
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")
	replacer := strings.NewReplacer(".", "_", ",", "_", " ", "_", "(", "_", ")", "_", "/", "_")
	for _, header := range headers {
		fixHeader := replacer.Replace(header)
		//for k, v := range tableType {
		sb.WriteString("\t" + fixHeader + " " + translator[tableType[header]] + ",\n")
	}
	data := strings.TrimSuffix(sb.String(), ",\n") + ");"
	return data
}

type Result struct {
	Userid          int       `json:"userid,omitempty"`
	Slice           int       `json:"slice,omitempty"`
	Tbl             int       `json:"tbl,omitempty"`
	Starttime       time.Time `json:"starttime,omitempty"`
	Session         int       `json:"session,omitempty"`
	Query           int       `json:"query,omitempty"`
	Filename        string    `json:"filename,omitempty"`
	Line_number     int       `json:"line_number,omitempty"`
	Colname         string    `json:"colname,omitempty"`
	Type            string    `json:"type,omitempty"`
	Col_length      string    `json:"col_length,omitempty"`
	Position        int       `json:"position,omitempty"`
	Raw_line        string    `json:"raw_line,omitempty"`
	Raw_field_value string    `json:"raw_field_value,omitempty"`
	Err_code        int       `json:"err_code,omitempty"`
	Err_reason      string    `json:"err_reason,omitempty"`
	Is_partial      string    `json:"is_partial,omitempty"`
	Start_offset    string    `json:"start_offset,omitempty"`
}

func (r *Result) Trim() {
	r.Filename = stringutils.Trim(r.Filename)
	r.Colname = stringutils.Trim(r.Colname)
	r.Type = stringutils.Trim(r.Type)
	r.Col_length = stringutils.Trim(r.Col_length)
	r.Raw_line = stringutils.Trim(r.Raw_line)
	r.Raw_field_value = stringutils.Trim(r.Raw_field_value)
	r.Err_reason = stringutils.Trim(r.Err_reason)
	r.Is_partial = stringutils.Trim(r.Is_partial)
	r.Start_offset = stringutils.Trim(r.Start_offset)
}

// GetCOPYErrors is delegated to retrieve the loading error related to the COPY commands, sorted by time
func GetCOPYErrors(connection *sql.DB) []Result {
	rows, err := connection.Query("select * from stl_load_errors order by starttime desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var errorsResult []Result
	for rows.Next() {
		var res Result
		if err = rows.Scan(&res.Userid,
			&res.Slice,
			&res.Tbl,
			&res.Starttime,
			&res.Session,
			&res.Query,
			&res.Filename,
			&res.Line_number,
			&res.Colname,
			&res.Type,
			&res.Col_length,
			&res.Position,
			&res.Raw_line,
			&res.Raw_field_value,
			&res.Err_code,
			&res.Err_reason,
			&res.Is_partial,
			&res.Start_offset); err != nil {
			panic(err)
		} else {
			res.Trim()
			errorsResult = append(errorsResult, res)
		}
	}
	return errorsResult
}

// SetAutoOptimization is delegated to scan all the tables in a redshift cluster and set the automatic diststyle and sortkey
func SetAutoOptimization(connection *sql.DB) error {
	query := `
select distinct t.table_name
from information_schema.tables t
where t.table_schema = 'public' and table_type = 'BASE TABLE'
order by t.table_name;`

	rows, err := connection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			return err
		}
		result = append(result, s)
	}

	// Convert a table to a diststyle AUTO table
	var dist = `ALTER TABLE %s ALTER DISTSTYLE AUTO;`
	// Convert a table to a sort key AUTO table
	var sort = `ALTER TABLE %s ALTER SORTKEY AUTO;`
	var encode = `ALTER TABLE %s ALTER ENCODE AUTO;`
	bar := progressbar.Default(int64(len(result)))
	for _, table := range result {
		bar.Add(1)
		bar.Describe(fmt.Sprintf(dist, table))
		sqlutils.ExecuteStatement(connection, fmt.Sprintf(dist, table))
		bar.Describe(fmt.Sprintf(sort, table))
		sqlutils.ExecuteStatement(connection, fmt.Sprintf(sort, table))
		bar.Describe(fmt.Sprintf(encode, table))
		sqlutils.ExecuteStatement(connection, fmt.Sprintf(encode, table))
	}
	return nil
}

// PhysicalDelete is delegated to perform the physical delete of the items marked as deleted
func PhysicalDelete(connection *sql.DB) error {
	query := `
	SELECT distinct tablename
	FROM pg_table_def
	WHERE schemaname = 'public'
	and "column" = 'flag_delete';`

	rows, err := connection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			return err
		}
		result = append(result, s)
	}
	var remove = `delete from %s where flag_delete=true;`
	bar := progressbar.Default(int64(len(result)))
	for _, table := range result {
		bar.Describe(table)
		bar.Add(1)
		if err = sqlutils.ExecuteStatement(connection, fmt.Sprintf(remove, table)); err != nil {
			log.Println(err)
			log.Println()
		}
	}
	return nil
}

func VACUUM(connection *sql.DB) error {
	query := `
	SELECT distinct tablename
	FROM pg_table_def
	WHERE schemaname = 'public';`
	rows, err := connection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			return err
		}
		result = append(result, s)
	}
	bar := progressbar.Default(int64(len(result)))
	for _, table := range result {
		bar.Describe(table)
		bar.Add(1)
		q := fmt.Sprintf("VACUUM FULL %s TO 100 PERCENT BOOST;", table)
		if _, err := connection.Exec(q); err != nil {
			log.Println(err)
		}
		q = fmt.Sprintf("VACUUM REINDEX %s;", table)
		if _, err := connection.Exec(q); err != nil {
			log.Println(err)
		}
	}
	return nil
}
