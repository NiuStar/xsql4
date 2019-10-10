package xsql4

import (
	"strings"
	"fmt"
	"github.com/NiuStar/reflect"

)

const DEFAULTENGINE  = "InnoDB"
type fieldTable struct {
	tableName string
	table interface{}
	comment string
	pkey []string
	ukey map[string][]string
	ikey map[string][]string
}

func (this *fieldTable) AddPrimaryKey(key string) {
	this.pkey = append(this.pkey, key)
}

func (this *fieldTable) AddUniqueKey(key string, keys []string) {
	if 0 == len(keys) {
		return
	} else if 0 == len(key) {
		if 1 == len(keys) {
			this.ukey[keys[0]] = keys
		} else {
			return
		}
	}
	this.ukey[key] = keys
}

func (this *fieldTable) AddIndexKey(key string, keys []string) {
	if 0 == len(keys) {
		return
	} else if 0 == len(key) {
		if 1 == len(keys) {
			this.ikey[keys[0]] = keys
		} else {
			return
		}
	}
	this.ikey[key] = keys
}

func (this *fieldTable) GetUniqueKey() map[string][]string {
	return this.ukey
}

func (this *fieldTable) GetIndexKey() map[string][]string {
	return this.ikey
}

func (this *fieldTable) GetPrimaryKey() []string {
	return this.pkey
}

func Register(tableName string,table interface{},comment string) *fieldTable {
	t := &fieldTable{tableName:tableName,table:table,comment:comment,ukey:make(map[string][]string),ikey:make(map[string][]string)}
	t.copyToMySQL()
	return t
}

func (this *fieldTable)copyToMySQL() {

	t := &Table{tableName:this.tableName,ukey:make(map[string][]string),ikey:make(map[string][]string),engine:DEFAULTENGINE,comment:this.comment}
	t.columns = this.createTable()
	t.pkey = this.pkey
	t.ukey = this.ukey
	t.ikey = this.ikey
	verificationUpdate(t)
	//fmt.Println(t)
}

func (this *fieldTable)createTable() (columList map[string]*Column) {

	columList = make(map[string]*Column)
	sql := "CREATE TABLE IF NOT EXISTS `" + this.tableName + "` (" + DB_SQL_ENTER_STRING
	tableType := reflect.GetReflectType(this.table)
	var columns_head = "`id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'id自增长'"
	var hadID = false
	var columns = ""
	for i:=0;i<tableType.NumField();i++ {

		fieldTag := tableType.Field(i).Tag

		if len(fieldTag.Get("json")) > 0 && len(fieldTag.Get("type")) > 0 {

			mark := fieldTag.Get("mark")
			name := fieldTag.Get("json")

			if strings.EqualFold(name,"id") {
				hadID = true
			}

			type_ := fieldTag.Get("type")
			default_ := fieldTag.Get("default")
			comment := fieldTag.Get("comment")
			canNull := !strings.EqualFold(strings.ToLower(fieldTag.Get("required")),"yes")

			if name == "bseUnitId" {
				//fmt.Println("测试数据：",strings.ToLower(fieldTag.Get("canNull")),canNull)
			}

			pri := strings.Contains(strings.ToUpper(mark), " PRIMARY KEY ")

			extra := ""
			if strings.Contains(strings.ToUpper(mark), " AUTO_INCREMENT ") {
				extra = "auto_increment"
			}


			col := &Column{Index:i,Name:name,Type:type_,Default:default_,Comment:comment,IsNull:canNull,Mark:mark,Extra:extra,PRIMARY:pri}
			columList[name] = col

			column := ""
			if 0 != len(columns) {
				column = "," + DB_SQL_ENTER_STRING
			}
			column += "\t`"
			column += name + "` "
			column += type_

			if !canNull {
				column += " NOT NULL "
			} else {
				column += " "
			}

			if 0 != len(mark) {
				column += mark + " "
			}
			if 0 == len(mark) ||
				-1 == strings.Index(strings.ToUpper(mark), "AUTO_INCREMENT") {
				if 0 != len(default_) {
					column += "DEFAULT " + default_ + " "
				}
			}
			column += "COMMENT '" + comment + "'"
			columns += column
		}
	}

	if !hadID {
		sql += columns_head + "," + DB_SQL_ENTER_STRING
	}
	fmt.Println("dbConfig.DB_charset:",dbConfig)
	sql += columns
	sql += DataBaseKey(this.pkey,this.ukey,this.ikey)
	sql += DB_SQL_ENTER_STRING
	sql += ") ENGINE=" + DEFAULTENGINE + " DEFAULT CHARSET=" +
		dbConfig.DB_charset + " COMMENT='" + this.comment + "';"
	xsql := CreateInstance(xs)
	xsql.Qurey(sql)
	fmt.Println("result:",xsql.ExecuteNoResult())
	return
}

func DataBaseKey(pkey []string,ukey map[string][]string,ikey map[string][]string) string {
	sql := ""
	//PRIMARY KEY
	if len(pkey) > 0 {
		sql += "," + DB_SQL_ENTER_STRING + "\tPRIMARY KEY ("
		for i := 0; i < len(pkey); i++ {
			if 0 != i {
				sql += ",`"
			}
			sql += "`" + pkey[i] + "`"
		}
		sql += ") COMMENT 'PRIMARY'"
	}

	//UNIQUE KEY
	if ukey != nil && len(ukey) > 0 {
		sql += "," + DB_SQL_ENTER_STRING

		count := 0
		for key, val := range ukey {
			if 0 != count {
				sql += "," + DB_SQL_ENTER_STRING
			}
			sql += "\tUNIQUE KEY `" + key + "` ("
			for i := 0; i < len(val); i++ {
				if 0 != i {
					sql += ","
				}
				sql += "`" + val[i] + "`"
			}
			sql += ") COMMENT 'UNIQUE'"
			//sql += DB_SQL_ENTER_STRING
			count += 1
		}
	}
	///sql += DB_SQL_ENTER_STRING

	//INDEX KEY
	if ikey != nil && len(ikey) > 0 {
		sql += "," + DB_SQL_ENTER_STRING

		count := 0
		for key, val := range ikey {
			if 0 != count {
				sql += "," + DB_SQL_ENTER_STRING
			}
			sql += "\tKEY `" + key + "` ("
			for i := 0; i < len(val); i++ {
				if 0 != i {
					sql += ","
				}
				sql += "`" + val[i] + "`"
			}
			sql += ") COMMENT 'KEY'"
			//sql += DB_SQL_ENTER_STRING
			count += 1
		}
	}
	return sql
}

func IsEXITSTable(tableNames... string) map[string]bool {

	sql := CreateInstance(GetServerDB())

	var sqlString string = "SELECT table_name FROM information_schema.TABLES WHERE table_name in("
	for index,name := range tableNames {
		sqlString += "'" + name + "'"
		if len(tableNames) - 1 > index {

			sqlString  += " , "
		}
	}

	sqlString  += ");"
	sql.Qurey(sqlString)
	l := sql.Execute()
	list_sqlName := make(map[string]bool)
	for _,obj := range l {
		list_sqlName[strings.ToLower(string(obj["table_name"].([]uint8)))] = true
	}
	return list_sqlName
}
