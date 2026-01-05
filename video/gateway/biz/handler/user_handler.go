package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"west2-video/common/errs"
	"west2-video/common/img"
	"west2-video/common/logs"
	"west2-video/common/upload"
	"west2-video/gateway/biz/model"

	"github.com/cloudwego/hertz/pkg/app"
	pbcommon "west2-video/api/common/v1"
	pbuser "west2-video/api/user/v1"
	"west2-video/gateway/biz/client"
)

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	// 直接从表单获取参数
	req := &pbuser.RegisterRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		//Email:    c.PostForm("email"),
	}

	// 参数验证
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, pbuser.RegisterReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: username 和 password 不能为空"},
		})
		return
	}

	// 调用 gRPC 服务
	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.Register(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.RegisterReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
func GetIp(c *app.RequestContext) string {
	// 获取客户端 IP 地址
	clientIP := c.ClientIP()
	if clientIP == "::1" {
		clientIP = "127.0.0.1"
	}
	return clientIP
}

// RefreshToken 刷新 Token
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	// 获取客户端 IP 地址
	clientIP := GetIp(c)

	// 从请求中获取 refreshToken（可以从 header 或 query 参数获取）
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		// 也可以从 header 获取
		refreshToken = string(c.Request.Header.Get("Refresh-Token"))
	}

	req := &pbuser.RefreshRequest{
		RefreshToken: refreshToken,
		Ip:           clientIP,
	}

	// 参数验证
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, pbuser.RefreshReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: refresh_token 不能为空"},
		})
		return
	}

	// 调用 gRPC 服务
	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.Refresh(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.RefreshReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Login 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	// 获取客户端 IP 地址
	clientIP := GetIp(c)
	req := &pbuser.LoginRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Ip:       clientIP,
	}

	// 参数验证
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, pbuser.LoginReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "参数错误: username 和 password 不能为空"},
		})
		return
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.Login(ctx, req)
	if err != nil {
		code, m := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.LoginReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: m},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Query("user_id")
	var userID int64 = 0

	if userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			userID = id
		}
	}

	req := &pbuser.UserInfoRequest{UserId: userID}
	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.GetUserInfo(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.UserInfoReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UploadAvatar 上传头像
func UploadAvatar(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "文件上传失败: " + err.Error()},
		})
		return
	}

	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "打开文件失败: " + err.Error()},
		})
		return
	}
	defer src.Close()

	// 使用 io.ReadAll 确保读取完整的文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "读取文件失败: " + err.Error()},
		})
		return
	}
	isImage, form := img.IsImage(data)
	if !isImage {
		c.JSON(http.StatusBadRequest, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "上传的不是图片: "},
		})
		return
	}
	logs.LG.Info("上传图片格式为:" + form + ", 文件大小: " + strconv.FormatInt(int64(len(data)), 10))
	id := c.GetInt64("user_id")
	idStr := strconv.FormatInt(id, 10)

	// 使用原始文件名（如果有），否则使用自定义文件名
	originalFilename := file.Filename
	logs.LG.Info("原始文件名: " + originalFilename)
	if originalFilename == "" {
		originalFilename = idStr + "img"
		logs.LG.Info("使用自定义文件名: " + originalFilename)
	}

	url, e := upload.Uploader.UploadImageWithContext(ctx, data, originalFilename)
	if e != nil {
		logs.LG.Error("上传图片失败: " + e.Error())
		c.JSON(http.StatusBadRequest, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "上传图片失败: " + e.Error()},
		})
		return
	}
	req := &pbuser.UploadAvatarRequest{Url: url, Id: id}
	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.UploadAvatar(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.UploadAvatarReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMfaQrcode 获取 MFA 二维码
func GetMfaQrcode(ctx context.Context, c *app.RequestContext) {
	req := &pbuser.GetMfaQrcodeRequest{
		Id:       c.GetInt64("user_id"),
		Username: c.GetString("username"),
	}
	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.GetMfaQrcode(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.GetMfaQrcodeReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BindMfa 绑定 MFA
func BindMfa(ctx context.Context, c *app.RequestContext) {
	// 直接从表单获取参数（根据 API 文档：code, secret）
	req := &pbuser.BindMfaRequest{
		Code:   c.PostForm("code"),
		Secret: c.PostForm("secret"),
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.BindMfa(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.BindMfaReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchByImage 以图搜图
func SearchByImage(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.SearchByImageReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "文件上传失败: " + err.Error()},
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.SearchByImageReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "打开文件失败: " + err.Error()},
		})
		return
	}
	defer src.Close()

	data := make([]byte, file.Size)
	_, err = src.Read(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, pbuser.SearchByImageReply{
			Base: &pbcommon.BaseResponse{Code: model.Failed, Msg: "读取文件失败: " + err.Error()},
		})
		return
	}

	// 解析分页参数（从表单获取）
	pageNum, _ := strconv.ParseInt(c.PostForm("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.PostForm("page_size"), 10, 32)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req := &pbuser.SearchByImageRequest{
		Data: data,
		Page: &pbcommon.PageRequest{
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		},
	}

	clientMgr := client.GetClientManager()
	resp, err := clientMgr.UserClient.SearchByImage(ctx, req)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusInternalServerError, pbuser.SearchByImageReply{
			Base: &pbcommon.BaseResponse{Code: code, Msg: msg},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
