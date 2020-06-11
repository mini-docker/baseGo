package conf

import (
	"fmt"
	"sync"

	"github.com/mini-docker/baseGo/src/fecho/golog"
	"github.com/mini-docker/baseGo/src/fecho/xorm"
	"github.com/mini-docker/baseGo/src/fecho/xorm/core"
)

var (
	engine   *xorm.Engine
	engineRW sync.RWMutex
)

// InitMysql初始化mysql
func InitMysql(c MysqlConfig) error {
	engineTemp, err := initMysql(&c)
	if err != nil {
		return err
	}
	engineRW.Lock()
	engine = engineTemp
	engineRW.Unlock()
	return nil
}

func initMysql(c *MysqlConfig) (*xorm.Engine, error) {
	var err error
	data := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=%s&interpolateParams=true",
		//data := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=%s&interpolateParams=true&allowNativePasswords=true",
		c.Username,
		c.Password,
		c.Host,
		c.DbName,
		c.Timeout,
	)
	var engineTemp *xorm.Engine
	engineTemp, err = xorm.NewEngine("mysql", data)
	if err != nil {
		return nil, err
	}
	engineTemp.ShowSQL(c.ShowSql == 1)
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, c.TablePrefix)
	engineTemp.SetTableMapper(tbMapper)
	engineTemp.SetLogger(new(golog.LoggingXorm))
	err = engineTemp.Ping()
	if err != nil {
		return nil, err
	}
	golog.Info("", "", "mysql", "InitMysql", "mysql connect success", "host", c.Host)
	return engineTemp, nil
}

// GetXormSession 获取xorm session
func GetXormSession() *xorm.Session {
	engineRW.RLock()
	defer engineRW.RUnlock()
	return engine.NewSession()
}
