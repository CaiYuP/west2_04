package service

import (
	"context"

	pb "west2-video/api/video/v1"
)

type VideoServiceService struct {
	pb.UnimplementedVideoServiceServer
}

func NewVideoServiceService() *VideoServiceService {
	return &VideoServiceService{}
}

func (s *VideoServiceService) Feed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedReply, error) {
    return &pb.FeedReply{}, nil
}
func (s *VideoServiceService) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishReply, error) {
    return &pb.PublishReply{}, nil
}
func (s *VideoServiceService) PublishList(ctx context.Context, req *pb.PublishListRequest) (*pb.PublishListReply, error) {
    return &pb.PublishListReply{}, nil
}
func (s *VideoServiceService) HotRanking(ctx context.Context, req *pb.HotRankingRequest) (*pb.HotRankingReply, error) {
    return &pb.HotRankingReply{}, nil
}
func (s *VideoServiceService) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchReply, error) {
    return &pb.SearchReply{}, nil
}
