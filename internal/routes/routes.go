package routes

import (
	"money-tracker/internal/database"
	"money-tracker/internal/handlers"
	"money-tracker/internal/middleware"
	"money-tracker/internal/repository"
	"money-tracker/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db := database.Connect()

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")

	authRepo := repository.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	accountRepo := repository.NewAccountRepository(db)
	accountService := services.NewAccountService(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	budgetRepo := repository.NewBudgetAllocationRepository(db)
	budgetService := services.NewBudgetAllocationService(budgetRepo)
	budgetHandler := handlers.NewBudgetAllocationHandler(budgetService)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, db)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	transferRepo := repository.NewTransferRepository(db)
	transferService := services.NewTransferService(transferRepo, db)
	transferHandler := handlers.NewTransferHandler(transferService)

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/google", authHandler.GoogleCallback)
	}

	v1.Use(middleware.AuthMiddleware())
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET("", accountHandler.GetAccounts)
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.PUT("/:id", accountHandler.UpdateAccount)
			accounts.DELETE("/:id", accountHandler.DeleteAccount)
		}
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}
		budgets := v1.Group("/budgets")
		{
			budgets.GET("", budgetHandler.GetBudgetAllocations)
			budgets.POST("", budgetHandler.CreateBudgetAllocation)
			budgets.GET("/:id", budgetHandler.GetBudgetAllocation)
			budgets.PUT("/:id", budgetHandler.UpdateBudgetAllocation)
			budgets.DELETE("/:id", budgetHandler.DeleteBudgetAllocation)
		}
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}
		transfers := v1.Group("/transfers")
		{
			transfers.GET("", transferHandler.GetTransfers)
			transfers.POST("", transferHandler.CreateTransfer)
			transfers.GET("/:id", transferHandler.GetTransfer)
			transfers.PUT("/:id", transferHandler.UpdateTransfer)
			transfers.DELETE("/:id", transferHandler.DeleteTransfer)
		}
	}
	return r
}
