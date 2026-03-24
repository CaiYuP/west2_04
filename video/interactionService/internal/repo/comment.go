package repo

import (
	"context"
	"interactionService/internal/data"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, comment *data.Comment) error
	IncryChildCount(ctx context.Context, id int64) error
	IncrLikeCount(ctx context.Context, id int64, isLike bool) error
	FindCommentComment(ctx context.Context, id int64, page int, size int) ([]*data.Comment, error)
	FindVideoComment(ctx context.Context, id int64, page int, size int) ([]*data.Comment, error)
	DeleteCommentByVideoId(ctx context.Context, vid int64) error
	DeleteCommentByCommentId(ctx context.Context, cid int64, uid int64) error
	FindUserIdByCommentId(ctx context.Context, cid int64) (uid int64, err error)
}
