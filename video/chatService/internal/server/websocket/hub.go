package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"west2-video/chatService/internal/biz"
	"west2-video/chatService/internal/model"
)

// Hub 维护所有活跃的客户端，并负责消息的广播
type Hub struct {
	clients      map[uint64]*Client
	register     chan *Client
	unregister   chan *Client
	chatUseCase  *biz.ChatUseCase
	groupUseCase *biz.GroupUseCase
	mu           sync.RWMutex
}

// 全局唯一的Hub实例
var GlobalHub *Hub

// InitHub 初始化全局Hub实例
func InitHub(chatUseCase *biz.ChatUseCase, groupUseCase *biz.GroupUseCase) {
	GlobalHub = &Hub{
		clients:      make(map[uint64]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		chatUseCase:  chatUseCase,
		groupUseCase: groupUseCase,
	}
}

// Run 启动Hub的消息处理循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if oldClient, ok := h.clients[client.userID]; ok {
				oldClient.Close()
			}
			h.clients[client.userID] = client
			h.mu.Unlock()
			hlog.Infof("客户端已连接: userID=%d", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
			}
			h.mu.Unlock()
			hlog.Infof("客户端已断开: userID=%d", client.userID)
		}
	}
}

// SendToUser 将消息发送给指定的用户
func (h *Hub) SendToUser(userID uint64, message interface{}) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		if err := client.SendMessage(message); err != nil {
			hlog.Errorf("向用户 %d 发送消息失败: %v", userID, err)
		}
	}
}

// InboundMessage 用于确定消息类型的通用结构
type InboundMessage struct {
	Type int `json:"type"`
}

// handleMessage 是所有传入WebSocket消息的中央路由器
func (h *Hub) handleMessage(client *Client, rawMsg []byte) {
	var msg InboundMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		hlog.Errorf("消息解析失败: %v", err)
		return
	}

	hlog.Infof("收到用户 %d 的消息, 类型: %d", client.userID, msg.Type)

	ctx := context.Background()

	switch msg.Type {
	case 1: // type1:向某个人发送消息(私聊)
		var req model.MessageRequest
		if err := json.Unmarshal(rawMsg, &req); err != nil {
			hlog.Errorf("私聊消息解析失败: %v", err)
			return
		}
		_, err := h.chatUseCase.SendPrivateMessage(ctx, client.userID, &req)
		if err != nil {
			hlog.Errorf("发送私聊消息失败: %v", err)
		}

	case 2: // type2:获取与某个人的全部历史记录(分页)
		var req model.HistoryRequest
		if err := json.Unmarshal(rawMsg, &req); err != nil {
			hlog.Errorf("历史消息请求解析失败: %v", err)
			return
		}
		messages, _, err := h.chatUseCase.GetHistoryMessages(ctx, client.userID, req.ConversationID, req.Page, req.PageSize)
		if err != nil {
			hlog.Errorf("获取历史消息失败: %v", err)
			return
		}
		client.SendMessage(messages)

	case 4: // type4:向某个聊天室发送消息(群组)
		var req model.MessageRequest
		if err := json.Unmarshal(rawMsg, &req); err != nil {
			hlog.Errorf("群聊消息解析失败: %v", err)
			return
		}
		_, err := h.chatUseCase.SendGroupMessage(ctx, client.userID, req.GroupID, &req)
		if err != nil {
			hlog.Errorf("发送群聊消息失败: %v", err)
		}

	default:
		hlog.Warnf("收到未知消息类型: %d", msg.Type)
	}
}
