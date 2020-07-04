package dao

import (
	"baseGo/src/fecho/echo"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	log "baseGo/src/fecho/golog"
	"baseGo/src/imserver/internal/logic/model"

	"github.com/zhenjl/cityhash"
)

const (
	_prefixMidServer    = "mid_%d"             // mid -> key:server
	_prefixKeyServer    = "%s"                 // key -> server
	_prefixServerOnline = "ol_%s"              // server -> online
	_jPush              = "jpush_%s"           // jpush
	_groupMemberAll     = "groupMemberAll_%s"  // _groupMemberAll
	_groupMemberRoom    = "groupMemberRoom_%s" // _groupMemberRoom
)

func keyMidServer(mid int64) string {
	return fmt.Sprintf(_prefixMidServer, mid)
}

func keyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func keyServerOnline(key string) string {
	return fmt.Sprintf(_prefixServerOnline, key)
}
func jPushKey(key string) string {
	return fmt.Sprintf(_jPush, key)
}
func groupMemberAll(lineId string) string {
	return fmt.Sprintf(_groupMemberAll, lineId)
}
func groupMemberRoom(lineId, RoomId string) string {
	return fmt.Sprintf(_groupMemberRoom, lineId+"_"+RoomId)
}

// AddMapping add a mapping.
// Mapping:
//	mid -> key_server
//	key -> server
func (d *Dao) AddMapping(c context.Context, mid int64, key, server string) (err error) {
	conn := model.GetRedis().Get()
	var n = 2
	if mid > 0 {
		_, err = conn.Do("HSet", keyMidServer(mid), key, server)
		if err != nil {
			log.Error("Dao", "AddMapping", "conn.Send(HSET %d,%s,%s) error(%v)", nil, mid, server, key, err)
			return
		}
		n += 2
	}
	// 之前的set没有给过期时间 所以给0
	_, err = conn.Do("Set", keyKeyServer(key), server, 0)
	if err != nil {
		log.Error("Dao", "AddMapping", "conn.Send(SET %d,%s,%s) error(%v)", nil, mid, server, key, err)
		return
	}
	return
}

// ExpireMapping expire a mapping.
func (d *Dao) ExpireMapping(c context.Context, mid int64, key string) (bool, error) {
	conn := model.GetRedis().Get()
	var (
		//n    = 1
		flag interface{}
		err  error
	)

	flag, err = conn.Do("Expire", keyKeyServer(key), time.Duration(d.redisExpire)*time.Second)

	if err != nil {
		log.Error("Dao", "ExpireMapping", "conn.Send(EXPIRE %d,%s) error(%v)", nil, mid, key, err)
		return flag.(bool), err
	}
	return flag.(bool), err
}

// DelMapping del a mapping.
func (d *Dao) DelMapping(c context.Context, mid int64, key, server string) (bool, error) {
	conn := model.GetRedis().Get()
	var result int64
	var err error
	n := 1
	if mid > 0 {
		res, err := conn.Do("HDel", keyMidServer(mid), key)
		result = res.(int64)
		if err != nil {
			log.Error("Dao", "DelMapping", "conn.Send(HDEL %d,%s,%s) error(%v)", nil, mid, key, server, err)
			return false, err
		}
		n++
	}
	res, err := conn.Do("Del", keyKeyServer(key))
	result = res.(int64)
	if err != nil {
		log.Error("Dao", "DelMapping", "conn.Send(HDEL %d,%s,%s) error(%v)", nil, mid, key, server, err)
		return false, err
	}
	return strconv.ParseBool(fmt.Sprintf("%v", result))
}

// ServersByKeys get a server by key.
func (d *Dao) ServersByKeys(c echo.Context, keys []string) ([]string, error) {
	conn := model.GetRedis().Get()
	var (
		args   []string
		resd   []interface{}
		err    error
		back   []string
		getArg interface{}
	)
	for _, key := range keys {
		// args = append(args, keyKeyServer(key))
		getArg, err = conn.Do("GET", keyKeyServer(key))
		resd = append(resd, getArg)
	}
	// resd, err = redis.Values(conn.Do("MGET", args...))

	if err != nil {
		log.Error("Dao", "ServersByKeys", "conn.Do(MGET %v) error(%v)", nil, args, err)
	}
	// val []interface{}
	for _, param := range resd {
		back = append(back, param.(string))
	}
	return back, err
}

