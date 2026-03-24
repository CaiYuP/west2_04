package service

import (
	"context"
	"go.uber.org/zap"
	"socialService/biz"

	"socialService/internal/domain"
	"time"
	v1 "west2-video/api/common/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"

	pb "west2-video/api/social/v1"
)

type SocialServiceService struct {
	socialDomain  *domain.SocialDomain
	userRpcDomain *domain.UserRpcDomain
	pb.UnimplementedSocialServiceServer
}

func NewSocialServiceService() *SocialServiceService {
	return &SocialServiceService{
		socialDomain:  domain.NewSocialDomain(),
		userRpcDomain: domain.NewUserRpcDomain(),
	}
}

func (s *SocialServiceService) FollowAction(ctx context.Context, req *pb.FollowActionRequest) (*pb.FollowActionReply, error) {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var err *errs.BError
	if req.ActionType == biz.IsFollow {
		err = s.socialDomain.CreateFollow(c, req.Id, req.ToUserId)
	} else if req.ActionType == biz.IsCancelFollow {
		err = s.socialDomain.DeleteFollow(c, req.Id, req.ToUserId)
	}
	if err != nil {
		logs.LG.Error("SocialServiceService.FollowAction.CreateFollow.DeleteFollow error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.FollowActionReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *SocialServiceService) FollowList(ctx context.Context, req *pb.FollowListRequest) (*pb.FollowListReply, error) {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	fids, total, err := s.socialDomain.FindFollowListByFollowerId(c, req.UserId, int(req.Page.PageNum), int(req.Page.PageSize))
	if err != nil {
		logs.LG.Error("SocialServiceService.FollowList.FindFollowListByFollowerId error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.FollowListReply{}
	rsp.Base = v1.Success()
	if len(fids) == 0 {
		return rsp, nil
	}
	users, e := s.userRpcDomain.FindUserInfoByIds(c, fids)
	if e != nil {
		logs.LG.Error("SocialServiceService.FollowList.FindUserInfoByIds error", zap.Error(err))
		return nil, e
	}
	us := make([]*pb.SocialUser, len(users))
	for i, u := range users {
		us[i] = &pb.SocialUser{
			Id:        u.Id,
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
		}
	}
	rsp.Users = us
	rsp.Total = total
	return rsp, nil
}
func (s *SocialServiceService) FollowerList(ctx context.Context, req *pb.FollowerListRequest) (*pb.FollowerListReply, error) {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	fids, total, err := s.socialDomain.FindFollowListByFolloweeId(c, req.UserId, int(req.Page.PageNum), int(req.Page.PageSize))
	if err != nil {
		logs.LG.Error("SocialServiceService.FollowList.FindFollowListByFollowerId error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.FollowerListReply{}
	rsp.Base = v1.Success()
	if len(fids) == 0 {
		return rsp, nil
	}
	users, e := s.userRpcDomain.FindUserInfoByIds(c, fids)
	if e != nil {
		logs.LG.Error("SocialServiceService.FollowList.FindUserInfoByIds error", zap.Error(err))
		return nil, e
	}
	us := make([]*pb.SocialUser, len(users))
	for i, u := range users {
		us[i] = &pb.SocialUser{
			Id:        u.Id,
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
		}
	}
	rsp.Users = us
	rsp.Total = total
	return rsp, nil
}
func (s *SocialServiceService) FriendList(ctx context.Context, req *pb.FriendListRequest) (*pb.FriendListReply, error) {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	//必须关注和被关注才是好友
	fids, total, err := s.socialDomain.FindFriend(c, req.Id, int(req.Page.PageNum), int(req.Page.PageSize))
	if err != nil {
		logs.LG.Error("SocialServiceService.FriendList.FindFriendList error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.FriendListReply{}
	rsp.Base = v1.Success()
	if len(fids) == 0 {
		return rsp, nil
	}
	users, e := s.userRpcDomain.FindUserInfoByIds(c, fids)
	if e != nil {
		logs.LG.Error("SocialServiceService.FollowList.FindUserInfoByIds error", zap.Error(err))
		return nil, e
	}
	us := make([]*pb.SocialUser, len(users))
	for i, u := range users {
		us[i] = &pb.SocialUser{
			Id:        u.Id,
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
		}
	}
	rsp.Users = us
	rsp.Total = total
	return rsp, nil
}
