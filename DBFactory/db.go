package DBFactory

import (
	"github.com/NiuStar/xsql4/Type"
	"fmt"
	"strconv"
	"github.com/NiuStar/xsql4"
	"bytes"
	"reflect"
	"strings"
	reflect2 "github.com/NiuStar/reflect"
)

/*
DBFactory.NewDBFactory().SetTable(bridge,test).
		Count().
		SetFields(&bridge.BseUnitId,&bridge.Len,&bridge.ID,&test.ID,&test.Name).
		SetConditions(test.ID.EqualFold(&bridge.ID),test.ID.GreaterEqual(2)).
			SetGroupBy(&bridge.ID).
				SetOrderBy(&bridge.ID).
					DESC().
						Limit(0,10).
						GetResults()
数据库使用方法，查询和更新两个方法需提前实例化对象，调用NewDBFactory()方法，
SetTable设置更新、查询的数据库对象
Count可以用来统计行数
SetFields设置查询内容，否则为全部查询，更新的时候不需要设置
SetConditions设置条件，目前支持=，>，>=，<，<=
SetGroupBy Group By的字段
SetOrderBy Order By的字段
DESC 默认为正序，否则为倒序
Limit 分页
GetResults 获取结果集
*/

/*
DBFactory.NewDBFactory().SetConditions(test.ID.EqualFold(&bridge.ID)).
		UpdateDB(bridge,test)
数据库使用方法，查询和更新两个方法需提前实例化对象，调用NewDBFactory()方法，
SetConditions设置条件，目前支持=，>，>=，<，<=
UpdateDB 更新哪几张表，里面有内容的字段将会被更新，没有内容的字段不会被更新
*/

type DBFactory struct {
	tables []Type.DBOperation
	fields []Type.TableType

	conditions []string
	conditionsValue []interface{}
	orConditions []string
	orConditionsValue []interface{}
	hadCount	 bool
	groupbys []Type.TableType
	orderbys []Type.TableType
	desc     bool
	limitStart,limitEnd int

	xsql     *xsql3.XSqlOrder
	transactionBegin bool

	unionAllDB []*DBFactory
	unionDISTINCTDB []*DBFactory

	joins []*Join
}

func NewDBFactory() *DBFactory {
	return &DBFactory{limitStart:-1,limitEnd:-1,xsql:xsql3.CreateInstance(xsql3.GetServerDB()),transactionBegin:false}
}
//如果起始值都为-1，则不限制条数
func (db *DBFactory)Limit(start,end int) *DBFactory {
	db.limitStart = start
	db.limitEnd = end
	return db
}

func (db *DBFactory)Count() *DBFactory {
	db.hadCount = true
	return db
}
//设置查询哪些表
func (db *DBFactory)SetTable(tables... Type.DBOperation) *DBFactory {
	db.tables = tables
	return db
}
//添加查询表
func (db *DBFactory)AddTable(tables... Type.DBOperation) *DBFactory {
	db.tables = append(db.tables,tables...)
	return db
}

//设置查询哪些字段
func (db *DBFactory)SetFields(fields... Type.TableType) *DBFactory {
	db.fields = fields
	return db
}
//设置查询哪些字段
func (db *DBFactory)AddFields(fields... Type.TableType) *DBFactory {
	db.fields = append(db.fields,fields...)
	return db
}
//设置AND查询条件或者插入条件
func (db *DBFactory)SetConditions(conditions... *Type.Condition) *DBFactory {

	db.conditions = []string{}
	db.conditionsValue = []interface{}{}

	for _,condition := range conditions {
		str,args := condition.String()
		db.conditions = append(db.conditions,str)
		db.conditionsValue = append(db.conditionsValue,args...)
	}
	return db
}
//添加AND查询条件或者插入条件
func (db *DBFactory)AddCondition(conditions... *Type.Condition) *DBFactory {

	for _,condition := range conditions {
		str,args := condition.String()
		db.conditions = append(db.conditions,str)
		db.conditionsValue = append(db.conditionsValue,args...)
	}
	//db.conditions = append(db.conditions,conditions...)
	return db
}

//设置OR查询条件或者插入条件
func (db *DBFactory)SetORConditions(conditions... *Type.Condition) *DBFactory {
	db.orConditions = []string{}
	db.orConditionsValue = []interface{}{}

	for _,condition := range conditions {
		str,args := condition.String()
		db.orConditions = append(db.orConditions,str)
		db.orConditionsValue = append(db.orConditionsValue,args...)
	}
	//db.orConditions = conditions
	return db
}
//添加OR查询条件或者插入条件
func (db *DBFactory)AddORCondition(conditions... *Type.Condition) *DBFactory {
	for _,condition := range conditions {
		str,args := condition.String()
		db.orConditions = append(db.orConditions,str)
		db.orConditionsValue = append(db.orConditionsValue,args...)
	}
	//db.orConditions = append(db.orConditions,conditions...)
	return db
}

