package xsql4

import (

	"fmt"
	"strings"
	//"strconv"
)

const SQLSPACE  =  " "
func copy2Struct() {

}
//表结构的映射结构体
type Table struct {
	tableName string
	columns map[string]*Column
	comment string
	engine string
	collation string
	pkey []string
	ukey map[string][]string
	ikey map[string][]string
}
//字段属性的映射结构体
type Column struct {
	Index int
	Name string
	Type string
	Comment string
	Extra string
	PRIMARY bool
	Default string
	IsNull bool
	Mark string
	oldName string    //在数据库中，也就是以前的字段名
}
//判断数据库中的字段名称及属性内容是否一致
func (c *Column)Equal(c2 *Column) bool {
/*
	if c.Name != c2.Name {
		fmt.Println("name 不一致：",c.Name,c2.Name)
	}
	if c.Type != c2.Type {
		fmt.Println("Type 不一致：",c.Type,c2.Type)
	}
	if c.Comment != c2.Comment {
		fmt.Println("Comment 不一致：",c.Comment,c2.Comment)
	}
	if c.Extra != c2.Extra {
		fmt.Println("Extra 不一致：",c.Extra,c2.Extra)
	}
	if c.PRIMARY != c2.PRIMARY {
		fmt.Println("PRIMARY 不一致：",c.PRIMARY,c2.PRIMARY)
	}
	if c.Default != c2.Default {
		fmt.Println("Default 不一致：",c.Default,c2.Default)
	}
	if c.IsNull != c2.IsNull {
		fmt.Println("IsNull 不一致：",c.IsNull,c2.IsNull)
	}
*/
	return c.Name == c2.Name && c.Type == c2.Type && c.Comment == c2.Comment && c.Extra == c2.Extra && c.PRIMARY == c2.PRIMARY && c.Default == c2.Default && c.IsNull == c2.IsNull
}
//备注一致的，认为是同一个字段改名称
func (c *Column)EqualByComment(c2 *Column) bool {
	return c.Comment == c2.Comment
}

