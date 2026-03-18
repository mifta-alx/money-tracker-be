package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, errs interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func FormatValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			fieldName := strings.ToLower(fe.Field())
			var msg string
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldName)
			case "email":
				msg = fmt.Sprintf("%s must be a valid email", fieldName)
			case "oneof":
				msg = fmt.Sprintf("%s must be one of: %s", fieldName, fe.Param())
			case "max":
				msg = fmt.Sprintf("%s must not exceed %s characters", fieldName, fe.Param())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", fieldName, fe.Param())
			default:
				msg = fmt.Sprintf("%s is invalid", fieldName)
			}
			errorsMap[strings.ToLower(fieldName)] = msg
		}
		return errorsMap
	}
	return nil
}
