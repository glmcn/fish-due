package handler

import (
	"fishserver/cluster/node/config"
	"fishserver/cluster/node/routers"
	"fishserver/cluster/node/service"
	"time"

	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
)

type CSEnterPublicRoomParam struct {
	RoomConfId int64
	Token      string
	Sign       string
}
type CSEnterPublicRoomResp struct {
	Errcode int
	ErrMsg  string
	Ip      string
	Port    int
	Mark    int

	Time  int64
	Sign  string
	Token string

	RoomId int64
}

func EnterPublicRoom(ctx node.Context) {
	req := &CSEnterPublicRoomParam{}

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
	}

	// 1. 根据token查询用户（通过RPC查询）
	// 游客方式进入
	pl := service.CreateGuestPlayer()
	token := req.Token
	sign := req.Sign

	// 1. 绑定用户和客户端对象
	cl := service.GetRoomMgr().FindOrAddClient(ctx)
	if cl.UserInfo != nil {
		return
	}
	cl.UserInfo = pl
	pl.SetClient(cl)

	var ret CSEnterPublicRoomResp
	defer cl.SendMsg(routers.SC_ENTER_PUBLIC_ROOM, ret)

	roomId, err := service.GetRoomMgr().EnterPublicRoom(pl)
	if err != nil {
		ret.Errcode = 0
		ret.ErrMsg = "ok"
		ret.Ip = config.GameConf.GameHost
		ret.Port = config.GameConf.GamePort
		ret.Sign = sign
		ret.RoomId = roomId
		ret.Time = time.Now().Unix() * 1000
		ret.Token = token
		ret.Mark = 1
	} else {
		ret.Errcode = 400
		ret.ErrMsg = "ok"
		ret.Ip = config.GameConf.GameHost
		ret.Port = config.GameConf.GamePort
		ret.Sign = sign
		ret.RoomId = roomId
		ret.Time = time.Now().Unix() * 1000
		ret.Token = token
		ret.Mark = 1
	}

}
