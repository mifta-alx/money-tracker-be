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

type BudgetAllocationHandler struct {
	service *services.BudgetAllocationService
}

func NewBudgetAllocationHandler(s *services.BudgetAllocationService) *BudgetAllocationHandler {
	return &BudgetAllocationHandler{
		service: s,
	}
}

func (h *BudgetAllocationHandler) CreateBudgetAllocation(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	userID := val.(uuid.UUID)
	var req struct {
		Name         string `json:"name" binding:"required"`
		Percentage   int64  `json:"percentage" binding:"required,min=0,max=100"`
		TargetAmount int64  `json:"target_amount" binding:"required,min=0"`
		Period       string `json:"period" binding:"required,oneof=weekly monthly yearly"`
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

	newBudget := &models.BudgetAllocation{
		UserID:       userID,
		Name:         req.Name,
		Percentage:   req.Percentage,
		TargetAmount: req.TargetAmount,
		Period:       req.Period,
	}

	budget, err := h.service.CreateBudgetAllocation(c.Request.Context(), newBudget)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}
	response := struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Percentage   int64     `json:"percentage"`
		TargetAmount int64     `json:"target_amount"`
		Period       string    `json:"period"`
		CreatedAt    time.Time `json:"created_at"`
	}{
		ID:           budget.ID,
		Name:         budget.Name,
		Percentage:   budget.Percentage,
		TargetAmount: budget.TargetAmount,
		Period:       budget.Period,
		CreatedAt:    budget.CreatedAt,
	}

	utils.JSON(c, http.StatusCreated, "Budget created successfully", response)
}

func (h *BudgetAllocationHandler) UpdateBudgetAllocation(c *gin.Context) {
	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	var req struct {
		Name         string `json:"name" binding:"required"`
		Percentage   int64  `json:"percentage" binding:"required,min=0,max=100"`
		TargetAmount int64  `json:"target_amount" binding:"required,min=0"`
		Period       string `json:"period" binding:"required,oneof=weekly monthly yearly"`
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

	updatedBudget := &models.BudgetAllocation{
		ID:           budgetID,
		UserID:       userID,
		Name:         req.Name,
		Percentage:   req.Percentage,
		TargetAmount: req.TargetAmount,
		Period:       req.Period,
	}

	budget, err := h.service.UpdateBudgetAllocation(c.Request.Context(), updatedBudget)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrBudgetNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	response := struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Percentage   int64     `json:"percentage"`
		TargetAmount int64     `json:"target_amount"`
		Period       string    `json:"period"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{
		ID:           budget.ID,
		Name:         budget.Name,
		Percentage:   budget.Percentage,
		TargetAmount: budget.TargetAmount,
		Period:       budget.Period,
		UpdatedAt:    budget.UpdatedAt,
	}

	utils.JSON(c, http.StatusOK, "Budget updated successfully", response)
}

func (h *BudgetAllocationHandler) DeleteBudgetAllocation(c *gin.Context) {
	budgetID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	err = h.service.DeleteBudgetAllocation(c.Request.Context(), budgetID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrBudgetNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Budget deleted successfully", nil)
}

func (h *BudgetAllocationHandler) GetBudgetAllocations(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	budget, err := h.service.GetBudgetAllocations(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Budgets retrieved successfully", budget)
}

func (h *BudgetAllocationHandler) GetBudgetAllocation(c *gin.Context) {
	budgetIDStr := c.Param("id")
	budgetID, err := uuid.Parse(budgetIDStr)

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrBudgetNotFound), nil)
		return
	}

	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	budget, err := h.service.GetBudgetAllocation(c.Request.Context(), budgetID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrBudgetNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		}

		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Budget retrieved successfully", budget)
}
