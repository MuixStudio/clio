package logger

import (
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

// Skipper is a function to skip logs based on provided Context
type Skipper func(c *gin.Context) bool
type LoggerConfig struct {

	// SkipPaths is a URL path array which logs are not written.
	// Optional.
	SkipPaths []string

	// Skip is a Skipper that indicates which logs should not be written.
	// Optional.
	Skip Skipper
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {

	notlogged := conf.SkipPaths

	//isTerm := true

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when it is not being skipped
		if _, ok := skip[path]; ok || (conf.Skip != nil && conf.Skip(c)) {
			return
		}
		//
		//param := LogFormatterParams{
		//	Request: c.Request,
		//	isTerm:  isTerm,
		//	Keys:    c.Keys,
		//}
		//
		//// Stop timer
		//param.TimeStamp = time.Now()
		//param.Latency = param.TimeStamp.Sub(start)
		//
		//param.ClientIP = c.ClientIP()
		//param.Method = c.Request.Method
		//param.StatusCode = c.Writer.Status()
		//param.ErrorMessage = c.Errors.String()
		//
		//param.BodySize = c.Writer.Size()
		//
		//if raw != "" {
		//	path = path + "?" + raw
		//}
		//
		//param.Path = path
		//klog.Ctx
		//klog.Control()
		klog.Infof("request: %s, keys: %v, start: %v, raw: %v", c.Request, c.Keys, start, raw)
	}
}
