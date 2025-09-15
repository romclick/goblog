package middleware

import (
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Log() gin.HandlerFunc {
	Logger, _ := zap.NewDevelopment()
	defer Logger.Sync()
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logger.Info(
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Any("user_id", c.MustGet("userID")),
		)
	}
}
