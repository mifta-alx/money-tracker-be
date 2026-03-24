package handlers

import (
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/pkg/utils"
	"money-tracker/internal/services"
	"net/http"
	"time"

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
		AllocationID *uuid.UUID `json:"allocation_id"`
		Name         string     `json:"name" binding:"required"`
		Type         string     `json:"type" binding:"required,oneof=expense income"`
		Icon         string     `json:"icon" binding:"required"`
		Color        string     `json:"color" binding:"required"`
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

	response := struct {
		ID           uuid.UUID  `json:"id"`
		AllocationID *uuid.UUID `json:"allocation_id"`
		Name         string     `json:"name"`
		Type         string     `json:"type"`
		Icon         string     `json:"icon"`
		Color        string     `json:"color"`
		CreatedAt    time.Time  `json:"created_at"`
	}{
		ID:           category.ID,
		AllocationID: category.AllocationID,
		Name:         category.Name,
		Type:         category.Type,
		Icon:         category.Icon,
		Color:        category.Color,
		CreatedAt:    category.CreatedAt,
	}
	utils.JSON(c, http.StatusCreated, "Category created successfully", response)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	var req struct {
		AllocationID *uuid.UUID `json:"allocation_id"`
		Name         string     `json:"name" binding:"required"`
		Type         string     `json:"type" binding:"required,oneof=expense income"`
		Icon         string     `json:"icon" binding:"required"`
		Color        string     `json:"color" binding:"required"`
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

	updatedCategory := &models.Category{
		ID:           categoryID,
		UserID:       userID,
		AllocationID: req.AllocationID,
		Name:         req.Name,
		Type:         req.Type,
		Icon:         req.Icon,
		Color:        req.Color,
	}

	category, err := h.service.UpdateCategory(c.Request.Context(), updatedCategory)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrCategoryNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	response := struct {
		ID           uuid.UUID  `json:"id"`
		AllocationID *uuid.UUID `json:"allocation_id"`
		Name         string     `json:"name"`
		Type         string     `json:"type"`
		Icon         string     `json:"icon"`
		Color        string     `json:"color"`
		UpdatedAt    time.Time  `json:"updated_at"`
	}{
		ID:           category.ID,
		AllocationID: category.AllocationID,
		Name:         category.Name,
		Type:         category.Type,
		Icon:         category.Icon,
		Color:        category.Color,
		UpdatedAt:    category.UpdatedAt,
	}

	utils.JSON(c, http.StatusOK, "Category updated successfully", response)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	err = h.service.DeleteCategory(c.Request.Context(), categoryID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrCategoryNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Category deleted successfully", nil)
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

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrCategoryNotFound), nil)
		return
	}

	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	category, err := h.service.GetCategory(c.Request.Context(), categoryID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrCategoryNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		}

		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Category retrieved successfully", category)
}
