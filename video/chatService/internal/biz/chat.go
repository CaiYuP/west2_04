package biz

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"west2-video/chatService/internal/data"
	"west2-video/chatService/internal/model"
)

type ChatUseCase struct {
	chatRepo  data.ChatRepo
	groupRepo data.GroupRepo
	redis     *data.RedisClient
	mu        sync.Mutex
}

func NewChatUseCase(chatRepo data.ChatRepo, groupRepo data.GroupRepo, redis *data.RedisClient) *ChatUseCase {
	return &ChatUseCase{
		chatRepo:  chatRepo,
		groupRepo: groupRepo,
		redis:     redis,
	}
}

// GetUserConversations 获取用户会话列表
func (uc *ChatUseCase) GetUserConversations(ctx context.Context, userID uint64) ([]*model.Conversation, error) {
	conversations, err := uc.chatRepo.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %w", err)
	}

	for _, conv := range conversations {
		lastMsg, err := uc.chatRepo.GetLastMessage(ctx, conv.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			hlog.CtxErrorf(ctx, "无法获取会话 %s 的最后一条消息: %v", conv.ID, err)
			continue
		}
		conv.LastMessage = lastMsg
	}

	return conversations, nil
}

// GetHistoryMessages 获取历史消息
func (uc *ChatUseCase) GetHistoryMessages(ctx context.Context, userID uint64, conversationID string, page, pageSize int) ([]*model.MessageResponse, int64, error) {
	if !uc.hasAccessToConversation(ctx, userID, conversationID) {
		return nil, 0, fmt.Errorf("无权访问此会话")
	}

	total, err := uc.chatRepo.CountMessages(ctx, conversationID)
	if err != nil {
		return nil, 0, fmt.Errorf("无法统计消息数量: %w", err)
	}

	messages, err := uc.chatRepo.GetMessages(ctx, conversationID, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("无法获取消息: %w", err)
	}

	var resp []*model.MessageResponse
	for _, msg := range messages {
		resp = append(resp, &model.MessageResponse{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			ReceiverID:  msg.ReceiverID,
			GroupID:     msg.GroupID,
			ContentType: msg.ContentType,
			Content:     msg.Content,
			IsRead:      msg.IsRead,
			CreatedAt:   msg.CreatedAt,
		})
	}

	return resp, total, nil
}

// SendPrivateMessage 发送私聊消息
func (uc *ChatUseCase) SendPrivateMessage(ctx context.Context, senderID uint64, req *model.MessageRequest) (*model.Message, error) {
	if req.ReceiverID == 0 {
		return nil, fmt.Errorf("接收者ID是必需的")
	}

	conversationID := generateConversationID(senderID, req.ReceiverID)

	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		ReceiverID:     req.ReceiverID,
		ContentType:    model.MessageType(req.ContentType), // 类型转换
		Content:        req.Content,
		IsRead:         false,
	}

	if err := uc.chatRepo.SaveMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("无法保存消息: %w", err)
	}

	if err := uc.updateConversation(ctx, conversationID, senderID, req.ReceiverID, msg); err != nil {
		hlog.CtxErrorf(ctx, "无法更新会话: %v", err)
	}

	return msg, nil
}

// SendGroupMessage 发送群聊消息
func (uc *ChatUseCase) SendGroupMessage(ctx context.Context, senderID, groupID uint64, req *model.MessageRequest) (*model.Message, error) {
	isMember, err := uc.groupRepo.IsGroupMember(ctx, groupID, senderID)
	if err != nil {
		return nil, fmt.Errorf("无法检查群组成员身份: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("用户不是该群组成员")
	}

	conversationID := fmt.Sprintf("group_%d", groupID)

	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		GroupID:        groupID,
		ContentType:    model.MessageType(req.ContentType), // 类型转换
		Content:        req.Content,
		IsRead:         false,
	}

	if err := uc.chatRepo.SaveMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("无法保存消息: %w", err)
	}

	members, err := uc.groupRepo.GetGroupMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("无法获取群组成员: %w", err)
	}

	for _, member := range members {
		if member.UserID == senderID {
			continue
		}

		if err := uc.updateConversation(ctx, conversationID, member.UserID, 0, msg); err != nil {
			hlog.CtxErrorf(ctx, "无法更新用户 %d 的会话: %v", member.UserID, err)
			continue
		}
	}

	return msg, nil
}

// MarkMessageAsRead 标记消息为已读
func (uc *ChatUseCase) MarkMessageAsRead(ctx context.Context, userID, messageID uint64) error {
	msg, err := uc.chatRepo.GetMessage(ctx, messageID)
	if err != nil {
		return fmt.Errorf("无法获取消息: %w", err)
	}

	if msg.ReceiverID != userID && !(msg.GroupID > 0) {
		return fmt.Errorf("未经授权")
	}

	if err := uc.chatRepo.MarkMessageAsRead(ctx, messageID); err != nil {
		return fmt.Errorf("无法标记消息为已读: %w", err)
	}

	if err := uc.chatRepo.DecrementUnreadCount(ctx, userID, msg.ConversationID); err != nil {
		hlog.CtxErrorf(ctx, "无法减少未读计数: %v", err)
	}

	return nil
}

// generateConversationID 生成会话ID（私聊）
func generateConversationID(userID1, userID2 uint64) string {
	ids := []uint64{userID1, userID2}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return fmt.Sprintf("%d_%d", ids[0], ids[1])
}

// hasAccessToConversation 检查用户是否有权限访问此会话
func (uc *ChatUseCase) hasAccessToConversation(ctx context.Context, userID uint64, conversationID string) bool {
	parts := strings.Split(conversationID, "_")
	if len(parts) == 2 {
		user1ID, _ := strconv.ParseUint(parts[0], 10, 64)
		user2ID, _ := strconv.ParseUint(parts[1], 10, 64)
		return userID == user1ID || userID == user2ID
	}

	var groupID uint64
	n, _ := fmt.Sscanf(conversationID, "group_%d", &groupID)
	if n == 1 {
		isMember, err := uc.groupRepo.IsGroupMember(ctx, groupID, userID)
		if err != nil {
			hlog.CtxErrorf(ctx, "无法检查群组成员身份: %v", err)
			return false
		}
		return isMember
	}

	return false
}

// updateConversation 更新会话
func (uc *ChatUseCase) updateConversation(ctx context.Context, conversationID string, userID, otherUserID uint64, msg *model.Message) error {
	conv := &model.Conversation{
		ID:           conversationID,
		Type:         1,
		TargetID:     otherUserID,
		LastMessage:  msg,
		UnreadCount:  1,
		UpdatedAt:    time.Now(),
	}

	if msg.GroupID > 0 {
		conv.Type = 2
		conv.TargetID = msg.GroupID
	}

	return uc.chatRepo.SaveConversation(ctx, userID, conv)
}