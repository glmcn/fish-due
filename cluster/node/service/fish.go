package service

type FishUtil struct {
	//ActiveFish []*Fish

	//Lock sync.Mutex
	CurrentFishId    int64
	BuildFishChan    chan *Fish
	StopBuildFish    chan bool //暂停出鱼
	RestartBuildFish chan bool //重新开始出鱼
	Exit             chan bool //接收信号
}
