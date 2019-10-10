package xsql4

import "fmt"

func (x *XSqlOrder) Begin() {
	var err error
	x.tx, err = x.xs.db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	x.txopen = 1
}

func (x *XSqlOrder) Commit() {
	var err error
	err = x.tx.Commit()
	x.txopen = 0
	if err != nil {
		x.RollBack()
		return
	}
}

func (x *XSqlOrder) RollBack() {
	var err error
	err = x.tx.Rollback()
	x.txopen = 0
	if err != nil {
		fmt.Println(err)
		return
	}
}