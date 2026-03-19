package router

import (
	"peer-link-server/internal/handler"
	"peer-link-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	Health *handler.HealthHandler
	User   *handler.UserHandler
}

func New(h *Handlers, log *zap.Logger) *gin.Engine {
	// 不使用 gin.Default()，完全自控中间件
	r := gin.New()

	// 全局中间件（顺序很重要）
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))
	r.Use(middleware.CORS())

	// 健康检查（不需要鉴权）
	r.GET("/health", h.Health.Ping)

	// API v1
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		users.GET("", h.User.List)
		users.POST("", h.User.Create)
		users.GET("/:id", h.User.Get)
		users.PUT("/:id", h.User.Update)
		users.DELETE("/:id", h.User.Delete)
	}

	return r
}
