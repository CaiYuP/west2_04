package service

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"time"
	"videoService/internal/data"
	"videoService/internal/domain"
	v1 "west2-video/api/common/v1"
	userpb "west2-video/api/user/v1"
	videoPb "west2-video/api/video/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"
)

type VideoServiceService struct {
	videoDomain   *domain.VideoDomain
	userRpcDomain *domain.UserRpcDomain
	videoPb.UnimplementedVideoServiceServer
}

func NewVideoServiceService() *VideoServiceService {
	return &VideoServiceService{
		videoDomain:   domain.NewVideoDomain(),
		userRpcDomain: domain.NewUserRpcDomain(),
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
	ids := make([]int64, 0)
	users := make(map[int64]*userpb.User, 0)
	for _, video := range videos {
		ids = append(ids, int64(video.UserID))
	}
	byIds, e := s.userRpcDomain.FindUserInfoByIds(c, ids)
	if e != nil {
		logs.LG.Error("VideoServiceService.PublishList.FindUserNameByIds error", zap.Error(err))
		return nil, e
	}
	for _, user := range byIds {
		uid, _ := strconv.ParseInt(user.Id, 10, 64)
		users[uid] = user
	}
	for _, video := range rsp.Items {
		uid, _ := strconv.ParseInt(video.UserId, 10, 64)
		video.Author = users[uid].Username
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) Publish(ctx context.Context, req *videoPb.PublishRequest) (*videoPb.PublishReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//1.先判断是否为视频，是的话上传到minio
	err := s.videoDomain.VerifyMp4(c, req.Data)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	idStr := strconv.FormatInt(req.Id, 10)
	url, err := s.videoDomain.ParseMinioUrl(c, req.Data, idStr)
	if err != nil {
		logs.LG.Error("VideoServiceService.Publish.ParseMinioUrl error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	//2.保存video
	err = s.videoDomain.CreateVideo(c, url, req.Id, req.Description, req.Title)
	if err != nil {
		logs.LG.Error("VideoServiceService.Publish.ParseMinioUrl error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &videoPb.PublishReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) PublishList(ctx context.Context, req *videoPb.PublishListRequest) (*videoPb.PublishListReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	videos, total, err := s.videoDomain.FindVideosByUserId(c, req.UserId, req.Page.PageSize, req.Page.PageNum)
	if err != nil {
		logs.LG.Error("VideoServiceService.PublishList.FindVideosByUserId error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &videoPb.PublishListReply{
		Total:     total,
		VideoList: data.CopierVideos(videos),
	}
	user, e := s.userRpcDomain.FindUserInfoById(c, req.UserId)
	if e != nil {
		logs.LG.Error("VideoServiceService.PublishList.FindUserInfoById error", zap.Error(e))
		return nil, err
	}
	userName := user.Username
	for _, v := range rsp.VideoList {
		v.Author = userName
	}

	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) HotRanking(ctx context.Context, req *videoPb.HotRankingRequest) (*videoPb.HotRankingReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ranking, total, err := s.videoDomain.GetHotRanking(c, int(req.Page.PageSize), int(req.Page.PageNum))
	if err != nil {
		logs.LG.Error("VideoServiceService.HotRanking.GetHotRanking error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	videos := data.CopierVideos(ranking)
	rsp := &videoPb.HotRankingReply{
		Total:     total,
		VideoList: videos,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) WatchVideo(ctx context.Context, req *videoPb.WatchVideoRequest) (*videoPb.WatchVideoReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//观看次数+1
	err := s.videoDomain.IncrementVisitCount(c, uint64(req.Id))
	if err != nil {
		logs.LG.Error("VideoServiceService.WatchVideo.IncrementVisitCount error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	videoById, err := s.videoDomain.FindVideosById(c, req.Id)
	if err != nil {
		logs.LG.Error("VideoServiceService.FindVideosById error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	video := data.CopierVideo(videoById)
	rsp := &videoPb.WatchVideoReply{
		Video: video,
		Base:  v1.Success(),
	}
	return rsp, nil
}
func (s *VideoServiceService) FindVideosByIds(ctx context.Context, req *videoPb.FindVideosByIdsRequest) (*videoPb.FindVideosByIdsReply, error) {
	videos, err := s.videoDomain.FindVideosByIds(ctx, req.Ids)
	if err != nil {
		logs.LG.Error("VideoServiceService.FindVideosByIds error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	copierVideos := data.CopierVideos(videos)
	rsp := &videoPb.FindVideosByIdsReply{
		Items: copierVideos,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) IncrLikeCount(ctx context.Context, req *videoPb.IncrLikeCountRequest) (*videoPb.IncrLikeCountReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.videoDomain.IncrLikeCount(c, req.Ids, req.IsLike)
	if err != nil {
		logs.LG.Error("VideoServiceService.IncrLikeCount error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &videoPb.IncrLikeCountReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) FindVideosById(ctx context.Context, req *videoPb.FindVideosByIdRequest) (*videoPb.FindVideosByIdReply, error) {
	video, err := s.videoDomain.FindVideosById(ctx, req.Ids)
	if err != nil {
		logs.LG.Error("VideoServiceService.FindVideosById error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	v := data.CopierVideo(video)
	rsp := &videoPb.FindVideosByIdReply{
		Items: v,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *VideoServiceService) Search(ctx context.Context, req *videoPb.SearchRequest) (*videoPb.SearchReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ft, tt, err := s.videoDomain.VerifyTime(c, req.FromDate, req.ToDate)
	if err != nil {
		logs.LG.Error("VideoServiceService.Search.VerifyTime error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	if req.ToDate == 0 {
		tt = time.Now()
	}
	userId := int64(0)
	var e error
	if req.Username != "" {
		userId, e = s.userRpcDomain.FindUserIdByUsername(c, req.Username)
		if e != nil {
			logs.LG.Error("VideoServiceService.Search.FindUserIdByUsername error", zap.Error(e))
			return nil, e
		}
	}
	itmes, total, err := s.videoDomain.FindVideosByTimeAndUserName(c, ft, tt, userId, req.Page.PageSize, req.Page.PageNum, req.Keywords)
	videos := data.CopierVideos(itmes)
	rsp := &videoPb.SearchReply{
		VideoList: videos,
		Total:     total,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
