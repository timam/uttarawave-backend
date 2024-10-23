package response

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type PaginationInfo struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
		Error:   err,
	})
}
