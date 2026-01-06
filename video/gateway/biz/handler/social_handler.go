package handler

import (
	"context"
	"net/http"
	"strconv"
	"west2-video/common/errs"
	"west2-video/gateway/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	pbcommon "west2-video/api/common/v1"
	pbsocial "west2-video/api/social/v1"
	"west2-video/gateway/biz/client"
)

// FollowAction 关注操作
func FollowAction(ctx context.Context, c *app.RequestContext) {
	// 直接从表单获取参数（根据 API 文档：to_user_id, action_type）
	toUserIDStr := c.PostForm("to_user_id")
	actionTypeStr := c.PostForm("action_type")

	// 参数验证
	if toUserIDStr == "" || actionTypeStr == "" {
		c.JSON(http.StatusBadRequest, pbsocial.FollowActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: to_user_id 和 action_type 不能为空"},
		})
		return
	}

	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbsocial.FollowActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: to_user_id 格式错误"},
		})
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbsocial.FollowActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: action_type 格式错误"},
		})
		return
	}

	// API 文档：0-关注, 1-取关；但 proto 中：1-关注, 2-取关
	// 需要转换：0->1(关注), 1->2(取关)
	var protoActionType int32
	if actionType == 0 {
		protoActionType = 1 // 关注
	} else if actionType == 1 {
		protoActionType = 2 // 取关
	} else {
		c.JSON(http.StatusBadRequest, pbsocial.FollowActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: action_type 必须是 0(关注) 或 1(取关)"},
		})
		return
	}

	req := &pbsocial.FollowActionRequest{
		ToUserId:   toUserID,
		ActionType: protoActionType,
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.SocialClient.FollowAction(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: nil,
	}
	c.JSON(http.StatusOK, httpResp)
}

// FollowList 关注列表
func FollowList(ctx context.Context, c *app.RequestContext) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbsocial.FollowListRequest{
		UserId: userID,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.SocialClient.FollowList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		Users []*pbsocial.SocialUser   `json:"users"`
		Page  *pbcommon.PageResponse `json:"page"`
	}{
		Users: resp.GetUsers(),
		Page:  resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// FollowerList 粉丝列表
func FollowerList(ctx context.Context, c *app.RequestContext) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbsocial.FollowerListRequest{
		UserId: userID,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.SocialClient.FollowerList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		Users []*pbsocial.SocialUser   `json:"users"`
		Page  *pbcommon.PageResponse `json:"page"`
	}{
		Users: resp.GetUsers(),
		Page:  resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// FriendList 好友列表
func FriendList(ctx context.Context, c *app.RequestContext) {
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbsocial.FriendListRequest{
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.SocialClient.FriendList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		Users []*pbsocial.SocialUser   `json:"users"`
		Page  *pbcommon.PageResponse `json:"page"`
	}{
		Users: resp.GetUsers(),
		Page:  resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

