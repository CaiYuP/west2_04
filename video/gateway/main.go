package main

import (
	"fmt"
	"log"
	"west2-video/gateway/biz/router/routerRepo"

	"github.com/cloudwego/hertz/pkg/app/server"
	"west2-video/gateway/biz/client"
	_ "west2-video/gateway/biz/router"
	"west2-video/gateway/config"
)

func main() {

	// 初始化 gRPC 客户端（使用 etcd 服务发现）
	serviceCfg := &client.ServiceConfig{
		UserServiceName:        config.C.Services.UserService,
		VideoServiceName:       config.C.Services.VideoService,
		InteractionServiceName: config.C.Services.InteractionService,
		SocialServiceName:      config.C.Services.SocialService,
	}

	_, err := client.InitClients(serviceCfg)
	if err != nil {
		log.Fatalf("初始化 gRPC 客户端失败: %v", err)
	}
	defer client.GetClientManager().Close()

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", config.C.Server.Host, config.C.Server.Port)
	h := server.Default(server.WithHostPorts(addr))

	// 注册路由
	routerRepo.InitRouters(h)

	log.Printf("Gateway 服务启动在 %s", addr)
	h.Spin()
}
