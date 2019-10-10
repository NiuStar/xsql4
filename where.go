package xsql4

import (
	"bytes"
	"reflect"
	reflect2 "github.com/NiuStar/reflect"
	"strings"
	"github.com/NiuStar/xsql4/Type"
)

func Where(list... Type.DBOperation) (sqlString string,args []interface{}) {

	if len(list) <= 0 {
		return "",args
	}

	var sqlOrder bytes.Buffer
	sqlOrder.Grow(8192)

	var fields []reflect.StructField

	type_ := reflect2.GetReflectType(list[0])

	for i:=0;i<type_.NumField();i++ {

		fieldTag := type_.Field(i).Tag

		if len(fieldTag.Get("json")) > 0 {
			fields = append(fields, type_.Field(i))
		}
	}

	for _,value := range list {

		rvalue := reflect2.GetReflectValue(value)
		sqlOrder .WriteString( "(")
		for _,field := range fields {

			fieldValue := rvalue.FieldByName(field.Name)

			if fieldValue.IsValid() && !fieldValue.IsNil() {

				jsonName := field.Tag.Get("json")
				sqlOrder .WriteString(  jsonName + "=? AND ")

				args = append(args,fieldValue.Interface().(Type.TableType).Value())
			}
		}

		sqlStr := sqlOrder.String()

		if strings.HasSuffix(sqlStr," AND ") {
			sqlOrder.Truncate(len(sqlStr) - 4)
		} else if strings.HasSuffix(sqlStr," OR ") {
			sqlOrder.Truncate(len(sqlStr) - 3)
		}
		sqlOrder .WriteString( ") OR ")

	}
	sqlStr := sqlOrder.String()

	if strings.HasSuffix(sqlStr," AND ") {
		sqlStr = sqlStr[:len(sqlStr) - 4]
	} else if strings.HasSuffix(sqlStr," OR ") {
		sqlStr = sqlStr[:len(sqlStr) - 3]
	}
	return sqlStr,args
}
