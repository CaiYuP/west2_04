package handler

import (
	"context"
	"net/http"
	"strconv"
	"west2-video/common/errs"
	"west2-video/gateway/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	pbcommon "west2-video/api/common/v1"
	pbvideo "west2-video/api/video/v1"
	"west2-video/gateway/biz/client"
)

// Feed 视频流
func Feed(ctx context.Context, c *app.RequestContext) {
	latestTime, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	//pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	//if pageSize == 0 {
	//	pageSize = 10
	//}

	req := &pbvideo.FeedRequest{
		LatestTime: latestTime,
		//PageSize:   int32(pageSize),
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.VideoClient.Feed(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}
	items := &videoListResp{
		Items: resp.Items,
	}
	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: items,
	}
	c.JSON(http.StatusOK, httpResp)
}

// Publish 投稿
func Publish(ctx context.Context, c *app.RequestContext) {
	title := c.PostForm("title")
	description := c.PostForm("description")

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, pbvideo.PublishReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "文件上传失败: " + err.Error()},
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, pbvideo.PublishReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "打开文件失败: " + err.Error()},
		})
		return
	}
	defer src.Close()

	data := make([]byte, file.Size)
	_, err = src.Read(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbvideo.PublishReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "读取文件失败: " + err.Error()},
		})
		return
	}

	req := &pbvideo.PublishRequest{
		Title:       title,
		Description: description,
		Data:        data,
		Id:          c.GetInt64("user_id"),
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.VideoClient.Publish(ctx, req)
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

// PublishList 发布列表
func PublishList(ctx context.Context, c *app.RequestContext) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbvideo.PublishListRequest{
		UserId: userID,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.VideoClient.PublishList(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		VideoList []*pbvideo.Video       `json:"video_list"`
		Page      *pbcommon.PageResponse `json:"page"`
	}{
		VideoList: resp.GetVideoList(),
		Page:      resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// HotRanking 热门排行榜
func HotRanking(ctx context.Context, c *app.RequestContext) {
	pageNum, _ := strconv.ParseInt(c.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbvideo.HotRankingRequest{
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.VideoClient.HotRanking(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		VideoList []*pbvideo.Video       `json:"video_list"`
		Page      *pbcommon.PageResponse `json:"page"`
	}{
		VideoList: resp.GetVideoList(),
		Page:      resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}

// SearchVideo 搜索视频
func SearchVideo(ctx context.Context, c *app.RequestContext) {
	// 根据 API 文档，POST 方法使用 multipart/form-data
	keywords := c.PostForm("keywords")
	pageNumStr := c.PostForm("page_num")
	pageSizeStr := c.PostForm("page_size")
	fromDateStr := c.PostForm("from_date")
	toDateStr := c.PostForm("to_date")
	username := c.PostForm("username")

	// 参数验证（根据 API 文档，keywords, page_size, page_num 是必需的）
	if pageNumStr == "" || pageSizeStr == "" {
		c.JSON(http.StatusBadRequest, pbvideo.SearchReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: page_num 和 page_size 不能为空"},
		})
		return
	}

	pageNum, err := strconv.ParseInt(pageNumStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbvideo.SearchReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: page_num 格式错误"},
		})
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbvideo.SearchReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: page_size 格式错误"},
		})
		return
	}

	// 注意：proto 中只有 keywords 和 page，没有 from_date, to_date, username
	// 这些字段可能需要扩展 proto 或通过其他方式传递
	req := &pbvideo.SearchRequest{
		Keywords: keywords,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	// TODO: 如果需要支持 from_date, to_date, username，需要扩展 proto 定义
	_ = fromDateStr
	_ = toDateStr
	_ = username

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.VideoClient.Search(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, HTTPResponse{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
			Data: nil,
		})
		return
	}

	data := struct {
		VideoList []*pbvideo.Video       `json:"video_list"`
		Page      *pbcommon.PageResponse `json:"page"`
	}{
		VideoList: resp.GetVideoList(),
		Page:      resp.GetPage(),
	}

	httpResp := &HTTPResponse{
		Base: resp.GetBase(),
		Data: data,
	}
	c.JSON(http.StatusOK, httpResp)
}
