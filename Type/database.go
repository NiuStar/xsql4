package Type

import (
	"github.com/NiuStar/reflect"
)

type ISQLDB interface {
	GetDBName() string
	String() (sqlString string,args []interface{})
}


//where条件语句的条件
type Condition struct {
	conditions []*Condition
	operation *Operation
	type_     int   //0是all 1 是 equalOperation 2 是 equalValues 3 是 greaterValue 4 是 greaterEqualValue 5 是 lessValue 6 是 lessEqualValue 7 是 notEqualOperation 8 是 notEqualValues 9 是 INDB 10 是 NOTINDB 11 like
					//12是UNION ALL 13 是 UNION DISTINCT

	conditionString string
	argvs []interface{}
}

func NewSelfDefineCondition(con string) *Condition {
	return &Condition{conditionString:con}
}

func (c *Condition)String() (string,[]interface{}) {

	if len(c.conditionString) <= 0 && c.operation != nil{
		c.conditionString,c.argvs = c.operation.Where(c.type_)
	}

	return c.conditionString,c.argvs
}

func (c *Condition)AddCondition(condition *Condition) *Condition {

	if len(c.conditionString) <= 0 {
		c.conditionString,c.argvs = c.operation.Where(c.type_)
	}

	str,argv := condition.String()
	c.conditionString += " AND " + str
	c.argvs = append(c.argvs,argv...)
	c.conditions = append(c.conditions,condition)
	return c
}


type Operation struct {
	parent TableType
	equalOperation TableType
	equalValues []interface{}

	notINDB ISQLDB
	INDB ISQLDB

	notEqualOperation TableType
	notEqualValues []interface{}

	greaterValue interface{}//大于计算的集合
	greaterEqualValue interface{}//大于等于计算的集合

	lessValue interface{}//小于计算的集合
	lessEqualValue interface{}//小于等于计算的集合

	likeValue interface{}

	unionALL ISQLDB
	unionDISTINCT ISQLDB
}
func (src *Operation)SetParent(parent TableType) {
	//fmt.Println("SetParent:",parent)
	src.parent = parent
}
//两种表之间通过关联关系对比，如：a.id=b.aid
func (src *Operation)NotEqualFold(dstFiled TableType) *Condition {
	src.notEqualOperation = dstFiled
	return &Condition{operation:src,type_:7}
}
//两种表之间通过关联关系对比，如：a.id=b.aid
func (src *Operation)EqualFold(dstFiled TableType) *Condition {
	src.equalOperation = dstFiled
	return &Condition{operation:src,type_:1}
}
//where id not in (1,2,3,4,5)
func (src *Operation)NotEqual(value... interface{}) *Condition {
	src.notEqualValues = src.checkValues(value...)
	return &Condition{operation:src,type_:8}
}

//where id in (1,2,3,4,5)
func (src *Operation)Equal(value... interface{}) *Condition {
	src.equalValues = src.checkValues(value...)
	return &Condition{operation:src,type_:2}
}

//大于
func (src *Operation)Greater(value interface{}) *Condition {
	src.greaterValue = src.checkValues(value)[0]
	return &Condition{operation:src,type_:3}
}

//大于等于
func (src *Operation)GreaterEqual(value interface{}) *Condition {
	src.greaterEqualValue = src.checkValues(value)[0]
	return &Condition{operation:src,type_:4}
}

//小于
func (src *Operation)Less(value interface{}) *Condition {
	src.lessValue = src.checkValues(value)[0]
	return &Condition{operation:src,type_:5}
}

//小于等于
func (src *Operation)LessEqual(value interface{}) *Condition {
	src.lessEqualValue = src.checkValues(value)[0]
	return &Condition{operation:src,type_:6}
}

//嵌套临时表的时候需要用到的，在一个select里面嵌套一个select
func (src *Operation)NotInDB(value ISQLDB) *Condition {
	src.notINDB = value
	return &Condition{operation:src,type_:9}
}

//嵌套临时表的时候需要用到的，在一个select里面嵌套一个select
func (src *Operation)InDB(value ISQLDB) *Condition {
	src.INDB = value
	return &Condition{operation:src,type_:10}
}

//小于
func (src *Operation)Like(value string) *Condition {
	src.likeValue = value
	return &Condition{operation:src,type_:11}
}
/*
//UNION的时候需要用到的
func (src *Operation)UNIONALLDB(value ISQLDB) *Condition {
	src.unionALL = value
	return &Condition{operation:src,type_:12}
}

//UNION DISTINCT的时候需要用到的
func (src *Operation)UNIONDISTINCTDB(value ISQLDB) *Condition {
	src.unionDISTINCT = value
	return &Condition{operation:src,type_:13}
}
*/
func (src *Operation)checkValues(value... interface{}) []interface{} {
	var list []interface{}
	for _,v := range value {
		if IsTabelType(reflect.GetReflectType(v)) {
			list = append(list, v.(TableType).Value())
		} else {
			list = append(list, v)
		}
	}
	return list
}
/*
func (src *Operation)Where() (string,[]interface{}) {
	if src.equalOperation != nil {
		return src.parent.TableName() + "." + src.parent.Name() + "=" + src.equalOperation.TableName() + "." + src.equalOperation.Name(),nil
	} else if len(src.equalValues) > 0 {
		str := src.parent.TableName() + "." + src.parent.Name() + " IN ("
		for index,_ := range src.equalValues {
			str += "?"
			if index != len(src.equalValues) - 1 {
				str += ","
			}
		}
		str += ")"
		return str,src.equalValues
	} else {
		return "",nil
	}
}*/

