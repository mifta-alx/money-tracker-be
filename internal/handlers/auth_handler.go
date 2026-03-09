package handlers

import (
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
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
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
		Email    string `json:"email"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		GoogleID string `json:"google_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.service.HandleLoginGoogle(c.Request.Context(), input.Email, input.Name, input.Avatar, input.GoogleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "user": user})
}
