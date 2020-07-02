package conf

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model"
)

// SetConfig 设置配置参数
func Init(conf *Config) error {
	mConfig = conf
	//日志初始化
	err := InitLog(mConfig.Log)
	if err != nil {
		golog.Error("initd", "InitLog", "config log error %v", err, "log", mConfig.Log)
		return err
	}

	//数据库初始化
	err = InitMysql(mConfig.Mysql)
	if err != nil {
		golog.Error("initd", "InitConfig", "config mysql %v", err)
		return err
	}

	ConfigRW.Lock()
	IdgenInit(uint16(mConfig.App.UniqueId))
	ConfigRW.Unlock()

	//redis初始化
	err = InitRedis(&mConfig.Redis)
	if err != nil {
		golog.Error("initd", "InitConfig", "config redis %v %v", err, "Host", mConfig.Redis.Addrs)
		return err
	}

	model.IdgenInit(uint16(mConfig.App.UniqueId))
	// err = InitEtcd(
	// 	registry.Addrs(strings.Split(strings.Trim(GetRegistryConfig().Addr, ","), ",")...),
	// 	registry.Timeout(time.Duration(GetRegistryConfig().TTL*2)*time.Second),
	// )
	// if err != nil {
	// 	golog.Error("initd", "InitConfig", "config etcd error:%v", err)
	// 	return err
	// }
	return nil
}
