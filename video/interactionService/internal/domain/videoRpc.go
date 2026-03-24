package domain

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"interactionService/internal/rpc/conn"
	pbvideo "west2-video/api/video/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

type VideoRpcDomain struct {
	// 不直接存储 client，而是延迟获取
}

func NewVideoRpcDomain() *VideoRpcDomain {
	return &VideoRpcDomain{}
}

//延迟获取用户客户端
func (u *VideoRpcDomain) getVideoClient() (pbvideo.VideoServiceClient, error) {
	clientManager := conn.GetClientManager()
	if clientManager == nil {
		return nil, fmt.Errorf("客户端管理器未初始化，请稍后重试")
	}
	return clientManager.VideoClient, nil
}

func (u *VideoRpcDomain) FindVideosByIds(ctx context.Context, vids []int64) ([]*pbvideo.Video, error) {
	client, err := u.getVideoClient()
	if err != nil {
		logs.LG.Error("VideoRpcDomain.FindVideosByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.RpcError)
	}
	req := &pbvideo.FindVideosByIdsRequest{
		Ids: vids,
	}
	videosByIds, err := client.FindVideosByIds(ctx, req)
	if err != nil {
		logs.LG.Error("VideoRpcDomain.FindVideosByIds error", zap.Error(err))
		return nil, err
	}
	return videosByIds.Items, nil
}

func (u *VideoRpcDomain) IncrLikeCount(ctx context.Context, vid int64, isLike bool) error {
	client, err := u.getVideoClient()
	if err != nil {
		logs.LG.Error("VideoRpcDomain.FindVideosByIds error", zap.Error(err))
		return errs.GrpcError(model.RpcError)
	}
	req := &pbvideo.IncrLikeCountRequest{
		Ids:    vid,
		IsLike: isLike,
	}

	_, err = client.IncrLikeCount(ctx, req)
	if err != nil {
		logs.LG.Error("VideoRpcDomain.IncrLikeCount error", zap.Error(err))
		return err
	}
	return nil
}

func (u *VideoRpcDomain) FindVideosById(ctx context.Context, id int64) (*pbvideo.Video, error) {
	client, err := u.getVideoClient()
	if err != nil {
		logs.LG.Error("VideoRpcDomain.FindVideosByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.RpcError)
	}
	req := &pbvideo.FindVideosByIdRequest{
		Ids: id,
	}
	video, err := client.FindVideosById(ctx, req)
	if err != nil {
		logs.LG.Error("VideoRpcDomain.FindVideosById error", zap.Error(err))
		return nil, err
	}
	return video.Items, nil
}
