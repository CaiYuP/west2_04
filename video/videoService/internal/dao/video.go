package dao

import (
	"context"
	"errors"
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

func (v *VideoDao) DecrLikeCount(ctx context.Context, ids int64) error {
	err := v.conn.Session(ctx).Model(&data.Video{}).Where("id = ?", ids).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	return err
}

func (v *VideoDao) IncrLikeCount(ctx context.Context, ids int64) error {
	err := v.conn.Session(ctx).Model(&data.Video{}).Where("id = ?", ids).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	return err
}

func (v *VideoDao) FindVideosByIds(ctx context.Context, ids []int64) (videos []*data.Video, err error) {
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("id IN (?)", ids).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (v *VideoDao) FindVideosByTimeAndUserNameWithKeyWord(ctx context.Context, ft time.Time, tt time.Time, userId int64, size int32, num int32, keyword string) (videoes []*data.Video, total int64, err error) {
	offset := (num - 1) * size
	kw := "%" + keyword + "%"
	base := v.conn.Session(ctx).Model(&data.Video{}).Where("created_at > ? and created_at < ? and deleted_at IS NULL ", ft, tt).
		Where("(title LIKE ? OR description LIKE ?)", kw, kw)
	if userId == 0 {
		err = base.
			Offset(int(offset)).
			Limit(int(size)).Find(&videoes).Error
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		base.Count(&total)
		return
	} else {
		err = base.Where("user_id = ?", userId).
			Offset(int(offset)).
			Limit(int(size)).Find(&videoes).Error
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		base.Where("user_id = ?", userId).Count(&total)
		return
	}
}

func (v *VideoDao) FindVideosByTimeAndUserName(ctx context.Context, ft time.Time, tt time.Time, userId int64, size int32, num int32) (videoes []*data.Video, total int64, err error) {
	offset := (num - 1) * size
	base := v.conn.Session(ctx).Model(&data.Video{}).Where("created_at > ? and created_at < ? and deleted_at IS NULL", ft, tt)
	if userId == 0 {
		err = base.
			Offset(int(offset)).
			Limit(int(size)).Find(&videoes).Error
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		base.Count(&total)
		return
	} else {
		err = base.Where("user_id=?", userId).
			Offset(int(offset)).
			Limit(int(size)).Find(&videoes).Error
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		base.Where("user_id = ?", userId).Count(&total)
		return
	}

}

func (v *VideoDao) FindVideosById(ctx context.Context, id int64) (item *data.Video, err error) {
	err = v.conn.Session(ctx).Where("id = ? and deleted_at IS NULL", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return item, err
}

func (v *VideoDao) BatchUpdateVisitCount(ctx context.Context, videoIDToCount map[uint64]int64) (map[uint64]int64, error) {
	conn := gorms.NewTran()
	conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			conn.Rollback()
		}
	}()

	// 使用循环逐个更新，因为gorm的批量更新对累加操作支持有限
	for videoID, increment := range videoIDToCount {
		err := conn.Tx(ctx).Model(&data.Video{}).
			Where("id = ? and deleted_at IS NULL", videoID).
			Update("visit_count", gorm.Expr("visit_count + ?", increment)).Error
		if err != nil {
			conn.Rollback()
			return nil, err
		}
	}

	err := conn.Commit()
	if err != nil {
		conn.Rollback()
		return nil, err
	}

	// 更新后，查询数据库获取总访问量（用于更新 Redis 排行榜）
	totalVisitCounts := make(map[uint64]int64)
	ids := make([]uint64, 0, len(videoIDToCount))
	for id := range videoIDToCount {
		ids = append(ids, id)
	}

	// 批量查询更新后的访问量
	var videos []data.Video
	err = v.conn.Session(ctx).Model(&data.Video{}).
		Where("id IN ? ", ids).
		Select("id", "visit_count").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	// 构建返回的 map
	for _, video := range videos {
		totalVisitCounts[video.ID] = int64(video.VisitCount)
	}

	return totalVisitCounts, nil
}

func (v *VideoDao) FindHotRankingVideos(ctx context.Context, size, page int) (items []*data.Video, err error) {
	offset := (page - 1) * size
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("and deleted_at IS NULL").Order("visit_count desc").Offset(offset).Limit(size).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return
}

func (v *VideoDao) FindVideosByIDsDesc(ctx context.Context, videoIDs []uint64) (items []*data.Video, err error) {
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("id in (?) and deleted_at IS NULL", videoIDs).Order("visit_count desc").Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return
}

func (v *VideoDao) FindVideosByUserId(ctx context.Context, id int64, size int32, page int32) (videoes []*data.Video, total int64, err error) {
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("user_id=? and deleted_at IS NULL", id).Count(&total).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		} else {
			return nil, 0, err
		}
	}
	offset := (page - 1) * size
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("user_id=? and deleted_at IS NULL", id).Offset(int(offset)).Find(&videoes).Error
	if err != nil {
		return nil, 0, err
	}
	return videoes, total, nil
}

func (v *VideoDao) CreateVideo(ctx context.Context, video *data.Video) error {
	err := v.conn.Session(ctx).Create(video).Error
	return err
}

func (v *VideoDao) FindVideosAfterTime(ctx context.Context, cur int64) (items []*data.Video, err error) {
	t := time.UnixMilli(cur)
	err = v.conn.Session(ctx).Model(&data.Video{}).Where("created_at > ? and deleted_at IS NULL", t).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
