package data

import (
	"github.com/jinzhu/copier"
	"strconv"
	"time"
	userpb "west2-video/api/user/v1"
	"west2-video/common/tms"
)

// User 用户表模型
type User struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:用户ID" json:"id"`
	Username     string    `gorm:"column:username;type:varchar(32);not null;uniqueIndex:uk_username;comment:用户名" json:"username"`
	Password     string    `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"-"`
	Email        string    `gorm:"column:email;type:varchar(128);comment:邮箱" json:"email"`
	Nickname     string    `gorm:"column:nickname;type:varchar(64);comment:昵称" json:"nickname"`
	AvatarURL    string    `gorm:"column:avatar_url;type:varchar(512);comment:头像URL" json:"avatar_url"`
	Description  string    `gorm:"column:description;type:varchar(512);comment:个人简介" json:"description"`
	IsMfaEnabled bool      `gorm:"column:is_mfa_enabled;default:0;comment:是否启用MFA" json:"is_mfa_enabled"`
	MfaSecret    string    `gorm:"column:mfa_secret;type:varchar(64);comment:MFA Secret" json:"mfa_secret"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func UserFormat(user *User) *userpb.User {
	u := &userpb.User{}
	copier.Copy(u, user)
	u.Id = strconv.FormatUint(user.ID, 10)
	u.CreatedAt = tms.Format(user.CreatedAt)
	u.UpdatedAt = tms.Format(user.UpdatedAt)
	u.DeletedAt = tms.Format(user.DeletedAt)
	return u
}