func (db *DBFactory)DESC() *DBFactory {
	db.desc = true
	return db
}

//设置Order By 条件
func (db *DBFactory)SetOrderBy(orderbys... Type.TableType) *DBFactory {
	db.orderbys = orderbys
	return db
}

//设置group By 条件
func (db *DBFactory)SetGroupBy(groupbys... Type.TableType) *DBFactory {
	db.groupbys = groupbys
	return db
}


//设置查询哪些字段
func (db *DBFactory)SetJoins(joins... *Join) *DBFactory {
	db.joins = joins
	return db
}
//设置查询哪些字段
func (db *DBFactory)AddJoin(joins... *Join) *DBFactory {
	db.joins = append(db.joins,joins...)
	return db
}


func (db *DBFactory)GetResultsOperation() (results []map[string]Type.DBOperation) {
	return db.ParseResults(db.GetResults())
}

func (db *DBFactory)getTableFields(operation Type.DBOperation) string {

	var str string = ""
	type_ := reflect2.GetReflectType(operation)

	for i:=0;i<type_.NumField();i++ {
		if len(type_.Field(i).Tag.Get("json")) > 0 && len(type_.Field(i).Tag.Get("type")) > 0 {

			jsonName := type_.Field(i).Tag.Get("json")

			str += operation.TableName() + "." + jsonName + " as " + operation.TableName() + "_" + jsonName + ","
		}

	}
	return str[:len(str)-1]
}

func (db *DBFactory)GetCount() int64 {

	list := db.Count().GetResults()
	switch list[0]["count"].(type) {

	case string:
		{
			count, err := strconv.ParseInt(list[0]["count"].(string), 10, 64)
			if err != nil {
				return -1
			}
			return count
		}
		default:
			return list[0]["count"].(int64)

	}
	return list[0]["count"].(int64)
}

func (db *DBFactory)GetResults() []map[string]interface{} {

	str,args := db.String()

	fmt.Println("strSql:",str,args)
	db.xsql.Qurey(str,args...)
	l := db.xsql.Execute()
	return l
}

func (db *DBFactory)GetMax(max Type.TableType) interface{} {


	str := "select max(" + max.Name() + ") as max from " + max.TableName()

	db.xsql.Qurey(str)
	l := db.xsql.Execute()
	return l[0]["max"]
}

func PrintStruct(data interface{}) {
	_type := reflect2.GetReflectType(data)
	_value := reflect2.GetReflectValue(data)

	for i:=0;i<_value.NumField();i++ {

		if _value.Field(i).IsValid() {
			if _value.Field(i).CanAddr() && !_value.Field(i).Addr().IsNil() {


				if  Type.IsTabelType(_value.Field(i).Type()) {
					if  _value.Field(i).Addr().MethodByName("IsNil").Call([]reflect.Value{})[0].Bool() {
						fmt.Println(_type.Field(i).Name,nil)
					} else {
						fmt.Println(_type.Field(i).Name,_value.Field(i).Addr().MethodByName("Value").Call([]reflect.Value{})[0].Interface())
					}
				} else {
					fmt.Println(_type.Field(i).Name, _value.Field(i))
				}


			} else {
				fmt.Println(_type.Field(i).Name,nil)
			}

		} else {
			fmt.Println(_type.Field(i).Name,nil)
		}

	}

}

func (db *DBFactory)InsertDB(list... Type.DBOperation) (latsid int64) {
	return xsql3.InsertDB(db.xsql,list...)
}

func (db *DBFactory)DeleteDB(obj Type.DBOperation) (num int64) {
	//return xsql3.DeleteDB(db.xsql,list...)

	str1,argvs :=  db.where()

	sqlOrder := "delete from " + obj.TableName() + str1
	db.xsql.Qurey(sqlOrder,argvs...)
	return db.xsql.ExecuteNoResult()
}

func (db *DBFactory)DeleteDBALL(tableName string) (num int64) {
	return xsql3.DeleteDBALL(db.xsql,tableName)
}

