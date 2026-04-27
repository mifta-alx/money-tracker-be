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

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(s *services.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	userID := val.(uuid.UUID)
	var req struct {
		Name       string `json:"name" binding:"required"`
		Type       string `json:"type" binding:"required,oneof=bank e-wallet cash credit-card investment"`
		Balance    *int64 `json:"balance" binding:"required,numeric,min=0"`
		Icon       string `json:"icon" binding:"required"`
		Color      string `json:"color" binding:"required"`
		IsExcluded bool   `json:"is_excluded"`
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

	newAccount := &models.Account{
		UserID:     userID,
		Name:       req.Name,
		Type:       req.Type,
		Balance:    *req.Balance,
		Icon:       req.Icon,
		Color:      req.Color,
		IsExcluded: req.IsExcluded,
	}

	account, err := h.service.CreateAccount(c.Request.Context(), newAccount)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}
	response := struct {
		ID         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Type       string    `json:"type"`
		Balance    int64     `json:"balance"`
		Icon       string    `json:"icon"`
		Color      string    `json:"color"`
		IsExcluded bool      `json:"is_excluded"`
		CreatedAt  time.Time `json:"created_at"`
	}{
		ID:         account.ID,
		Name:       account.Name,
		Type:       account.Type,
		Balance:    account.Balance,
		Icon:       account.Icon,
		Color:      account.Color,
		IsExcluded: account.IsExcluded,
		CreatedAt:  account.CreatedAt,
	}

	utils.JSON(c, http.StatusCreated, "Account created successfully", response)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	var req struct {
		Name       string `json:"name" binding:"required"`
		Type       string `json:"type" binding:"required,oneof=bank e-wallet cash credit-card investment"`
		Balance    *int64 `json:"balance" binding:"required,numeric"`
		Icon       string `json:"icon" binding:"required"`
		Color      string `json:"color" binding:"required"`
		IsExcluded bool   `json:"is_excluded"`
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

	updatedAccount := &models.Account{
		ID:         accountID,
		UserID:     userID,
		Name:       req.Name,
		Type:       req.Type,
		Icon:       req.Icon,
		Color:      req.Color,
		IsExcluded: req.IsExcluded,
	}

	if req.Balance != nil {
		updatedAccount.Balance = *req.Balance
	} else {
		existingAccount, err := h.service.GetAccount(c.Request.Context(), accountID, userID)
		if err == nil {
			updatedAccount.Balance = existingAccount.Balance
		}
	}

	account, err := h.service.UpdateAccount(c.Request.Context(), updatedAccount)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrAccountNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Account updated successfully", account)
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	err = h.service.DeleteAccount(c.Request.Context(), accountID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrAccountNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Account deleted successfully", nil)
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	accounts, total, err := h.service.GetAccounts(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	response := gin.H{
		"total_balance": total,
		"accounts":      accounts,
	}
	utils.JSON(c, http.StatusOK, "Accounts retrieved successfully", response)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrAccountNotFound), nil)
		return
	}

	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	account, err := h.service.GetAccount(c.Request.Context(), accountID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrAccountNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		}

		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Account retrieved successfully", account)
}
