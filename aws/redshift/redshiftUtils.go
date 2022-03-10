package redshiftutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	S3utils "github.com/alessiosavi/GoGPUtils/aws/S3"
	"github.com/alessiosavi/GoGPUtils/helper"
	sqlutils "github.com/alessiosavi/GoGPUtils/sql"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	_ "github.com/lib/pq"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"path"
	"strings"
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
// headers: List of headers necessary to preserve orders
// tableType: Map of headers:type for the given table
func CreateTableByType(tableName string, headers []string, tableType map[string]string) string {
	var sb strings.Builder

	//for i := range headers {
	//	headers[i] = stringutils.ExtractUpperBlock(headers[i])
	//}

	translator := sqlutils.GetRedshiftTranslator()
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")
	replacer := strings.NewReplacer(".", "", ",", "", " ", "", "(", "", ")", "")
	for _, header := range headers {
		fixHeader := replacer.Replace(header)
		//for k, v := range tableType {
		sb.WriteString("\t" + stringutils.ExtractUpperBlock(fixHeader, nil) + " " + translator[tableType[header]] + ",\n")
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
select t.table_name
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
	bar := progressbar.Default(int64(len(result)))
	for _, table := range result {
		bar.Describe(dist)
		sqlutils.ExecuteStatement(connection, fmt.Sprintf(dist, table))
		bar.Describe(sort)
		sqlutils.ExecuteStatement(connection, fmt.Sprintf(sort, table))
		bar.Add(1)
	}
	return nil
}

// PhysicalDelete is delegated to perform the physical delete of the items marked as deleted
func PhysicalDelete(connection *sql.DB) error {
	query := `
	SELECT tablename
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
		if err = sqlutils.ExecuteStatement(connection, fmt.Sprintf(remove, table)); err != nil {
			log.Println(fmt.Sprintf("Error with table: %s | %s", table, err.Error()))
		}
		bar.Add(1)
	}
	return nil
}

// Check the number of rows for every table
//select table_name, diststyle, sortkey1, tbl_rows
//from svv_tables t1,
//svv_table_info t2
//where t1.table_schema = t2.schema
//and t1.table_name = t2.table
//and t1.table_schema not in ('pg_catalog', 'information_schema')
//and t1.table_type = 'BASE TABLE'
//order by tbl_rows desc;

// UnloadDB is delegated to dump all the DATABASE data into an S3 folder
func UnloadDB(conn *sql.DB, schemaName, bucket, prefix, role string) {
	tables, err := GetAllTables(conn, schemaName)
	if err != nil {
		panic(err)
	}

	log.Println(helper.MarshalIndent(tables))
	bar := progressbar.Default(int64(len(tables)))
	for _, table := range tables {
		bar.Describe(table)
		var q string = fmt.Sprintf(`unload ('select * from %[1]s')
				to 's3://%[2]s/%[3]s/%[1]s_'
				iam_role '%[4]s'
				header
				parallel off
				CSV DELIMITER AS '|' ;`, table, bucket, prefix, role)

		_, err := conn.Exec(q)
		if err != nil {
			log.Println(fmt.Sprintf("Error with table: %s | %s", table, err.Error()))
		}
		bar.Add(1)
	}
	bar.Close()
}

func LoadDB(conn *sql.DB, schemaName, bucket, prefix, role string) {
	// Get all tables from the DB
	tables, err := GetAllTables(conn, schemaName)
	if err != nil {
		panic(err)
	}

	// Get all table dump
	objects, err := S3utils.ListBucketObjects(bucket, prefix)
	if err != nil {
		panic(err)
	}

	for i := range objects {
		split := strings.Split(path.Base(objects[i]), "_")
		objects[i] = stringutils.JoinSeparator("_", split[:len(split)-1]...)
	}

	tablesMap := stringutils.ArrayToMap(tables)
	//objectsMap := stringutils.ArrayToMap(objects)

	var qTest string = `copy %s
    from '%s'
    iam_role '%s'
    DELIMITER '|'
	NOLOAD
    IGNOREHEADER 1
    IGNOREBLANKLINES
    REMOVEQUOTES
    EMPTYASNULL;`

	var q string = `copy %s
    from '%s'
    iam_role '%s'
    DELIMITER '|'
    IGNOREHEADER 1
    IGNOREBLANKLINES
    REMOVEQUOTES
    EMPTYASNULL;`

	bar := progressbar.Default(int64(len(tablesMap)))

	for k := range tablesMap {
		bar.Describe(k)
		bar.Add(1)
		count := stringutils.Count(objects, k)
		var sb strings.Builder
		for i := 0; i < count; i++ {
			fname := fmt.Sprintf("s3://%s/%s/%s_00%d", bucket, prefix, k, i)
			sb.WriteString(fmt.Sprintf(qTest, k, fname, role))
		}
		_, err := conn.Exec(sb.String())
		if err != nil {
			log.Println("QUERY:\n", sb.String())
			log.Println(fmt.Sprintf("Error with table: %s | %s", k, err.Error()))
		} else {
			sb.Reset()
			sb.WriteString(fmt.Sprintf("truncate %s;\n", k))
			for i := 0; i < count; i++ {
				fname := fmt.Sprintf("s3://%s/%s/%s_00%d", bucket, prefix, k, i)
				sb.WriteString(fmt.Sprintf(q, k, fname, role))
			}
			_, err := conn.Exec(sb.String())
			if err != nil {
				log.Println("QUERY:\n", sb.String())
				log.Println(fmt.Sprintf("Error with table: %s | %s", k, err.Error()))
			}
		}
	}
	bar.Finish()
	log.Println()

}

func GetAllTables(conn *sql.DB, schema string) ([]string, error) {
	rows, err := conn.Query(`
		select t.table_name
		from information_schema.tables t
		where t.table_schema = $1
		  and t.table_type = 'BASE TABLE'
		order by t.table_name;`, schema)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
