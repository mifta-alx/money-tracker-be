package routes

import (
	"money-tracker/internal/database"
	"money-tracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db := database.Connect()

	r := gin.Default()

	transactionHandler := handlers.NewTransactionHandler(db)

	transactions := r.Group("/transactions")
	{
		transactions.GET("", transactionHandler.GetAllTransactions)
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.GET("/:id", transactionHandler.GetTransactionById)
		transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
	}
	return r
}
