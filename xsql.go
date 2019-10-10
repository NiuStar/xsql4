package xsql4

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NiuStar/log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
	"github.com/NiuStar/Basic"
	"reflect"
)

var (
	TAG string = " <XSQL>"
	xs *XSql = nil
	dbConfig *Basic.DBConfig = nil
)

func init() {

	serConfig := Basic.GetServerConfig()
	if serConfig == nil {
		return 
	}
	c := serConfig.DBConfig

	if c == nil {
		return
	}
	xs = InitSql(c.DB_user,c.DB_password,c.DB_server,c.DB_port,c.DB_name,c.DB_charset)

	CreateDB(c.DB_name,c.DB_charset,c.DB_charset + "_general_ci")

	UseDataBase(c.DB_name)

	dbConfig = c
}

func GetServerDB() *XSql {

	return xs
}

type XSql struct {

	db        *sql.DB
	name      string
	password  string
	ip        string
	port      string
	sqlName   string
	charset	  string
	stmts	  map[string]*sql.Stmt
}

type XSqlOrder struct {
	txopen    int
	tx        *sql.Tx
	xs        *XSql
	reqString string
	args []interface{}
}

func CreateInstance(xs *XSql) *XSqlOrder {
	o := new(XSqlOrder)
	o.xs = xs
	return o
}

func InitSql(name string, password string, ip string, port string, sqlName string,charset string) *XSql {
	db := connectDB(name, password, ip, port, sqlName,charset)
	fmt.Println("初始化数据库成功")
	s := new(XSql)
	s.db = db
	s.name = name
	s.password = password
	s.ip = ip
	s.port = port
	s.sqlName = sqlName
	s.stmts = make(map[string]*sql.Stmt)

	return s
}

func connectDB(name string, password string, ip string, port string, sqlName string,charset string) *sql.DB {
	db, err := sql.Open("mysql", name+":"+password+"@tcp("+ip+":"+port+")/"+sqlName+"?charset=" + charset)

	checkErr(err)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(10 * time.Minute)
	err = db.Ping()

	checkErr(err)

	return db
}


func (s *XSqlOrder) Qurey(suffixes string,args... interface{}) { //执行sql语句
	s.reqString = suffixes
	s.args = args
}

func checkErr(err error) {
	if err != nil {
		log.Error(TAG,err)
	}
}
func (s *XSqlOrder) ExecuteForJson() string { //执行sql语句得到json

	body, err := json.Marshal(s.Execute())
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func (s *XSqlOrder) GetSQLString() string {
	return s.reqString
}

//受影响条数，为-1时是出错了
func (s *XSqlOrder) ExecuteNoResult() (num int64) {
	//SQL
	fmt.Println("ExecuteNoResult执行sql语句: " + s.reqString)
	var result sql.Result
	var err error

	if len(s.args) > 0 {
		if s.tx != nil {
			result, err = s.tx.Exec(s.reqString,s.args...)
			fmt.Println("rows, err = s.tx.Exec(s.reqString,s.args...) ")
		} else {
			fmt.Println("rows, err = s.xs.db.Exec(s.reqString,s.args...)")
			result, err = s.xs.db.Exec(s.reqString,s.args...)
		}
	} else {
		if s.tx != nil {
			result, err = s.tx.Exec(s.reqString)
			fmt.Println("rows, err = s.tx.Exec(s.reqString) ")
		} else {
			fmt.Println("rows, err = s.xs.db.Exec(s.reqString)")
			result, err = s.xs.db.Exec(s.reqString)
		}
	}

	if err != nil {
		panic(err)
		return -1
	}
	if result == nil {
		return 0
	}
	num ,err = result.RowsAffected()
	if err != nil {
		panic(err)
		return -1
	}
	return num

	//RowsAffected
}
//插入的最后一条的id
func (s *XSqlOrder) InsertExecute() (lastid int64) {
	defer func() {
		if err := recover(); err != nil {

			fmt.Println("数据库执行错误：", err)
			panic(err)
		}
	}()

	var result sql.Result
	var err error

	if len(s.args) > 0 {
		if s.tx != nil {
			result, err = s.tx.Exec(s.reqString,s.args...)
			fmt.Println("rows, err = s.tx.Exec(s.reqString,s.args...) ")
		} else {
			fmt.Println("rows, err = s.xs.db.Exec(s.reqString,s.args...)")
			result, err = s.xs.db.Exec(s.reqString,s.args...)
		}
	} else {
		if s.tx != nil {
			result, err = s.tx.Exec(s.reqString)
			fmt.Println("rows, err = s.tx.Exec(s.reqString) ")
		} else {
			fmt.Println("rows, err = s.xs.db.Exec(s.reqString)")
			result, err = s.xs.db.Exec(s.reqString)
		}
	}


	if err != nil {
		return -1
	}
	lastid ,err = result.LastInsertId()
	if err != nil {
		return -1
	}
	return lastid
}

func (s *XSqlOrder) Execute() (results []map[string]interface{}) { //SQL

	defer func() {
		if err := recover(); err != nil {

			fmt.Println("数据库执行错误：", err)
			panic(err)
		}

	}()

	fmt.Println("Execute执行sql语句: " + s.reqString)
	var rows *sql.Rows
	var err error
	if len(s.args) > 0 {

		if s.tx != nil {
			rows, err = s.tx.Query(s.reqString,s.args...)
			fmt.Println("rows, err = s.tx.Query(s.reqString,s.args...) ")
		} else {
			fmt.Println("rows, err = s.xs.db.Query(s.reqString,s.args...)")
			rows, err = s.xs.db.Query(s.reqString,s.args...)
		}
	} else {
		if s.tx != nil {
			fmt.Println("rows, err = s.tx.Query(s.reqString)")
			rows, err = s.tx.Query(s.reqString)
		} else {
			fmt.Println("rows, err = s.xs.db.Query(s.reqString)")
			rows, err = s.xs.db.Query(s.reqString)
		}
	}

	if err != nil {
		fmt.Println("error: ", err)
		checkErr(err)
		return nil
	}

	defer rows.Close()

	columns, err2 := rows.Columns()
	if err2 != nil {
		log.Error(TAG,err2) // proper error handling instead of panic in your app
		return nil
	}

	if len(columns) <= 0 {
		return nil
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]

	}
	//var results []map[string]interface{}

	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		t := make(map[string]interface{})

		for i, col := range values {
			//fmt.Println("columns:",columns[i],string(col.([]byte)))
			if col == nil {
				t[columns[i]] = nil
			} else {
				switch reflect.ValueOf(col).Type().String() {
				case "[]uint8":{
						t[columns[i]] = string(col.([]uint8))
					}
				default:

					t[columns[i]] = col
				}
			}

		}
		results = append(results, t)

	}
	return results

}
func byte2Int(value []byte) int64 {

	result, err := strconv.ParseInt(string(value), 10, 64)
	checkErr(err)
	return result
}
func byte2Float(value []byte) float64 {

	result, err := strconv.ParseFloat(string(value), 64)
	checkErr(err)
	return result
}

func byte2String(value []byte) string {
	return string(value)
}