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
		Name    string `json:"name" binding:"required"`
		Type    string `json:"type" binding:"required,oneof=Bank E-Wallet Cash"`
		Balance *int64 `json:"balance" binding:"required"`
		Icon    string `json:"icon" binding:"required"`
		Color   string `json:"color" binding:"required"`
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
		UserID:  userID,
		Name:    req.Name,
		Type:    req.Type,
		Balance: *req.Balance,
		Icon:    req.Icon,
		Color:   req.Color,
	}

	account, err := h.service.CreateAccount(c.Request.Context(), newAccount)
	errorMessage := utils.TranslateError(err)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errorMessage, nil)
		return
	}

	utils.JSON(c, http.StatusCreated, "Account created successfully", account)
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	val, exists := c.Get(utils.UserIDKey)
	userID, ok := val.(uuid.UUID)
	if !exists || !ok {
		utils.Error(c, http.StatusUnauthorized, utils.TranslateError(services.ErrUnauthorized), nil)
		return
	}
	accounts, err := h.service.GetAccounts(c.Request.Context(), userID)
	if err != nil {
		errorMessage := utils.TranslateError(err)
		utils.Error(c, http.StatusInternalServerError, errorMessage, nil)
		return
	}
	utils.JSON(c, http.StatusOK, "Accounts retrieved successfully", accounts)
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
