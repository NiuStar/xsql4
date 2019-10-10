package xsql4

import (
	"reflect"
	"fmt"
	"github.com/NiuStar/xsql4/Type"
	//reflect2 "github.com/NiuStar/reflect"
	//"strings"
)


type QueryTable struct {
	links []Type.DBOperation
}

func (link *QueryTable)Link(list... Type.DBOperation) {
	link.links = list
}

func (link *QueryTable)Equal(srcField reflect.StructField,dstField reflect.StructField) {
	fmt.Println("srcField.PkgPath:",srcField.PkgPath)
}

