package server

import (
	"baseGo/src/red_robot/conf"
	"fmt"
)

const (
	_red_packet_collect = "red_packet_collect_%s" // _red_packet_collect

)

func redPacketCollect(lineId string) string {
	return fmt.Sprintf(_red_packet_collect, lineId)
}

func InsertRedPacketCollect(lineId string, data string) (err error) {
	conn := conf.GetRedis()
	redisClient := conn.Get()
	_, err = redisClient.Do("LPush", redPacketCollect(lineId), data)
	if err != nil {
		return
	}

	return
}
