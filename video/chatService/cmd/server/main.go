package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"west2-video/chatService/configs"
	"west2-video/chatService/internal/biz"
	"west2-video/chatService/internal/data"
	http_server "west2-video/chatService/internal/server/http"
	"west2-video/chatService/internal/server/http/handler"
	"west2-video/chatService/internal/server/websocket"
)

func main() {
	// 1. 加载配置
	workDir, _ := os.Getwd()
	cfg, err := configs.LoadConfig(workDir + "/configs/config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 2. 初始化日志
	initLogger(cfg)

	// 3. 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		hlog.Fatalf("无法初始化数据库: %v", err)
	}

	// 4. 初始化Redis
	redisClient := data.NewRedisClient(cfg.Redis)
	defer redisClient.Close()

	// 5. 依赖注入：自底向上创建所有组件
	// - 数据层 (data)
	chatRepo := data.NewChatRepo(db)
	groupRepo := data.NewGroupRepo(db)

	// - 业务逻辑层 (biz)
	chatUseCase := biz.NewChatUseCase(chatRepo, groupRepo, redisClient)
	groupUseCase := biz.NewGroupUseCase(groupRepo, chatRepo)

	// - WebSocket 核心
	websocket.InitHub(chatUseCase, groupUseCase) // 初始化Hub
	go websocket.GlobalHub.Run()                 // 启动Hub

	// - 处理器层 (handler)
	wsHandler := websocket.NewWebSocketHandler(chatUseCase, groupUseCase)
	chatHandler := handler.NewChatHandler(chatUseCase)
	groupHandler := handler.NewGroupHandler(groupUseCase)

	// 6. 创建HTTP服务器
	h := server.New(server.WithHostPorts(cfg.Server.Addr))

	// 7. 注册路由
	http_server.RegisterRoutes(h, wsHandler, chatHandler, groupHandler)

	// 8. 启动服务器
	go func() {
		hlog.Infof("服务器正在运行于 %s", cfg.Server.Addr)
		if err := h.Run(); err != nil {
			hlog.Fatalf("无法启动服务器: %v", err)
		}
	}()

	// 9. 等待中断信号以实现优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 10. 优雅关闭
	hlog.Info("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		hlog.Errorf("服务器被强制关闭: %v", err)
	}

	hlog.Info("服务器已退出")
}

func initLogger(cfg *configs.Config) {
	// 此处可添加更复杂的日志配置，例如写入文件
	hlog.SetLevel(hlog.LevelInfo)
}

func initDB(cfg *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("无法获取 sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	return db, nil
}
