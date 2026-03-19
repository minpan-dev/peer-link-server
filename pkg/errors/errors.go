package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 携带 HTTP 状态码 + 业务错误码 + 用户可读信息
type AppError struct {
	HTTPCode int
	Code     int
	Message  string
	Err      error // 原始错误，不对外暴露
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
func (e *AppError) Unwrap() error { return e.Err }

func New(httpCode, code int, message string) *AppError {
	return &AppError{HTTPCode: httpCode, Code: code, Message: message}
}
func Wrap(httpCode, code int, message string, err error) *AppError {
	return &AppError{HTTPCode: httpCode, Code: code, Message: message, Err: err}
}

// 预定义常用错误，可在 service 层直接返回
var (
	ErrNotFound   = New(http.StatusNotFound, 40400, "resource not found")
	ErrBadRequest = New(http.StatusBadRequest, 40000, "bad request")
	ErrUnauth     = New(http.StatusUnauthorized, 40100, "unauthorized")
	ErrForbidden  = New(http.StatusForbidden, 40300, "forbidden")
	ErrConflict   = New(http.StatusConflict, 40900, "resource already exists")
	ErrInternal   = New(http.StatusInternalServerError, 50000, "internal server error")
)

func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	return ae, errors.As(err, &ae)
}
