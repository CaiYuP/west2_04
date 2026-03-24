package data

import (
	"context"
	"gorm.io/gorm"
	"west2-video/chatService/internal/model"
)

// ChatRepo 聊天数据访问接口
type ChatRepo interface {
	// SaveMessage 保存消息
	SaveMessage(ctx context.Context, message *model.Message) error
	// GetMessage 获取消息
	GetMessage(ctx context.Context, id uint64) (*model.Message, error)
	// GetMessages 获取消息列表
	GetMessages(ctx context.Context, conversationID string, offset, limit int) ([]*model.Message, error)
	// CountMessages 统计消息数量
	CountMessages(ctx context.Context, conversationID string) (int64, error)
	// MarkMessageAsRead 标记消息为已读
	MarkMessageAsRead(ctx context.Context, id uint64) error
	// SaveConversation 保存会话
	SaveConversation(ctx context.Context, userID uint64, conv *model.Conversation) error
	// GetConversation 获取会话
	GetConversation(ctx context.Context, userID uint64, conversationID string) (*model.Conversation, error)
	// GetUserConversations 获取用户会话列表
	GetUserConversations(ctx context.Context, userID uint64) ([]*model.Conversation, error)
	// GetLastMessage 获取最后一条消息
	GetLastMessage(ctx context.Context, conversationID string) (*model.Message, error)
	// IncrementUnreadCount 增加未读计数
	IncrementUnreadCount(ctx context.Context, userID uint64, conversationID string) error
	// DecrementUnreadCount 减少未读计数
	DecrementUnreadCount(ctx context.Context, userID uint64, conversationID string) error
}

type chatRepo struct {
	db *gorm.DB
}

// NewChatRepo 创建聊天仓库
func NewChatRepo(db *gorm.DB) ChatRepo {
	return &chatRepo{db: db}
}

func (r *chatRepo) SaveMessage(ctx context.Context, message *model.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *chatRepo) GetMessage(ctx context.Context, id uint64) (*model.Message, error) {
	var msg model.Message
	err := r.db.WithContext(ctx).First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *chatRepo) GetMessages(ctx context.Context, conversationID string, offset, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	// 反转顺序，使最新的消息在最后
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *chatRepo) CountMessages(ctx context.Context, conversationID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error

	return count, err
}

func (r *chatRepo) MarkMessageAsRead(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Message{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *chatRepo) SaveConversation(ctx context.Context, userID uint64, conv *model.Conversation) error {
	// 这里简化处理，实际项目中应该使用事务
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO user_conversations (user_id, conversation_id, unread_count, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE 
		 unread_count = unread_count + VALUES(unread_count),
		 updated_at = VALUES(updated_at)`,
		userID, conv.ID, conv.UnreadCount, conv.UpdatedAt,
	).Error
}

func (r *chatRepo) GetConversation(ctx context.Context, userID uint64, conversationID string) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.WithContext(ctx).Raw(
		`SELECT conversation_id as id, unread_count, updated_at 
		 FROM user_conversations 
		 WHERE user_id = ? AND conversation_id = ?`,
		userID, conversationID,
	).Scan(&conv).Error

	if err != nil {
		return nil, err
	}

	return &conv, nil
}

func (r *chatRepo) GetUserConversations(ctx context.Context, userID uint64) ([]*model.Conversation, error) {
	var convs []*model.Conversation

	err := r.db.WithContext(ctx).Raw(
		`SELECT conversation_id as id, unread_count, updated_at 
		 FROM user_conversations 
		 WHERE user_id = ? 
		 ORDER BY updated_at DESC`,
		userID,
	).Scan(&convs).Error

	if err != nil {
		return nil, err
	}

	return convs, nil
}

func (r *chatRepo) GetLastMessage(ctx context.Context, conversationID string) (*model.Message, error) {
	var msg model.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		First(&msg).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &msg, nil
}

func (r *chatRepo) IncrementUnreadCount(ctx context.Context, userID uint64, conversationID string) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO user_conversations (user_id, conversation_id, unread_count, updated_at)
		 VALUES (?, ?, 1, NOW())
		 ON DUPLICATE KEY UPDATE 
		 unread_count = unread_count + 1,
		 updated_at = NOW()`,
		userID, conversationID,
	).Error
}

func (r *chatRepo) DecrementUnreadCount(ctx context.Context, userID uint64, conversationID string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE user_conversations 
		 SET unread_count = GREATEST(0, unread_count - 1),
		     updated_at = NOW()
		 WHERE user_id = ? AND conversation_id = ?`,
		userID, conversationID,
	).Error
}