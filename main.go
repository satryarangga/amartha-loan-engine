package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/satryarangga/amartha-loan-engine/config"
	"github.com/satryarangga/amartha-loan-engine/controllers"
	"github.com/satryarangga/amartha-loan-engine/repositories"
	"github.com/satryarangga/amartha-loan-engine/services"

	"github.com/gin-gonic/gin"
)

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

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Borrower routes
		api.GET("/borrowers", borrowerController.GetBorrowers)
		api.GET("/borrowers/:id", borrowerController.GetBorrowerByID)

		// Loan routes
		api.POST("/loans", loanController.CreateLoan)
		api.GET("/loans", loanController.GetLoans)
		api.GET("/loans/:id", loanController.GetLoanByID)

		// Payment routes
		api.POST("/loans/:id/payment-link", paymentController.GeneratePaymentLink)
		api.POST("/webhook/payment", paymentController.HandlePaymentWebhook)
	}

	// Get port from environment
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
