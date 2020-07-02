package grpc

import (
	"context"
	"net"
	"strconv"
	"time"

	pb "baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/imserver/internal/comet"
	"baseGo/src/imserver/internal/comet/conf"
	"baseGo/src/imserver/internal/comet/errors"

	"fecho/golog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// New comet grpc server.
func New(c *conf.RPCServer, s *comet.Server) *grpc.Server {
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(c.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(c.ForceCloseWait),
		Time:                  time.Duration(c.KeepAliveInterval),
		Timeout:               time.Duration(c.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(c.MaxLifeTime),
	})
	srv := grpc.NewServer(keepParams)
	pb.RegisterCometServer(srv, &server{s})
	lis, err := net.Listen(c.Network, ":"+c.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	srv *comet.Server
}

var _ pb.CometServer = &server{}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	// TODO: some graceful close
	return &pb.Empty{}, nil
}

// PushMsg push a message to specified sub keys.
func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
	golog.Info("comet", "PushMsg", "push msg:", req.Keys)
	if len(req.Keys) == 0 || req.Proto == nil {
		return nil, errors.ErrPushMsgArg
	}
	for _, v := range req.Keys {
		golog.Info("bucket", "PushMsg", "---------------private msg push:%s", v)
		bkt := s.srv.Bucket(req.Proto.LineId, req.Proto.AgencyId)
		if bkt == nil {
			continue
		}
		if channel := bkt.Channel(v); channel != nil {
			if err = channel.Push(req.Proto); err != nil {
				golog.Error("server", "PushMsg", "err:", err)
				continue
			}
		} else {
			golog.Info("server", "PushMsg", "channel didnt find :%s---------------", v)
			continue
		}
		golog.Info("server", "PushMsg", "private msg push success:%s---------------", v)
	}
	return &pb.PushMsgReply{}, nil
}

// Broadcast broadcast msg to all user.
func (s *server) Broadcast(ctx context.Context, req *pb.BroadcastReq) (*pb.BroadcastReply, error) {
	if req.Proto == nil {
		return nil, errors.ErrBroadCastArg
	}
	// TODO use broadcast queue
	go func() {
		//for _, bucket := range s.srv.Buckets() {
		bucket := s.srv.Bucket(req.Proto.LineId, req.Proto.AgencyId)
		if nil == bucket {
			return
		}
		bucket.Broadcast(req.GetProto(), req.ProtoOp)
		if req.Speed > 0 {
			t := bucket.ChannelCount() / int(req.Speed)
			time.Sleep(time.Duration(t) * time.Second)
		}
		//}
	}()
	return &pb.BroadcastReply{}, nil
}

// BroadcastRoom broadcast msg to specified room.
func (s *server) BroadcastRoom(ctx context.Context, req *pb.BroadcastRoomReq) (*pb.BroadcastRoomReply, error) {

	//golog.Info( "comet", "PushMsg", "push room:", req.RoomID)
	if req.Proto == nil || req.RoomID == "" {
		return nil, errors.ErrBroadCastRoomArg
	}
	//golog.Info( "comet", "BroadcastRoom", "req", string(req.Proto.Body))
	for _, bucket := range s.srv.Buckets() {
		room := bucket.Room(req.RoomID)
		if nil == room {
			continue
		}
		bucket.BroadcastRoom(req)
	}
	//if nil != room {
	//	s.srv.Bucket(req.Proto.LineId, req.Proto.AgencyId).BroadcastRoom(req)
	//}

	//}
	//golog.Info( "comet", "BroadcastRoom", string(req.Proto.Body), req.Proto.Op)
	return &pb.BroadcastRoomReply{}, nil
}

// Rooms gets all the room ids for the server.
func (s *server) Rooms(ctx context.Context, req *pb.RoomsReq) (*pb.RoomsReply, error) {
	var (
		roomIds = make(map[string]bool)
	)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.Rooms() {
			roomIds[roomID] = true
		}
	}
	return &pb.RoomsReply{Rooms: roomIds}, nil
}

// 进入房间
func (s *server) JoinRoom(ctx context.Context, in *pb.JoinRoomReq) (rs *pb.JoinRoomeResp, err error) {
	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)
	if btk == nil {
		//golog.Info( "bucket", "JoinRoom", "joinroom false :%+v", in.Proto)
		return &pb.JoinRoomeResp{Result: false}, nil
	}
	if ch := btk.Channel(in.Proto.SessionKey); nil != ch {
		err = btk.ChangeRoom(in.RoomID, ch)
	}
	return &pb.JoinRoomeResp{Result: true}, err
}

// 退出房间
func (s *server) LeaveRoom(ctx context.Context, in *pb.LeaveRoomReq) (*pb.LeaveRoomResp, error) {
	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)
	if btk == nil {
		//golog.Info( "bucket", "LeaveRoom", "leave false :%+v", in.Proto)
		return &pb.LeaveRoomResp{Result: false}, nil
	}
	err := btk.LeaveRoom(btk.Channel(in.Proto.SessionKey))
	return &pb.LeaveRoomResp{Result: true}, err
}

// 注销房间
func (s *server) CancelRoom(ctx context.Context, in *pb.CancelRoomReq) (*pb.CancelRoomResp, error) {

	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)
	if btk == nil {
		return &pb.CancelRoomResp{Result: false}, nil
	}

	if rm := btk.Room(in.RoomID); nil != rm {
		//btk.DelRoom(rm)
	}

	return &pb.CancelRoomResp{Result: true}, nil

}

// 房间信息
func (s *server) RoomInfolRoom(ctx context.Context, in *pb.RoomInfoReq) (*pb.RoomInfoResp, error) {
	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)
	if nil != btk {
		room := btk.Room(in.RoomID)
		if nil != room {
			return &pb.RoomInfoResp{Result: true, Mid: room.MIds}, nil
		}
	}

	return &pb.RoomInfoResp{Result: false, Mid: []int64{}}, nil
}

// 进入房间
func (s *server) JoinRooms(ctx context.Context, in *pb.JoinRoomsReq) (*pb.JoinRoomesResp, error) {
	//new(comet.Bucket).ChangeRoom()
	return nil, nil
}

//房间人数
func (s *server) RoomUserCount(ctx context.Context, in *pb.RoomUserCountReq) (*pb.RoomUserCountResp, error) {
	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)

	if nil != btk {
		roomMap := btk.RoomUserCount()
		rucs := make([]*pb.RoomUserCount, 0)
		for k, v := range roomMap {
			rucs = append(rucs, &pb.RoomUserCount{
				RoomId: k,
				Count:  int64(v),
			})
		}
		return &pb.RoomUserCountResp{
			RoomUserCount: rucs,
			Result:        true,
		}, nil
	}

	return &pb.RoomUserCountResp{Result: false}, nil
}

// VeeGebruikerSessieUit 会员是否在线
func (s *server) VeeGebruikerSessieUit(ctx context.Context, in *pb.RoomUserCountReq) (*pb.RoomUserCountResp, error) {
	btk := s.srv.Bucket(in.Proto.LineId, in.Proto.AgencyId)
	if btk == nil {
		return &pb.RoomUserCountResp{Result: false}, nil
	}
	var keyStr string
	// 发送人是管理员，接收人就是会员
	keyStr = strconv.Itoa(int(in.Proto.UserId))

	if ch := btk.Channel(keyStr); nil != ch {
		return &pb.RoomUserCountResp{Result: true}, nil
	}
	return &pb.RoomUserCountResp{Result: false}, nil
}
