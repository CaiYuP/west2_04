package domain

import (
	"context"
	"go.uber.org/zap"
	"socialService/internal/dao"
	"socialService/internal/data"
	"socialService/internal/repo"
	"time"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/gateway/biz/model"
)

func NewSocialDomain() *SocialDomain {
	return &SocialDomain{
		socialRepo: dao.NewSocialDao(),
	}
}

type SocialDomain struct {
	socialRepo repo.SocialRepo
}

func (d *SocialDomain) CreateFollow(ctx context.Context, id int64, toId int64) *errs.BError {
	if id == toId {
		return model.CannotFollowSelf
	}
	b, err := d.socialRepo.FindIsExist(ctx, id, toId)
	if err != nil {
		logs.LG.Error("SocialServiceService.CreateFollow error", zap.Error(err))
		return model.DBError
	}
	if b {
		return model.ConnotRepeatFollow
	}
	follow := &data.Follow{
		FollowerID: uint64(id),
		FolloweeID: uint64(toId),
		CreatedAt:  time.Now(),
	}
	err = d.socialRepo.CreateFollow(ctx, follow)
	if err != nil {
		logs.LG.Error("SocialDomain.CreateFollow.socialRepo.CreateFollow error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *SocialDomain) DeleteFollow(ctx context.Context, id int64, toId int64) *errs.BError {
	if id == toId {
		return model.CannotCancelFollowSelf
	}
	err := d.socialRepo.DeleteFollow(ctx, id, toId)
	if err != nil {
		logs.LG.Error("SocialDomain.DeleteFollow.deleteFollow error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *SocialDomain) FindFollowListByFollowerId(ctx context.Context, id int64, page int, size int) ([]int64, int64, *errs.BError) {
	fs, total, err := d.socialRepo.FindFollowListByFollowerId(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("SocialDomain.FindFollowListByFollowerId error", zap.Error(err))
		return nil, 0, model.DBError
	}
	if total == 0 {
		return make([]int64, 0), 0, nil
	}
	fids := make([]int64, len(fs))
	for i, f := range fs {
		fids[i] = int64(f.FolloweeID)
	}
	return fids, total, nil
}

func (d *SocialDomain) FindFollowListByFolloweeId(ctx context.Context, id int64, page int, size int) ([]int64, int64, *errs.BError) {
	fs, total, err := d.socialRepo.FindFollowListByFolloweeId(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("SocialDomain.FindFollowListByFolloweeId error", zap.Error(err))
		return nil, 0, model.DBError
	}
	if total == 0 {
		return make([]int64, 0), 0, nil
	}
	fids := make([]int64, len(fs))
	for i, f := range fs {
		fids[i] = int64(f.FollowerID)
	}
	return fids, total, nil
}

func (d *SocialDomain) FindFriend(ctx context.Context, id int64, page int, size int) ([]int64, int64, *errs.BError) {
	fs, total, err := d.socialRepo.FindFriend(ctx, id, page, size)
	if err != nil {
		logs.LG.Error("SocialDomain.FindFriend error", zap.Error(err))
		return nil, 0, model.DBError
	}
	if total == 0 {
		return make([]int64, 0), 0, nil
	}
	fids := make([]int64, len(fs))
	for i, f := range fs {
		fids[i] = int64(f.FollowerID)
	}
	return fids, total, nil
}
