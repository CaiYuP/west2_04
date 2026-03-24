package rpc

import (
	"context"
	"fmt"
	"time"

	"userService/config"

	etcdregistry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// RegisterEtcd 注册服务到 etcd
func RegisterEtcd(serviceName, version string, id string) (registry.Registrar, error) {
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.C.Ec.Addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 etcd 客户端失败: %w", err)
	}

	// 创建 etcd 注册器
	etcdRegistrar := etcdregistry.New(etcdClient)

	// 从配置获取 gRPC 地址和权重
	grpcAddr := config.C.Gc.Addr
	weight := config.C.Gc.Weight

	// 创建服务实例信息
	serviceInstance := &registry.ServiceInstance{
		ID:      id,
		Name:    serviceName,
		Version: version,
		Metadata: map[string]string{
			"weight": fmt.Sprintf("%d", weight),
		},
		Endpoints: []string{
			fmt.Sprintf("grpc://%s", grpcAddr),
		},
	}

	// 注册服务到 etcd
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := etcdRegistrar.Register(ctx, serviceInstance); err != nil {
		return nil, fmt.Errorf("注册服务到 etcd 失败: %w", err)
	}

	return etcdRegistrar, nil
}
