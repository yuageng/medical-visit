package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response 统一 API 响应结构
type Response struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      "SUCCESS",
		Message:   "操作成功",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// Created 返回创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      "CREATED",
		Message:   "创建成功",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

// BadRequest 返回 400 错误
func BadRequest(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:      "BAD_REQUEST",
		Message:   message,
		Details:   details,
		RequestID: getRequestID(c),
	})
}

// NotFound 返回 404 错误
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:      "NOT_FOUND",
		Message:   message,
		RequestID: getRequestID(c),
	})
}

// Conflict 返回 409 错误
func Conflict(c *gin.Context, code, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Code:      code,
		Message:   message,
		RequestID: getRequestID(c),
	})
}

// UnprocessableEntity 返回 422 错误
func UnprocessableEntity(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Code:      "UNPROCESSABLE_ENTITY",
		Message:   message,
		Details:   details,
		RequestID: getRequestID(c),
	})
}

// InternalServerError 返回 500 错误
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:      "INTERNAL_SERVER_ERROR",
		Message:   message,
		RequestID: getRequestID(c),
	})
}

// getRequestID 获取或生成请求 ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return uuid.New().String()
}
