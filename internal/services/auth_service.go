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

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(r *repository.AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) (*models.User, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInternal
	}

	if existingUser != nil {
		return nil, ErrEmailExist
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	pepper := os.Getenv("AUTH_PEPPER")
	passwordWithPepper := password + pepper
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithPepper), bcrypt.DefaultCost)

	if err != nil {
		return nil, ErrInternal
	}

	newUser := &models.User{
		Name:  name,
		Email: email,
	}

	err = s.repo.CreateUser(ctx, newUser, string(hashPassword))
	if err != nil {
		return nil, ErrInternal
	}

	return newUser, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if strings.Contains(password, " ") {
		return ErrPasswordContainsSpace
	}

	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	if !hasNumber {
		return ErrPasswordNoNumber
	}

	matchUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !matchUpper {
		return ErrPasswordNoUpper
	}

	matchLower, _ := regexp.MatchString(`[a-z]`, password)
	if !matchLower {
		return ErrPasswordNoLower
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *models.User, error) {
	user, hashPassword, err := s.repo.GetUserWithPassword(ctx, email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	pepper := os.Getenv("AUTH_PEPPER")
	passwordWithPepper := password + pepper

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordWithPepper))
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, ErrInternal
	}

	err = s.repo.UpdateLastSign(context.Background(), user.ID)
	return accessToken, refreshToken, user, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, email, name, avatar, googleID string) (string, string, *models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, sql.ErrNoRows) {
		newUser := &models.User{
			Email:     email,
			Name:      name,
			AvatarURL: &avatar,
		}

		newProvider := &models.AuthProvider{
			Provider:       "google",
			ProviderUserID: googleID,
		}

		err := s.repo.CreateUserOAuth(ctx, newUser, newProvider)
		if err != nil {
			return "", "", nil, ErrFailedCreateOAuth
		}

		return s.generateAuthTokens(ctx, newUser)
	}

	if err != nil {
		return "", "", nil, err
	}

	exist, err := s.repo.CheckProviderExists(ctx, user.ID, "google")
	if err != nil {
		return "", "", nil, err
	}

	if !exist {
		err = s.repo.AddAuthProvider(ctx, user.ID, "google", googleID)
		if err != nil {
			return "", "", nil, ErrFailedLinkGoogle
		}
	}

	return s.generateAuthTokens(ctx, user)
}

func (s *AuthService) generateAuthTokens(ctx context.Context, user *models.User) (string, string, *models.User, error) {
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, ErrInternal
	}
	_ = s.repo.UpdateLastSign(ctx, user.ID)
	return accessToken, refreshToken, user, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
