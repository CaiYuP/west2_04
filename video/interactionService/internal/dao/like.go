package dao

import (
	"context"
	"gorm.io/gorm"
	"interactionService/internal/data"
	"interactionService/internal/database/gorms"
)

func NewLikeDao() *LikeDao {
	return &LikeDao{
		conn: gorms.New(),
	}
}

type LikeDao struct {
	conn *gorms.GormConn
}

func (l *LikeDao) FindLikeList(ctx context.Context, id int64, page int, size int) (vids []int64, err error) {
	offset := (page - 1) * size
	err = l.conn.Session(ctx).Model(&data.Like{}).Where("user_id=? and is_comment=false", id).Select("video_id").Offset(offset).Limit(size).Find(&vids).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (l *LikeDao) CreateVideoLike(ctx context.Context, like *data.Like) error {
	err := l.conn.Session(ctx).Model(&data.Like{}).Create(like).Error
	return err
}

func (l *LikeDao) DeleteVideoLike(ctx context.Context, uid int64, vid int64) error {
	err := l.conn.Session(ctx).Where("user_id=? and video_id=?", uid, vid).Delete(&data.Like{}).Error
	return err
}

func (l *LikeDao) DeleteCommentLike(ctx context.Context, uid int64, commentId int64) error {
	err := l.conn.Session(ctx).Where("user_id=? and comment_id=?", uid, commentId).Delete(&data.Like{}).Error
	return err
}

func (l *LikeDao) CreateCommentLike(ctx context.Context, like *data.Like) error {
	// 使用 Select 明确指定要插入的字段，包括 video_id=0 和 comment_id
	// 因为数据库字段是 NOT NULL，必须插入值（即使是 0）
	err := l.conn.Session(ctx).Model(&data.Like{}).
		Select("user_id", "video_id", "comment_id", "is_comment", "created_at").
		Create(like).Error
	return err
}

func (l *LikeDao) IsLikeExist(ctx context.Context, like *data.Like) (bool, error) {
	var count int64
	err := l.conn.Session(ctx).Model(&data.Like{}).Where("video_id=? and user_id=? and comment_id=?", like.VideoID, like.UserID, like.CommentID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