//如果数据库及现有结构体中字段存在，则进行对字段名及相关属性进行更新操作
func (c *Column)updateColumn(tableName string) {

	sql := "alter table " + tableName + " change " + c.oldName + SQLSPACE + c.Name + SQLSPACE + c.Type

	if !c.IsNull {
		sql += " NOT NULL "
	} else {
		sql += SQLSPACE
	}

	if 0 != len(c.Mark) {
		if -1 == strings.Index(strings.ToUpper(c.Mark), "PRIMARY KEY") {
			sql += c.Mark + " "
		} else {
			sql += strings.Replace(strings.ToUpper(c.Mark), "PRIMARY KEY","",-1) + " "
		}


	}
	if 0 == len(c.Mark) ||
		-1 == strings.Index(strings.ToUpper(c.Mark), "AUTO_INCREMENT") {
		if 0 != len(c.Default) {
			sql += "DEFAULT " + c.Default + SQLSPACE
		}
	}
	sql += "COMMENT '" + c.Comment + "'"
	order := CreateInstance(xs)
	order.Qurey(sql)
	fmt.Println("alert result:",order.ExecuteNoResult())
}
//如果字段在现结构体中存在、数据库中不存在，则进行对字段名进行添加操作
func (c *Column)addColumn(tableName,afterName string) {
	sql := "alter table " + tableName + " add " + c.Name + SQLSPACE + c.Type

	if !c.IsNull {
		sql += " NOT NULL "
	} else {
		sql += SQLSPACE
	}

	if 0 != len(c.Mark) {
		sql += c.Mark + " "
	}
	if 0 == len(c.Mark) ||
		-1 == strings.Index(strings.ToUpper(c.Mark), "AUTO_INCREMENT") {
		if 0 != len(c.Default) {
			sql += "DEFAULT " + c.Default + SQLSPACE
		}
	}

	if !strings.EqualFold(c.Name , afterName) {
		sql += "COMMENT '" + c.Comment + "'" + " after " + afterName
	} else {
		sql += "COMMENT '" + c.Comment + "'"
	}

	order := CreateInstance(xs)
	order.Qurey(sql)
	order.ExecuteNoResult()

}
//如果字段在现结构体中不存在、数据库中存在，则进行对字段名进行删除操作
func (c *Column)delColumn(tableName string) {
	sql := "alter table " + tableName + " DROP " + c.Name + SQLSPACE
	order := CreateInstance(xs)
	order.Qurey(sql)
	order.ExecuteNoResult()

}
//node：现在的问题还剩下一个，字段名称和备注不能同时修改，否则会认为是新增的，会把以前的数据清除
func verificationUpdate(table *Table) {
	table2 := getTableByName(table.tableName)

	//fmt.Println("table2:",*table2.columns["id"])
	//fmt.Println("table:",*table.columns["id"])

	indexList := make(map[int]*Column)  //需要更新的字段属性
	addList := make(map[string]*Column)  //需要添加的字段属性

	for name,col := range table.columns {
		if table2.columns[name] != nil {
			if !col.Equal(table2.columns[name]) {
				fmt.Println("字段："+name+"数据表中与结构体中存在不一致，需检查!!!")
			}
			delete(table2.columns,name)
		} else {
			addList[name] = col
		}
		indexList[col.Index] = col
	}




/*把原来的自动更新功能去掉
	for _ , col := range updateList {
		col.updateColumn(table.tableName)
	}
	for _ , col := range addList {
		afterName := "id" //:= strconv.FormatInt(index - 1,64,10)
		if indexList[col.Index - 1] != nil {
			afterName = indexList[col.Index - 1].Name
		}
		col.addColumn(table.tableName,afterName)
	}

	for _ , col := range table2.columns {

		if strings.EqualFold(table.tableName,"id") {
			col.delColumn(table.tableName)
		}
	}*/
	/*增加Model与数据库不一致的提示功能
	*/


	for name , col := range addList {

		var ok string
		fmt.Print("数据表中缺少："+name+"字段，是否插入(插入该字段请输入y，否则输入任意字符，回车确认):")
		fmt.Scanf("%s",&ok)
		//fmt.Println("输入了",b)

		ok = strings.ToLower(ok)
		if ok == "y" {
			afterName := "id" //:= strconv.FormatInt(index - 1,64,10)
			if indexList[col.Index - 1] != nil {
				afterName = indexList[col.Index - 1].Name
			}
			col.addColumn(table.tableName,afterName)
		}


	}
	for name , _ := range table2.columns {

		fmt.Println(`结构体中缺少：`+name+`字段，请详细检查`)
	}
}
//去数据库中把相应的表信息和表结构映射到代码里来
func getTableByName(name string) *Table {

	table := &Table{tableName:name,columns:make(map[string]*Column),ukey:make(map[string][]string),ikey:make(map[string][]string)}

	order := CreateInstance(xs)

	{//从数据库里面查找表相关的字符集及相关备注
		order.Qurey(`select engine,table_collation,table_comment from information_schema.tables where table_name = '` + name + `'`)
		list := order.Execute()

		if len(list) <= 0 {
			return nil
		}
		{
			fmt.Println("list[0]:",list[0])
			table.engine = string(list[0]["engine"].(string))
			table.collation = string(list[0]["table_collation"].(string))
			table.comment = string(list[0]["table_comment"].(string))
		}
	}

	{//从数据库
		order.Qurey(`show index from ` + name)
		list := order.Execute()

		//fmt.Println("list:",list)
		if len(list) <= 0 {
			return nil
		}

		for _,l := range list {
			key_name := string(l["Key_name"].(string))
			column_name := string(l["Column_name"].(string))
			index_comment := string(l["Index_comment"].(string))

			//if !strings.EqualFold(column_name,"id") || !strings.EqualFold(key_name,"PRIMARY")
			{

				if strings.EqualFold(key_name,"PRIMARY") {
					table.pkey = append(table.pkey,column_name)
				} else if strings.EqualFold(index_comment,"UNIQUE") {
					table.ukey[key_name] = append(table.ukey[key_name],column_name)
				} else if strings.EqualFold(index_comment,"KEY") {
					table.ikey[key_name] = append(table.ikey[key_name],column_name)
				}
			}
		}
	}

	{
		order.Qurey(`select column_name,column_type,extra,column_comment,column_default,is_nullable,column_key from information_schema.columns where table_name = '` + name + "'")
		list := order.Execute()

		if len(list) <= 0 {
			return nil
		}

		for index,l := range list {
			column_name := string(l["column_name"].(string))

			//if !strings.EqualFold(column_name,"id")
			{
				column_type := string(l["column_type"].(string))
				column_comment := string(l["column_comment"].(string))
				extra := string(l["extra"].(string))

				column_default := ""
				if l["column_default"] != nil {
					column_default = l["column_default"].(string)
				}


				is_nullable := strings.ToLower(string(l["is_nullable"].(string))) == "yes"
				PRIMARY := strings.ToUpper(string(l["column_key"].(string))) == "PRI"

				table.columns[column_name] = &Column{Index:index - 1,Name: column_name, Type: column_type, Comment: column_comment, Extra: extra,
				Default: column_default,IsNull:is_nullable,PRIMARY:PRIMARY}
			}
			if strings.EqualFold(column_name,"id") {
				fmt.Println("table.columns[COLUMN_NAME]:",table.columns[column_name])
			}
		}
	}
	return table
}