package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func FormatValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			var msg string
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("Field %s is required", fe.Field())
			case "email":
				msg = fmt.Sprintf("Field %s must be a valid email", fe.Field())
			case "max":
				msg = fmt.Sprintf("Field %s must not exceed %s characters", fe.Field(), fe.Param())
			case "min":
				msg = fmt.Sprintf("Field %s must be at least %s characters", fe.Field(), fe.Param())
			default:
				msg = fmt.Sprintf("Field %s is invalid", fe.Field())
			}
			errorsMap[strings.ToLower(fe.Field())] = msg
		}
	}
	return errorsMap
}
