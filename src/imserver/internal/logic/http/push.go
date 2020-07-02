package http

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/model/structs"
)

type Data struct {
	Code int              `json:"code"`
	Info *structs.MsgResp `json:"info"`
}

func (s *Server) pushMid(c echo.Context) error {
	var arg struct {
		Operation   int32  `json:"operation"`
		RoomId      int    `json:"room_id"`
		Mid         int32  `json:"mid"`
		Msg         string `json:"msg" valid:"Html"`
		MsgType     int32  `json:"msg_type"`
		SendId      int32  `json:"send_id"`
		Account     string `json:"account"`
		ReceiveType int32  `json:"receive_type"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	info, err := s.logic.PushMid(arg.Operation, arg.SendId, arg.Mid, arg.MsgType, []byte(arg.Msg), arg.Account, arg.ReceiveType)
	if err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	data.Info = info
	return c.JSON4Item(data)
}

func (s *Server) pushRoom(c echo.Context) error {
	var arg struct {
		Op          int32  `json:"operation"`
		RoomId      int32  `json:"roomId"`
		LineId      string `json:"lineId"`
		AgencyId    string `json:"agencyId"`
		Msg         string `json:"msg" valid:"Html"`
		MsgType     int32  `json:"msgType"`
		SendId      int32  `json:"sendId"`
		Key         string `json:"key"`
		ReceiveType int32  `json:"receiveType"`
	}

	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	info, err := s.logic.PushRoom(arg.Op, arg.RoomId, arg.SendId, arg.MsgType, []byte(arg.Msg), arg.ReceiveType, arg.LineId, arg.AgencyId)
	if err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	data.Info = info
	return c.JSON4Item(data)
}

func (s *Server) PushRoomSimulation(c echo.Context) error {
	var arg struct {
		Op        int32  `json:"operation"`
		RoomId    int32  `json:"roomId"`
		Msg       string `json:"msg" valid:"Html"`
		MsgType   int32  `json:"msgType"`
		SendId    int32  `json:"sendId"`
		Keys      string `json:"key"`
		IsShillId int32  `json:"isShillId"`
		RoomType  int    `json:"roomType"`
	}

	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	_, err := s.logic.PushRoomSimulation(arg.Keys, arg.Op, arg.RoomId, arg.IsShillId, arg.MsgType, arg.SendId, []byte(arg.Msg), arg.RoomType)
	if err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	//data.MsgId = id
	return c.JSON4Item(data)
}

func (s *Server) pushAll(c echo.Context) error {
	var arg struct {
		Op    int32  `json:"operation"`
		Speed int32  `json:"speed"`
		Msg   string `json:"msg"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}

	if err := s.logic.PushAll(arg.Op, arg.Speed, []byte(arg.Msg)); err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	return c.JSON4Item(data)
}

func (s *Server) sendNotificationMessage(c echo.Context) error {
	var arg struct {
		Op         int32  `json:"operation"`
		Speed      int32  `json:"speed"`
		Msg        string `json:"msg" valid:"Html"`
		ReceiverId int32  `json:"receiverId"`
		SendTime   int32  `json:"sendTime"`
		MsgType    int32  `json:"msgType"`
		PushCrowd  int32  `json:"pushCrowd"` // 推送类型 1单人 2所有
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}

	if err := s.logic.SendNotificationMessage(arg.Op, arg.ReceiverId, arg.MsgType, arg.SendTime, arg.Speed, arg.PushCrowd, []byte(arg.Msg)); err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	return c.JSON4Item(data)
}

func (s *Server) pushGroup(c echo.Context) error {
	var arg struct {
		Op          int32  `json:"operation"`
		RoomId      int32  `json:"room_id"`
		Msg         string `json:"msg" valid:"Html"`
		MsgType     int32  `json:"msg_type"`
		SendId      int32  `json:"send_id"`
		Keys        string `json:"key"`
		RoomType    int    `json:"room_type"`
		ReceiveType int32  `json:"receive_type"`
		LineId      string `json:"line_id"`
		AgencyId    string `json:"agency_id"`
	}

	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	_, err := s.logic.PushRoom(arg.Op, arg.RoomId, arg.SendId, arg.MsgType, []byte(arg.Msg), arg.ReceiveType, arg.LineId, arg.AgencyId)
	if err != nil {
		return err
	}
	data := new(Data)
	data.Code = 200
	return c.JSON4Item(data)
}
