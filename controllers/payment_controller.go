package controllers

import (
	"amartha/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// GeneratePaymentLink handles POST /api/v1/loans/:id/payment-link
func (c *PaymentController) GeneratePaymentLink(ctx *gin.Context) {
	loanID := ctx.Param("id")
	if loanID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	paymentData, err := c.paymentService.GeneratePaymentLink(loanID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to generate payment link",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    paymentData,
		"message": "Payment link generated successfully",
	})
}

// HandlePaymentWebhook handles POST /api/v1/webhook/payment
func (c *PaymentController) HandlePaymentWebhook(ctx *gin.Context) {
	var paymentData map[string]interface{}
	if err := ctx.ShouldBindJSON(&paymentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid webhook payload",
			"details": err.Error(),
		})
		return
	}

	if err := c.paymentService.HandlePaymentWebhook(paymentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to process payment webhook",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Payment processed successfully",
	})
}

// GetPaymentHistory handles GET /api/v1/loans/:id/payments
func (c *PaymentController) GetPaymentHistory(ctx *gin.Context) {
	loanID := ctx.Param("id")
	if loanID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	payments, err := c.paymentService.GetPaymentHistory(loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch payment history",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  payments,
		"count": len(payments),
	})
}
