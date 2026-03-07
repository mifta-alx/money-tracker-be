package handlers

import (
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, title, amount, category, created_at FROM transactions ORDER BY created_at DESC`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var transactions []map[string]interface{}

	for rows.Next() {
		var id int
		var title string
		var amount float64
		var category string
		var createdAt time.Time

		err := rows.Scan(&id, &title, &amount, &category, &createdAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		transactions = append(transactions, gin.H{"id": id, "title": title, "amount": amount, "category": category, "created_at": createdAt})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	row := h.DB.QueryRow(`SELECT id, title, amount, category, created_at FROM transactions WHERE id = $1`, id)

	var t models.Transactions

	err = row.Scan(
		&t.ID,
		&t.Title,
		&t.Amount,
		&t.Category,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category string  `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO transactions (title, amount, category) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := h.DB.QueryRow(query, req.Title, req.Amount, req.Category).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category string  `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.DB.Exec(`UPDATE transactions SET title=$1, amount=$2, category=$3 WHERE id = $4`, req.Title, req.Amount, req.Category, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_, err = h.DB.Exec(`DELETE FROM transactions WHERE id = $1`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
