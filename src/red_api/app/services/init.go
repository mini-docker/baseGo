package services

import "baseGo/src/model/bo"

var (
	RoomBo             = new(bo.Room)
	UserBo             = new(bo.User)
	RedPacketBo        = new(bo.RedPacket)
	RedPacketLogBo     = new(bo.RedPacketLog)
	MemberCashRecordBo = new(bo.MemberCashRecord)
	AgencyBo           = new(bo.Agency)
	// AgencyCashRecordBo = new(bo.AgencyCashRecord)
	// RedPacketCollectBo = new(bo.RedPacketCollect)
	PostBo = new(bo.Post)
)
