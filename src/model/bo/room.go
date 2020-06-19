package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
)

type Room struct{}

// 查询房间列表
func (*Room) QueryRoomList(sess *xorm.Session, lineId string, agencyId string, startTime, endTime, gameType, status int, page, pagSize int) (int64, []*structs.RoomResp, error) {
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}

	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}

	if startTime != 0 {
		sess.Where("create_time >= ?", startTime)
	}

	if endTime != 0 {
		sess.Where("create_time <= ?", endTime)
	}

	if gameType != 0 {
		sess.Where("game_type = ? ", gameType)
	}

	if status != 0 {
		sess.Where("status = ? ", status)
	}
	sess.Where("delete_time = ? ", model.UNDEL)
	rooms := make([]*structs.RoomResp, 0)
	count, err := sess.Table(new(structs.Room).TableName()).Limit(pagSize, (page-1)*pagSize).
		OrderBy("room_sort asc,id desc").FindAndCount(&rooms)
	return count, rooms, err
}

// 添加房间
func (*Room) SaveRoom(sess *xorm.Session, room *structs.Room) error {
	_, err := sess.Insert(room)
	return err
}

// 查询单个房间
func (*Room) GetOne(sess *xorm.Session, id int) (bool, *structs.Room) {
	sess.Where("id = ?", id)
	room := new(structs.Room)
	has, _ := sess.Get(room)
	return has, room
}

// 查询单个房间
func (*Room) GetOneByRoomNo(sess *xorm.Session, lineId string, agencyId string, roomNo int) (bool, *structs.Room) {
	sess.Where("room_no = ? ", roomNo)
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	sess.Where("delete_time = ? and room_type = ?", 0, 2)
	room := new(structs.Room)
	has, _ := sess.Get(room)
	return has, room
}

// 根据ids查询多个房间
func (*Room) FindRoomsByIds(sess *xorm.Session, ids []int) ([]*structs.RoomResp, error) {
	if len(ids) > 0 {
		sess.In("id  ", ids)
	}
	sess.Where("delete_time = ?", 0)
	rooms := make([]*structs.RoomResp, 0)
	err := sess.Table(new(structs.Room).TableName()).Find(&rooms)
	return rooms, err
}

// 修改房间
func (*Room) ModifyRoom(sess *xorm.Session, room *structs.Room) error {
	_, err := sess.Where("id = ? and delete_time = ?", room.Id, 0).
		Cols("room_name", "game_type", "max_money", "min_money",
			"odds", "game_play", "red_num", "red_min_num", "royalty", "game_time", "room_sort",
			"room_type", "free_from_death", "robot_send_packet", "robot_send_packet_time", "robot_grab_packet", "control_kill").Update(room)
	return err
}

// 修改房间状态
func (*Room) ModifyRoomStatus(sess *xorm.Session, room *structs.Room) error {
	_, err := sess.Where("id = ? and delete_time = ?", room.Id, 0).Cols("status").Update(room)
	return err
}

// 删除房间
func (*Room) DelRoom(sess *xorm.Session, room *structs.Room) error {
	_, err := sess.Table(new(structs.Room).TableName()).
		ID(room.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("delete_time").
		Update(room)
	return err
}

// 查询游戏房间列表
func (*Room) GetRoomList(sess *xorm.Session, lineId, agencyId string, gameType, gamePlay int) ([]structs.Room, error) {
	data := make([]structs.Room, 0)
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("game_type = ?", gameType)
	if gameType != model.MINESWEEPER_RED_PACKET {
		sess.Where("game_play = ?", gamePlay)
	}
	sess.Where("status = 1")
	sess.Where("delete_time = 0 and room_type = 1")
	err := sess.OrderBy("room_sort asc").Find(&data)
	return data, err

}

func (*Room) UpdateRoom(sess *xorm.Session, req *structs.Room, fields string) (int, error) {
	sess.Where("delete_time = ? ", model.UNDEL)
	count, err := sess.ID(req.Id).Cols(fields).Update(req)
	return int(count), err
}

// 获取最大房间号
func (*Room) GetMaxRoomNo(sess *xorm.Session) (int, error) {
	room := new(structs.Room)
	res := sess.OrderBy("room_no desc").GetFirst(room)
	if res.Has && res.Error == nil && room.RoomNo != 0 {
		return room.RoomNo, nil
	}
	if !res.Has && res.Error == nil {
		return 99999, nil
	}
	return 0, res.Error
}

// 群枚举
func (*Room) RoomCode(sess *xorm.Session, lineId, agencyId string, gameType int) ([]*structs.RoomCode, error) {
	roomCodes := make([]*structs.RoomCode, 0)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	sess.Where("delete_time = 0")
	err := sess.Table(new(structs.Room).TableName()).Find(&roomCodes)
	return roomCodes, err
}
