package main

import (
	"log"
	"os"

	"amartha/config"
	"amartha/controllers"
	"amartha/repositories"
	"amartha/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

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
