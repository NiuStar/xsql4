package DBFactory


//UNION的时候需要用到的
func (src *DBFactory)UNIONALLDB(value... *DBFactory) *DBFactory {
	src.unionAllDB = value
	return src
}

//UNION DISTINCT的时候需要用到的
func (src *DBFactory)UNIONDISTINCTDB(value... *DBFactory) *DBFactory {
	src.unionDISTINCTDB = value
	return src
}

//UNION的时候需要用到的
func (src *Join)UNIONALLDB(value... *DBFactory) *Join {
	src.unionAllDB = value
	return src
}

//UNION DISTINCT的时候需要用到的
func (src *Join)UNIONDISTINCTDB(value... *DBFactory) *Join {
	src.unionDISTINCTDB = value
	return src
}