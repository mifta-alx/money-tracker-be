package handlers

import (
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/pkg/utils"
	"money-tracker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(s *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: s,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	userID := val.(uuid.UUID)
	var req struct {
		AllocationID uuid.UUID `json:"allocation_id"`
		Name         string    `json:"name" binding:"required"`
		Type         string    `json:"type" binding:"required,oneof=expense income"`
		Icon         string    `json:"icon" binding:"required"`
		Color        string    `json:"color" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			utils.Error(c, http.StatusUnprocessableEntity, utils.TranslateError(services.ErrValidation), utils.FormatValidationError(ve))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	newCategory := &models.Category{
		UserID:       userID,
		AllocationID: req.AllocationID,
		Name:         req.Name,
		Type:         req.Type,
		Icon:         req.Icon,
		Color:        req.Color,
	}

	category, err := h.service.CreateCategory(c.Request.Context(), newCategory)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	categories, err := h.service.GetCategories(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Categories retrieved successfully", categories)
}
