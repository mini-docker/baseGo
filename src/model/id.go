package model

import (
	"sync"

	"strconv"

	"github.com/mini-docker/baseGo/src/fecho/utility/sonyflake"
)

var sf *sonyflake.Sonyflake
var idKey uint16
var mu sync.Mutex

func IdgenInit(key uint16) {
	var st sonyflake.Settings
	// 初始化key .
	idKey = key
	st.MachineID = genId
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func genId() (uint16, error) {
	return idKey, nil
}

// GetId 获取唯一id
func GetId() (int, error) {
	mu.Lock()
	defer mu.Unlock()

	id, err := sf.NextID()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func GetIdString() (string, error) {
	id, err := GetId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}
