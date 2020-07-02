package comet

import (
	"fecho/golog"
	"sync"
	"sync/atomic"

	"baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/imserver/internal/comet/conf"
)

// Bucket is a channel holder.
type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex        // protect the channels for chs
	chs   map[string]*Channel // map sub key to a channel
	// room
	rooms       map[string]*Room // bucket room channels
	routines    []chan *grpc.BroadcastRoomReq
	routinesNum uint64

	ipCnts map[string]int32
}

//创建Bucket实体
// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.ipCnts = make(map[string]int32)
	b.c = c
	b.rooms = make(map[string]*Room, c.Room)
	//为广播房间数创建对应数量的gorounite
	b.routines = make([]chan *grpc.BroadcastRoomReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan *grpc.BroadcastRoomReq, c.RoutineSize)
		b.routines[i] = c
		go b.roomproc(c)
	}
	return
}

// ChannelCount channel count in the bucket
//bucket中Channel个数
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// ChangeRoom change ro room
func (b *Bucket) ChangeRoom(nrid string, ch *Channel) (err error) {
	var (
		nroom *Room
		ok    bool
		oroom *Room
	)
	if nil != ch.Room {
		oroom = ch.Room
	}
	// change to no room
	if nrid == "" {
		if oroom != nil && oroom.Del(ch) {
			//b.DelRoom(oroom)
		}
		ch.Room = nil
		return
	}
	b.cLock.Lock()
	if nroom, ok = b.rooms[nrid]; !ok {
		nroom = NewRoom(nrid)
		b.rooms[nrid] = nroom
	}
	b.cLock.Unlock()
	if oroom != nil && oroom.Del(ch) {
		//b.DelRoom(oroom)
	}

	if err = nroom.Put(ch); err != nil {
		golog.Error("Bucket", "ChangeRoom", "error(%v)", nil, err)
		return
	}
	ch.Room = nroom
	return
}

// Put put a channel according with sub key.
//sub key为key，ch为channel，rid是Room号，一并初始化并put到Bucket中
func (b *Bucket) Put(rids []string, ch *Channel) (err error) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	for _, rid := range rids {

		var (
			room *Room
			ok   bool
		)

		if rid != "" {
			if room, ok = b.rooms[rid]; !ok {
				room = NewRoom(rid)
				b.rooms[rid] = room
			}
			ch.Room = room
		}

		if room != nil {
			err = room.Put(ch)
		}
	}

	b.ipCnts[ch.IP]++
	return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(dch *Channel) {
	var (
		ok   bool
		ch   *Channel
		room *Room
	)
	b.cLock.Lock()
	if ch, ok = b.chs[dch.Key]; ok {
		room = ch.Room
		if ch == dch {
			delete(b.chs, ch.Key)
		}
		// ip counter
		if b.ipCnts[ch.IP] > 1 {
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.cLock.Unlock()
	if room != nil && room.Del(ch) {
		// if empty room, must delete from bucket
		//b.DelRoom(room)
	}
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

// Broadcast push msgs to all channels in the bucket.
func (b *Bucket) Broadcast(p *grpc.Proto, op int32) {
	var ch *Channel
	b.cLock.RLock()
	for _, ch = range b.chs {
		//if !ch.NeedPush(op) {
		//	continue
		//}
		_ = ch.Push(p)
	}
	//fmt.Println("-----------------广播消息：", string(p.Body))
	b.cLock.RUnlock()
}

// Room get a room by roomid.
func (b *Bucket) Room(rid string) (room *Room) {
	b.cLock.RLock()
	room = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

// BroadcastRoom broadcast a message to specified room
func (b *Bucket) BroadcastRoom(arg *grpc.BroadcastRoomReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// Rooms get all room id where online number > 0.
func (b *Bucket) Rooms() (res map[string]struct{}) {
	var (
		roomID string
		room   *Room
	)
	res = make(map[string]struct{})
	b.cLock.RLock()
	for roomID, room = range b.rooms {
		if len(room.Users) > 0 {
			res[roomID] = struct{}{}
		}
	}
	b.cLock.RUnlock()
	return
}

// roomproc
func (b *Bucket) roomproc(c chan *grpc.BroadcastRoomReq) {
	for {
		arg := <-c
		var push_mids []string
		if room := b.Room(arg.RoomID); room != nil {
			//golog.Info( "bucket", "roomproc", "push room : ", arg.RoomID)
			for _, user := range room.Users {
				ch, ok := b.chs[user.Key]
				if ok && nil != ch {
					push_mids = append(push_mids, user.Key)
					ch.Push(arg.Proto)
				}
			}
		}
		//golog.Info( "room", "Push", "push mids :%s", push_mids)
	}
}

// LeaveRoom leave a room
func (b *Bucket) LeaveRoom(ch *Channel) (err error) {
	var (
		oroom *Room
	)
	if nil == ch {
		return
	}
	if nil != ch.Room {
		oroom = ch.Room
	}
	// change to no room
	if oroom != nil && oroom.Del(ch) {
		//b.DelRoom(oroom)
	}
	ch.Room = nil
	return
	b.cLock.Lock()

	b.cLock.Unlock()
	if oroom != nil && oroom.Del(ch) {
		//b.DelRoom(oroom)
	}

	return
}

func (b *Bucket) RoomUserCount() (res map[string]int32) {

	res = make(map[string]int32)
	b.cLock.RLock()
	var not_in_room_ch []int64
	var in_room_ch_user []int64
	for _, chs := range b.chs {
		if nil != chs.Room {
			res[chs.Room.ID] += 1
			in_room_ch_user = append(in_room_ch_user, chs.Mid)
		} else {
			not_in_room_ch = append(not_in_room_ch, chs.Mid)
		}

	}
	//golog.Info( "bucket", "RoomUserCount", "in room user ch :%s", in_room_ch_user)
	//golog.Info( "bucket", "RoomUserCount", "in room admin ch :%s", in_room_ch_admin)
	//golog.Info( "bucket", "RoomUserCount", "not in room ch :%s", not_in_room_ch)

	b.cLock.RUnlock()
	return
}
func (b *Bucket) RefreshRoom(rids []string, ch *Channel) (err error) {
	b.cLock.Lock()
	defer b.cLock.Unlock()

	for _, rid := range rids {

		var (
			room *Room
			ok   bool
		)

		if rid != "" {
			if room, ok = b.rooms[rid]; !ok {
				room = NewRoom(rid)
				b.rooms[rid] = room
			}
			ch.Room = room
		}

		if room != nil {
			err = room.Put(ch)
		}
	}

	return
}
