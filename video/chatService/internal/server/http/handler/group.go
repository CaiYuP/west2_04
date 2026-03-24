package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"west2-video/chatService/internal/biz"
	"west2-video/chatService/internal/model"
)

type GroupHandler struct {
	groupUseCase *biz.GroupUseCase
}

func NewGroupHandler(groupUseCase *biz.GroupUseCase) *GroupHandler {
	return &GroupHandler{
		groupUseCase: groupUseCase,
	}
}

// CreateGroup 创建群组
func (h *GroupHandler) CreateGroup(c context.Context, ctx *app.RequestContext) {
	userID, _ := ctx.Get("user_id")

	var req model.CreateGroupRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的请求参数: " + err.Error(),
		})
		return
	}

	group, err := h.groupUseCase.CreateGroup(c, userID.(uint64), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "创建群组失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": group,
	})
}

// GetGroup 获取群组信息
func (h *GroupHandler) GetGroup(c context.Context, ctx *app.RequestContext) {
	groupID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的群组ID",
		})
		return
	}

	group, err := h.groupUseCase.GetGroup(c, groupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "获取群组信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": group,
	})
}

// AddMember 添加群组成员
func (h *GroupHandler) AddMember(c context.Context, ctx *app.RequestContext) {
	groupID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的群组ID",
		})
		return
	}

	var req model.AddGroupMemberRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的请求参数: " + err.Error(),
		})
		return
	}

	if err := h.groupUseCase.AddMember(c, groupID, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "添加成员失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
}

// RemoveMember 移除群组成员
func (h *GroupHandler) RemoveMember(c context.Context, ctx *app.RequestContext) {
	groupID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的群组ID",
		})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	if err := h.groupUseCase.RemoveMember(c, groupID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "移除成员失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
}
