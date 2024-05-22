package routers

// CS 客户端-服务端
// SC 服务端-客户端
const (
	CS_LOGIN              = iota + 1001
	CS_ENTER_PUBLIC_ROOM  // 自动进房
	SC_ENTER_PUBLIC_ROOM  // 自动进房回复
	CS_CREATE_PUBLIC_ROOM // 建房
	SC_BUILD_FISH_REPLY
	SC_BUILD_FISHARRAY_REPLY
	SC_FISH_ARRAY // 鱼群
)

/*
连上socket之后 自动进入房间 进房需要发送一个token，用于查询用户信息
如果当前房间列表里没有空房间 则自动创建新的房间
 （应该是HTTP请求：进房，准备，连socket，找到这个客户端对应的账户）



房间负责管理

进入房间后，每隔一段时间给客服端发送当前房间的状态（鱼的状态）
*/

// func EnterPublicRoom(ctx node.Context) {

// }
// func FindOrAddClient(ctx node.Context) {
// 	// 根据uid来找
// 	if _, ok := RoomMgrIns.Clients[ctx.UID()]; ok {

// 	} else {
// 		RoomMgrIns.Clients[ctx.UID()] = &Client{
// 			GID: ctx.GID(),
// 			UID: ctx.UID(),
// 			NID: ctx.NID(),
// 			CID: ctx.CID(),
// 		}
// 	}
// }