// KeysByMids get a key server by mid.
func (d *Dao) KeysByMids(c echo.Context, mids []int64) (map[string]string, []int64, error) {
	conn := model.GetRedis().Get()
	ress := make(map[string]string)
	var olMids []int64
	for k, mid := range mids {
		resd, err := conn.Do("HGetAll", keyMidServer(mid))
		if err != nil {
			log.Error("Dao", "KeysByMids", "conn.Do(HGETALL %d) error(%v)", nil, mid, err)
			return ress, olMids, err
		}
		fResults := string(resd.([]byte))
		var fs map[string]string
		json.Unmarshal([]byte(fResults), &fs)

		if len(fs) > 0 {
			olMids = append(olMids, mids[k])
		}
		for index, v := range fs {
			ress[index] = v
		}
	}
	return ress, olMids, nil
}

// AddServerOnline add a server online.
func (d *Dao) AddServerOnline(c context.Context, server string, online *model.Online) (err error) {
	roomsMap := map[uint32]map[string]int32{}
	for room, count := range online.RoomCount {
		rMap := roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64]
		if rMap == nil {
			rMap = make(map[string]int32)
			roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64] = rMap
		}
		rMap[room] = count
	}
	key := keyServerOnline(server)
	for hashKey, value := range roomsMap {
		err = d.addServerOnline(c, key, strconv.FormatInt(int64(hashKey), 10), &model.Online{RoomCount: value, Server: online.Server, Updated: online.Updated})
		if err != nil {
			return
		}
	}
	return
}

func (d *Dao) addServerOnline(c context.Context, key string, hashKey string, online *model.Online) (err error) {
	conn := model.GetRedis().Get()
	b, _ := json.Marshal(online)
	_, err = conn.Do("HSet", key, hashKey, b)
	if err != nil {
		log.Error("Dao", "addServerOnline", "conn.Send(SET %s,%s) error(%v)", nil, key, hashKey, err)
		return
	}
	_, err = conn.Do("Expire", key, time.Duration(d.redisExpire)*time.Second)
	if err != nil {
		log.Error("Dao", "addServerOnline", "conn.Send(EXPIRE %s) error(%v)", nil, key, err)
		return
	}
	return
}

// ServerOnline get a server online.
func (d *Dao) ServerOnline(c context.Context, server string) (online *model.Online, err error) {
	online = &model.Online{RoomCount: map[string]int32{}}
	key := keyServerOnline(server)
	for i := 0; i < 64; i++ {
		ol, err := d.serverOnline(c, key, strconv.FormatInt(int64(i), 10))
		if err == nil && ol != nil {
			online.Server = ol.Server
			if ol.Updated > online.Updated {
				online.Updated = ol.Updated
			}
			for room, count := range ol.RoomCount {
				online.RoomCount[room] = count
			}
		}
	}
	return
}

func (d *Dao) serverOnline(c context.Context, key string, hashKey string) (online *model.Online, err error) {
	conn := model.GetRedis().Get()
	back, err := conn.Do("HGet", key, hashKey)
	fResult := string(back.([]byte))
	if err != nil {
		log.Error("Dao", "serverOnline", "conn.Do(HGET %s %s) error(%v)", nil, key, hashKey, err)
		return
	}
	online = new(model.Online)
	if len(fResult) > 0 {
		if err = json.Unmarshal([]byte(fResult), online); err != nil {
			log.Error("Dao", "serverOnline", "serverOnline json.Unmarshal(%s) error(%v)", nil, back, err)
			return
		}
	}
	return
}

// DelServerOnline del a server online.
func (d *Dao) DelServerOnline(c context.Context, server string) (err error) {
	conn := model.GetRedis().Get()
	key := keyServerOnline(server)
	_, err = conn.Do("Del", key)
	if err != nil {
		log.Error("Dao", "DelServerOnline", "conn.Do(DEL %s) error(%v)", nil, key, err)
	}
	return
}

func (d *Dao) SessionByKey(c context.Context, key string) (session *model.LogicSession, err error) {
	conn := model.GetRedis().Get()
	back, err := conn.Do("Get", keyKeyServer(key))
	fResult := string(back.([]byte))
	if err != nil {

		log.Error("Dao", "SessionByKey", "conn.Do(GET %s) error(%v)", nil, key, err)
		return
	}
	session = new(model.LogicSession)
	if err = json.Unmarshal([]byte(fResult), session); err != nil {
		log.Error("Dao", "SessionByKey", "session json.Unmarshal(%s) error(%v)", nil, back, err)
		return
	}
	return
}

