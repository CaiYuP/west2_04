package main

import (
	"flag"
	"fmt"
	"os"
	"userService/config"

	"github.com/go-kratos/kratos/v2"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"userService/internal/service"
	pb "west2-video/api/user/v1"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "user-service"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	id, _ = os.Hostname()
)

func newApp(logger kratoslog.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := kratoslog.NewStdLogger(os.Stdout)

	grpcAddr := config.C.Gc.Addr

	// 创建 gRPC 服务器
	grpcSrv := grpc.NewServer(
		grpc.Address(grpcAddr),
		grpc.Middleware(
			recovery.Recovery(), // 恢复中间件，防止 panic
		),
	)

	// 创建服务实例
	userService := service.NewUserServiceService()

	// 注册 gRPC 服务
	pb.RegisterUserServiceServer(grpcSrv, userService)

	// 创建 Kratos 应用
	app := newApp(logger, grpcSrv)

	// 启动服务
	logger.Log(kratoslog.LevelInfo, "msg", fmt.Sprintf("用户服务启动在 %s", grpcAddr))
	if err := app.Run(); err != nil {
		panic(err)
	}
}
