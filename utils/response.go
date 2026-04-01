package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Error(code, message))
}
