package xsql4

import (
	"github.com/NiuStar/xsql4/Type"
	"reflect"
	//"strings"
)

func ScanStructInterface(value interface{}) Type.DBOperation {

	scanStructInterface(reflect.ValueOf(value),reflect.ValueOf(value.(Type.DBOperation).TableName()))
	return value.(Type.DBOperation)
}

//实例化对象的时候对TableType进行变量名定义
func scanStructInterface(valueType reflect.Value,tableName reflect.Value) bool {

	switch valueType.Kind() {
	case reflect.Struct:

		if Type.IsTabelType(valueType.Type()) {
			return true
		} else {
			for i:=0;i<valueType.NumField();i++ {
				if scanStructInterface(valueType.Field(i),tableName) {

					jsonName := valueType.Type().Field(i).Name
					if len(valueType.Type().Field(i).Tag.Get("json")) > 0 {
						jsonName = valueType.Type().Field(i).Tag.Get("json")
					}
					valueType.Field(i).FieldByName("Names").SetString(jsonName)
					valueType.Field(i).Addr().MethodByName("SetParent").Call([]reflect.Value{valueType.Field(i).Addr()})
					valueType.Field(i).Addr().MethodByName("SetTableName").Call([]reflect.Value{tableName})
				}
			}
		}

	case reflect.Ptr:
		for ;reflect.Ptr == valueType.Kind(); {
			valueType = valueType.Elem()
		}
		return scanStructInterface(valueType,tableName)
	case reflect.Interface:
		for ;reflect.Interface == valueType.Kind(); {
			valueType = valueType.Elem()
		}
		return scanStructInterface(valueType,tableName)
	default:
		return false
	}
	return false
}