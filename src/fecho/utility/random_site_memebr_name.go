package utility

import (
	"fmt"
	"sync"
	"time"
)

var randomMemberNameBySiteIndexIdMutex sync.Mutex
var lastUnixTime int64
var unixTimeMakeupCount int64

//
// RandomMemberNameBySiteIndexId 用尽量少的字符，为每个子站生成一个不重复的随机用户名
//
func RandomMemberNameBySiteIndexId(siteIndexId string) string {
	randomMemberNameBySiteIndexIdMutex.Lock()
	defer randomMemberNameBySiteIndexIdMutex.Unlock()
	unixTime := time.Now().Unix()
	if lastUnixTime != unixTime {
		lastUnixTime = unixTime
		unixTimeMakeupCount = 0
		return siteIndexId + fmt.Sprintf("%x", unixTime)
	}

	result := siteIndexId + fmt.Sprintf("%x%x", unixTime, unixTimeMakeupCount)
	unixTimeMakeupCount++

	return result
}
