package handler

import (
	"fishserver/cluster/node/service"
	"fishserver/pkg/panictool"

	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
)

type CSUserFireParam struct {
	UserId     string  `json:"userId"` //int64->string
	ChairId    int     `json:"chairId"`
	BulletKind int     `json:"bulletKind"`
	BulletId   int     `json:"bulletId"` // int64->string 由客户端生成
	Angle      float64 `json:"angle"`
	Sign       string  `json:"sign"`
	LockFishId string  `json:"lockFishId"` //int64->string
}

func UserFire(ctx node.Context) {

	req := &CSUserFireParam{}
	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
	}

	cl := service.GetRoomMgr().FindOrAddClient(ctx)
	if cl.UserInfo == nil {
		log.Errorf("UserFire: client is not bind up with player")
		return
	}

	bullet := service.Bullet{
		UserId:     panictool.GetInt64(req.UserId),
		ChairId:    req.ChairId,
		BulletKind: req.BulletKind,
		BulletId:   panictool.GetString(req.BulletId),
		Angle:      req.Angle,
	}

	cl.Fire(&bullet)
}
