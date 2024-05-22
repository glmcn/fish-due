package service

import (
	"fishserver/cluster/node/routers"
	"fmt"
	"math/rand"
	"time"

	"github.com/dobyte/due/v2/log"
)

type Room struct {
	RoomId        int64
	ActiveFish    []*Fish //待激活的鱼
	CreateTime    time.Time
	Users         map[int64]*Player
	Conf          *RoomConf
	FrozenEndTime time.Time
	//FormationEndTime  time.Time
	Status            int
	AliveFish         map[int64]*Fish
	AliveBullet       map[string]*Bullet
	Utils             *FishUtil
	fishArrayEndTimer <-chan time.Time
	frozenEndTimer    <-chan time.Time
	Exit              chan bool
	ClientReqChan     chan *clientReqData //todo 客户端的请求通过chan传递，省去加锁的写法.包括加入房间
	HttpReqChan       chan *HttpReqData
}

type Bullet struct {
	UserId     int64   `json:"userId"`
	ChairId    int     `json:"chairId"`
	BulletKind int     `json:"bulletKind"`
	BulletId   string  `json:"bulletId"`
	Angle      float64 `json:"angle"`
	Sign       string  `json:"sign"`
	LockFishId int64   `json:"lockFishId"`
}

type RoomConf struct {
	BaseScore    int   `json:"base_score"`     // 该房间的基础分数（子弹加呗按照该分数计算）
	MinHaveScore int   `json:"min_have_score"` // 最小携带金币 暂未使用
	MaxHaveScore int   `json:"max_have_score"` // 最大携带金币 暂未使用
	TaxRatio     int   `json:"tax_ratio"`      // 抽水 千分之x
	Creator      int64 `json:"creator"`
}

type RoomInfo struct {
	UserInfo    []int64
	HttpReqChan chan *HttpReqData
	BaseScore   int //差点忘了请求的房间类型要一致才能进入
}

type HttpReqData struct {
	UserInfo Player
	ErrChan  chan error
}

const (
	GameStatusWaitBegin = iota
	GameStatusFree
	GameStatusPlay
	GameStatusFormation
	GameStatusFrozen
)

// 普通鱼
type Fish struct {
	FishKind        int       `json:"fishKind"`
	Trace           [][]int   `json:"trace"`
	Speed           int       `json:"speed"`
	FishId          int64     `json:"fishId"`
	ActiveTime      time.Time `json:"-"`
	FrontActiveTime int64     `json:"activeTime"` //给客户端的时间
}

// 组合鱼
type ArrayFish struct {
	FishKind  int   `json:"fishKind"`
	TraceKind int   `json:"traceKind"`
	FishId    int64 `json:"fishId"`
	Speed     int   `json:"speed"`
}

func (room *Room) begin() {
	log.Debug("room %d begin", room.RoomId)
	buildNormalFishTicker := time.NewTicker(time.Second * 1)     //普通鱼每秒刷一次
	buildGroupFishTicker := time.NewTicker(time.Second * 5 * 60) //鱼群
	flushTimeOutFishTicker := time.NewTicker(time.Second * 5)    //清理过期鱼

	go func() {
		defer func() {
			// log.Trace("room %v exit...", room.RoomId)
			buildNormalFishTicker.Stop()
			buildGroupFishTicker.Stop()
			flushTimeOutFishTicker.Stop()
			room.Utils.Exit <- true
			go func() { //启动协程取数据，防止utils阻塞在出鱼阶段导致无法退出 :)
				for range room.Utils.BuildFishChan {
				}
			}()
			close(room.Exit)
			close(room.HttpReqChan)
			close(room.ClientReqChan)

			roomMgrIns.RoomLock.Lock()
			log.Info("exit room goroutine get lock...")
			defer roomMgrIns.RoomLock.Unlock()
			defer log.Info("exit room goroutine set free lock...")

			delete(roomMgrIns.Rooms, room.RoomId)
		}()
		//defer room.Wg.Done()
		for {
			// 每次循环只走一个case
			select {
			case <-buildNormalFishTicker.C:
				room.buildNormalFish()
			case <-buildGroupFishTicker.C:
				if room.Status != GameStatusFree {
					continue
				}
				room.Utils.StopBuildFish <- true
				room.AliveFish = make(map[int64]*Fish) //清理鱼
				room.buildFormation()
			case <-flushTimeOutFishTicker.C:
				now := time.Now()
				AliveFishCheck := make(map[int64]*Fish)
				for _, fish := range room.AliveFish {
					if now.Sub(fish.ActiveTime) < 60*2*time.Second {
						AliveFishCheck[fish.FishId] = fish
					}
				}
				room.AliveFish = AliveFishCheck
			case fish := <-room.Utils.BuildFishChan:
				room.ActiveFish = append(room.ActiveFish, fish)
			// case clientReq := <-room.ClientReqChan:
			// 这里就不需要了
			//logs.Debug("room [%d] receive client message %v", room.RoomId, clientReq.reqData)
			// handleUserRequest(clientReq)
			case httpReq := <-room.HttpReqChan: // 进房操作（
				httpReq.ErrChan <- room.EnterRoom(&httpReq.UserInfo)
				close(httpReq.ErrChan)
			case <-room.Exit:
				return
			case <-room.fishArrayEndTimer:
				room.Status = GameStatusFree
				//room.AliveFish = make(map[FishId]*Fish) //清理鱼,因为时间有可能不同步，所以结束也不清理鱼
				room.Utils.RestartBuildFish <- true
			case <-room.frozenEndTimer:
				room.Status = GameStatusFree
				room.Utils.RestartBuildFish <- true
			}
		}
	}()
}

