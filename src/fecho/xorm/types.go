package xorm

import (
	"reflect"

	"fecho/xorm/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
