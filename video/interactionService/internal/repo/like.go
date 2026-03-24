package repo

import (
	"context"
	"interactionService/internal/data"
)

type LikeRepo interface {
	CreateCommentLike(ctx context.Context, l *data.Like) error
	DeleteCommentLike(ctx context.Context, uid int64, id int64) error
	IsLikeExist(ctx context.Context, like *data.Like) (bool, error)
	CreateVideoLike(ctx context.Context, like *data.Like) error
	DeleteVideoLike(ctx context.Context, uid int64, vid int64) error
	FindLikeList(ctx context.Context, id int64, page int, size int) (vids []int64, err error)
}
