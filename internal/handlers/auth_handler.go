package handlers

import (
	"errors"
	"money-tracker/internal/pkg/utils"
	"money-tracker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email           string `json:"email" binding:"required,email"`
		Name            string `json:"name" binding:"required"`
		Password        string `json:"password" binding:"required,min=8,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			utils.Error(c, http.StatusUnprocessableEntity, utils.TranslateError(services.ErrValidation), utils.FormatValidationError(err))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	user, err := h.service.Register(c.Request.Context(), input.Email, input.Name, input.Password)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}

	utils.JSON(c, http.StatusCreated, "User created successfully", gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			utils.Error(c, http.StatusUnprocessableEntity, utils.TranslateError(services.ErrValidation), utils.FormatValidationError(ve))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	accessToken, refreshToken, user, err := h.service.Login(c.Request.Context(), input.Email, input.Password)

	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, services.ErrInternal) {
			status = http.StatusInternalServerError
		}
		utils.Error(c, status, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Login successful", gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"avatar_url":   user.AvatarURL,
		"last_sign_in": user.LastSignInAt,
	}})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrTokenRequired), nil)
		return
	}

	googleUser, err := utils.VerifyGoogleToken(input.Token)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrInvalidGoogleToken), nil)
		return
	}

	accessToken, refreshToken, user, err := h.service.LoginWithGoogle(c.Request.Context(),
		googleUser.Email,
		googleUser.Name,
		googleUser.Picture,
		googleUser.Sub,
	)

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Login successful", gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": user})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	userID := val.(uuid.UUID)
	user, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		utils.Error(c, status, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "User profile retrieved successfully", user)
}
