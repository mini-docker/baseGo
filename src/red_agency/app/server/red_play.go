package server

import (
	"math/rand"
	"model"
	"strconv"
	"time"
)

// 红包玩法计算
type RedPlay struct{}

//
type RedSettlement struct {
	RedId   int // 红包ID
	roomId  int // 房间ID
	RedType int // 红包类型 1 牛牛红包 2
	RedPlay int // 红包玩法 1
	Data    []UserRed
}

// 红包结果结构体
type UserRed struct {
	UserId    int     // 会员ID
	Identity  int     //	会员身份 1庄家 2闲家
	RedAmount float64 //	红包金额
	Win       float64 // 结果
}

// 计算尾部3位数字的合的个位数
func NiuNiuCalculation(money float64) int {
	num := int(money * 100)
	n1 := num % 10         // 个位数字
	n2 := (num / 10) % 10  //拾位
	n3 := (num / 100) % 10 //百位
	return (n1 + n2 + n3) % 10
}

// 订单号生成
func OderNo(gameType int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var order string
	if gameType == model.NIUNIU_RED_ENVELOPE {
		order = order + "slhb"
	} else if gameType == model.MINESWEEPER_RED_PACKET {
		order = order + "nnhb"
	}

	order += time.Now().Format("20060102150405") + strconv.Itoa(100000+r.Intn(899999))
	return order
}
