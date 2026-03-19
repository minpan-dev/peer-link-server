package response

import (
	"net/http"
	"peer-link-server/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Envelope 是所有 API 响应的统一格式
type Envelope struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Code: 0, Message: "ok", Data: data, RequestID: rid(c)})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Code: 0, Message: "created", Data: data, RequestID: rid(c)})
}

// Error 自动识别 AppError，其余统一返回 500
func Error(c *gin.Context, err error) {
	ae, ok := errors.AsAppError(err)
	if !ok {
		ae = errors.ErrInternal
	}
	c.AbortWithStatusJSON(ae.HTTPCode, Envelope{Code: ae.Code, Message: ae.Message, RequestID: rid(c)})
}

func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Envelope{Code: code, Message: message, RequestID: rid(c)})
}

func rid(c *gin.Context) string {
	id, _ := c.Get("request_id")
	s, _ := id.(string)
	return s
}
