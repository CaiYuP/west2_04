package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"socialService/config"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"socialService/internal/rpc"
	"socialService/internal/service"
	pb "west2-video/api/social/v1"
)

var (
	Name    = "social-service"
	Version = "v1.0.0"
	id, _   = os.Hostname()
)

func newApp(logger kratoslog.Logger, gs *grpc.Server, registrar registry.Registrar) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs),
	}

	if registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	return kratos.New(opts...)
}

func main() {
	flag.Parse()
	logger := kratoslog.NewStdLogger(os.Stdout)

	grpcAddr := config.C.Gc.Addr
	serviceName := config.C.Gc.Name

	// 创建 gRPC 服务器
	grpcSrv := grpc.NewServer(
		grpc.Address(grpcAddr),
		grpc.Middleware(
			recovery.Recovery(), // 恢复中间件，防止 panic
		),
	)

	// 创建服务实例
	videoService := service.NewSocialServiceService()

	// 注册 gRPC 服务
	pb.RegisterSocialServiceServer(grpcSrv, videoService)

	// 注册到 etcd
	var etcdRegistrar registry.Registrar
	var err error
	etcdRegistrar, err = rpc.RegisterEtcd(serviceName, Version, id)
	if err != nil {
		logger.Log(kratoslog.LevelError, "msg", fmt.Sprintf("注册服务到 etcd 失败: %v", err))
		// 不阻止服务启动，但记录错误
	} else {
		logger.Log(kratoslog.LevelInfo, "msg", fmt.Sprintf("服务已注册到 etcd: %s", serviceName))
	}

	// 创建 Kratos 应用
	app := newApp(logger, grpcSrv, etcdRegistrar)

	logger.Log(kratoslog.LevelInfo, "msg", fmt.Sprintf("视频服务启动在 %s", grpcAddr))

	go func() {
		if err := app.Run(); err != nil {
			logger.Log(kratoslog.LevelError, "msg", fmt.Sprintf("服务运行错误: %v", err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Log(kratoslog.LevelInfo, "msg", "收到关闭信号，开始优雅关闭...")

	if etcdRegistrar != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		serviceInstance := &registry.ServiceInstance{
			ID:   id,
			Name: serviceName,
		}
		if err := etcdRegistrar.Deregister(ctx, serviceInstance); err != nil {
			logger.Log(kratoslog.LevelError, "msg", fmt.Sprintf("从 etcd 注销服务失败: %v", err))
		} else {
			logger.Log(kratoslog.LevelInfo, "msg", "服务已从 etcd 注销")
		}
	}

	if err := app.Stop(); err != nil {
		logger.Log(kratoslog.LevelError, "msg", fmt.Sprintf("停止服务失败: %v", err))
	}
}
