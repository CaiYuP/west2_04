package biz

import (
	"context"
	"fmt"
	"time"
	"west2-video/chatService/internal/data"
	"west2-video/chatService/internal/model"
)

type GroupUseCase struct {
	groupRepo data.GroupRepo
	chatRepo  data.ChatRepo
}

func NewGroupUseCase(groupRepo data.GroupRepo, chatRepo data.ChatRepo) *GroupUseCase {
	return &GroupUseCase{
		groupRepo: groupRepo,
		chatRepo:  chatRepo,
	}
}

// CreateGroup 创建群组
func (uc *GroupUseCase) CreateGroup(ctx context.Context, creatorID uint64, req *model.CreateGroupRequest) (*model.Group, error) {
	group := &model.Group{
		Name:      req.Name,
		CreatorID: creatorID,
	}

	creatorMember := &model.GroupMember{
		UserID:   creatorID,
		Role:     2, // 群主
		JoinTime: time.Now(),
	}

	if err := uc.groupRepo.CreateGroup(ctx, group, creatorMember); err != nil {
		return nil, fmt.Errorf("无法创建群组: %w", err)
	}

	// 邀请其他成员
	for _, userID := range req.UserIDs {
		if userID == creatorID {
			continue
		}
		member := &model.GroupMember{
			GroupID:  group.ID,
			UserID:   userID,
			Role:     0, // 普通成员
			JoinTime: time.Now(),
		}
		if err := uc.groupRepo.AddMember(ctx, member); err != nil {
			// 忽略单个成员添加失败的情况，但记录日志
			fmt.Printf("无法添加成员 %d 到群组 %d: %v\n", userID, group.ID, err)
		}
	}

	return group, nil
}

// GetGroup 获取群组信息
func (uc *GroupUseCase) GetGroup(ctx context.Context, groupID uint64) (*model.Group, error) {
	return uc.groupRepo.GetGroup(ctx, groupID)
}

// AddMember 添加群组成员
func (uc *GroupUseCase) AddMember(ctx context.Context, groupID uint64, req *model.AddGroupMemberRequest) error {
	for _, userID := range req.UserIDs {
		member := &model.GroupMember{
			GroupID:  groupID,
			UserID:   userID,
			Role:     0,
			JoinTime: time.Now(),
		}
		if err := uc.groupRepo.AddMember(ctx, member); err != nil {
			return fmt.Errorf("无法添加成员 %d: %w", userID, err)
		}
	}
	return nil
}

// RemoveMember 移除群组成员
func (uc *GroupUseCase) RemoveMember(ctx context.Context, groupID, userID uint64) error {
	return uc.groupRepo.RemoveMember(ctx, groupID, userID)
}