func (room *Room) buildNormalFish() {
	if room.Status != GameStatusFree {
		return
	}
	newFish := make([]*Fish, 0)
	for _, fish := range room.ActiveFish {
		if _, ok := room.AliveFish[fish.FishId]; ok {
			continue
		}
		//if len(room.AliveFish) < 30 {
		room.AliveFish[fish.FishId] = fish
		newFish = append(newFish, fish)
		//} else {
		//	break
		//}
	}
	room.ActiveFish = make([]*Fish, 0)
	if len(newFish) > 0 {
		room.broadcast(routers.SC_BUILD_FISH_REPLY, newFish)
	}
}

func (room *Room) broadcast(route int32, data interface{}) {
	// if dataByte, err := json.Marshal(data); err != nil {
	// 	log.Error("broadcast [%v] json marshal err :%v ", data, err)
	// } else {
	// 	dataByte = append([]byte{'4', '2'}, dataByte...)
	// 	for _, userInfo := range room.Users {
	// 		if userInfo.client != nil {
	// 			userInfo.client.SendMsg(route,dataByte)
	// 		}
	// 	}
	// }

	for _, userInfo := range room.Users {
		if userInfo.client != nil {
			userInfo.client.SendMsg(route, data)
		}
	}
}

func (room *Room) EnterRoom(userInfo *Player) (err error) {
	log.Debug("user %d request enter room %v", userInfo.UserId, room.RoomId)

	userCount := len(room.Users)
	if userCount >= 4 {
		log.Error("enterRoom err: room [%v] is full", room.RoomId)
		return
	}
	seatIndex := -1
out:
	for i := 0; i < 4; i++ {
		for _, roomUserInfo := range room.Users {
			if roomUserInfo.SeatIndex == i {
				continue out
			}
		}
		seatIndex = i
		break
	}

	if seatIndex == -1 {
		return fmt.Errorf("enterRoom  roomId [%v] failed", room.RoomId)
	}
	userInfo.SeatIndex = seatIndex
	room.Users[userInfo.UserId] = userInfo
	return
}

func (room *Room) buildFormation() {
	if room.Status != GameStatusFree {
		room.Utils.RestartBuildFish <- true
		return
	}
	room.Status = GameStatusFormation
	fishArrayData := BuildFishArray()
	activeTime := time.Now()
	for _, fishArray := range fishArrayData.FishArray {
		for _, arrayFish := range fishArray {
			room.AliveFish[arrayFish.FishId] = &Fish{
				FishId:     arrayFish.FishId,
				FishKind:   arrayFish.FishKind,
				Speed:      0,
				ActiveTime: activeTime,
			}
		}
	}
	room.fishArrayEndTimer = time.After(fishArrayData.EndTime.Sub(time.Now()))
	room.broadcast(routers.SC_BUILD_FISHARRAY_REPLY, fishArrayData)
}

type FishArrayRet struct {
	FormationKind int            `json:"formationKind"`
	FishArray     [][]*ArrayFish `json:"fishArray"`
	EndTime       time.Time      `json:"-"`
	EndTimeStamp  int64          `json:"endTime"`
}

// 启动鱼阵
func BuildFishArray() (ret *FishArrayRet) {
	var fishId int64
	var generateFishId = func() int64 {
		fishId++
		return fishId
	}
	var fishArray = make([][]*ArrayFish, 0)
	var duration = 0
	//直线鱼阵
	var buildFormationLine = func() {
		duration = 60
		fishArray = append(fishArray, make([]*ArrayFish, 0))
		fishArray = append(fishArray, make([]*ArrayFish, 0))
		var kind = 14
		for i := 0; i < 30; i++ {
			kind = i/3 + 10
			fishArray[0] = append(fishArray[0], &ArrayFish{
				FishKind:  kind,
				TraceKind: 0,
				FishId:    generateFishId(),
				Speed:     0,
			})
			fishArray[1] = append(fishArray[1], &ArrayFish{
				FishKind:  kind,
				TraceKind: 0,
				FishId:    generateFishId(),
				Speed:     0,
			})
		}
	}

	//环形鱼阵
	var buildCircleGroupFish = func() {
		duration = 60
		kind, fishNum := 1, 20
		for i := 0; i < 10; i++ {
			kind += 2
			fishArray = append(fishArray, make([]*ArrayFish, 0))
			if i > 20 {
				fishNum = 10
			}
			for j := 0; j < fishNum; j++ {
				fishArray[i] = append(fishArray[i], &ArrayFish{
					FishKind:  kind,
					TraceKind: 0,
					FishId:    generateFishId(),
					Speed:     0,
				})
			}
		}
	}

	// 两个螺旋形数组
	var buildSpiralGroupFish = func() {
		duration = 60
		fishArray = append(fishArray, make([]*ArrayFish, 0))
		fishArray = append(fishArray, make([]*ArrayFish, 0))
		kind := 1
		for i := 1; i <= 30; i++ {
			kind = ((i-1)/10 + 1) * 5
			fishArray[0] = append(fishArray[0], &ArrayFish{
				FishKind:  kind,
				TraceKind: 0,
				FishId:    generateFishId(),
				Speed:     0,
			})
			fishArray[1] = append(fishArray[1], &ArrayFish{
				FishKind:  kind,
				TraceKind: 0,
				FishId:    generateFishId(),
				Speed:     0,
			})
		}
	}
	ret = &FishArrayRet{}
	ret.FormationKind = rand.Intn(3) + 1
	//ret.FormationKind = 1
	switch ret.FormationKind {
	case 1:
		buildFormationLine()
	case 2:
		buildCircleGroupFish()
	case 3:
		buildSpiralGroupFish()
	}
	ret.FishArray = fishArray
	ret.EndTime = time.Now().Add(time.Second * time.Duration(duration))
	ret.EndTimeStamp = ret.EndTime.Unix() * 1e3
	return
}
