package main

import (
	"context"
	"log"
	"os"

	"github.com/satryarangga/amartha-loan-engine/config"
	"github.com/satryarangga/amartha-loan-engine/controllers"
	"github.com/satryarangga/amartha-loan-engine/repositories"
	"github.com/satryarangga/amartha-loan-engine/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/satryarangga/amartha-loan-engine/docs"
)

// @title           Amartha Loan Management API
// @version         1.0
// @description     A Golang backend application for managing loans, borrowers, and payments.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	// Initialize logger
	logger := config.NewLogger()
	ctx := context.Background()

	conf, err := config.NewConfig()
	if err != nil {
		logger.Errorf(ctx, "Unable to initialize config. Error: %v", err)
	}
	config.Config = conf

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	borrowerRepo := repositories.NewBorrowerRepository(db)
	loanRepo := repositories.NewLoanRepository(db)
	loanScheduleRepo := repositories.NewLoanScheduleRepository(db)
	loanPaymentRepo := repositories.NewLoanPaymentRepository(db)

	// Initialize services
	borrowerService := services.NewBorrowerService(borrowerRepo)
	loanService := services.NewLoanService(loanRepo, loanScheduleRepo)
	paymentService := services.NewPaymentService(loanPaymentRepo, loanScheduleRepo)

	// Initialize controllers
	borrowerController := controllers.NewBorrowerController(borrowerService)
	loanController := controllers.NewLoanController(loanService)
	paymentController := controllers.NewPaymentController(paymentService)

	// Setup router
	r := gin.Default()

	// Debug route to check if docs are accessible
	r.GET("/docs.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		// Borrower routes
		api.GET("/borrowers", borrowerController.GetBorrowers)
		api.GET("/borrowers/:id", borrowerController.GetBorrowerByID)
		api.POST("/borrowers", borrowerController.CreateBorrower)

		// Loan routes
		api.POST("/loans", loanController.CreateLoan)
		api.GET("/loans", loanController.GetLoans)
		api.GET("/loans/:id", loanController.GetLoanByID)
		api.PUT("/loans/:id", loanController.UpdateLoan)
		api.DELETE("/loans/:id", loanController.DeleteLoan)

		// Payment routes
		api.POST("/loans/:id/payment-link", paymentController.GeneratePaymentLink)
		api.POST("/webhook/payment", paymentController.HandlePaymentWebhook)
		api.GET("/loans/:id/payments", paymentController.GetPaymentHistory)
	}

	// Get port from environment
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	log.Printf("Direct docs.json available at http://localhost:%s/doc.json", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
