package service

import (
	"context"

	pb "west2-video/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{}
}

func (s *UserServiceService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
    return &pb.RegisterReply{}, nil
}
func (s *UserServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
    return &pb.LoginReply{}, nil
}
func (s *UserServiceService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
    return &pb.UserInfoReply{}, nil
}
func (s *UserServiceService) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarReply, error) {
    return &pb.UploadAvatarReply{}, nil
}
func (s *UserServiceService) GetMfaQrcode(ctx context.Context, req *pb.GetMfaQrcodeRequest) (*pb.GetMfaQrcodeReply, error) {
    return &pb.GetMfaQrcodeReply{}, nil
}
func (s *UserServiceService) BindMfa(ctx context.Context, req *pb.BindMfaRequest) (*pb.BindMfaReply, error) {
    return &pb.BindMfaReply{}, nil
}
func (s *UserServiceService) SearchByImage(ctx context.Context, req *pb.SearchByImageRequest) (*pb.SearchByImageReply, error) {
    return &pb.SearchByImageReply{}, nil
}
