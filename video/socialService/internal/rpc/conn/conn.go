package conn

import (
	"context"
	"fmt"
	etcdregistry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
	"videoService/config"
	pbinteraction "west2-video/api/interaction/v1"
	pbuser "west2-video/api/user/v1"
	pbvideo "west2-video/api/video/v1"
)

type ClientManager struct {
	UserClient        pbuser.UserServiceClient
	InteractionClient pbinteraction.InteractionServiceClient
	VideoClient       pbvideo.VideoServiceClient

	conns []*grpc.ClientConn
}

var (
	GlobalClientManager *ClientManager
	initMutex           sync.RWMutex
	etcdClient          *clientv3.Client
	initAttempted       bool // 标记是否已经尝试过初始化
)

func initEtcdClient() (*clientv3.Client, error) {
	// 创建 etcd 客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.C.Ec.Addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 etcd 客户端失败: %w", err)
	}

	return client, nil
}

// InitClients 初始化客户端管理器（支持重试）
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

	// 初始化用户服务客户端（使用服务发现，支持重试）
	userConn, err := createConnectionWithDiscoveryWithRetry(cfg.UserServiceName, 3, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("初始化用户服务客户端失败: %w", err)
	}
	manager.UserClient = pbuser.NewUserServiceClient(userConn)
	manager.conns = append(manager.conns, userConn)

	// 初始化互动服务客户端（使用服务发现，支持重试）
	interactionConn, err := createConnectionWithDiscoveryWithRetry(cfg.InteractionServiceName, 3, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("初始化互动服务客户端失败: %w", err)
	}
	manager.InteractionClient = pbinteraction.NewInteractionServiceClient(interactionConn)
	manager.conns = append(manager.conns, interactionConn)

	// 初始化video服务客户端（使用服务发现，支持重试）
	videoConn, err := createConnectionWithDiscoveryWithRetry(cfg.VideoServiceName, 3, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("初始化video服务客户端失败: %w", err)
	}
	manager.VideoClient = pbvideo.NewVideoServiceClient(videoConn)
	manager.conns = append(manager.conns, videoConn)

	log.Println("所有 gRPC 客户端初始化成功（使用 etcd 服务发现）")
	return manager, nil
}

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

// （支持重试）
func createConnectionWithDiscoveryWithRetry(serviceName string, maxRetries int, retryInterval time.Duration) (*grpc.ClientConn, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		conn, err := createConnectionWithDiscovery(serviceName)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		if i < maxRetries-1 {
			log.Printf("连接服务 %s 失败，%v 后重试... (尝试 %d/%d)", serviceName, retryInterval, i+1, maxRetries)
			time.Sleep(retryInterval)
		}
	}
	return nil, fmt.Errorf("连接服务 %s 失败，已重试 %d 次: %w", serviceName, maxRetries, lastErr)
}

// GetClientManager 获取客户端管理器（单例模式，延迟初始化，支持重试）
func GetClientManager() *ClientManager {
	// 先尝试读取已初始化的管理器
	initMutex.RLock()
	manager := GlobalClientManager
	initMutex.RUnlock()

	if manager != nil {
		return manager
	}

	// 如果未初始化，尝试初始化（支持重试）
	initMutex.Lock()
	defer initMutex.Unlock()

	// 双重检查，避免并发时重复初始化
	if GlobalClientManager != nil {
		return GlobalClientManager
	}

	cfg := &ServiceConfig{
		UserServiceName:        config.C.ESc.UserService,
		InteractionServiceName: config.C.ESc.InteractionService,
		VideoServiceName:       config.C.ESc.VideoService,
	}

	var err error
	manager, err = InitClients(cfg)
	if err != nil {
		if !initAttempted {
			log.Printf("初始化客户端管理器失败: %v，将在下次使用时重试", err)
			initAttempted = true
		}
		return nil
	}

	GlobalClientManager = manager
	initAttempted = true
	log.Println("客户端管理器初始化成功")
	return manager
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
	InteractionServiceName string
	VideoServiceName       string
}
