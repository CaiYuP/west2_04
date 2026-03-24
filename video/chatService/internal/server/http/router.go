package http

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"west2-video/chatService/internal/server/http/handler"
	"west2-video/chatService/internal/server/websocket"
)

// RegisterRoutes 注册所有HTTP和WebSocket路由
func RegisterRoutes(h *server.Hertz, wsHandler *websocket.WebSocketHandler, chatHandler *handler.ChatHandler, groupHandler *handler.GroupHandler) {
	// WebSocket 路由
	h.GET("/ws/chat", wsHandler.Upgrade)

	// API 路由组
	api := h.Group("/api/v1")
	{
		// 在实际应用中，这里应该有一个JWT认证中间件
		// api.Use(jwt.Middleware())

		// 聊天相关路由
		chat := api.Group("/chat")
		{
			chat.GET("/conversations", chatHandler.GetConversations)
			chat.GET("/history", chatHandler.GetHistory)
			chat.POST("/read", chatHandler.MarkAsRead)
		}

		// 群组相关路由
		group := api.Group("/group")
		{
			group.POST("", groupHandler.CreateGroup)
			group.GET("/:id", groupHandler.GetGroup)
			group.POST("/:id/members", groupHandler.AddMember)
			group.DELETE("/:id/members/:user_id", groupHandler.RemoveMember)
		}
	}
}
