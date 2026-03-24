package websocket

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"west2-video/chatService/internal/biz"
)

// WebSocketHandler 负责处理WebSocket连接的升级
type WebSocketHandler struct {
	chatUseCase  *biz.ChatUseCase
	groupUseCase *biz.GroupUseCase
}

// NewWebSocketHandler 创建一个新的WebSocket处理器
func NewWebSocketHandler(chatUseCase *biz.ChatUseCase, groupUseCase *biz.GroupUseCase) *WebSocketHandler {
	return &WebSocketHandler{
		chatUseCase:  chatUseCase,
		groupUseCase: groupUseCase,
	}
}

// Upgrade 将HTTP连接升级到WebSocket
func (h *WebSocketHandler) Upgrade(ctx context.Context, c *app.RequestContext) {
	// 正确的 Upgrader 初始化方式
	upgrader := websocket.HertzUpgrader{
		CheckOrigin: func(c *app.RequestContext) bool {
			// 在生产环境中，这里应该有更严格的来源域名白名单检查
			return true
		},
	}

	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			hlog.Error("WebSocket连接失败: 缺少 user_id")
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil || userID == 0 {
			hlog.Error("WebSocket连接失败: 无效的用户ID")
			return
		}

		// NewClient 调用已修正，不再需要 handler 参数
		client := NewClient(GlobalHub, conn, userID)
		GlobalHub.register <- client

		go client.WritePump()
		// 调用不带参数的 ReadPump
		client.ReadPump()
	})

	if err != nil {
		hlog.Errorf("WebSocket升级失败: %v", err)
	}
}
