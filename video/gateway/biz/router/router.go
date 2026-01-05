package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"west2-video/gateway/biz/handler"
	"west2-video/gateway/biz/midd"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz) {
	// 全局中间件
	h.Use(midd.LogMiddleware())

	// API 路由组
	api := h.Group("/api")
	{
		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/register", handler.Register)                                 // 注册
			user.POST("/login", handler.Login)                                       // 登录
			user.POST("/refresh", handler.RefreshToken)                               // 刷新 Token
			user.GET("/info", midd.AuthMiddleware(), handler.GetUserInfo)            // 用户信息
			user.POST("/avatar", midd.AuthMiddleware(), handler.UploadAvatar)        // 上传头像
			user.GET("/mfa/qrcode", midd.AuthMiddleware(), handler.GetMfaQrcode)     // 获取 MFA qrcode
			user.POST("/mfa/bind", midd.AuthMiddleware(), handler.BindMfa)           // 绑定 MFA
			user.POST("/search/image", midd.AuthMiddleware(), handler.SearchByImage) // 以图搜图
		}

		// 视频相关路由
		video := api.Group("/video")
		{
			video.GET("/feed", handler.Feed)                               // 视频流
			video.POST("/publish", midd.AuthMiddleware(), handler.Publish) // 投稿
			video.GET("/publish/list", handler.PublishList)                // 发布列表
			video.GET("/hot/ranking", handler.HotRanking)                  // 热门排行榜
			video.GET("/search", handler.SearchVideo)                      // 搜索视频
		}

		// 互动相关路由
		interaction := api.Group("/interaction")
		{
			interaction.POST("/like", midd.AuthMiddleware(), handler.LikeAction)         // 点赞操作
			interaction.GET("/like/list", handler.LikeList)                              // 点赞列表
			interaction.POST("/comment", midd.AuthMiddleware(), handler.CommentAction)   // 评论
			interaction.GET("/comment/list", handler.CommentList)                        // 评论列表
			interaction.DELETE("/comment", midd.AuthMiddleware(), handler.DeleteComment) // 删除评论
		}

		// 社交相关路由
		social := api.Group("/social")
		{
			social.POST("/follow", midd.AuthMiddleware(), handler.FollowAction)   // 关注操作
			social.GET("/follow/list", handler.FollowList)                        // 关注列表
			social.GET("/follower/list", handler.FollowerList)                    // 粉丝列表
			social.GET("/friend/list", midd.AuthMiddleware(), handler.FriendList) // 好友列表
		}
	}

	// 健康检查
	h.GET("/ping", handler.Ping)
}



