package utils

import (
	"errors"

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

	fieldLabels := map[string]string{
		"account_id":      "Account",
		"category_id":     "Category",
		"amount":          "Amount",
		"title":           "Title",
		"type":            "Type",
		"date":            "Date",
		"notes":           "Notes",
		"name":            "Name",
		"email":           "Email",
		"password":        "Password",
		"from_account_id": "Source account",
		"to_account_id":   "Destination account",
	}

	tagMessages := map[string]string{
		"required": "is required",
		"oneof":    "is invalid",
		"gt":       "must be greater than 0",
		"gte":      "must be at least 0",
		"email":    "must be a valid email address",
		"uuid":     "is not a valid identifier",
		"min":      "is too short",
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			label, ok := fieldLabels[field]
			if !ok {
				label = field
			}
			if msg, ok := tagMessages[tag]; ok {
				errorsMap[field] = label + " " + msg
			} else {
				errorsMap[field] = label + " is invalid"
			}
		}
		return errorsMap
	}
	return nil
}
