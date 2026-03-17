package routes

import (
	"money-tracker/internal/database"
	"money-tracker/internal/handlers"
	"money-tracker/internal/middleware"
	"money-tracker/internal/repository"
	"money-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db := database.Connect()

	r := gin.Default()
	v1 := r.Group("/api/v1")

	authRepo := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	v1.POST("/register", authHandler.Register)
	v1.POST("/login", authHandler.Login)

	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", authHandler.GetProfile)
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
