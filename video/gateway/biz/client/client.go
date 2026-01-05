package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbinteraction "west2-video/api/interaction/v1"
	pbsocial "west2-video/api/social/v1"
	pbuser "west2-video/api/user/v1"
	pbvideo "west2-video/api/video/v1"
)

// ClientManager 管理所有 gRPC 客户端
type ClientManager struct {
	UserClient       pbuser.UserServiceClient
	VideoClient      pbvideo.VideoServiceClient
	InteractionClient pbinteraction.InteractionServiceClient
	SocialClient     pbsocial.SocialServiceClient
	
	conns []*grpc.ClientConn
}

var globalClientManager *ClientManager

// InitClients 初始化所有 gRPC 客户端
func InitClients(cfg *ServiceConfig) (*ClientManager, error) {
	manager := &ClientManager{
		conns: make([]*grpc.ClientConn, 0),
	}

	// 初始化用户服务客户端
	userConn, err := createConnection(cfg.UserServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("连接用户服务失败: %w", err)
	}
	manager.UserClient = pbuser.NewUserServiceClient(userConn)
	manager.conns = append(manager.conns, userConn)

	// 初始化视频服务客户端
	videoConn, err := createConnection(cfg.VideoServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("连接视频服务失败: %w", err)
	}
	manager.VideoClient = pbvideo.NewVideoServiceClient(videoConn)
	manager.conns = append(manager.conns, videoConn)

	// 初始化互动服务客户端
	interactionConn, err := createConnection(cfg.InteractionServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("连接互动服务失败: %w", err)
	}
	manager.InteractionClient = pbinteraction.NewInteractionServiceClient(interactionConn)
	manager.conns = append(manager.conns, interactionConn)

	// 初始化社交服务客户端
	socialConn, err := createConnection(cfg.SocialServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("连接社交服务失败: %w", err)
	}
	manager.SocialClient = pbsocial.NewSocialServiceClient(socialConn)
	manager.conns = append(manager.conns, socialConn)

	globalClientManager = manager
	log.Println("所有 gRPC 客户端初始化成功")
	return manager, nil
}

// createConnection 创建 gRPC 连接
func createConnection(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// GetClientManager 获取全局客户端管理器
func GetClientManager() *ClientManager {
	if globalClientManager == nil {
		log.Fatal("客户端管理器未初始化")
	}
	return globalClientManager
}

// Close 关闭所有连接
func (m *ClientManager) Close() error {
	for _, conn := range m.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	UserServiceAddr       string
	VideoServiceAddr      string
	InteractionServiceAddr string
	SocialServiceAddr     string
}

