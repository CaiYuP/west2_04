package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"west2-video/chatService/internal/biz"
)

type ChatHandler struct {
	chatUseCase *biz.ChatUseCase
}

func NewChatHandler(chatUseCase *biz.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUseCase,
	}
}

// GetConversations 获取会话列表
func (h *ChatHandler) GetConversations(c context.Context, ctx *app.RequestContext) {
	userID, _ := ctx.Get("user_id")

	conversations, err := h.chatUseCase.GetUserConversations(c, userID.(uint64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "获取会话列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": conversations,
	})
}

// GetHistory 获取历史消息
func (h *ChatHandler) GetHistory(c context.Context, ctx *app.RequestContext) {
	userID, _ := ctx.Get("user_id")

	conversationID := ctx.Query("conversation_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	messages, total, err := h.chatUseCase.GetHistoryMessages(c, userID.(uint64), conversationID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "获取历史消息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"list":  messages,
			"total": total,
			"page":  page,
		},
	})
}

// MarkAsRead 标记消息为已读
func (h *ChatHandler) MarkAsRead(c context.Context, ctx *app.RequestContext) {
	userID, _ := ctx.Get("user_id")

	messageID, err := strconv.ParseUint(ctx.Query("message_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的消息ID",
		})
		return
	}

	if err := h.chatUseCase.MarkMessageAsRead(c, userID.(uint64), messageID); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "标记已读失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
}
