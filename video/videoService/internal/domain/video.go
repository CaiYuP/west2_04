package domain

import (
	"context"
	"go.uber.org/zap"
	"time"
	"videoService/internal/dao"
	"videoService/internal/data"
	"videoService/internal/repo"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

func NewVideoDomain() *VideoDomain {
	return &VideoDomain{
		videoRepo: dao.NewVideoDao(),
	}
}

type VideoDomain struct {
	videoRepo repo.VideoRepo
}

func (d *VideoDomain) VerifySeed(ctx context.Context, cur int64) *errs.BError {
	t := time.Now().UnixMilli()
	if t < cur {
		return model.SeedError
	}
	return nil
}

func (d *VideoDomain) FindVideosAfterTime(ctx context.Context, latestTime int64) ([]*data.Video, *errs.BError) {
	items, err := d.videoRepo.FindVideosAfterTime(ctx, latestTime)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosAfterTime.FindVideosAfterTime error", zap.Error(err))
		return nil, model.DBError
	}
	if items == nil {
		return make([]*data.Video, 0), nil
	}
	return items, nil
}
