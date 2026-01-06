package service

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"time"
	userData "userService/internal/data"
	"userService/internal/domain"
	v1 "west2-video/api/common/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"
	userpb "west2-video/api/user/v1"
)

type UserServiceService struct {
	userpb.UnimplementedUserServiceServer
	userDomain *domain.UserDomain
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{
		userDomain: domain.NewUserDomain(),
	}
}

func (s *UserServiceService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterReply, error) {

	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := s.userDomain.Register(c, req.Username, req.Password)
	if err != nil {
		logs.LG.Error("UserServiceService userDomain Register error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &userpb.RegisterReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *UserServiceService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginReply, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	u, err := s.userDomain.Login(c, req.Username, req.Password)
	if err != nil {
		logs.LG.Error("UserServiceService userDomain Login error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	token, err := s.userDomain.CreateToken(c, u.ID, u.Username, req.Ip)
	rsp := &userpb.LoginReply{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		AccessExpiresIn:  token.AccessExp,
		RefreshExpiresIn: token.RefreshExp,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *UserServiceService) Refresh(ctx context.Context, req *userpb.RefreshRequest) (*userpb.RefreshReply, error) {
	token, err := s.userDomain.RefreshToken(req.RefreshToken, req.Ip)
	if err != nil {
		logs.LG.Error("UserServiceService RefreshToken error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &userpb.RefreshReply{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		AccessExpiresIn:  token.AccessExp,
		RefreshExpiresIn: token.RefreshExp,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *UserServiceService) GetUserInfo(ctx context.Context, req *userpb.UserInfoRequest) (*userpb.UserInfoReply, error) {
	uid := req.UserId
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	userById, err := s.userDomain.FindUserById(c, uid)
	if err != nil {
		logs.LG.Error("UserServiceService userDomain FindUserById error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	u := userData.UserFormat(userById)
	rsp := &userpb.UserInfoReply{
		User: u,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}

func (s *UserServiceService) UploadAvatar(ctx context.Context, req *userpb.UploadAvatarRequest) (*userpb.UploadAvatarReply, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	id := req.Id
	err := s.userDomain.UpdateAvatar(c, id, req.Url)
	if err != nil {
		logs.LG.Error("UserServiceService userDomain UpdateAvatar error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &userpb.UploadAvatarReply{}
	rsp.Base = v1.Success()
	rsp.AvatarUrl = req.Url
	return rsp, nil
}
func (s *UserServiceService) GetMfaQrcode(ctx context.Context, req *userpb.GetMfaQrcodeRequest) (*userpb.GetMfaQrcodeReply, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	qcode, url, err := s.userDomain.CreateMfaQrcode(c, req.Id, req.Username)
	if err != nil {
		logs.LG.Error("UserServiceService CreateMfaQrcode error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	// 4. 持久化 secret（绑定前你可以先暂存在某张表，如 user_mfa_temp 或 user 表的 mfa_secret 字段）
	// 典型做法：根据 userId 更新用户的 mfa_secret（但此时 is_mfa_enabled=false）
	err = s.userDomain.SaveMFASecret(c, req.Id, qcode)
	rsp := &userpb.GetMfaQrcodeReply{
		Secret: qcode,
		Qrcode: url,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *UserServiceService) BindMfa(ctx context.Context, req *userpb.BindMfaRequest) (*userpb.BindMfaReply, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	secret, err := s.userDomain.FindSecretById(c, req.Id)
	if err != nil {
		logs.LG.Error("UserServiceService FindSecretById error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	err = s.userDomain.VerifySecret(c, secret, req.Secret, req.Code)
	if err != nil {
		logs.LG.Error("UserServiceService VerifySecret error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	err = s.userDomain.UpdateIsSecretEnabled(c, req.Id)
	if err != nil {
		logs.LG.Error("UserServiceService UpdateIsSecretEnabled error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &userpb.BindMfaReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *UserServiceService) SearchByImage(ctx context.Context, req *userpb.SearchByImageRequest) (*userpb.SearchByImageReply, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	idStr := strconv.FormatInt(req.Id, 10)
	url, err := s.userDomain.CreatePngUrl(c, req.Data, idStr)
	if err != nil {
		logs.LG.Error("UserServiceService CreatePngUrl error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &userpb.SearchByImageReply{
		Url: url,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
