package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
		}
		if rid, ok := c.Get("request_id"); ok {
			fields = append(fields, zap.Any("request_id", rid))
		}
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error(e, fields...)
			}
		} else {
			log.Info("access", fields...)
		}
	}
}
