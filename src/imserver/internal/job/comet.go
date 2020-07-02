package job

import (
	"context"
	registry_module "fecho/registry/registry-module"
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/connectivity"

	comet "baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/imserver/internal/job/conf"

	log "fecho/golog"

	"google.golang.org/grpc"
)

var (
	// grpc options
	grpcKeepAliveTime    = time.Duration(10) * time.Second
	grpcKeepAliveTimeout = time.Duration(3) * time.Second
	grpcBackoffMaxDelay  = time.Duration(3) * time.Second
	grpcMaxSendMsgSize   = 1 << 24
	grpcMaxCallMsgSize   = 1 << 24
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
)

var (
	cometGrpcAddr map[string]string
	cometClients  map[string]*grpc.ClientConn
)

func newCometClient() map[string]*grpc.ClientConn {
	clients := make(map[string]*grpc.ClientConn)
	if len(cometGrpcAddr) == 0 {
		cometGrpcAddr = make(map[string]string, 0)
	}
	if len(cometClients) == 0 {
		cometClients = make(map[string]*grpc.ClientConn, 0)
	}
	if len(registry_module.CometGrpcAddr) > 0 {
		// 检查连接是否有变化
		for k, v := range registry_module.CometGrpcAddr {
			if addr, ok := cometGrpcAddr[k]; !ok {
				// 地址变化更新地址
				cometGrpcAddr[k] = v
				// 更新连接
				conn, err := grpc.Dial(v, grpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := comet.NewCometClient(conn)
				clients[v] = conn
			} else {
				if addr == v {
					if cometClients[v] != nil {
						clients[v] = cometClients[v]
						continue
					}
				}
				// 地址变化更新地址
				cometGrpcAddr[k] = v
				// 更新连接
				conn, err := grpc.Dial(v, grpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := comet.NewCometClient(conn)
				clients[v] = conn
			}
		}
	} else {
		//  清空缓存
		for k, old := range cometClients {
			if old.GetState() != connectivity.Ready {
				fmt.Println("---------清空过期连接")
				old.Close()
			} else {
				clients[k] = old
			}
		}
		cometClients = make(map[string]*grpc.ClientConn)
		cometGrpcAddr = make(map[string]string)
	}
	for key, old := range cometClients {
		if _, ok := clients[key]; !ok {
			if old.GetState() != connectivity.Ready {
				fmt.Println("---------清空过期连接")
				old.Close()
			} else {
				clients[key] = old
			}
		}
	}
	cometClients = clients

	return cometClients
}

// Comet is a comet.
type Comet struct {
	serverID      string
	client        comet.CometClient
	pushChan      []chan *comet.PushMsgReq
	roomChan      []chan *comet.BroadcastRoomReq
	broadcastChan chan *comet.BroadcastReq
	pushChanNum   uint64
	roomChanNum   uint64
	routineSize   uint64

	ctx    context.Context
	cancel context.CancelFunc
}

// NewComet new a comet.
func NewComet(c *conf.Comet) map[string]*Comet {
	cmts := make(map[string]*Comet, 0)
	clients := newCometClient()
	if len(clients) > 0 {
		for addr, conn := range clients {
			cmt := &Comet{
				serverID:      "im-job",
				pushChan:      make([]chan *comet.PushMsgReq, c.RoutineSize),
				roomChan:      make([]chan *comet.BroadcastRoomReq, c.RoutineSize),
				broadcastChan: make(chan *comet.BroadcastReq, c.RoutineSize),
				routineSize:   uint64(c.RoutineSize),
			}
			cmt.client = comet.NewCometClient(conn)
			cmt.ctx, cmt.cancel = context.WithCancel(context.Background())

			for i := 0; i < c.RoutineSize; i++ {
				cmt.pushChan[i] = make(chan *comet.PushMsgReq, c.RoutineChan)
				cmt.roomChan[i] = make(chan *comet.BroadcastRoomReq, c.RoutineChan)
				go cmt.process(cmt.pushChan[i], cmt.roomChan[i], cmt.broadcastChan)
			}
			cmts[addr] = cmt
		}
	}

	return cmts
}

// Push push a user message.
func (c *Comet) Push(arg *comet.PushMsgReq) (err error) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.routineSize
	c.pushChan[idx] <- arg
	return
}

// BroadcastRoom broadcast a room message.
func (c *Comet) BroadcastRoom(arg *comet.BroadcastRoomReq) (err error) {
	idx := atomic.AddUint64(&c.roomChanNum, 1) % c.routineSize
	c.roomChan[idx] <- arg
	return
}

// Broadcast broadcast a message.
func (c *Comet) Broadcast(arg *comet.BroadcastReq) (err error) {
	c.broadcastChan <- arg
	return
}

func (c *Comet) process(pushChan chan *comet.PushMsgReq, roomChan chan *comet.BroadcastRoomReq, broadcastChan chan *comet.BroadcastReq) {
	for {
		select {
		case broadcastArg := <-broadcastChan:
			_, err := c.client.Broadcast(context.Background(), &comet.BroadcastReq{
				Proto:   broadcastArg.Proto,
				ProtoOp: broadcastArg.ProtoOp,
				Speed:   broadcastArg.Speed,
			})
			if err != nil {
				log.Error("Comet", "process", "c.client.Broadcast(%s, reply) serverId:%s error(%v)", nil, broadcastArg, c.serverID, err)
			}
		case roomArg := <-roomChan:
			_, err := c.client.BroadcastRoom(context.Background(), &comet.BroadcastRoomReq{
				RoomID: roomArg.RoomID,
				Proto:  roomArg.Proto,
			})
			if err != nil {
				log.Error("Comet", "process", "c.client.BroadcastRoom(%s, reply) serverId:%s error(%v)", nil, roomArg, c.serverID, err)
			}
		case pushArg := <-pushChan:
			_, err := c.client.PushMsg(context.Background(), &comet.PushMsgReq{
				Keys:    pushArg.Keys,
				Proto:   pushArg.Proto,
				ProtoOp: pushArg.ProtoOp,
			})
			if err != nil {
				log.Error("Comet", "process", "c.client.PushMsg(%s, reply) serverId:%s error(%v)", nil, pushArg, c.serverID, err)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Close close the resouces.
func (c *Comet) Close() (err error) {
	finish := make(chan bool)
	go func() {
		for {
			n := len(c.broadcastChan)
			for _, ch := range c.pushChan {
				n += len(ch)
			}
			for _, ch := range c.roomChan {
				n += len(ch)
			}
			if n == 0 {
				finish <- true
				return
			}
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-finish:
		log.Info("Comet", "Close", "close comet finish")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("close comet(server:%s push:%d room:%d broadcast:%d) timeout", c.serverID, len(c.pushChan), len(c.roomChan), len(c.broadcastChan))
	}
	c.cancel()
	return
}
