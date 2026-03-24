package domain

import (
	"context"
	"go.uber.org/zap"
	"interactionService/biz"
	"interactionService/internal/dao"
	"interactionService/internal/data"
	"interactionService/internal/repo"
	"time"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

func NewLikeDomain() *LikeDomain {
	return &LikeDomain{
		likeRepo: dao.NewLikeDao(),
	}
}

type LikeDomain struct {
	likeRepo repo.LikeRepo
}

func (d *LikeDomain) LikeComment(ctx context.Context, uid int64, commentId int64, actionType int32) *errs.BError {
	var isLike bool
	var err error
	isLike = actionType == biz.Like
	if isLike {
		like := &data.Like{
			UserID:    uint64(uid),
			VideoID:   0, // 评论点赞时 video_id 设为 0
			CommentID: uint64(commentId),
			CreatedAt: time.Now(),
			IsComment: isLike,
		}

		// 先检查是否已经点过赞，避免重复插入
		var exist bool
		exist, err = d.likeRepo.IsLikeExist(ctx, like)
		if err != nil {
			logs.LG.Error("LikeDomain.LikeComment.IsLikeExist error", zap.Error(err))
			return model.DBError
		}
		if exist {
			// 已经点过赞，视为成功，直接返回（幂等）
			return model.IsAlreadyLike
		}

		// 未点赞时才真正创建记录
		err = d.likeRepo.CreateCommentLike(ctx, like)
	} else {
		err = d.likeRepo.DeleteCommentLike(ctx, uid, commentId)
	}
	if err != nil {
		logs.LG.Error("LikeDomain.LikeComment.likeRepo.CreateCommentLike error", zap.Error(err))
		return model.DBError
	}

	return nil
}

func (d *LikeDomain) LikeVideo(ctx context.Context, uid int64, vid int64, actionType int32) *errs.BError {
	var err error
	isLike := actionType == biz.Like
	if isLike {
		like := &data.Like{
			UserID:    uint64(uid),
			VideoID:   uint64(vid),
			CreatedAt: time.Now(),
			IsComment: false,
		}

		// 先检查是否已经点过赞，避免重复插入
		var exist bool
		exist, err = d.likeRepo.IsLikeExist(ctx, like)
		if err != nil {
			logs.LG.Error("LikeDomain.LikeVideo.IsLikeExist error", zap.Error(err))
			return model.DBError
		}
		if exist {
			return model.IsAlreadyLike
		}

		// 未点赞时才真正创建记录
		err = d.likeRepo.CreateVideoLike(ctx, like)
	} else {
		err = d.likeRepo.DeleteVideoLike(ctx, uid, vid)
	}
	if err != nil {
		logs.LG.Error("LikeDomain.LikeVideo.DeleteCommentLike error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *LikeDomain) FindLikeList(ctx context.Context, id int64, page int, size int) ([]int64, *errs.BError) {
	vids, err := d.likeRepo.FindLikeList(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("LikeDomain.FindLikeList error", zap.Error(err))
		return nil, model.DBError
	}
	return vids, nil
}