/**
group_member_count
*/
func (d *Dao) UpdateMemberAll(c context.Context, lineId, roomId, count string) (err error) {
	conn := model.GetRedis().Get()
	_, err = conn.Do("HSet", groupMemberAll(lineId), roomId, count)
	if err != nil {
		log.Error("Dao", "UpdateMemberAll", "%v", err)
		return
	}

	return
}
func (d *Dao) UpdateMemberRoom(c context.Context, lineId, roomId, memberInfo, mid string) (err error) {
	conn := model.GetRedis().Get()
	_, err = conn.Do("HSet", groupMemberRoom(lineId, roomId), mid, memberInfo)
	if err != nil {
		log.Error("Dao", "UpdateMemberRoom", "%v", err)
		return
	}

	return
}
func (d *Dao) GetMemberCount(c context.Context, lineId, roomId string) (map[string]string, error) {
	conn := model.GetRedis().Get()
	obj, err := conn.Do("HGetAll", groupMemberRoom(lineId, roomId))
	var cm map[string]string
	fResults := string(obj.([]byte))
	json.Unmarshal([]byte(fResults), &cm)
	if err != nil {
		log.Error("Dao", "GetMemberCount", "%v", err)
		return nil, err
	} else {
		return cm, nil
	}

}
func (d *Dao) GetGroupMemberAll(c context.Context, lineId string) (map[string]string, error) {
	conn := model.GetRedis().Get()

	objs, err := conn.Do("HGetAll", groupMemberAll(lineId))
	var cm map[string]string
	fResults := string(objs.([]byte))
	json.Unmarshal([]byte(fResults), &cm)
	if err != nil {
		log.Error("Dao", "GetGroupMemberAll", "%v", err)
		return nil, err
	} else {

		return cm, nil
	}

}
func (d *Dao) UpdateGroupMemberCount(c context.Context, lineId, RoomId string, mid string, addCount bool, memberInfo string, online bool, isDel bool) error {
	conn := model.GetRedis().Get()

	fs, err := conn.Do("HGet", groupMemberAll(lineId), RoomId)
	obj := string(fs.([]byte))
	if err != nil {
		if addCount && err == nil {
			if online {
				_, err = conn.Do("HSet", groupMemberAll(lineId), RoomId, 1)
				if err != nil {
					log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
				}
			}

		}
		if isDel {
			_, err = conn.Do("HDel", groupMemberRoom(lineId, RoomId), mid)
			if err != nil {
				log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
			}
		} else {
			_, err = conn.Do("HSet", groupMemberRoom(lineId, RoomId), mid, memberInfo)
			if err != nil {
				log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
			}
		}
	} else {
		objInt, err := strconv.Atoi(obj)
		if err != nil {
			log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
		}
		if addCount {
			objInt += 1
			if online {
				_, err = conn.Do("HSet", groupMemberAll(lineId), RoomId, objInt)
				if err != nil {
					log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
				}
			}
			_, err = conn.Do("HSet", groupMemberRoom(lineId, RoomId), mid, memberInfo)
			if err != nil {
				log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
			}
		} else {
			objInt -= 1
			if online {
				if objInt <= 0 {
					_, err = conn.Do("HDel", groupMemberAll(lineId), RoomId)
					if err != nil {
						log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
					}

				} else {
					_, err = conn.Do("HSet", groupMemberAll(lineId), RoomId, objInt)
					if err != nil {
						log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
					}
				}
			}
			if isDel {
				_, err = conn.Do("HDel", groupMemberRoom(lineId, RoomId), mid)
				if err != nil {
					log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
				}
			} else {
				_, err = conn.Do("HSet", groupMemberRoom(lineId, RoomId), mid, memberInfo)
				if err != nil {
					log.Error("Dao", "UpdateGroupMemberCount", "%v", err)
				}
			}
		}

	}
	return nil
}
func (d *Dao) KickGroupMemberCount(c context.Context, lineId, RoomId string, mid string, online bool) error {
	conn := model.GetRedis().Get()

	if obj, err := conn.Do("HGet", groupMemberAll(lineId), RoomId); err != nil {

		_, err = conn.Do("HDel", groupMemberRoom(lineId, RoomId), mid)
		if err != nil {
			log.Error("Dao", "KickGroupMemberCount", "%v", err)
		}
	} else {
		fs := string(obj.([]byte))
		objInt, err := strconv.Atoi(fs)
		if err != nil {
			log.Error("Dao", "KickGroupMemberCount", "%v", err)
		}
		objInt -= 1
		if online {
			if _, err = conn.Do("HSet", groupMemberAll(lineId), RoomId, objInt); err != nil {
				log.Error("Dao", "KickGroupMemberCount", "%v", err)
			}
		}
		if _, err = conn.Do("HDel", groupMemberRoom(lineId, RoomId), mid); err != nil {
			log.Error("Dao", "KickGroupMemberCount", "%v", err)
		}

	}
	return nil
}
