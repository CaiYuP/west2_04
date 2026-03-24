package domain

import (
	"context"
	"go.uber.org/zap"
	"interactionService/internal/dao"
	"interactionService/internal/data"
	"interactionService/internal/repo"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

func NewCommentDomain() *CommentDomain {
	return &CommentDomain{
		commentRepo: dao.NewCommentDao(),
	}
}

type CommentDomain struct {
	commentRepo repo.CommentRepo
}

func (d *CommentDomain) CommentComment(ctx context.Context, id int64, commentId int64, content string) *errs.BError {
	comment := &data.Comment{
		UserID:    uint64(id),
		ParentID:  uint64(commentId),
		Content:   content,
		IsComment: true,
	}
	err := d.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		logs.LG.Error("LikeDomain.CommentComment error", zap.Error(err))
		return model.DBError
	}
	err = d.commentRepo.IncryChildCount(ctx, commentId)
	if err != nil {
		logs.LG.Error("LikeDomain.CommentComment error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *CommentDomain) CommentVideo(ctx context.Context, id int64, vid int64, content string) *errs.BError {
	comment := &data.Comment{
		UserID:    uint64(id),
		VideoID:   uint64(vid),
		Content:   content,
		IsComment: false,
	}
	err := d.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		logs.LG.Error("LikeDomain.CommentComment error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *CommentDomain) GetCommentComment(ctx context.Context, id int64, page int, size int) ([]*data.Comment, *errs.BError) {
	comments, err := d.commentRepo.FindCommentComment(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("LikeDomain.GetCommentComment error", zap.Error(err))
		return nil, model.DBError
	}
	return comments, nil
}

func (d *CommentDomain) GetVideoComment(ctx context.Context, id int64, page int, size int) ([]*data.Comment, *errs.BError) {
	comments, err := d.commentRepo.FindVideoComment(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("LikeDomain.FindVideoComment error", zap.Error(err))
		return nil, model.DBError
	}
	return comments, nil
}

func (d *CommentDomain) DeleteCommentByVideoId(ctx context.Context, vid int64, uid, authId int64) *errs.BError {
	if uid != authId {
		return model.IsNotAuthor
	}
	err := d.commentRepo.DeleteCommentByVideoId(ctx, vid)
	if err != nil {
		logs.LG.Error("LikeDomain.DeleteCommentByVideoId error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *CommentDomain) DeleteCommentByCommentId(ctx context.Context, cid int64, uid int64) *errs.BError {
	userId, err := d.commentRepo.FindUserIdByCommentId(ctx, cid)
	if err != nil {
		logs.LG.Error("LikeDomain.DeleteCommentByCommentId error", zap.Error(err))
		return model.DBError
	}
	if userId != uid {
		return model.ConnotDeleteOtherComment
	}
	err = d.commentRepo.DeleteCommentByCommentId(ctx, cid, uid)
	if err != nil {
		logs.LG.Error("LikeDomain.DeleteCommentByCommentId error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *CommentDomain) IncrLikeCount(ctx context.Context, id int64, like bool) *errs.BError {
	err := d.commentRepo.IncrLikeCount(ctx, id, like)
	if err != nil {
		logs.LG.Error("LikeDomain.IncrLikeCount error", zap.Error(err))
		return model.DBError
	}
	return nil
}
