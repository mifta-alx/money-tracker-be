package handlers

import (
	"errors"
	"money-tracker/internal/pkg/utils"
	"money-tracker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		prettyError := utils.FormatValidationError(err)
		utils.Error(c, http.StatusBadRequest, "Validation failed", prettyError)
		return
	}

	user, err := h.service.Register(c.Request.Context(), input.Email, input.Name, input.Password)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}

	utils.JSON(c, http.StatusCreated, "User created successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		prettyError := utils.FormatValidationError(err)
		utils.Error(c, http.StatusBadRequest, "Validation failed", prettyError)
		return
	}

	accessToken, refreshToken, user, err := h.service.Login(c.Request.Context(), input.Email, input.Password)

	if err != nil {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(err), nil)
	}

	responseUser := struct {
		ID        string  `json:"id"`
		Email     string  `json:"email"`
		Name      string  `json:"name"`
		AvatarURL *string `json:"avatar_url"`
	}{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}

	utils.JSON(c, http.StatusOK, "Login successful", gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": responseUser})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Token is required", nil)
		return
	}

	googleUser, err := utils.VerifyGoogleToken(input.Token)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid Google token", nil)
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
	userID, exist := c.Get("user_id")
	if !exist {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(errors.New("unauthorized")), nil)
		return
	}

	user, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "User profile retrieved successfully", user)
}
