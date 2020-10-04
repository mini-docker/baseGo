package logic

import (
	registry_module "baseGo/src/fecho/registry/registry-module"
	"fmt"
	"sync"
	"time"

	"github.com/jasonlvhit/gocron"
	"google.golang.org/grpc/connectivity"

	comet "baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/internal/logic/conf"
	"baseGo/src/imserver/internal/logic/dao"
	"baseGo/src/imserver/pkg/etcd"
	log "fecho/golog"
	"strings"

	"github.com/bilibili/discovery/naming"
	cluster "github.com/bsm/sarama-cluster"
	igrpc "google.golang.org/grpc"
)

// Logic struct
type Logic struct {
	c *conf.Config
	//dis             *naming.Discovery
	offlineConsumer *cluster.Consumer
	//kconsumer kafka.Consumer
	dao *dao.Dao
	// online
	totalIPs   int64
	totalConns int64
	roomCount  map[string]int32
	// load balancer
	nodes        []*naming.Instance
	loadBalancer *LoadBalancer
	regions      map[string]string // province -> region
	lock         sync.Mutex
}

// New init
func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		c:               c,
		dao:             dao.New(c),
		offlineConsumer: newOffLineKafkaSub(c.Kafka),
		loadBalancer:    NewLoadBalancer(),
	}
	return l
}

var (
	cometGrpcAddr map[string]string
	cometClients  map[string]*igrpc.ClientConn
)

func UpdateCometConn() {
	l := new(Logic)
	gocron.Every(3).Seconds().Do(l.updateClient)
	gocron.NextRun()
	<-gocron.Start()
	s := gocron.NewScheduler()
	s.Every(3).Seconds().Do(l.updateClient)
	<-s.Start()
}

