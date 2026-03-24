package repo

import (
	"context"
	"socialService/internal/data"
)

type SocialRepo interface {
	CreateFollow(ctx context.Context, follow *data.Follow) error
	DeleteFollow(ctx context.Context, id int64, toId int64) error
	FindFollowListByFollowerId(ctx context.Context, id int64, page int, size int) ([]*data.Follow, int64, error)
	FindIsExist(ctx context.Context, id int64, toId int64) (bool, error)
	FindFollowListByFolloweeId(ctx context.Context, id int64, page int, size int) ([]*data.Follow, int64, error)
	FindFriend(ctx context.Context, id int64, page int, size int) ([]*data.Follow, int64, error)
}
