package repo

import (
	"context"
	"time"
	"videoService/internal/data"
)

type VideoRepo interface {
	FindVideosAfterTime(ctx context.Context, time int64) (items []*data.Video, err error)
	CreateVideo(ctx context.Context, v *data.Video) error
	FindVideosByUserId(ctx context.Context, id int64, size int32, page int32) (videoes []*data.Video, total int64, err error)
	BatchUpdateVisitCount(ctx context.Context, videoIDToCount map[uint64]int64) (map[uint64]int64, error)
	FindHotRankingVideos(ctx context.Context, size, page int) (items []*data.Video, err error)
	FindVideosByIDsDesc(ctx context.Context, videoIDs []uint64) (items []*data.Video, err error)
	FindVideosById(ctx context.Context, id int64) (item *data.Video, err error)
	FindVideosByTimeAndUserName(ctx context.Context, ft time.Time, tt time.Time, userId int64, size int32, num int32) (videoes []*data.Video, total int64, err error)
	FindVideosByTimeAndUserNameWithKeyWord(ctx context.Context, ft time.Time, tt time.Time, id int64, size int32, num int32, keyword string) ([]*data.Video, int64, error)
	FindVideosByIds(ctx context.Context, ids []int64) (videos []*data.Video, err error)
	IncrLikeCount(ctx context.Context, ids int64) error
	DecrLikeCount(ctx context.Context, ids int64) error
}
