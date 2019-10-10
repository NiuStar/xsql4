package Type

import (
	"strings"
	"reflect"
)
//数据库操作所需
type DBOperation interface {
	TableName() string
	NewInterface() DBOperation
	//Count() TableType
}
//json转换的时候需要用到
type IHandler interface {
}
//数据库操作需要，不能用基本类型
type TableType interface {
	TableName() string
	Type() string
	Value() interface{}
	String() string
	Int() int64
	Float() float64
	IsNil() bool
	SetValue(i interface{})
	Name() string
}

func IsTabelType(type_ reflect.Type) bool {
	return strings.HasPrefix(type_.String(), "*Type.") || strings.HasPrefix(type_.String(), "Type.")
}

func IsMultipleTabelType(type_ reflect.Type) bool {
	return strings.HasPrefix(type_.String(), "[]*Type.") || strings.HasPrefix(type_.String(), "[]Type.")
}

//string 1 int 2 float 3
func GetTabelType(type_ reflect.Type) int {
	switch type_.String() {
	case "[]Type.String":{
		return 1
	}
	case  "[]Type.Int":{
		return 2
	}
	case  "[]Type.Float":{
		return 3
	}
	}
	return 0
//	return strings.HasPrefix(type_.String(), "[]*Type.") || strings.HasPrefix(type_.String(), "[]Type.")
}

