package logic

import (
	"fmt"
	"model/bo"
)

var (
	RoomBo           = new(bo.Room)
	User             = new(bo.User)
	_groupMemberRoom = "groupMemberRoom_%s" // _groupMemberRoom
	_groupMemberAll  = "groupMemberAll"
)

func groupMemberRoom(RoomId string) string {
	return fmt.Sprintf(_groupMemberRoom, RoomId)
}
