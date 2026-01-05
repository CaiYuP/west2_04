package midd

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// LogMiddleware 日志中间件
func LogMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Path())
		method := string(c.Method())

		c.Next(ctx)

		latency := time.Since(start)
		statusCode := c.Response.StatusCode()

		// 日志记录（可以在这里添加日志输出）
		_ = utils.H{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  c.ClientIP(),
			"user_agent": string(c.UserAgent()),
		}
	}
}
