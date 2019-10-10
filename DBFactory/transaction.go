package DBFactory


func (x *DBFactory) Begin() {
	if x.xsql != nil && x.transactionBegin {
		x.xsql.RollBack()
	}

	//x.xsql = xsql3.CreateInstance(xsql3.GetServerDB())
	x.xsql.Begin()
	x.transactionBegin = true
}

func (x *DBFactory) Commit() {
	//xr := xsql3.CreateInstance(xsql3.GetServerDB())
	if x.xsql != nil {
		x.xsql.Commit()
		//x.xsql = nil
	}
	x.transactionBegin = false

}

func (x *DBFactory) RollBack() {
	if x.xsql != nil {
		x.xsql.RollBack()
		//x.xsql = nil
	}
	x.transactionBegin = false

}