func (db *DBFactory)UpdateDB(list... Type.DBOperation) (num int64) {

	var args []interface{}
	var where string
	var whereArgs []interface{}

	for index,condition := range db.conditions {

		str1 := condition
		where += str1
		if index != len(db.conditions) - 1 {
			where += " AND "
		}
	}
	whereArgs = append(whereArgs,db.conditionsValue...)

	var sqlOrder bytes.Buffer
	sqlOrder.Grow(8192)
	var tableName string
	for index,object := range list {

		tableName += object.TableName()

		if index != len(list) - 1 {
			tableName += ","
		}

		var fields []reflect.StructField

		type_ := reflect2.GetReflectType(object)

		for i:=0;i<type_.NumField();i++ {

			fieldTag := type_.Field(i).Tag

			if len(fieldTag.Get("json")) > 0 && !strings.Contains(strings.ToUpper(fieldTag.Get("mark")),"AUTO_INCREMENT") && !(len(fieldTag.Get("type")) <= 0) {

				//fmt.Println("mark:", type_.Field(i).Name,fieldTag.Get("comment"))
				fields = append(fields, type_.Field(i))
			}
		}

		{
			rvalue := reflect2.GetReflectValue(object)
			for _,field := range fields {

				fieldValue := rvalue.FieldByName(field.Name)

				if fieldValue.IsValid() && Type.IsTabelType(fieldValue.Type()) {

					obj := fieldValue.Addr().Interface().(Type.TableType)

					if !obj.IsNil() {
						jsonName := field.Tag.Get("json")

						sqlOrder .WriteString( object.TableName()+ "." + jsonName + "=? , ")
						args = append(args,fieldValue.Addr().Interface().(Type.TableType).Value())
					}

				}
			}
		}
	}

	args = append(args,whereArgs...)
	sqlStr := sqlOrder.String()

	if strings.HasSuffix(sqlStr," , ") {
		sqlOrder.Truncate(len(sqlStr) - 3)
	}

	sqlStr = "UPDATE " + tableName + " SET " + sqlOrder.String()
	if len(where) > 0 {
		sqlStr += " where " + where
	}

	db.xsql.Qurey(sqlStr,args...)
	return db.xsql.ExecuteNoResult()
}

func (db *DBFactory)ParseResults(datas []map[string]interface{}) (results []map[string]Type.DBOperation) {

	var fields = make(map[string]map[string]string)
	var tables = make(map[string]reflect.Type)

	for _,table := range db.tables {
		tables[table.TableName()] = reflect2.GetReflectType(table)

		type_ := reflect2.GetReflectType(table)
		for i:=0;i<type_.NumField();i++ {
			jsonName := type_.Field(i).Name
			if len(type_.Field(i).Tag.Get("json")) > 0 {
				jsonName = type_.Field(i).Tag.Get("json")
			}
			if fields[table.TableName()] == nil {
				fields[table.TableName()] = make(map[string]string)
			}
			fields[table.TableName()][jsonName] = type_.Field(i).Name
		}
	}

	for index,data := range datas {
		var result_list map[string]reflect.Value = make(map[string]reflect.Value)
		for name,value := range data {

			if strings.Index(name,"_") < 0 || value == nil {
				continue
			}
			tableName := name[:strings.LastIndex(name,"_")]
			fieldName := name[len(tableName) + 1:]

			var v reflect.Value
			if !result_list[tableName].IsValid() {
				v = reflect.New(tables[tableName]).Elem()

				v.Addr().MethodByName("")

				result_list[tableName] = v
			} else {
				v = result_list[tableName]
			}

			if Type.IsTabelType(v.FieldByName(fields[tableName][fieldName]).Type()) {

				if v.FieldByName(fields[tableName][fieldName]).Type().Name() == "Int" {
					fmt.Println("Int fieldName:",fieldName,value)
					value_ := reflect.ValueOf(value)
					if value_.Kind() == reflect.String {

						value_o,_ := strconv.ParseInt(value.(string),10,64)

						v.FieldByName(fields[tableName][fieldName]).Addr().MethodByName("SetValue").Call([]reflect.Value{reflect.ValueOf(value_o)})
					} else {
						v.FieldByName(fields[tableName][fieldName]).Addr().MethodByName("SetValue").Call([]reflect.Value{value_})
					}

				} else if v.FieldByName(fields[tableName][fieldName]).Type().Name() == "Float" {
					fmt.Println("Float fieldName:",fieldName,value)
					value_ := reflect.ValueOf(value)
					if value_.Kind() == reflect.String {
						value_o,_ := strconv.ParseFloat(value.(string),64)
						v.FieldByName(fields[tableName][fieldName]).Addr().MethodByName("SetValue").Call([]reflect.Value{reflect.ValueOf(value_o)})
					} else {
						v.FieldByName(fields[tableName][fieldName]).Addr().MethodByName("SetValue").Call([]reflect.Value{value_})
					}
				} else {
					fmt.Println("String fieldName:",fieldName,value)
					v.FieldByName(fields[tableName][fieldName]).Addr().MethodByName("SetValue").Call([]reflect.Value{reflect.ValueOf(value)})
				}

			}
		}

		if index >= len(results) {
			results = append(results,make(map[string]Type.DBOperation))
		}
		for name,value := range result_list {
			results[index][name] = xsql3.ScanStructInterface(value.Addr().Interface().(Type.DBOperation))
		}
	}
	return results
}