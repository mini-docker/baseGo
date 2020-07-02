package comet

import (
	"fecho/utility"
	"sync"
)

// Room is a room and store channel room info.
type Room struct {
	ID    string
	rLock sync.RWMutex
	Users map[int64]*User
	drop  bool
	//Online    int32 // dirty read is ok
	//AllOnline int32
	MIds []int64
	AIds []int64
}
type User struct {
	Key string
}

// NewRoom new a room struct, store channel room info.
func NewRoom(id string) (r *Room) {
	r = new(Room)
	r.ID = id
	r.drop = false
	r.Users = make(map[int64]*User, 0)
	//r.Online = 0
	return
}

// Put put channel into the room.
func (r *Room) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	//if !r.drop {
	//rch, ok := r.chs[ch.Mid]
	//if !ok || nil == rch {
	r.Users[ch.Mid] = &User{
		Key: ch.Key,
	}

	//}

	mids := append(r.MIds, ch.Mid)
	mids = utility.IntSliceDeduplicationInt64(mids)
	r.MIds = mids

	r.rLock.Unlock()
	return
}

// Del delete channel from the room.
func (r *Room) Del(ch *Channel) bool {
	r.rLock.Lock()
	delete(r.Users, ch.Mid)

	// 去掉mid
	var ids []int64
	for _, v := range r.MIds {
		if v != ch.Mid {
			ids = append(ids, v)
		}
	}

	r.MIds = ids

	//r.drop = (len(r.chs) == 0)
	r.rLock.Unlock()
	return r.drop
}
