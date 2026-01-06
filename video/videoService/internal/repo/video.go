package repo

import (
	"context"
	"videoService/internal/data"
)

type VideoRepo interface {
	FindVideosAfterTime(ctx context.Context, time int64) (items []*data.Video, err error)
}
