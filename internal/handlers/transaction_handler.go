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

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(s *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	userID := val.(uuid.UUID)
	var req struct {
		AccountID  uuid.UUID  `json:"account_id" binding:"required"`
		CategoryID *uuid.UUID `json:"category_id"`
		Title      string     `json:"title" binding:"required"`
		Type       string     `json:"type" binding:"required,oneof=expense income"`
		Amount     *int64     `json:"amount" binding:"required,min=0"`
		Date       time.Time  `json:"date" binding:"required"`
		Notes      *string    `json:"notes"`
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

	newTransaction := &models.Transaction{
		UserID:     userID,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Type:       req.Type,
		Amount:     *req.Amount,
		Date:       req.Date,
		Notes:      req.Notes,
	}

	transaction, err := h.service.CreateTransaction(c.Request.Context(), newTransaction)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}

	response := struct {
		ID         uuid.UUID  `json:"id"`
		AccountID  uuid.UUID  `json:"account_id"`
		CategoryID *uuid.UUID `json:"category_id"`
		Title      string     `json:"title"`
		Type       string     `json:"type"`
		Amount     int64      `json:"amount"`
		Date       time.Time  `json:"date"`
		Notes      *string    `json:"notes"`
		CreatedAt  time.Time  `json:"created_at"`
	}{
		ID:         transaction.ID,
		AccountID:  transaction.AccountID,
		CategoryID: transaction.CategoryID,
		Title:      transaction.Title,
		Type:       transaction.Type,
		Amount:     transaction.Amount,
		Date:       transaction.Date,
		Notes:      transaction.Notes,
		CreatedAt:  transaction.CreatedAt,
	}
	utils.JSON(c, http.StatusCreated, "Transaction created successfully", response)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	var req struct {
		AccountID  uuid.UUID  `json:"account_id" binding:"required"`
		CategoryID *uuid.UUID `json:"category_id"`
		Title      string     `json:"title" binding:"required"`
		Type       string     `json:"type" binding:"required,oneof=expense income"`
		Amount     *int64     `json:"amount" binding:"required,min=0"`
		Date       time.Time  `json:"date" binding:"required"`
		Notes      *string    `json:"notes"`
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

	updatedTransaction := &models.Transaction{
		ID:         transactionID,
		UserID:     userID,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Type:       req.Type,
		Amount:     *req.Amount,
		Date:       req.Date,
		Notes:      req.Notes,
	}

	transaction, err := h.service.UpdateTransaction(c.Request.Context(), updatedTransaction)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransactionNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transaction updated successfully", transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	err = h.service.DeleteTransaction(c.Request.Context(), transactionID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransactionNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	transactions, err := h.service.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)

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

	transaction, err := h.service.GetTransaction(c.Request.Context(), transactionID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransactionNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		}

		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}
