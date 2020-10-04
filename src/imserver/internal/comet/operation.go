package comet

import (
	registry_module "baseGo/src/fecho/registry/registry-module"
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "baseGo/src/imserver/api/comet/grpc"
	logic "baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/pkg/strings"
	log "fecho/golog"

	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

// Connect connected a connection.
func (s *Server) Connect(c context.Context, p *model.Proto, cookie string) (rs bool, mid int64, key, rid string, accepts []int32, heartbeat time.Duration, roomids []string, err error) {
	//var bodyData struct {
	//	Op   int32  `json:"op"`
	//	Body []byte `json:"body"`
	//}
	var params struct {
		UserId  int64  `json:"userId"`
		Account string `json:"account"`
	}
	//if err = json.Unmarshal(p.Body, &bodyData); err != nil {
	//	return
	//}
	if err = json.Unmarshal(p.Body, &params); err != nil {
		return
	}
	p.UserId = params.UserId
	p.SessionKey = fmt.Sprintf("%s-%d", params.Account, params.UserId)

	client := s.GetLogicClient()
	if nil == client {
		err = errors.New("auth error")
		return
	}
	reply, err := client.EstablishConn(c, &logic.EstablishConnReq{
		UserId: p.UserId,
		Base: &logic.Base{
			Server: registry_module.GetCometRpcUrl(),
			Tokken: p.SessionKey,
		},
	})
	rs = false
	if err != nil {
		return
	}
	rs = reply.Result
	return rs, reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), reply.Rooms, nil
}

// Disconnect disconnected a connection.
func (s *Server) Disconnect(c context.Context, mid int64, key string) (err error) {
	client := s.GetLogicClient()
	if nil == client {
		return
	}
	_, err = client.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: registry_module.GetCometRpcUrl(),
		Mid:    mid,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
func (s *Server) Heartbeat(ctx context.Context, mid int64, key string) (userRoom string, adminRooms []string, err error) {
	fmt.Println("接收心跳：", key)
	client := s.GetLogicClient()
	if nil == client {
		return
	}
	_, err = client.Heartbeat(ctx, &logic.HeartbeatReq{
		Server: registry_module.GetCometRpcUrl(),
		Mid:    mid,
		Key:    key,
	})
	//reply, err := s.rpcClient.EstablishConn(ctx, &logic.EstablishConnReq{
	//	UserId: mid,
	//	Base: &logic.Base{
	//		Server:      s.serverID,
	//	},
	//})
	//if nil != reply {
	//	if nil != reply.Rooms {
	//		adminRooms = reply.Rooms
	//	}
	//	if common.IsNotBlank(reply.RoomID) {
	//		userRoom = reply.RoomID
	//	}
	//}

	return
}

// RenewOnline renew room online.
func (s *Server) RenewOnline(ctx context.Context, serverID string, rommCount map[string]int32) (allRoom map[string]int32, err error) {
	client := s.GetLogicClient()
	if nil == client {
		return
	}
	reply, err := client.RenewOnline(ctx, &logic.OnlineReq{
		Server:    registry_module.GetCometRpcUrl(),
		RoomCount: rommCount,
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return
	}
	return reply.AllRoomCount, nil
}

// Receive receive a message.
func (s *Server) Receive(ctx context.Context, mid int64, p *model.Proto) (err error) {
	client := s.GetLogicClient()
	if nil == client {
		return
	}
	_, err = client.Receive(ctx, &logic.ReceiveReq{Mid: mid, Proto: p})
	return
}

// Operate operate.
func (s *Server) Operate(ctx context.Context, p *model.Proto, ch *Channel, b *Bucket) error {
	switch p.Op {
	case model.OpChangeRoom:
		var params struct {
			UserId int64 `json:"userId"`
		}
		if err := json.Unmarshal(p.Body, &params); err != nil {
			return nil
		}
		p.UserId = params.UserId
		p.SessionKey = fmt.Sprint(params.UserId)

		client := s.GetLogicClient()
		if nil == client {
			return nil
		}
		reply, err := client.JoinRoom(ctx, &logic.JoinRoomReq{
			UserId: p.UserId,
			Base: &logic.Base{
				Server: registry_module.GetCometRpcUrl(),
			},
		})
		if err != nil {
			return err
		}
		if reply.Result == false {
			return nil
		}
		if err := b.ChangeRoom(reply.RoomID, ch); err != nil {
			log.Error("Server", "Operate", "error :b.ChangeRoom(%v, %v) error(%v)", nil, reply.RoomID, ch, err)
		}
		p.Op = model.OpChangeRoomReply
	case model.OpSub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = model.OpSubReply
	case model.OpUnsub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = model.OpUnsubReply
	default:
		// TODO ack ok&failed
		if err := s.Receive(ctx, ch.Mid, p); err != nil {
			log.Info("Server", "Operate", "s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		}
		p.Body = nil
	}
	return nil
}

// WsConnect connected a connection.
func (s *Server) WsConnect(c context.Context, p *model.Proto, cookie string) (result bool, mid int64, key, rid string, accepts []int32, heartbeat time.Duration, lineId, agencyId string, roomids []string, err error) {

	var params struct {
		LineId   string `json:"lineId"`
		AgencyId string `json:"agencyId"`
		UserId   int64  `json:"userId"`
	}
	if err = json.Unmarshal(p.Body, &params); err != nil {
		return
	}
	p.UserId = params.UserId
	p.AgencyId = params.AgencyId
	p.LineId = params.LineId
	p.SessionKey = params.LineId + "-" + params.AgencyId + "-" + fmt.Sprint(params.UserId)

	client := s.GetLogicClient()
	if nil == client {
		err = errors.New("auth error")
		return
	}
	reply, err := client.EstablishConn(c, &logic.EstablishConnReq{
		UserId: p.UserId,
		Base: &logic.Base{
			LineId:   p.LineId,
			AgencyId: p.AgencyId,
			Server:   registry_module.GetCometRpcUrl(),
		},
	})

	result = false
	if err != nil {
		return
	}
	if nil == reply || !reply.Result {
		err = errors.New("seesion not found")
		return
	}

	result = reply.Result
	return result, reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), params.LineId, params.AgencyId, reply.Rooms, nil
}
