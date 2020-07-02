package http

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/imserver/internal/logic/model"
	"context"
)

func (s *Server) CacheSession(ctx echo.Context) error {
	req := new(model.CacheSessionReq)
	if err := ctx.ValidRequest(req); err != nil {
		return ctx.JSON4Error(RequestErr)
	}

	mid, key, room, accepts, hb, server, err := s.logic.CacheSession(context.Background(), req.Cookie, req.Token, req.UserId, req.Rooms)
	if err != nil {
		return err
	}

	return ctx.JSON4Item(model.ConnectReply{
		Mid:       mid,
		Key:       key,
		RoomID:    room,
		Accepts:   accepts,
		Heartbeat: hb,
		Server:    server,
	})
}

func (s *Server) DelRoom(ctx echo.Context) error {
	req := new(model.DelRoomReq)
	if err := ctx.ValidRequest(req); err != nil {
		return ctx.JSON4Error(RequestErr)
	}

	err := s.logic.DelRoom(context.Background(), req.RoomId)
	if err != nil {
		return err
	}

	return ctx.JSON4Item(200)
}
func (s *Server) IntiveOrKickRoom(ctx echo.Context) error {
	req := new(model.IntiveOrKickRoomReq)
	if err := ctx.ValidRequest(req); err != nil {
		return ctx.JSON4Error(RequestErr)
	}
	if req.IsInvite {
		err := s.logic.IntiveRoom(context.Background(), req.UserKeys, req.RoomId, req.NoticeType, req.SenderId, req.LineId)
		if err != nil {
			return err
		}
	} else {
		err := s.logic.KickRoom(context.Background(), req.UserKeys, req.RoomId, req.NoticeType, req.SenderId, req.LineId)
		if err != nil {
			return err
		}
	}
	data := new(Data)
	data.Code = 200

	return ctx.JSON4Item(data)
}

func (s *Server) JoinRoom(ctx echo.Context) error {
	req := new(model.ChangeRoomReq)
	if err := ctx.ValidRequest(req); err != nil {
		return ctx.JSON4Error(RequestErr)
	}

	for _, userId := range req.UserId {
		err := s.logic.JoinRoom(context.Background(), req.LineId, req.AgencyId, int64(userId), req.RoomId, req.NoticeType, req.SenderId, req.Content, req.RoomType)
		if err != nil {
			return err
		}
	}

	return ctx.JSON4Item(200)
}

func (s *Server) OutRoom(ctx echo.Context) error {
	req := new(model.OutRoomReq)
	if err := ctx.ValidRequest(req); err != nil {
		return ctx.JSON4Error(RequestErr)
	}

	err := s.logic.OutRoom(context.Background(), req.LineId, req.AgencyId, int64(req.UserId), req.RoomId, req.NoticeType, req.RoomType)
	if err != nil {
		return err
	}

	return ctx.JSON4Item(200)
}
