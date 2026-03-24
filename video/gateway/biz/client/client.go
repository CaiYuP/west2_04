package client

import (
	"context"
	"fmt"
	"log"
	"time"

	etcdregistry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbinteraction "west2-video/api/interaction/v1"
	pbsocial "west2-video/api/social/v1"
	pbuser "west2-video/api/user/v1"
	pbvideo "west2-video/api/video/v1"
	"west2-video/gateway/config"
)

type ClientManager struct {
	UserClient        pbuser.UserServiceClient
	VideoClient       pbvideo.VideoServiceClient
	InteractionClient pbinteraction.InteractionServiceClient
	SocialClient      pbsocial.SocialServiceClient

	conns []*grpc.ClientConn
}

var globalClientManager *ClientManager
var etcdClient *clientv3.Client

func initEtcdClient() (*clientv3.Client, error) {
	// 创建 etcd 客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.C.Etcd.Addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 etcd 客户端失败: %w", err)
	}

	return client, nil
}

func InitClients(cfg *ServiceConfig) (*ClientManager, error) {
	// 初始化 etcd 客户端
	var err error
	etcdClient, err = initEtcdClient()
	if err != nil {
		return nil, fmt.Errorf("初始化 etcd 客户端失败: %w", err)
	}

	manager := &ClientManager{
		conns: make([]*grpc.ClientConn, 0),
	}

	// 初始化用户服务客户端（使用服务发现）
	userConn, err := createConnectionWithDiscovery(cfg.UserServiceName)
	if err != nil {
		return nil, fmt.Errorf("连接用户服务失败: %w", err)
	}
	manager.UserClient = pbuser.NewUserServiceClient(userConn)
	manager.conns = append(manager.conns, userConn)

	// 初始化视频服务客户端（使用服务发现）
	videoConn, err := createConnectionWithDiscovery(cfg.VideoServiceName)
	if err != nil {
		return nil, fmt.Errorf("连接视频服务失败: %w", err)
	}
	manager.VideoClient = pbvideo.NewVideoServiceClient(videoConn)
	manager.conns = append(manager.conns, videoConn)

	// 初始化互动服务客户端（使用服务发现）
	interactionConn, err := createConnectionWithDiscovery(cfg.InteractionServiceName)
	if err != nil {
		return nil, fmt.Errorf("连接互动服务失败: %w", err)
	}
	manager.InteractionClient = pbinteraction.NewInteractionServiceClient(interactionConn)
	manager.conns = append(manager.conns, interactionConn)

	// 初始化社交服务客户端（使用服务发现）
	socialConn, err := createConnectionWithDiscovery(cfg.SocialServiceName)
	if err != nil {
		return nil, fmt.Errorf("连接社交服务失败: %w", err)
	}
	manager.SocialClient = pbsocial.NewSocialServiceClient(socialConn)
	manager.conns = append(manager.conns, socialConn)

	globalClientManager = manager
	log.Println("所有 gRPC 客户端初始化成功（使用 etcd 服务发现）")
	return manager, nil
}

// createConnectionWithDiscovery 使用 etcd 服务发现创建 gRPC 连接
func createConnectionWithDiscovery(serviceName string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建 etcd registry
	registry := etcdregistry.New(etcdClient)

	// 从 etcd 获取服务实例
	instances, err := registry.GetService(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("从 etcd 获取服务 %s 失败: %w", serviceName, err)
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("服务 %s 在 etcd 中未找到可用实例", serviceName)
	}

	// 使用第一个可用实例的地址（可以根据负载均衡策略选择）
	var target string
	for _, instance := range instances {
		if len(instance.Endpoints) > 0 {
			// 提取 gRPC 地址（格式：grpc://host:port）
			endpoint := instance.Endpoints[0]
			if len(endpoint) > 7 && endpoint[:7] == "grpc://" {
				target = endpoint[7:] // 移除 "grpc://" 前缀
				break
			}
		}
	}

	if target == "" {
		return nil, fmt.Errorf("服务 %s 没有有效的 gRPC 端点", serviceName)
	}

	// 创建 gRPC 连接
	conn, err := grpc.DialContext(ctx, target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("连接服务 %s 失败: %w", serviceName, err)
	}

	log.Printf("成功连接到服务: %s (通过 etcd 服务发现，地址: %s)", serviceName, target)
	return conn, nil
}

func GetClientManager() *ClientManager {
	if globalClientManager == nil {
		log.Fatal("客户端管理器未初始化")
	}
	return globalClientManager
}

func (m *ClientManager) Close() error {
	for _, conn := range m.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

type ServiceConfig struct {
	UserServiceName        string
	VideoServiceName       string
	InteractionServiceName string
	SocialServiceName      string
}
