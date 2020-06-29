package conf

import (
	"sync"
)

var (
	sessionConfig   SessionConfig //session配置信息
	sessionConfigRW sync.RWMutex  //session配置信息因为会有变动并且会有外部访问,所以加锁
)

var memberSessionMap sync.Map

func GetSessionKey(siteId, siteIndexId string) string {
	v, ok := memberSessionMap.Load(siteId + siteIndexId)
	if ok {
		return v.(string)
	}

	key := "pk_" + siteId + "_" + siteIndexId + "_member_session_list_key"
	memberSessionMap.Store(siteId+siteIndexId, key)
	return key
}

// SetSessionConfig 设置sessionConfig
func SetSessionConfig(sessionConfigTemp SessionConfig) {
	sessionConfigRW.Lock()
	sessionConfig = sessionConfigTemp
	sessionConfigRW.Unlock()
}

// GetSessionConfig 获取session配置信息
func GetSessionConfig() SessionConfig {
	sessionConfigRW.RLock()
	defer sessionConfigRW.RUnlock()
	return sessionConfig
}
