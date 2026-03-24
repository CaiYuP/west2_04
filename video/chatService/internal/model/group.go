package model

import "time"

// Group 聊天群组
type Group struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:64;not null"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	CreatorID uint64    `json:"creator_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GroupMember 群组成员
type GroupMember struct {
	ID       uint64    `json:"id" gorm:"primaryKey"`
	GroupID  uint64    `json:"group_id" gorm:"uniqueIndex:uk_group_user"`
	UserID   uint64    `json:"user_id" gorm:"uniqueIndex:uk_group_user"`
	Role     int       `json:"role" gorm:"default:0"` // 0-普通成员 1-管理员 2-群主
	JoinTime time.Time `json:"join_time"`
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Name    string   `json:"name" binding:"required,min=2,max=20"`
	UserIDs []uint64 `json:"user_ids"` // 邀请的成员ID
}

// AddGroupMemberRequest 添加群组成员请求
type AddGroupMemberRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required"`
}
