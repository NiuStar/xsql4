package xsql4

import(
	"bytes"
	"reflect"
	reflect2 "github.com/NiuStar/reflect"
	"strings"
	"fmt"
	"github.com/NiuStar/xsql4/Type"
)


/*
b := "BseUnitId"
	c := "2"
	xsql3.InsertDB(t,&BridgeBaseBasic{BseUnitId:&b,Node:&c})
*/

func InsertDB(xs *XSqlOrder, list... Type.DBOperation) (latsid int64) {


	if len(list) <= 0 {
		return
	}

	if xs == nil {
		xs = CreateInstance(GetServerDB())
	}

	var fields []string
	var fields2 []string

	type_ := reflect2.GetReflectType(list[0])

	for i:=0;i<type_.NumField();i++ {

		fieldTag := type_.Field(i).Tag
		if len(fieldTag.Get("json")) > 0  && len(fieldTag.Get("type")) > 0 && !strings.Contains(strings.ToUpper(fieldTag.Get("mark")),"AUTO_INCREMENT"){

			jsonName := type_.Field(i).Name
			if len(type_.Field(i).Tag.Get("json")) > 0 {
				jsonName = type_.Field(i).Tag.Get("json")
			}
			fields = append(fields, type_.Field(i).Name)
			fields2 = append(fields2,jsonName)
		}
	}


	var sqlOrder bytes.Buffer
	sqlOrder.Grow(8192)
	sqlOrder.WriteString("insert into " + list[0].TableName())

	{
		sqlOrder .WriteString( "(")

		for index1,_ := range fields {



			sqlOrder .WriteString(  fields2[index1])

			if index1 != len(fields) - 1 {
				sqlOrder .WriteString(  " , ")
			}
		}

		sqlOrder .WriteString(  ")")
	}

	sqlOrder .WriteString( " VALUES ")

	var args []interface{}
	for index,value := range list {
		rvalue := reflect2.GetReflectValue(value)
		sqlOrder .WriteString( "(")
		for index1,name := range fields {

			field := rvalue.FieldByName(name)
			fmt.Println("args:",field,field.Kind())

			sqlOrder .WriteString(  "?")

			if index1 != len(fields) - 1 {
				sqlOrder .WriteString(  " , ")
			}

			if field.Kind() == reflect.Struct {

				if Type.IsTabelType(field.Type()) {
					if field.IsValid() && !field.Addr().Interface().(Type.TableType).IsNil() {
						args = append(args,field.Addr().Interface().(Type.TableType).Value())

					} else {

						fmt.Println("Struct:",field.String(),field.IsValid(),field.Addr().Interface(),field.Addr().Interface().(Type.TableType).IsNil())
						args = append(args,nil)
					}
				} else {
					//该处需添加子关联表逻辑
				}


			} else {
				if field.IsValid() && !field.IsNil() {
					args = append(args,field.Interface().(Type.TableType).Value())

				} else {
					fmt.Println("Interface:",field)
					args = append(args,nil)
				}
			}

		}

		sqlOrder .WriteString(  ")")
		if index != len(list) - 1 {
			sqlOrder .WriteString( " , ")
		}
	}

	fmt.Println("sqlOrder.String():",sqlOrder.String(),args)

	xs.Qurey(sqlOrder.String(),args...)
	return xs.InsertExecute()
}

func DeleteDBALL(xsc *XSqlOrder,tableName string) (num int64) {
	if xsc == nil {
		xsc = CreateInstance(GetServerDB())
	}
	xsc.Qurey("delete from " + tableName)
	return xsc.ExecuteNoResult()
}
/*
b := "BseUnitId"
c := "2"
xsql3.DeleteDB(&BridgeBaseBasic{BseUnitId:&b},&BridgeBaseBasic{Node:&c})
同一个对象内的是and关联，不同的对象内的是OR关联
*/
func DeleteDB(xsc *XSqlOrder,list... Type.DBOperation) (num int64) {

	if len(list) <= 0 {
		return
	}

	if xsc == nil {
		xsc = CreateInstance(GetServerDB())
	}

	sqlStr,args := Where(list...)
	/*xs.Begin()
	{

		var sqlString = "use " + Basic.GetServerConfig().DBConfig.DB_name + ";"
		xs.Qurey(sqlString)
		xs.ExecuteNoResult()
	}*/
	sqlOrder := "delete from " + list[0].TableName() + " where " + sqlStr
	xsc.Qurey(sqlOrder,args...)
	return xsc.ExecuteNoResult()
	//xs.Commit()
}

func UpdateDB(xsc *XSqlOrder,list Type.DBOperation,where string,argvs []interface{}) (num int64) {


	if xsc == nil {
		xsc = CreateInstance(GetServerDB())
	}


	var sqlOrder bytes.Buffer
	sqlOrder.Grow(8192)

	var fields []reflect.StructField

	type_ := reflect2.GetReflectType(list)

	for i:=0;i<type_.NumField();i++ {

		fieldTag := type_.Field(i).Tag

		if len(fieldTag.Get("json")) > 0 {
			fields = append(fields, type_.Field(i))
		}
	}

	var args []interface{}
	{
		rvalue := reflect2.GetReflectValue(list)
		for _,field := range fields {

			fieldValue := rvalue.FieldByName(field.Name)

			if fieldValue.IsValid() && !fieldValue.IsNil() {
				jsonName := field.Tag.Get("json")
				sqlOrder .WriteString(  jsonName + "=? , ")
				args = append(args,fieldValue.Interface().(Type.TableType).Value())
			}
		}
	}

	args = append(args,argvs...)
	sqlStr := sqlOrder.String()

	if strings.HasSuffix(sqlStr," , ") {
		sqlOrder.Truncate(len(sqlStr) - 3)
	}

	sqlStr = "UPDATE " + list.TableName() + " SET " + sqlOrder.String() + " where " + where
	xsc.Qurey(sqlStr,args...)
	return xsc.ExecuteNoResult()
}

func SelectDB(list... Type.DBOperation) []map[string]interface{} {
	if len(list) <= 0 {
		return []map[string]interface{}{}
	}

	xs := CreateInstance(GetServerDB())
	sqlStr,args := Where(list...)
	sqlOrder := "select * from " + list[0].TableName() + " where " + sqlStr
	xs.Qurey(sqlOrder,args...)
	return xs.Execute()
}
