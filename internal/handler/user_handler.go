package handler

import (
	"net/http"
	"peer-link-server/internal/service"
	apperr "peer-link-server/pkg/errors"
	"peer-link-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ svc service.UserService }

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	users, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"list": users, "total": total, "page": page, "page_size": pageSize})
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := parseUint(c, "id")
	if err != nil {
		response.Error(c, apperr.ErrBadRequest)
		return
	}
	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	user, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseUint(c, "id")
	if err != nil {
		response.Error(c, apperr.ErrBadRequest)
		return
	}
	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	user, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseUint(c, "id")
	if err != nil {
		response.Error(c, apperr.ErrBadRequest)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

func parseUint(c *gin.Context, key string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(v), err
}
