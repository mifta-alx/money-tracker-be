package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uuid.UUID) (string, string, error) {
	accessExpiryStr := os.Getenv("JWT_ACCESS_EXPIRY")
	if accessExpiryStr == "" {
		accessExpiryStr = "15m"
	}

	refreshExpiryStr := os.Getenv("JWT_REFRESH_EXPIRY")
	if refreshExpiryStr == "" {
		refreshExpiryStr = "168h"
	}

	accessDuration, _ := time.ParseDuration(accessExpiryStr)
	refreshDuration, _ := time.ParseDuration(refreshExpiryStr)

	accessTokenClaims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	at, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	rt, err := refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}
