package comet

import (
	registry_module "fecho/registry/registry-module"
	"fmt"
	"math/rand"
	"time"

	"github.com/jasonlvhit/gocron"
	"google.golang.org/grpc/connectivity"

	//"github.com/zhenjl/cityhash"
	"google.golang.org/grpc"

	logicClient "baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/internal/comet/conf"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
)

var (
	logicGrpcAddr map[string]string
	logicClients  map[string]*grpc.ClientConn
)

func UpdateLogicConn() {
	gocron.Every(3).Seconds().Do(updateClient)
	gocron.NextRun()
	<-gocron.Start()
	s := gocron.NewScheduler()
	s.Every(3).Seconds().Do(updateClient)
	<-s.Start()
}

// Buckets return all buckets.
func (s *Server) Buckets() map[string]*Bucket {
	return s.buckets
}

func updateClient() {
	clients := make(map[string]*grpc.ClientConn)
	if len(logicGrpcAddr) == 0 {
		logicGrpcAddr = make(map[string]string, 0)
	}
	if len(logicClients) == 0 {
		logicClients = make(map[string]*grpc.ClientConn, 0)
	}
	if len(registry_module.LogicGrpcAddr) > 0 {
		// 检查连接是否有变化
		for k, v := range registry_module.LogicGrpcAddr {
			if addr, ok := logicGrpcAddr[k]; !ok {
				// 地址变化更新地址
				logicGrpcAddr[k] = v
				// 更新连接
				conn, err := grpc.Dial(v, grpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := logicClient.NewLogicClient(conn)
				clients[k] = conn
			} else {
				if addr == v {
					if logicClients[k] != nil {
						// 没有变化继续执行
						clients[k] = logicClients[k]
						continue
					}
				}
				// 地址变化更新地址
				logicGrpcAddr[k] = v
				// 更新连接
				conn, err := grpc.Dial(v, grpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := logicClient.NewLogicClient(conn)
				clients[k] = conn
			}
		}
	} else {
		//  清空缓存
		for k, old := range logicClients {
			if old.GetState() != connectivity.Ready {
				fmt.Println("---------清空过期连接")
				old.Close()
			} else {
				clients[k] = old
			}
		}
		logicClients = make(map[string]*grpc.ClientConn)
		logicGrpcAddr = make(map[string]string)
	}
	for key, old := range logicClients {
		if _, ok := clients[key]; !ok {
			if old.GetState() != connectivity.Ready {
				fmt.Println("---------清空过期连接")
				old.Close()
			} else {
				clients[key] = old
			}
		}
	}
	logicClients = clients
}

func (s *Server) GetLogicClient() logicClient.LogicClient {
	for _, v := range logicClients {
		client := logicClient.NewLogicClient(v)
		return client
	}
	return nil
}

// Server is comet server.
type Server struct {
	c         *conf.Config
	round     *Round             // accept round store
	buckets   map[string]*Bucket // subkey bucket
	bucketIdx uint32

	serverID string
}

// NewServer returns a new Server.
func NewServer(c *conf.Config) *Server {
	s := &Server{
		c:     c,
		round: NewRound(c),
	}
	// init bucket
	s.buckets = make(map[string]*Bucket, 0)
	s.bucketIdx = uint32(c.Bucket.Size)
	s.serverID = c.Env.Host
	return s
}

// Bucket get the bucket by subkey.
func (s *Server) Bucket(lineId, agencyId string) *Bucket {
	bck, ok := s.buckets[lineId+"-"+agencyId]
	if ok {
		return bck
	} else {
		return nil
	}

}

// RandServerHearbeat rand server heartbeat.
func (s *Server) RandServerHearbeat() time.Duration {
	return (minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat))))
}

// Close close the server.
func (s *Server) Close() (err error) {
	return
}
