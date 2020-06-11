// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package modules

import (
	"container/list"
	"os"
	"time"

	"github.com/mini-docker/baseGo/fecho/utility/atomicbool"
)

var modules *list.List
var addModule chan *Module
var GlobalShutdown chan bool
var loggingActive bool

type Module struct {
	Name  string
	Order uint8

	Start         chan bool
	Active        *atomicbool.AtomicBool
	startComplete chan bool

	Stop         chan bool
	Stopped      *atomicbool.AtomicBool
	stopComplete chan bool
}

func Register(name string, order uint8) *Module {
	newModule := &Module{
		name,
		order,
		make(chan bool),
		atomicbool.NewBool(true),
		make(chan bool),

		make(chan bool),
		atomicbool.NewBool(false),
		make(chan bool),
	}
	addModule <- newModule
	return newModule
}

func (module *Module) addToList() {
	if loggingActive {
		logger.Info("", "", "modules", "addToList", "Modules: starting %s", module.Name)
	}
	for e := modules.Back(); e != nil; e = e.Prev() {
		if module.Order > e.Value.(*Module).Order {
			modules.InsertAfter(module, e)
			return
		}
	}
	modules.PushFront(module)
}

func (module *Module) stop() {
	module.Active.UnSet()
	defer module.Stopped.Set()
	for {
		select {
		case module.Stop <- true:
		case <-module.stopComplete:
			return
		case <-time.After(1 * time.Second):
			if loggingActive {
				logger.Debug("", "", "modules", "stop", "Modules: waiting for %s to stop...", module.Name)
			}
		}
	}
}

func (module *Module) StopComplete() {
	if loggingActive {
		logger.Debug("", "", "modules", "StopComplete", "Modules: stopped %s", module.Name)
	}
	module.stopComplete <- true
}

func (module *Module) start() {
	module.Stopped.UnSet()
	defer module.Active.Set()
	for {
		select {
		case module.Start <- true:
		case <-module.startComplete:
			return
		}
	}
}

func (module *Module) StartComplete() {
	if loggingActive {
		logger.Info("", "", "modules", "StartComplete", "Modules: starting %s", module.Name)
	}
	module.startComplete <- true
}

func InitiateFullShutdown() {
	close(GlobalShutdown)
}

func fullStop() {
	for e := modules.Back(); e != nil; e = e.Prev() {
		if e.Value.(*Module).Active.IsSet() {
			e.Value.(*Module).stop()
		}
	}
}

func run() {
	select {
	case <-loggerRegistered:
		logger.Info("", "", "modules", "run", "Modules: starting ...")
		loggingActive = true
	case <-time.After(1 * time.Second):
	}

	for {
		select {
		case <-GlobalShutdown:
			if loggingActive {
				logger.Debug("", "", "modules", "run", "Modules: stopping")
			}
			fullStop()
			os.Exit(0)
		case m := <-addModule:
			m.addToList()
		}
	}
}

func init() {

	modules = list.New()
	addModule = make(chan *Module, 10)
	GlobalShutdown = make(chan bool)
	loggerRegistered = make(chan bool, 1)
	loggingActive = false

	go run()

}
