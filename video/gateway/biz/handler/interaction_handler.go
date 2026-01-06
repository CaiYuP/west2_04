package handler

import (
	"context"
	"net/http"
	"strconv"
	"west2-video/common/errs"
	"west2-video/gateway/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	pbcommon "west2-video/api/common/v1"
	pbinteraction "west2-video/api/interaction/v1"
	"west2-video/gateway/biz/client"
)

// LikeAction 点赞操作
func LikeAction(ctx context.Context, c *app.RequestContext) {
	// 直接从表单获取参数（根据 API 文档：video_id, comment_id, action_type）
	videoIDStr := c.PostForm("video_id")
	commentIDStr := c.PostForm("comment_id")
	actionTypeStr := c.PostForm("action_type")

	// 参数验证：video_id 和 comment_id 必须存在其中一个
	if videoIDStr == "" && commentIDStr == "" {
		c.JSON(http.StatusBadRequest, pbinteraction.LikeActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: video_id 和 comment_id 必须存在其中一个"},
		})
		return
	}

	// 解析 video_id（如果存在）
	var videoID int64
	if videoIDStr != "" {
		var err error
		videoID, err = strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, pbinteraction.LikeActionReply{
				Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: video_id 格式错误"},
			})
			return
		}
	}

	// 解析 action_type
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 32)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusBadRequest, pbinteraction.LikeActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: action_type 必须是 1(点赞) 或 2(取消点赞)"},
		})
		return
	}

	req := &pbinteraction.LikeActionRequest{
		VideoId:    videoID,
		ActionType: int32(actionType),
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.InteractionClient.LikeAction(ctx, req)
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

// LikeList 点赞列表
func LikeList(ctx context.Context, c *app.RequestContext) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbinteraction.LikeListRequest{
		UserId: userID,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.InteractionClient.LikeList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	// LikeListReply 中的 Videos 字段类型是 []*v12.Video (即 west2-video/api/video/v1.Video)
	// 这里直接使用 interface{} 来避免类型问题
	data := struct {
		Videos interface{}           `json:"videos"`
		Page   *pbcommon.PageResponse `json:"page"`
	}{
		Videos: resp.GetVideos(),
		Page:   resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// CommentAction 评论
func CommentAction(ctx context.Context, c *app.RequestContext) {
	// 直接从表单获取参数（根据 API 文档：video_id, comment_id, content）
	videoIDStr := c.PostForm("video_id")
	commentIDStr := c.PostForm("comment_id")
	content := c.PostForm("content")

	// 参数验证
	if content == "" {
		c.JSON(http.StatusBadRequest, pbinteraction.CommentActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: content 不能为空"},
		})
		return
	}

	// video_id 和 comment_id 必须存在其中一个
	if videoIDStr == "" && commentIDStr == "" {
		c.JSON(http.StatusBadRequest, pbinteraction.CommentActionReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: video_id 和 comment_id 必须存在其中一个"},
		})
		return
	}

	// 解析 video_id（如果存在）
	var videoID int64
	if videoIDStr != "" {
		var err error
		videoID, err = strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, pbinteraction.CommentActionReply{
				Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: video_id 格式错误"},
			})
			return
		}
	}

	// 解析 comment_id（如果存在，用于回复评论，但本次作业不要求）
	var commentID int64
	if commentIDStr != "" {
		var err error
		commentID, err = strconv.ParseInt(commentIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, pbinteraction.CommentActionReply{
				Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: comment_id 格式错误"},
			})
			return
		}
	}

	req := &pbinteraction.CommentActionRequest{
		VideoId:   videoID,
		Content:   content,
		CommentId: commentID,
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.InteractionClient.CommentAction(ctx, req)
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
		Data: resp.GetComment(),
	}
	c.JSON(http.StatusOK, httpResp)
}

// CommentList 评论列表
func CommentList(ctx context.Context, c *app.RequestContext) {
	videoID, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbinteraction.CommentListRequest{
		VideoId: videoID,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.InteractionClient.CommentList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		Comments []*pbinteraction.Comment `json:"comments"`
		Page     *pbcommon.PageResponse  `json:"page"`
	}{
		Comments: resp.GetComments(),
		Page:     resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	// 根据 API 文档，DELETE 方法也使用 multipart/form-data
	commentIDStr := c.PostForm("comment_id")
	// videoIDStr := c.PostForm("video_id") // API 文档中有，但 proto 中不需要

	if commentIDStr == "" {
		c.JSON(http.StatusBadRequest, pbinteraction.DeleteCommentReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: comment_id 不能为空"},
		})
		return
	}

	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbinteraction.DeleteCommentReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: comment_id 格式错误"},
		})
		return
	}

	req := &pbinteraction.DeleteCommentRequest{
		CommentId: commentID,
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.InteractionClient.DeleteComment(ctx, req)
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
