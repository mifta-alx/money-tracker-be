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

	fieldLabels := map[string]string{
		"account_id":       "Account",
		"category_id":      "Category",
		"amount":           "Amount",
		"title":            "Title",
		"type":             "Type",
		"date":             "Date",
		"notes":            "Notes",
		"name":             "Name",
		"email":            "Email",
		"password":         "Password",
		"balance":          "Balance",
		"icon":             "Icon",
		"color":            "Color",
		"confirm_password": "Confirm password",
		"from_account_id":  "Source account",
		"to_account_id":    "Destination account",
	}

	tagMessages := map[string]string{
		"required":    "is required",
		"oneof":       "is invalid",
		"gt":          "must be greater than 0",
		"gte":         "must be at least 0",
		"email":       "must be a valid email address",
		"uuid":        "is not a valid identifier",
		"min":         "must be at least %s characters",
		"max":         "must be at most %s characters",
		"eqfield":     "must match %s",
		"containsany": "must contain at least one uppercase, lowercase, and number",
		"numeric":     "must be a valid number",
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			param := fe.Param()

			label, ok := fieldLabels[field]
			if !ok {
				label = field
			}
			msg, ok := tagMessages[tag]
			if !ok {
				msg = "is invalid"
			}

			if strings.Contains(msg, "%s") {
				if tag == "eqfield" {
					if pLabel, exists := fieldLabels[param]; exists {
						param = strings.ToLower(pLabel)
					}
				}
				errorsMap[field] = fmt.Sprintf("%s %s", label, fmt.Sprintf(msg, param))
			} else {
				errorsMap[field] = fmt.Sprintf("%s %s", label, msg)
			}
		}
		return errorsMap
	}
	return nil
}
