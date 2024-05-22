package service

import (
	"context"
	"fishserver/pkg/snowflake"

	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/session"
)

type Client struct {
	// conn      *websocket.Conn
	UserInfo  *Player
	Room      *Room
	msgChan   chan []byte
	closeChan chan bool

	// GID 获取网关ID
	GID string
	// NID 获取节点ID
	NID string
	// CID 获取连接ID
	CID int64
	// UID 获取用户ID
	UID int64
}

// func (c *Client) sendMsg(msg []byte) {
// 	c.SendMsg(1, msg)
// 	if c.UserInfo != nil {
// 		//logs.Debug("user [%v] send msg %v", c.UserInfo.UserId, string(msg))
// 	}
// 	c.msgChan <- msg //为什么此处不担心发送数据到一个已关闭的chan？  因为channel没有手动关闭而是交给gc处理  :)
// }

func (c *Client) SendMsg(route int32, data interface{}) {
	roomMgrIns.proxy.Push(context.Background(), &cluster.PushArgs{
		GID:    c.GID,
		Kind:   session.User,
		Target: c.CID,
		Message: &cluster.Message{
			// Seq:   1,
			Route: route,
			Data:  data,
		},

		// 	Seq   int32       // 序列号
		// Route int32       // 路由ID
		// Data  interface{} // 消息数据，接收json、proto、[]byte
	})

}

// UserInfo->Player
type Player struct {
	// UserID int64
	UserId          int64   `json:"userId"`
	Score           int     `json:"-"`
	Bill            int     `json:"-"` //账单
	ConversionScore float64 `json:"score"`
	Name            string  `json:"name"`
	Ready           bool    `json:"ready"`
	SeatIndex       int     `json:"seatIndex"`
	Vip             int     `json:"vip"`
	CannonKind      int     `json:"cannonKind"`
	Power           float64 `json:"power"`
	LockFishId      int64   `json:"lockFishId"`
	Online          bool    `json:"online"`
	client          *Client `json:"-"`
	Ip              string  `json:"ip"`

	BaseScore int `json:"base_score"`
}

func (p *Player) SetClient(c *Client) {
	p.client = c
}

func CreateGuestPlayer() *Player {
	return &Player{
		UserId:    snowflake.GetID(),
		Score:     1000,
		BaseScore: 1,
	}
}

type clientReqData struct {
	client  *Client
	reqData []string
}
