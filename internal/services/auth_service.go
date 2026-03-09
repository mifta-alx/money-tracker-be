package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/pkg/utils"
	"money-tracker/internal/repository"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(r *repository.AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) (*models.User, error) {
	existingUser, _ := s.repo.GetUserByEmail(ctx, email)

	if existingUser != nil {
		return nil, errors.New("email_exists")
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	pepper := os.Getenv("AUTH_PEPPER")
	passwordWithPepper := password + pepper
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithPepper), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("process_password_failed")
	}

	now := time.Now()
	newUser := &models.User{
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.CreateUser(ctx, newUser, string(hashPassword))
	if err != nil {
		return nil, errors.New("create_account_failed")
	}

	return newUser, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password_too_short")
	}

	if strings.Contains(password, " ") {
		return errors.New("password_contains_space")
	}

	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	if !hasNumber {
		return errors.New("password_no_number")
	}

	matchUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !matchUpper {
		return errors.New("password_no_upper")
	}

	matchLower, _ := regexp.MatchString(`[a-z]`, password)
	if !matchLower {
		return errors.New("password_no_lower")
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *models.User, error) {
	user, hashPassword, err := s.repo.GetUserWithPassword(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("invalid_credentials")
	}

	pepper := os.Getenv("AUTH_PEPPER")
	passwordWithPepper := password + pepper

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordWithPepper))
	if err != nil {
		return "", "", nil, errors.New("invalid_credentials")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, errors.New("internal_server_error")
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) HandleLoginGoogle(ctx context.Context, email, name, avatar, googleID string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, sql.ErrNoRows) {
		newUser := &models.User{
			Email:     email,
			Name:      name,
			AvatarURL: avatar,
		}

		newProvider := &models.AuthProvider{
			Provider:       "google",
			ProviderUserID: googleID,
		}

		err := s.repo.CreateUserOAuth(ctx, newUser, newProvider)

		if err != nil {
			return nil, err
		}

		return newUser, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
