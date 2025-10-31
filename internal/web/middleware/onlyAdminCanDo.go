package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 在处理请求前
		path := c.Request.URL.Path
		method := c.Request.Method

		// 调用后续中间件和最终的 handler
		c.Next()

		// 在处理请求后
		latency := time.Since(start)
		status := c.Writer.Status()

		fmt.Printf("[GIN-TEST] %s %s | %d | %v", method, path, status, latency)
	}
}
