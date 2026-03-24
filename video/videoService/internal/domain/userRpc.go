package domain

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"videoService/internal/rpc/conn"
	pbuser "west2-video/api/user/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

type UserRpcDomain struct {
	// 不直接存储 client，而是延迟获取
}

func NewUserRpcDomain() *UserRpcDomain {
	return &UserRpcDomain{}
}

//延迟获取用户客户端
func (u *UserRpcDomain) getUserClient() (pbuser.UserServiceClient, error) {
	clientManager := conn.GetClientManager()
	if clientManager == nil {
		return nil, fmt.Errorf("客户端管理器未初始化，请稍后重试")
	}
	return clientManager.UserClient, nil
}
func (u *UserRpcDomain) FindUserInfoById(ctx context.Context, id int64) (*pbuser.User, error) {
	client, err := u.getUserClient()
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoById.getUserClient error", zap.Error(err))
		return nil, errs.GrpcError(model.RpcError)
	}

	info, err := client.GetUserInfo(ctx, &pbuser.UserInfoRequest{UserId: id})
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoById.GetUserInfo error", zap.Error(err))
		return nil, err
	}
	return info.User, nil
}

func (u *UserRpcDomain) FindUserInfoByIds(ctx context.Context, ids []int64) ([]*pbuser.User, error) {
	client, err := u.getUserClient()
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoByIds.getUserClient error", zap.Error(err))
		return nil, errs.GrpcError(model.RpcError)
	}

	info, err := client.GetUserInfos(ctx, &pbuser.UserInfosRequest{UserId: ids})
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoByIds.GetUserInfos error", zap.Error(err))
		return nil, err
	}
	return info.User, nil
}

func (u *UserRpcDomain) FindUserIdByUsername(ctx context.Context, username string) (int64, error) {
	client, err := u.getUserClient()
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoByIds.getUserClient error", zap.Error(err))
		return 0, errs.GrpcError(model.RpcError)
	}

	info, err := client.GetUserInfoByUserName(ctx, &pbuser.UserInfoUserNameRequest{Username: username})
	if err != nil {
		logs.LG.Error("UserRpcDomain.FindUserInfoByIds.GetUserInfoByUsername error", zap.Error(err))
		return 0, err
	}
	id, _ := strconv.ParseInt(info.User.Id, 10, 64)
	return id, nil
}
