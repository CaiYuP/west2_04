package dao

import (
	"context"
	"gorm.io/gorm"
	"interactionService/internal/data"
	"interactionService/internal/database/gorms"
)

func NewCommentDao() *CommentDao {
	return &CommentDao{
		conn: gorms.New(),
	}
}

type CommentDao struct {
	conn *gorms.GormConn
}

func (c *CommentDao) FindUserIdByCommentId(ctx context.Context, cid int64) (uid int64, err error) {
	err = c.conn.Session(ctx).Model(&data.Comment{}).Where("id = ?", cid).Pluck("user_id", &uid).Error
	return
}

func (c *CommentDao) DeleteCommentByVideoId(ctx context.Context, vid int64) error {
	// 先查询出要删除的评论ID列表
	var commentIDs []uint64
	err := c.conn.Session(ctx).Model(&data.Comment{}).
		Where("video_id=? ", vid).
		Pluck("id", &commentIDs).Error
	if err != nil {
		return err
	}

	if len(commentIDs) == 0 {
		return nil
	}

	// 删除这些评论及其所有子评论（递归删除）
	err = c.deleteCommentsAndChildren(ctx, commentIDs)
	if err != nil {
		return err
	}

	return nil
}

//递归删除评论及其所有子评论
func (c *CommentDao) deleteCommentsAndChildren(ctx context.Context, parentIDs []uint64) error {
	if len(parentIDs) == 0 {
		return nil
	}

	// 删除这些父评论
	err := c.conn.Session(ctx).Where("id IN ?", parentIDs).Delete(&data.Comment{}).Error
	if err != nil {
		return err
	}

	// 查询这些评论的所有子评论ID
	var childIDs []uint64
	err = c.conn.Session(ctx).Model(&data.Comment{}).
		Where("parent_id IN ?", parentIDs).
		Pluck("id", &childIDs).Error
	if err != nil {
		return err
	}

	// 如果有子评论，递归删除
	if len(childIDs) > 0 {
		return c.deleteCommentsAndChildren(ctx, childIDs)
	}

	return nil
}

func (c *CommentDao) DeleteCommentByCommentId(ctx context.Context, cid int64, uid int64) error {
	d := &data.Comment{}
	err := c.conn.Session(ctx).Where("id=? and  user_id=?", cid, uid).First(d).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	ids := make([]uint64, 0)
	ids = append(ids, uint64(cid))
	err = c.deleteCommentsAndChildren(ctx, ids)
	return err
}

func (c *CommentDao) FindVideoComment(ctx context.Context, id int64, page int, size int) (items []*data.Comment, err error) {
	offset := (page - 1) * size
	err = c.conn.Session(ctx).Model(&data.Comment{}).Where("video_id=?", id).Find(&items).Offset(offset).Limit(size).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (c *CommentDao) FindCommentComment(ctx context.Context, id int64, page int, size int) (items []*data.Comment, err error) {
	offset := (page - 1) * size
	err = c.conn.Session(ctx).Model(&data.Comment{}).Where("parent_id=?", id).Find(&items).Offset(offset).Limit(size).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (c *CommentDao) IncrLikeCount(ctx context.Context, id int64, isLike bool) error {
	if isLike {
		c.conn.Session(ctx).Model(&data.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	} else {
		c.conn.Session(ctx).Model(&data.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}
	return nil
}

func (c *CommentDao) IncryChildCount(ctx context.Context, id int64) error {
	err := c.conn.Session(ctx).Model(&data.Comment{}).Where("id = ?", id).
		UpdateColumn("child_count", gorm.Expr("child_count + ?", 1)).Error
	return err
}

func (c *CommentDao) CreateComment(ctx context.Context, comment *data.Comment) error {
	err := c.conn.Session(ctx).Create(comment).Error
	return err
}
