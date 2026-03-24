package service

import (
	"context"
	"go.uber.org/zap"
	"interactionService/biz"
	"interactionService/internal/dao"
	"interactionService/internal/data"
	"interactionService/internal/database"
	"interactionService/internal/database/tran"
	"interactionService/internal/domain"
	"strconv"
	"time"
	v1 "west2-video/api/common/v1"
	pb "west2-video/api/interaction/v1"
	"west2-video/common/errs"
	"west2-video/common/logs"
)

type InteractionServiceService struct {
	likeDomain     *domain.LikeDomain
	commentDomain  *domain.CommentDomain
	videoRpcDomain *domain.VideoRpcDomain
	tran           tran.Transaction
	pb.UnimplementedInteractionServiceServer
}

func NewInteractionServiceService() *InteractionServiceService {
	return &InteractionServiceService{
		likeDomain:     domain.NewLikeDomain(),
		videoRpcDomain: domain.NewVideoRpcDomain(),
		commentDomain:  domain.NewCommentDomain(),
		tran:           dao.NewTransactionDao(),
	}
}

func (s *InteractionServiceService) LikeAction(ctx context.Context, req *pb.LikeActionRequest) (*pb.LikeActionReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var berr *errs.BError
	isLike := req.ActionType == biz.Like
	isLikeComment := req.CommentId != biz.IsCommentVideo
	err := s.tran.Action(func(conn database.DbConn) error {
		if isLikeComment {
			berr = s.likeDomain.LikeComment(c, req.Id, req.CommentId, req.ActionType)
		} else {
			berr = s.likeDomain.LikeVideo(c, req.Id, req.VideoId, req.ActionType)
		}
		if berr != nil {
			logs.LG.Error("InteractionServiceService.LikeAction.LikeCommentOrLikeVideo error", zap.Error(berr))
			return errs.GrpcError(berr)
		}
		var e error
		if isLikeComment {
			berr = s.commentDomain.IncrLikeCount(ctx, req.CommentId, isLike)
		} else {
			e = s.videoRpcDomain.IncrLikeCount(ctx, req.VideoId, isLike)
		}
		if e != nil {
			logs.LG.Error("InteractionServiceService.LikeAction.commentDomain.IncrLikeCount error", zap.Error(e))
			return e
		}
		if berr != nil {
			logs.LG.Error("InteractionServiceService.LikeAction.LikeCommentOrLikeVideo error", zap.Error(berr))
			return errs.GrpcError(berr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	rsp := &pb.LikeActionReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *InteractionServiceService) LikeList(ctx context.Context, req *pb.LikeListRequest) (*pb.LikeListReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	vids, err := s.likeDomain.FindLikeList(c, req.UserId, int(req.Page.PageNum), int(req.Page.PageSize))
	if err != nil {
		logs.LG.Error("InteractionServiceService.LikeAction.LikeList error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	videos, e := s.videoRpcDomain.FindVideosByIds(c, vids)
	if e != nil {
		logs.LG.Error("InteractionServiceService.FindVideos error", zap.Error(e))
		return nil, e
	}
	rsp := &pb.LikeListReply{
		Videos: videos,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *InteractionServiceService) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (*pb.CommentActionReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var err *errs.BError
	if req.CommentId != biz.IsCommentVideo {
		err = s.commentDomain.CommentComment(c, req.Id, req.CommentId, req.Content)
	} else {
		err = s.commentDomain.CommentVideo(c, req.Id, req.VideoId, req.Content)
	}
	if err != nil {
		logs.LG.Error("InteractionServiceService.CommentCommentOrLikeVideo error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.CommentActionReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *InteractionServiceService) CommentList(ctx context.Context, req *pb.CommentListRequest) (*pb.CommentListReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var err *errs.BError
	items := make([]*data.Comment, 0)
	if req.CommentId != biz.IsCommentVideo {
		items, err = s.commentDomain.GetCommentComment(c, req.CommentId, int(req.Page.PageNum), int(req.Page.PageSize))
	} else {
		items, err = s.commentDomain.GetVideoComment(c, req.VideoId, int(req.Page.PageNum), int(req.Page.PageSize))
	}
	if err != nil {
		logs.LG.Error("InteractionServiceService.CommentList error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	comments := data.ConverToComment(items)
	rsp := &pb.CommentListReply{
		Comments: comments,
	}
	rsp.Base = v1.Success()
	return rsp, nil
}
func (s *InteractionServiceService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentReply, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var err *errs.BError
	if req.VideoId != biz.IsDeleteCommentComment {
		v, e := s.videoRpcDomain.FindVideosById(ctx, req.VideoId)
		if e != nil {
			logs.LG.Error("LikeDomain.DeleteCommentByVideoId error", zap.Error(e))
			return nil, e
		}
		authID, _ := strconv.ParseInt(v.UserId, 10, 64)
		err = s.commentDomain.DeleteCommentByVideoId(c, req.VideoId, req.Id, authID)
	} else {
		err = s.commentDomain.DeleteCommentByCommentId(c, req.CommentId, req.Id)
	}
	if err != nil {
		logs.LG.Error("InteractionServiceService.DeleteComment error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	rsp := &pb.DeleteCommentReply{}
	rsp.Base = v1.Success()
	return rsp, nil
}