func (src *Operation)Where(type_ int) (string,[]interface{}) {

	strSql := ""
	var args []interface{}

	switch type_ {

	case 1:
		{
			if src.equalOperation != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + "=" + src.equalOperation.TableName() + "." + src.equalOperation.Name()
			}
		}
	case 2:
		{
			if len(src.equalValues) > 0 {
				str := src.parent.TableName() + "." + src.parent.Name() + " IN ("
				for index,_ := range src.equalValues {
					str += "?"
					if index != len(src.equalValues) - 1 {
						str += ","
					}
				}
				str += ")"
				strSql = str
				args = append(args,src.equalValues...)
			}
		}
	case 3:
		{
			if src.greaterValue != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + " > ? "
				args = append(args,src.greaterValue)
			}
		}
	case 4:
		{
			if src.greaterEqualValue != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + " >= ? "
				args = append(args,src.greaterEqualValue)
			}
		}
	case 5:
		{
			if src.lessValue != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + " < ? "
				args = append(args,src.lessValue)
			}
		}
	case 6:
		{
			if src.lessEqualValue != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + " <= ? "
				args = append(args,src.lessEqualValue)
			}
		}
	case 7:
		{
			if src.notEqualOperation != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + "!=" + src.notEqualOperation.TableName() + "." + src.notEqualOperation.Name()
			}
		}
	case 8:
		{
			if len(src.notEqualValues) > 0 {
				str := src.parent.TableName() + "." + src.parent.Name() + " NOT IN ("
				for index,_ := range src.notEqualValues {
					str += "?"
					if index != len(src.notEqualValues) - 1 {
						str += ","
					}
				}
				str += ")"
				strSql = str
				args = append(args,src.notEqualValues...)
			}
		}
	case 9:
		{
			if src.notINDB != nil {
				str , argvs := src.notINDB.String()
				args = append(args,argvs...)

				strSql = src.parent.TableName() + "." + src.parent.Name() + " NOT IN(" + str + ")"
			}
		}
	case 10:
		{
			if src.INDB != nil {
				str , argvs := src.INDB.String()
				args = append(args,argvs...)
				strSql = src.parent.TableName() + "." + src.parent.Name() + " IN(" + str + ")"
			}
		}
	case 11:
		{
			if src.likeValue != nil {
				args = append(args,src.likeValue)
				strSql = src.parent.TableName() + "." + src.parent.Name() + " Like ? "
			}
		}
	case 12:
		{
			if src.unionALL != nil {
				str , argvs := src.unionALL.String()
				args = append(args,argvs...)
				strSql = src.parent.TableName() + "." + src.parent.Name() + " UNION ALL (" + str + ")"
			}
		}
	case 13:
		{
			if src.unionDISTINCT != nil {
				str , argvs := src.unionDISTINCT.String()
				args = append(args,argvs...)
				strSql = src.parent.TableName() + "." + src.parent.Name() + " UNION DISTINCT (" + str + ")"
			}
		}
	default:
		{
			if src.equalOperation != nil {
				strSql = src.parent.TableName() + "." + src.parent.Name() + "=" + src.equalOperation.TableName() + "." + src.equalOperation.Name()
			} else if len(src.equalValues) > 0 {
				str := src.parent.TableName() + "." + src.parent.Name() + " IN ("
				for index,_ := range src.equalValues {
					str += "?"
					if index != len(src.equalValues) - 1 {
						str += ","
					}
				}
				str += ")"
				strSql = str
				args = append(args,src.equalValues...)
			}

			start := false
			if len(strSql) > 0 {
				start = true
			}

			if start && src.greaterValue != nil {
				strSql += " AND " + src.parent.TableName() + "." + src.parent.Name() + " > ? "
				args = append(args,src.greaterValue)
			}


			if start && src.greaterEqualValue != nil {
				strSql += " AND " + src.parent.TableName() + "." + src.parent.Name() + " >= ? "
				args = append(args,src.greaterEqualValue)
			}

			if start && src.lessValue != nil {
				strSql += " AND " + src.parent.TableName() + "." + src.parent.Name() + " < ? "
				args = append(args,src.lessValue)
			}

			if start && src.lessEqualValue != nil {
				strSql += " AND " + src.parent.TableName() + "." + src.parent.Name() + " <= ? "
				args = append(args,src.lessEqualValue)
			}
		}
	}

	//fmt.Println("strSql:",strSql)
	return strSql,args
}

