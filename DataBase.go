package xsql4

import (
	"strings"
)

const (
	DB_SQL_ENTER_STRING = "\r\n"
)

func IsEXITSDB(dbNames... string) map[string]bool {

	sql := CreateInstance(GetServerDB())

	var sqlString string = "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA where SCHEMA_NAME in("
	for index,name := range dbNames {
		sqlString += "'" + name + "'"
		if len(dbNames) - 1 > index {

			sqlString  += " , "
		}
	}

	sqlString  += ");"
	sql.Qurey(sqlString)
	l := sql.Execute()
	list_sqlName := make(map[string]bool)
	for _,obj := range l {
		list_sqlName[strings.ToLower(obj["SCHEMA_NAME"].(string))] = true
	}
	return list_sqlName
}

func CreateDB(dbName,charset,COLLATE string) {

	//utf8_general_ci
	sql := CreateInstance(GetServerDB())
	var sqlString = "CREATE DATABASE IF NOT EXISTS " + dbName + " default charset " + charset + " COLLATE " + COLLATE + ";"
	sql.Qurey(sqlString)
	sql.ExecuteNoResult()
}

func UseDataBase(dbName string) {
	sql := CreateInstance(GetServerDB())
	var sqlString = "use " + dbName + ";"
	sql.Qurey(sqlString)
	sql.ExecuteNoResult()
}
