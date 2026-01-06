package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"videoService/internal/data"
	"videoService/internal/database/gorms"
)

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: gorms.New(),
	}
}

type VideoDao struct {
	conn *gorms.GormConn
}

func (v *VideoDao) FindVideosAfterTime(ctx context.Context, cur int64) (items []*data.Video, err error) {
	t := time.UnixMilli(cur)
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("created_at > ?", t).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
