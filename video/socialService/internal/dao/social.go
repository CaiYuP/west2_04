package dao

import (
	"context"
	"gorm.io/gorm"
	"socialService/internal/data"
	"socialService/internal/database/gorms"
)

func NewSocialDao() *SocialDao {
	return &SocialDao{
		conn: gorms.New(),
	}
}

type SocialDao struct {
	conn *gorms.GormConn
}

func (s *SocialDao) FindFriend(ctx context.Context, id int64, page int, size int) (fs []*data.Follow, total int64, err error) {
	offset := (page - 1) * size
	ids := make([]int64, 0)
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id=?", id).Offset(offset).Limit(size).Pluck("followee_id", &ids).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id in(?) and followee_id=?", ids, id).Find(&fs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id=?", id).Pluck("followee_id", &ids).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id in(?) and followee_id=?", ids, id).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return fs, total, err
}

func (s *SocialDao) FindFollowListByFolloweeId(ctx context.Context, id int64, page int, size int) (fs []*data.Follow, total int64, err error) {
	offset := (page - 1) * size
	err = s.conn.Session(ctx).Where("followee_id=?", id).Offset(offset).Limit(size).Find(&fs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("followee_id=?", id).Count(&total).Error
	return fs, total, err
}

func (s *SocialDao) FindIsExist(ctx context.Context, id int64, toId int64) (bool, error) {
	var count int64
	err := s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id = ? and followee_id = ?", id, toId).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SocialDao) FindFollowListByFollowerId(ctx context.Context, id int64, page int, size int) (fs []*data.Follow, total int64, err error) {
	offset := (page - 1) * size
	err = s.conn.Session(ctx).Where("follower_id=?", id).Offset(offset).Limit(size).Find(&fs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	err = s.conn.Session(ctx).Model(&data.Follow{}).Where("follower_id=?", id).Count(&total).Error
	return fs, total, err
}

func (s *SocialDao) CreateFollow(ctx context.Context, follow *data.Follow) error {
	err := s.conn.Session(ctx).Create(follow).Error
	return err
}

func (s *SocialDao) DeleteFollow(ctx context.Context, id int64, toId int64) error {
	err := s.conn.Session(ctx).Where("follower_id = ? and followee_id = ?", id, toId).Delete(&data.Follow{}).Error
	return err
}
