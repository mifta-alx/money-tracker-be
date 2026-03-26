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

type TransferHandler struct {
	service *services.TransferService
}

func NewTransferHandler(s *services.TransferService) *TransferHandler {
	return &TransferHandler{service: s}
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	userID := val.(uuid.UUID)
	var req struct {
		FromAccountID uuid.UUID `json:"from_account_id" binding:"required"`
		ToAccountID   uuid.UUID `json:"to_account_id" binding:"required"`
		Amount        *int64    `json:"amount" binding:"required,min=0"`
		Notes         *string   `json:"notes"`
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
	newTransfer := &models.Transfer{
		UserID:        userID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        *req.Amount,
		Notes:         req.Notes,
	}

	transfer, err := h.service.CreateTransfer(c.Request.Context(), newTransfer)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}

	response := struct {
		ID            uuid.UUID `json:"id"`
		FromAccountID uuid.UUID `json:"from_account_id"`
		ToAccountID   uuid.UUID `json:"to_account_id"`
		Amount        int64     `json:"amount"`
		Notes         *string   `json:"notes"`
		CreatedAt     time.Time `json:"created_at"`
	}{
		ID:            transfer.ID,
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		Notes:         transfer.Notes,
		CreatedAt:     transfer.CreatedAt,
	}
	utils.JSON(c, http.StatusCreated, "Transfer created successfully", response)
}

func (h *TransferHandler) UpdateTransfer(c *gin.Context) {
	transferID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	var req struct {
		FromAccountID uuid.UUID `json:"from_account_id" binding:"required"`
		ToAccountID   uuid.UUID `json:"to_account_id" binding:"required"`
		Amount        *int64    `json:"amount" binding:"required,min=0"`
		Notes         *string   `json:"notes"`
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

	updatedTransfer := &models.Transfer{
		ID:            transferID,
		UserID:        userID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        *req.Amount,
		Notes:         req.Notes,
	}

	transfer, err := h.service.UpdateTransfer(c.Request.Context(), updatedTransfer)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransferNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transfer updated successfully", transfer)
}

func (h *TransferHandler) DeleteTransfer(c *gin.Context) {
	transferID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrMalformedRequest), nil)
		return
	}

	val, _ := c.Get(utils.UserIDKey)
	userID := val.(uuid.UUID)

	err = h.service.DeleteTransfer(c.Request.Context(), transferID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransferNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Transfer deleted successfully", nil)
}

func (h *TransferHandler) GetTransfers(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	transfers, err := h.service.GetTransfers(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transfers retrieved successfully", transfers)
}

func (h *TransferHandler) GetTransfer(c *gin.Context) {
	transferIDStr := c.Param("id")
	transferID, err := uuid.Parse(transferIDStr)

	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.TranslateError(services.ErrTransferNotFound), nil)
		return
	}

	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}

	transfer, err := h.service.GetTransfer(c.Request.Context(), transferID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrTransferNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, services.ErrUnauthorized) {
			statusCode = http.StatusForbidden
		}

		utils.Error(c, statusCode, utils.TranslateError(err), nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Transfer retrieved successfully", transfer)
}
