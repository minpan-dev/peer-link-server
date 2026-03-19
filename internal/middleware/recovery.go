package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"peer-link-server/pkg/response"
	"runtime/debug"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered",
					zap.Any("error", rec),
					zap.ByteString("stack", debug.Stack()),
				)
				response.Fail(c, http.StatusInternalServerError, 50000, "internal server error")
			}
		}()
		c.Next()
	}
}
