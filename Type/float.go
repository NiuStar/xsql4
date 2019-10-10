package Type

import (
	"strconv"
)

type Float struct {
	Operation
	isValid bool
	value float64
	Names string
	tableName string
}

func (s *Float)Type() string {
	return "float64"
}

func (s *Float)Value() interface{} {
	return s.value
}

func (s *Float)Name() string {
	return s.Names
}

func (s *Float)SetValue(i interface{}) {
	s.value = i.(float64)
	s.isValid = true
}

func (s *Float)String() string{
	return strconv.FormatFloat(s.value,'b', 6, 64)
}

func (s *Float)SetTableName(tableName string) {
	s.tableName = tableName
}

func (s *Float)TableName() string {
	return s.tableName
}

func (s *Float)Int() int64 {

	return int64(s.value)
}

func (s *Float)Float() float64 {

	return s.value
}

func (s *Float)IsNil() bool {
	return !s.isValid
}