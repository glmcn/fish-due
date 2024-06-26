package main

import (
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/network/tcp/v2"
	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/registry/consul/v2"
	"github.com/dobyte/due/transport/rpcx/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/gate"
	"github.com/dobyte/due/v2/flag"
	"github.com/dobyte/due/v2/network"
)

func main() {
	// 创建容器
	container := due.NewContainer()

	var server network.Server
	// 创建服务器
	switch flag.String("protocol", "ws") {
	case "tcp":
		server = tcp.NewServer()
	default:
		server = ws.NewServer()
	}
	// 创建用户定位器
	locator := redis.NewLocator()
	// 创建服务发现
	registry := consul.NewRegistry()
	// 创建RPC传输器
	transporter := rpcx.NewTransporter()
	// 创建网关组件
	component := gate.NewGate(
		gate.WithServer(server),
		gate.WithLocator(locator),
		gate.WithRegistry(registry),
		gate.WithTransporter(transporter),
	)
	// 添加网关组件
	container.Add(component)
	// 启动容器
	container.Serve()
}
