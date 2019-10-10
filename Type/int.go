package Type

import (
	"strconv"
)

type Int struct {
	Operation
	isValid bool
	value int64
	Names string
	tableName string
}

func (s *Int)Type() string {
	return "int64"
}

func (s *Int)Value() interface{} {
	return s.value
}

func (s *Int)Name() string {
	return s.Names
}


func (s *Int)SetValue(i interface{}) {

	switch i.(type) {
	case float64:
		s.value = int64(i.(float64))
	case int:
		s.value = int64(i.(int))
	case int32:
		s.value = int64(i.(int32))
	default:
		s.value = i.(int64)
	}
	//s.value = i.(int64)
	s.isValid = true
}


func (s *Int)String() string{
	return strconv.FormatInt(s.value,10)
}

func (s *Int)Int() int64 {

	return s.value
}

func (s *Int)Float() float64 {
	return float64(s.value)
}

func (s *Int)IsNil() bool {
	return !s.isValid
}


func (s *Int)SetTableName(tableName string) {
	s.tableName = tableName
}

func (s *Int)TableName() string {
	return s.tableName
}