func (l *Logic) updateClient() {
	clients := make(map[string]*igrpc.ClientConn)
	if len(cometGrpcAddr) == 0 {
		cometGrpcAddr = make(map[string]string, 0)
	}
	if len(cometClients) == 0 {
		cometClients = make(map[string]*igrpc.ClientConn, 0)
	}
	if len(registry_module.CometGrpcAddr) > 0 {
		// 检查连接是否有变化
		for k, v := range registry_module.CometGrpcAddr {
			if addr, ok := cometGrpcAddr[k]; !ok {
				// 地址变化更新地址
				cometGrpcAddr[k] = v
				// 更新连接
				conn, err := igrpc.Dial(v, igrpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := comet.NewCometClient(conn)
				clients[k] = conn
			} else {
				if addr == v {
					if cometClients[k] != nil {
						clients[k] = cometClients[k]
						continue
					}
				}
				// 地址变化更新地址
				cometGrpcAddr[k] = v
				// 更新连接
				conn, err := igrpc.Dial(v, igrpc.WithInsecure())
				if err != nil {
					continue
				}
				//client := comet.NewCometClient(conn)
				clients[k] = conn
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
		cometClients = make(map[string]*igrpc.ClientConn)
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
}

func (l *Logic) GetCometClient() comet.CometClient {
	for _, v := range cometClients {
		client := comet.NewCometClient(v)
		return client
	}
	return nil
}

func (l *Logic) GetCometClinets() []comet.CometClient {
	clients := make([]comet.CometClient, 0)
	if len(cometClients) > 0 {
		for _, v := range cometClients {
			client := comet.NewCometClient(v)
			clients = append(clients, client)
		}
	}
	return clients
}

// Close close resources.
func (l *Logic) Close() {
	l.dao.Close()
}

func (l *Logic) initEtcd(c *conf.RegistryConfig) {
	m, err := etcd.OnWatch(strings.Split(strings.Trim(c.Addr, ","), ","), "/services")
	if err != nil {
		log.Error("logic", "initEtcd", "initEtcd error(%+v)", err)
	}
	for {
		m.Nodes.Range(func(k, v interface{}) bool {
			//log.Info( "logic", "initEtcd", "Etcd monitoring data ip->%v, ext->%v, key->%v", v.(*etcd.Node).Meta.IP, v.(*etcd.Node).Meta.Ext, v.(*etcd.Node).Key)
			return true
		})
		time.Sleep(time.Second * time.Duration(c.TTL) * time.Second)
	}
}

// 初始化离线消息数据的kafka
func newOffLineKafkaSub(c *conf.Kafka) *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(c.Brokers, "group-1", []string{"offlineMessage"}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

type OfflineMessageReq struct {
	Type             grpc.PushMsg_Type `protobuf:"varint,1,opt,name=type,proto3,enum=goim.logic.PushMsg_Type" json:"type,omitempty"` // 发送类型
	RoomId           int32             `json:"RoomId,omitempty"`                                                                     // 房间id
	Mid              int32             `json:"mid,omitempty"`                                                                        // 发送人id
	SendTime         int32             `json:"sendTime,omitempty"`                                                                   // 发送时间
	MsgType          int32             `json:"msgType,omitempty"`                                                                    // 消息类型
	Msg              []byte            `json:"Msg,omitempty"`                                                                        // 消息内容
	IsDel            int32             ` json:"isDel,omitempty"`                                                                     // 是否删除
	MsgId            int32             `json:"MsgId,omitempty"`                                                                      // 消息ID
	Speed            int32             `json:"speed,omitempty"`                                                                      // 禁言时间
	Server           string            `json:"server,omitempty"`
	ReceiverId       int32             `json:"receiverId,omitempty"`       // 接收人ID
	SenderName       string            `json:"senderName,omitempty"`       // 发送人昵称
	SenderHead       string            `json:"senderHead,omitempty"`       // 发送人头像
	SenderLevel      int64             `json:"senderLevel,omitempty"`      // 发送人会员等级
	SenderLevelPhoto string            `json:"senderLevelPhoto,omitempty"` // 发送人会员等级头像
	ReceiveType      int32             `json:"receiveType,omitempty"`      // 接收人类型
	SiteName         string            `json:"siteName,omitempty"`         // 站点名
	Account          string            `json:"account,omitempty"`          // 会员帐号
}

// Consume messages, watch signals
//func (l *Logic) Consume() {
//	for {
//		select {
//		case err := <-l.offlineConsumer.Errors():
//			log.Error("Job", "Consume", "consumer error(%v)", err)
//		case n := <-l.offlineConsumer.Notifications():
//			log.Info("Job", "Consume", "consumer rebalanced(%v)", n)
//		case msg, ok := <-l.offlineConsumer.Messages():
//			if !ok {
//				return
//			}
//			l.offlineConsumer.MarkOffset(msg, "")
//			// process push message
//			pushMsg := make([]OfflineMessageReq, 0)
//			if err := json.Unmarshal(msg.Value, &pushMsg); err != nil {
//				log.Error("Job", "Consume", "proto.Unmarshal(%v) error(%v)", nil, msg, err)
//				continue
//			}
//			item := make(map[int32]bool)
//			for _, v := range pushMsg {
//				if !item[v.MsgId] && v.MsgType < bo.IM_MSG_TYPE_NOTIFICSTION && v.Type == grpc.PushMsg_PUSH {
//					offInfo := new(po.ImOfflineMessage)
//					offInfo.Id = int(v.MsgId) // 直接强制同步ID 方便离线消息读取
//					offInfo.Message = string(v.Msg)
//					offInfo.SenderId = int(v.Mid)
//					offInfo.ReceiverId = int(v.ReceiverId)
//					offInfo.SendTime = int(v.SendTime)
//
//					if err := l.dao.OfflinePush(offInfo); err != nil {
//						log.Error("Job", "Consume", "j.push(%v) error(%v)", nil, pushMsg, err)
//					}
//					item[v.MsgId] = true
//				}
//
//			}
//		}
//	}
//}
