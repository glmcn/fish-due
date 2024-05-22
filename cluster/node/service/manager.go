package service

import (
	"fishserver/pkg/snowflake"
	"sync"
	"time"

	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
)

// 房间管理，后面可加入client一起维护
type RoomMgr struct {
	RoomLock   sync.Mutex
	RoomsInfo  map[int64]*RoomInfo //room暴露出去的信息和channel
	Rooms      map[int64]*Room
	RoomIdChan <-chan int64

	Clients map[int64]*Client

	proxy *node.Proxy
}

var (
	// RoomMgr->roomMgrIns
	roomMgrIns = &RoomMgr{
		RoomLock:   sync.Mutex{}, // 这个做成私有的，避免其他模块操作
		RoomsInfo:  make(map[int64]*RoomInfo),
		Rooms:      make(map[int64]*Room), //只在room的协程里操作
		RoomIdChan: make(<-chan int64),    //用于雪花算法生成ID 可以不用
	}
)

func GetRoomMgr() *RoomMgr {
	return roomMgrIns
}

func (c *RoomMgr) SetProxy(p *node.Proxy) {
	c.proxy = p
}

func (c *RoomMgr) FindOrAddClient(ctx node.Context) (ret *Client) {

	// 根据uid来找
	if p, ok := c.Clients[ctx.UID()]; ok {
		ret = p
	} else {
		ret = &Client{
			GID: ctx.GID(),
			UID: ctx.UID(),
			NID: ctx.NID(),
			CID: ctx.CID(),
		}
		c.Clients[ctx.UID()] = ret
	}

	return
}
func (r *RoomMgr) EnterPublicRoom(pl *Player) (roomId int64, err error) {

	userId := pl.UserId
	GetRoomMgr().RoomLock.Lock()
	log.Info("EnterPublicRoom get lock...")
	defer GetRoomMgr().RoomLock.Unlock()
	defer log.Info("EnterPublicRoom set free lock...")
	for RoomId, RoomInfo := range GetRoomMgr().RoomsInfo {
		for _, roomUserId := range RoomInfo.UserInfo {
			if userId == roomUserId { //已在房间内
				return
			}
		}
		// 已找到空房间
		if roomId == 0 && len(RoomInfo.UserInfo) < 4 && RoomInfo.BaseScore == pl.BaseScore {
			roomId = RoomId
			// break // 需要break吗
		}
	}

	if roomId == 0 { //房间全满，建新房
		roomId = GetRoomMgr().CreatePublicRoom(pl.BaseScore, userId)
	}
	cannonKindVip := map[int]int{
		0: 1,
		1: 4,
		2: 7,
		3: 10,
		4: 13,
		5: 16,
		6: 19,
	}
	pl.CannonKind = cannonKindVip[pl.Vip]

	if roomInfo, ok := GetRoomMgr().RoomsInfo[roomId]; ok {
		resChan := make(chan error)

		// 这里主要是让房间主协程去处理用户的进房操作
		roomInfo.HttpReqChan <- &HttpReqData{
			UserInfo: *pl, ErrChan: resChan,
		}
		timeOut := time.After(time.Second)
		select {
		case <-timeOut: //超时处理： 如果一秒钟没处理完就退出了
			return
		case ec := <-resChan:
			err = ec
			if err != nil {
				log.Error("EnterPublicRoom enter room [%d] err: %v", roomId, err)
			} else {
				exists := false
				for _, roomUserId := range GetRoomMgr().RoomsInfo[roomId].UserInfo {
					if roomUserId == userId {
						exists = true
					}
				}
				if !exists {
					roomInfo.UserInfo = append(GetRoomMgr().RoomsInfo[roomId].UserInfo, userId)
				}

				return
			}
		}
	}

	return
}

// 建新房（需要在外面加锁）
func (r *RoomMgr) CreatePublicRoom(baseScore int, userId int64) (roomId int64) {
	// roomId = RoomId(<-RoomMgr.RoomIdChan)
	roomId = snowflake.GetID()
	r.Rooms[roomId] = &Room{
		RoomId:     roomId,
		ActiveFish: make([]*Fish, 0),
		CreateTime: time.Now(),
		Users:      make(map[int64]*Player, 4),
		Conf: &RoomConf{
			BaseScore:    baseScore,
			MinHaveScore: MinHaveScore,
			MaxHaveScore: MaxHaveScore,
			TaxRatio:     TaxRatio,
			Creator:      userId,
		},
		FrozenEndTime: time.Time{},
		//FormationEndTime: time.Time{},
		Status:      GameStatusWaitBegin,
		AliveFish:   make(map[int64]*Fish),
		AliveBullet: make(map[string]*Bullet),
		Utils: &FishUtil{
			//ActiveFish: make([]*Fish, 0),
			//Lock:       sync.Mutex{},
			BuildFishChan:    make(chan *Fish, 10),
			StopBuildFish:    make(chan bool),    //暂停出鱼
			RestartBuildFish: make(chan bool),    //重新开始出鱼
			Exit:             make(chan bool, 1), //接收信号
		},
		fishArrayEndTimer: make(<-chan time.Time),
		frozenEndTimer:    make(<-chan time.Time),
		Exit:              make(chan bool, 1),
		ClientReqChan:     make(chan *clientReqData),
		HttpReqChan:       make(chan *HttpReqData),
	}
	r.RoomsInfo[roomId] = &RoomInfo{
		UserInfo: make([]int64, 0),
		//ClientReqChan: make(chan *clientReqData),
		HttpReqChan: r.Rooms[roomId].HttpReqChan,
		BaseScore:   baseScore,
	}

	r.Rooms[roomId].begin()
	return
}
