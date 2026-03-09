package routes

import (
	"money-tracker/internal/database"
	"money-tracker/internal/handlers"
	"money-tracker/internal/repository"
	"money-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db := database.Connect()

	r := gin.Default()

	authRepo := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

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
