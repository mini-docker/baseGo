package xorm

import (
	"reflect"

	"github.com/mini-docker/baseGo/fecho/xorm/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
