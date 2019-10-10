package DBFactory

import (
	"fmt"
	"strconv"
)

func (db *DBFactory)String() (string,[]interface{}) {

	var args []interface{}

	str := "select "

	for index,field := range db.fields {
		fmt.Println("field.TableName():",field.TableName())
		if len(field.TableName()) > 0 {
			str += field.TableName() + "."
		}
		str += field.Name() + " as "

		if len(field.TableName()) > 0 {
			str += field.TableName() + "_"
		}
		str += field.Name()
		if index != len(db.fields) - 1 {
			str += ","
		}
	}
	if len(db.fields) > 0 && db.hadCount {
		str += ","
	}

	if db.hadCount {
		str += " count(*) as count"
	}

	if len(str) == len("select ") {

		for index,table := range db.tables {
			str += db.getTableFields(table)
			if index != len(db.tables) - 1 {
				str += ","
			}
		}
		//str += "*"
	}
	str += " from "
	for index,table := range db.tables {
		str += table.TableName()
		if index != len(db.tables) - 1 {
			str += ","
		}
	}

	str1,argvs := db.where()

	str += str1
	args = append(args,argvs...)


	return str,args
}

func (db *DBFactory)where() (string,[]interface{}) {
	var args []interface{}

	str := ""

	if len(db.conditions) > 0 {
		str += " where "
		for index,condition := range db.conditions {
			str1 := condition
			str += str1
			if index != len(db.conditions) - 1 {
				str += " AND "
			}
		}
		args = append(args,db.conditionsValue...)
	}

	if len(db.orConditions) > 0 {
		if len(db.conditions) <= 0 {
			str += " where "
		}
		for index,condition := range db.orConditions {

			str1 := condition
			str += str1
			if index != len(db.orConditions) - 1 {
				str += " OR "
			}
			//args = append(args,argv...)
		}
		args = append(args,db.orConditionsValue...)
	}

	if len(db.groupbys) > 0 {
		str += " Group By "
		for index,groupby := range db.groupbys {

			if len(groupby.TableName()) > 0 {
				str += groupby.TableName() + "."
			}
			str += groupby.Name()
			if index != len(db.groupbys) - 1 {
				str += " , "
			}
		}
	}

	if len(db.orderbys) > 0 {
		str += " Order By "
		for index,orderby := range db.orderbys {
			if len(orderby.TableName()) > 0 {
				str += orderby.TableName() + "_"
			}
			str += orderby.Name()
			if index != len(db.orderbys) - 1 {
				str += " , "
			}
		}

		if db.desc {
			str += " DESC"
		}
	}

	if db.limitStart != db.limitEnd && db.limitEnd != -1 {
		str += " limit " + strconv.FormatInt(int64(db.limitStart),10) + "," + strconv.FormatInt(int64(db.limitEnd),10)
	}

	for _,value := range db.unionAllDB {
		str1,argvs := value.String()
		str += " UNION ALL (" + str1 + ")"
		args = append(args,argvs...)
	}
	for _,value := range db.unionDISTINCTDB {
		str1,argvs := value.String()
		str += " UNION DISTINCT (" + str1 + ")"
		args = append(args,argvs...)
	}

	for _,value := range db.joins {
		str1,argvs := value.String()
		str += str1
		args = append(args,argvs...)
	}
	return str,args
}

func (join *Join)String() (string,[]interface{}) {

	var args []interface{}

	str := ""

	switch join.type_ {
	case AUTO_JOIN:
		{
			str += " JOIN "
		}
	case INNER_JOIN:
		{
			str += " INNER JOIN "
		}
	case RIGHT_JOIN:
		{
			str += " RIGHT JOIN "
		}
	case LEFT_JOIN:
		{
			str += " LEFT JOIN "
		}
	}

	if join.table != nil {
		str += join.table.TableName() + " ON "
	} else {
		return "",[]interface{}{}
	}

	if len(join.conditions) > 0 {
		for index,condition := range join.conditions {
			str1 := condition
			str += str1
			if index != len(join.conditions) - 1 {
				str += " AND "
			}
		}
		args = append(args,join.conditionsValue...)
	}

	if len(join.orConditions) > 0 {

		for index,condition := range join.orConditions {

			str1 := condition
			str += str1
			if index != len(join.orConditions) - 1 {
				str += " OR "
			}
			//args = append(args,argv...)
		}
		args = append(args,join.orConditionsValue...)
	}

	if len(join.groupbys) > 0 {
		str += " Group By "
		for index,groupby := range join.groupbys {

			if len(groupby.TableName()) > 0 {
				str += groupby.TableName() + "_"
			}
			str += groupby.Name()
			if index != len(join.groupbys) - 1 {
				str += " , "
			}
		}
	}

	if len(join.orderbys) > 0 {
		str += " Order By "
		for index,orderby := range join.orderbys {
			if len(orderby.TableName()) > 0 {
				str += orderby.TableName() + "_"
			}
			str += orderby.Name()
			if index != len(join.orderbys) - 1 {
				str += " , "
			}
		}

		if join.desc {
			str += " DESC"
		}
	}

	if join.limitStart != join.limitEnd && join.limitEnd != -1 {
		str += " limit " + strconv.FormatInt(int64(join.limitStart),10) + "," + strconv.FormatInt(int64(join.limitEnd),10)
	}

	for _,value := range join.unionAllDB {
		str1,argvs := value.String()
		str += " UNION ALL (" + str1 + ")"
		args = append(args,argvs...)
	}
	for _,value := range join.unionDISTINCTDB {
		str1,argvs := value.String()
		str += " UNION DISTINCT (" + str1 + ")"
		args = append(args,argvs...)
	}

	return str + " ",args
}
