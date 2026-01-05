package service

import (
	"context"

	pb "west2-video/api/social/v1"
)

type SocialServiceService struct {
	pb.UnimplementedSocialServiceServer
}

func NewSocialServiceService() *SocialServiceService {
	return &SocialServiceService{}
}

func (s *SocialServiceService) FollowAction(ctx context.Context, req *pb.FollowActionRequest) (*pb.FollowActionReply, error) {
    return &pb.FollowActionReply{}, nil
}
func (s *SocialServiceService) FollowList(ctx context.Context, req *pb.FollowListRequest) (*pb.FollowListReply, error) {
    return &pb.FollowListReply{}, nil
}
func (s *SocialServiceService) FollowerList(ctx context.Context, req *pb.FollowerListRequest) (*pb.FollowerListReply, error) {
    return &pb.FollowerListReply{}, nil
}
func (s *SocialServiceService) FriendList(ctx context.Context, req *pb.FriendListRequest) (*pb.FriendListReply, error) {
    return &pb.FriendListReply{}, nil
}
