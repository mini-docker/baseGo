package services

import "baseGo/src/model/bo"

var (
	RoomBo             = new(bo.Room)
	UserBo             = new(bo.User)
	RedPacketBo        = new(bo.RedPacket)
	RedPacketLogBo     = new(bo.RedPacketLog)
	MemberCashRecordBo = new(bo.MemberCashRecord)
	AgencyBo           = new(bo.Agency)
	// PostBo             = new(bo.Post)
	// ActivePictureBo    = new(bo.ActivePicture)
	// ImOnlineMessageBo  = new(bo.ImOnlineMessage)
	MessageHistoryBo = new(bo.MessageHistory)
	// SystemLineBo       = new(bo.SystemLineBo)
	RedLogBo = new(bo.RedLog)
)
