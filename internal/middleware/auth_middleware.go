package middleware

import (
	"net/http"
	"os"
	"strings"

	"money-tracker/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "Invalid token format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userIDRaw, ok := claims["user_id"].(string)
			if !ok {
				utils.Error(c, http.StatusUnauthorized, "Invalid user ID in token", nil)
				c.Abort()
				return
			}

			userID, err := uuid.Parse(userIDRaw)
			if err != nil {
				utils.Error(c, http.StatusUnauthorized, "User ID format in token is not a valid UUID", nil)
				c.Abort()
				return
			}
			c.Set(utils.UserIDKey, userID)
		}

		c.Next()
	}
}
