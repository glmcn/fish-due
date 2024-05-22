package handler

import (
	"fishserver/cluster/node/routers"
	"fishserver/cluster/node/service"

	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
)

func InitListen(proxy *node.Proxy) {
	proxy.AddEventHandler(cluster.Connect, connectHandler)
	// proxy.AddRouteHandler(CS_CREATE_PUBLIC_ROOM, false, CreatePublicRoom)
	proxy.AddRouteHandler(routers.CS_ENTER_PUBLIC_ROOM, false, EnterPublicRoom)
	// proxy.Broadcast()
	// RoomMgrIns.Proxy = proxy
	service.GetRoomMgr().SetProxy(proxy)
}

func connectHandler(ctx node.Context) {
	// doPushMessage(conn)
	// c := &Client{
	// 	GID: ctx.GID(),
	// 	UID: ctx.UID(),
	// 	NID: ctx.NID(),
	// 	CID: ctx.CID(),
	// }

	// RoomMgrIns.FindOrAddClient(ctx)

}
