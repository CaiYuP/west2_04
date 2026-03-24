package data

import (
	"context"
	"gorm.io/gorm"
	"west2-video/chatService/internal/model"
)

// GroupRepo 群组数据访问接口
type GroupRepo interface {
	CreateGroup(ctx context.Context, group *model.Group, member *model.GroupMember) error
	GetGroup(ctx context.Context, id uint64) (*model.Group, error)
	AddMember(ctx context.Context, member *model.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID uint64) error
	IsGroupMember(ctx context.Context, groupID, userID uint64) (bool, error)
	GetGroupMembers(ctx context.Context, groupID uint64) ([]*model.GroupMember, error)
}

type groupRepo struct {
	db *gorm.DB
}

// NewGroupRepo 创建群组仓库
func NewGroupRepo(db *gorm.DB) GroupRepo {
	return &groupRepo{db: db}
}

func (r *groupRepo) CreateGroup(ctx context.Context, group *model.Group, member *model.GroupMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		member.GroupID = group.ID
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *groupRepo) GetGroup(ctx context.Context, id uint64) (*model.Group, error) {
	var group model.Group
	err := r.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepo) AddMember(ctx context.Context, member *model.GroupMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *groupRepo) RemoveMember(ctx context.Context, groupID, userID uint64) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&model.GroupMember{}).Error
}

func (r *groupRepo) IsGroupMember(ctx context.Context, groupID, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *groupRepo) GetGroupMembers(ctx context.Context, groupID uint64) ([]*model.GroupMember, error) {
	var members []*model.GroupMember
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
