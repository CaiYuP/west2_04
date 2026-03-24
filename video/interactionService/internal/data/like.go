package data

import "time"

// Like 表模型，对应 `likes` 表
type Like struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:点赞ID" json:"id"`
	UserID    uint64    `gorm:"column:user_id;not null;uniqueIndex:uk_user_video;comment:用户ID" json:"user_id"`
	VideoID   uint64    `gorm:"column:video_id;not null;uniqueIndex:uk_user_video;comment:视频ID" json:"video_id"`
	CommentID uint64    `gorm:"column:comment_id;not null;uniqueIndex:uk_user_video;comment:评论ID" json:"comment_id"`
	IsComment bool      `gorm:"column:is_comment;not null;comment:是否为评论点赞" json:"is_comment"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
