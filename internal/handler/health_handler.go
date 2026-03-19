package handler

import (
	"github.com/gin-gonic/gin"
	"peer-link-server/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Ping(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}
