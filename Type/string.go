package Type

import (
	"strconv"
)

type String struct {
	Operation
	isValid bool
	value string
	Names string
	tableName string
}

func (s *String)Type() string {
	return "string"
}

func (s *String)Value() interface{} {
	return s.value
}

func (s *String)String() string{

	return s.value
}

func (s *String)Int() int64 {
	i,err := strconv.ParseInt(s.value,10,64)
	if err != nil {
		return -1
	}
	return i
}

func (s *String)Float() float64 {

	i,err := strconv.ParseFloat(s.value,64)
	if err != nil {
		return -1
	}
	return i
}

func (s *String)IsNil() bool {

	return !s.isValid
}

func (s *String)Name() string {
	return s.Names
}

func (s *String)SetValue(i interface{}) {

	s.value = i.(string)
	s.isValid = true
}


func (s *String)SetTableName(tableName string) {
	s.tableName = tableName
}

func (s *String)TableName() string {
	return s.tableName
}

