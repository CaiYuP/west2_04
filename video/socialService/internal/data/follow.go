package data

import "time"

// Follow 表模型，对应 `follows`
// follower_id 关注者（粉丝），followee_id 被关注者

type Follow struct {
    ID         uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:关注ID" json:"id"`
    FollowerID uint64    `gorm:"column:follower_id;not null;comment:关注者ID" json:"follower_id"`
    FolloweeID uint64    `gorm:"column:followee_id;not null;comment:被关注者ID" json:"followee_id"`
    CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
}

func (Follow) TableName() string { return "follows" }

