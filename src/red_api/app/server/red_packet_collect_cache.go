package server

import (
	"baseGo/src/red_api/conf"
	"fmt"
)

const (
	_red_packet_collect           = "red_packet_collect_%s"           // _red_packet_collect
	_red_packet_collect_last_time = "red_packet_collect_last_time_%s" // _red_packet_collect_last_time

)

func redPacketCollect(lineId string) string {
	return fmt.Sprintf(_red_packet_collect, lineId)
}
func redPacketCollectLastTime(lastTime string) string {
	return fmt.Sprintf(_red_packet_collect_last_time, lastTime)
}

func QueryRedPacketCollect(lineId string) ([]string, error) {
	conn := conf.GetRedis()
	redisClient := conn.Get()
	count, err := redisClient.Do("LLen", redPacketCollect(lineId))
	if err != nil {
		return nil, err
	}
	counts := count.(int)
	rpc, err := redisClient.Do("LRange", redPacketCollect(lineId), 0, counts-1)
	if err != nil {
		return nil, err
	}
	// 删除已采集的数据
	redisClient.Do("Del", redPacketCollect(lineId))
	rpcs := rpc.([]string)
	return rpcs, nil
}

func QueryRedPacketCollectLastTime(lineId string) (string, error) {
	conn := conf.GetRedis()
	redisClient := conn.Get()
	lastTime, err := redisClient.Do("Get", redPacketCollectLastTime(lineId))
	lastTimes := string(lastTime.([]byte))
	if err != nil {
		return lastTimes, err
	}

	return lastTimes, nil
}
func SetRedPacketCollectLastTime(lineId, lastTime string) error {
	conn := conf.GetRedis()
	redisClient := conn.Get()
	_, err := redisClient.Do("Set", redPacketCollectLastTime(lineId), lastTime, -1)
	if err != nil {
		return err
	}

	return nil
}
