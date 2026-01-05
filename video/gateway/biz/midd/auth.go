package midd

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strconv"
	"strings"
	"west2-video/common/jwts"
	"west2-video/gateway/biz/model"
	"west2-video/gateway/config"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从 Header 获取 Token
		token := string(c.Request.Header.Get("Access-Token"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code": model.NoToken,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 解析 Bearer Token
		if strings.Contains(token, "bearer ") {
			token = strings.Replace(token, "bearer ", "", 1)
		}

		userID, username, err := jwts.ParseToken(token, config.C.JWT.AccessSecret, config.C.JWT.RefreshSecret, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code": model.TokenIsError,
				"msg":  "认证令牌无效: " + err.Error(),
			})
			c.Abort()
			return
		}
		id, _ := strconv.ParseInt(userID, 10, 64)
		// 将用户信息存储到上下文
		c.Set("user_id", id)
		c.Set("username", username)
		c.Next(ctx)
	}
}
