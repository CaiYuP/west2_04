package service

import (
	"context"
	"time"
	"videoService/internal/data"
	"videoService/internal/domain"
	v1 "west2-video/api/common/v1"
	videoPb "west2-video/api/video/v1"
	"west2-video/common/errs"
)

type VideoServiceService struct {
	videoDomain *domain.VideoDomain
	videoPb.UnimplementedVideoServiceServer
}

func NewVideoServiceService() *VideoServiceService {
	return &VideoServiceService{
		videoDomain: domain.NewVideoDomain(),
	}
}

func (s *VideoServiceService) Feed(ctx context.Context, req *videoPb.FeedRequest) (*videoPb.FeedReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.videoDomain.VerifySeed(c, req.LatestTime)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	videos, err := s.videoDomain.FindVideosAfterTime(c, req.LatestTime)
	items := data.CopierVideos(videos)
	rsp := &videoPb.FeedReply{
		Items: items,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) Publish(ctx context.Context, req *videoPb.PublishRequest) (*videoPb.PublishReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return &videoPb.PublishReply{}, nil
}
func (s *VideoServiceService) PublishList(ctx context.Context, req *videoPb.PublishListRequest) (*videoPb.PublishListReply, error) {
	return &videoPb.PublishListReply{}, nil
}
func (s *VideoServiceService) HotRanking(ctx context.Context, req *videoPb.HotRankingRequest) (*videoPb.HotRankingReply, error) {
	return &videoPb.HotRankingReply{}, nil
}
func (s *VideoServiceService) Search(ctx context.Context, req *videoPb.SearchRequest) (*videoPb.SearchReply, error) {
	return &videoPb.SearchReply{}, nil
}
