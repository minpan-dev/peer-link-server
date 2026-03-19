package handler

import (
	"peer-link-server/internal/service"
	"peer-link-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type LiveKitHandler struct {
	svc service.LiveKitService
}

func NewLiveKitHandler(svc service.LiveKitService) *LiveKitHandler {
	return &LiveKitHandler{svc: svc}
}

// POST /api/v1/rooms
func (h *LiveKitHandler) CreateRoom(c *gin.Context) {
	var req service.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	room, err := h.svc.CreateRoom(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, room)
}

// GET /api/v1/rooms/:roomName/token
func (h *LiveKitHandler) GetToken(c *gin.Context) {
	roomName := c.Param("roomName")
	// identity 从 JWT 中间件注入，退化到 query 参数
	identity, exists := c.Get("user_id")
	if !exists {
		identity = c.Query("identity")
	}
	token, err := h.svc.GetToken(c.Request.Context(), roomName, identity.(string))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, token)
}

// GET /api/v1/rooms
func (h *LiveKitHandler) ListRooms(c *gin.Context) {
	page, pageSize := parsePagination(c)
	rooms, total, err := h.svc.ListRooms(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"total": total, "items": rooms})
}

// DELETE /api/v1/rooms/:roomName
func (h *LiveKitHandler) DeleteRoom(c *gin.Context) {
	if err := h.svc.DeleteRoom(c.Request.Context(), c.Param("roomName")); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// GET /api/v1/rooms/:roomName/participants
func (h *LiveKitHandler) ListParticipants(c *gin.Context) {
	participants, err := h.svc.ListParticipants(c.Request.Context(), c.Param("roomName"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, participants)
}

// DELETE /api/v1/rooms/:roomName/participants/:identity
func (h *LiveKitHandler) RemoveParticipant(c *gin.Context) {
	err := h.svc.RemoveParticipant(c.Request.Context(),
		c.Param("roomName"), c.Param("identity"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

func parsePagination(c *gin.Context) (int, int) {
	page, pageSize := 1, 20
	// 复用你项目里已有的分页逻辑即可
	return page, pageSize
}
