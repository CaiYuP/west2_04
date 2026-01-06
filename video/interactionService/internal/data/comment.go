package data

import "time"

// Comment 表模型，对应 `comments` 表
// 包含父子层级、点赞数、子评论数等字段
type Comment struct {
    ID          uint64     `gorm:"column:id;primaryKey;autoIncrement;comment:评论ID" json:"id"`
    VideoID     uint64     `gorm:"column:video_id;not null;comment:视频ID" json:"video_id"`
    UserID      uint64     `gorm:"column:user_id;not null;comment:发表者ID" json:"user_id"`
    ParentID    uint64     `gorm:"column:parent_id;default:0;comment:父评论ID" json:"parent_id"`
    LikeCount   uint64     `gorm:"column:like_count;default:0;comment:点赞数" json:"like_count"`
    ChildCount  uint64     `gorm:"column:child_count;default:0;comment:子评论数" json:"child_count"`
    Content     string     `gorm:"column:content;type:text;not null;comment:评论内容" json:"content"`
    CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
    DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
