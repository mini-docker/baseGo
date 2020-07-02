package http

import (
	"baseGo/src/fecho/echo"
	"context"
)

func (s *Server) onlineTop(c echo.Context) error {
	var arg struct {
		Type  string `form:"type" valid:"Must;ErrorCode()"`
		Limit int    `form:"limit" valid:"Must;ErrorCode()"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}

	res, err := s.logic.OnlineTop(c, arg.Type, arg.Limit)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func (s *Server) onlineRoom(c echo.Context) error {
	var arg struct {
		Type  string   `form:"type" valid:"Must;ErrorCode()"`
		Rooms []string `form:"rooms" valid:"Must;ErrorCode()"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	res, err := s.logic.OnlineRoom(c, arg.Type, arg.Rooms)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func (s *Server) onlineTotal(c echo.Context) error {
	ipCount, connCount := s.logic.OnlineTotal(context.TODO())
	res := map[string]interface{}{
		"ip_count":   ipCount,
		"conn_count": connCount,
	}
	return c.JSON(200, res)
}

// 查询房间在线人数
func (s *Server) RoomsOnlineCount(c echo.Context) error {
	var arg struct {
		Ids      string `form:"ids"` // roomIds 房间号
		LineId   string `form:"lineId"`
		AgencyId string `form:"agencyId"`
	}
	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}

	resp, err := s.logic.RoomCount(context.Background(), arg.LineId, arg.AgencyId)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}

// 查询房间在线会员id
// 单个房间的id
func (s *Server) RoomsOnlineMids(c echo.Context) error {

	var arg struct {
		Id       string `form:"id"` // roomIds 房间号
		LineId   string `form:"lineId"`
		AgencyId string `form:"agencyId"`
	}

	if err := c.ValidRequest(&arg); err != nil {
		return c.JSON4Error(RequestErr)
	}
	resp, err := s.logic.RoomInfo(context.Background(), arg.Id, arg.LineId, arg.AgencyId)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}
