package model

import "time"

// MessageType 消息类型
type MessageType int

const (
	MessageTypeText  MessageType = 1 // 文本
	MessageTypeImage MessageType = 2 // 图片
	MessageTypeVoice MessageType = 3 // 语音
	MessageTypeVideo MessageType = 4 // 视频
	MessageTypeFile  MessageType = 5 // 文件
)

// Message 聊天消息
type Message struct {
	ID             uint64      `json:"id" gorm:"primaryKey"`
	ConversationID string      `json:"conversation_id" gorm:"size:64;index"`
	SenderID       uint64      `json:"sender_id" gorm:"index"`
	ReceiverID     uint64      `json:"receiver_id" gorm:"index"`
	GroupID        uint64      `json:"group_id" gorm:"index"`
	ContentType    MessageType `json:"content_type"`
	Content        string      `json:"content" gorm:"type:text"`
	IsRead         bool        `json:"is_read"`
	CreatedAt      time.Time   `json:"created_at"`
}

// Conversation 会话
type Conversation struct {
	ID           string    `json:"id"`
	Type         int       `json:"type"` // 1-私聊 2-群聊
	TargetID     uint64    `json:"target_id"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	LastMessage  *Message  `json:"last_message,omitempty"`
	UnreadCount  int       `json:"unread_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MessageRequest 发送消息请求
type MessageRequest struct {
	Type        MessageType `json:"type" binding:"required"` // 消息类型
	Content     string      `json:"content" binding:"required"`
	ReceiverID  uint64      `json:"receiver_id"`  // 接收者ID（私聊）
	GroupID     uint64      `json:"group_id"`     // 群组ID（群聊）
	ContentType int         `json:"content_type"` // 内容类型
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID          uint64      `json:"id"`
	SenderID    uint64      `json:"sender_id"`
	ReceiverID  uint64      `json:"receiver_id"`
	GroupID     uint64      `json:"group_id"`
	ContentType MessageType `json:"content_type"`
	Content     string      `json:"content"`
	IsRead      bool        `json:"is_read"`
	CreatedAt   time.Time   `json:"created_at"`
}

// HistoryRequest 历史消息请求
type HistoryRequest struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// UnreadRequest 未读消息请求
type UnreadRequest struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	LastMessageID  uint64 `form:"last_message_id" binding:"required"`
}