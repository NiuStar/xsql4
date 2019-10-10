package DBFactory

import (
	"github.com/NiuStar/xsql4/Type"
)

type JoinType int

const (
	AUTO_JOIN = iota
	INNER_JOIN
	LEFT_JOIN
	RIGHT_JOIN
)

type Join struct {
	table Type.DBOperation

	conditions []string
	conditionsValue []interface{}
	orConditions []string
	orConditionsValue []interface{}

	groupbys []Type.TableType
	orderbys []Type.TableType
	desc     bool
	limitStart,limitEnd int

	unionAllDB []*DBFactory
	unionDISTINCTDB []*DBFactory

	type_ JoinType
}


func NewJoinFactory(type_ JoinType,table Type.DBOperation) *Join {
	return &Join{limitStart:-1,limitEnd:-1,type_:type_,table:table}
}
//如果起始值都为-1，则不限制条数
func (join *Join)Limit(start,end int) *Join {
	join.limitStart = start
	join.limitEnd = end
	return join
}

//设置查询哪些表
func (join *Join)SetTable(table Type.DBOperation) *Join {
	join.table = table
	return join
}
//设置AND查询条件或者插入条件
func (join *Join)SetConditions(conditions... *Type.Condition) *Join {

	join.conditions = []string{}
	join.conditionsValue = []interface{}{}

	for _,condition := range conditions {
		str,args := condition.String()
		join.conditions = append(join.conditions,str)
		join.conditionsValue = append(join.conditionsValue,args...)
	}
	return join
}
//添加AND查询条件或者插入条件
func (join *Join)AddCondition(conditions... *Type.Condition) *Join {

	for _,condition := range conditions {
		str,args := condition.String()
		join.conditions = append(join.conditions,str)
		join.conditionsValue = append(join.conditionsValue,args...)
	}
	return join
}

//设置OR查询条件或者插入条件
func (join *Join)SetORConditions(conditions... *Type.Condition) *Join {
	join.orConditions = []string{}
	join.orConditionsValue = []interface{}{}

	for _,condition := range conditions {
		str,args := condition.String()
		join.orConditions = append(join.orConditions,str)
		join.orConditionsValue = append(join.orConditionsValue,args...)
	}
	//join.orConditions = conditions
	return join
}
//添加OR查询条件或者插入条件
func (join *Join)AddORCondition(conditions... *Type.Condition) *Join {
	for _,condition := range conditions {
		str,args := condition.String()
		join.orConditions = append(join.orConditions,str)
		join.orConditionsValue = append(join.orConditionsValue,args...)
	}
	//join.orConditions = append(join.orConditions,conditions...)
	return join
}

func (join *Join)DESC() *Join {
	join.desc = true
	return join
}

//设置Order By 条件
func (join *Join)SetOrderBy(orderbys... Type.TableType) *Join {
	join.orderbys = orderbys
	return join
}

//设置group By 条件
func (join *Join)SetGroupBy(groupbys... Type.TableType) *Join {
	join.groupbys = groupbys
	return join
}