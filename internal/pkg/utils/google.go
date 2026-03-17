package utils

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

type GoogleUser struct {
	Email   string
	Name    string
	Picture string
	Sub     string
}

func VerifyGoogleToken(token string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(context.Background(), token, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	return &GoogleUser{
		Email:   payload.Claims["email"].(string),
		Name:    payload.Claims["name"].(string),
		Picture: payload.Claims["picture"].(string),
		Sub:     payload.Subject,
	}, nil
}
