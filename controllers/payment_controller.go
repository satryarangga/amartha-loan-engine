package controllers

import (
	"net/http"

	"github.com/satryarangga/amartha-loan-engine/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *services.PaymentServiceImpl
}

func NewPaymentController(paymentService *services.PaymentServiceImpl) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// GeneratePaymentLink godoc
// @Summary Generate payment link
// @Description Generate a payment link for a specific loan
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /loans/{id}/payment-link [post]
func (c *PaymentController) GeneratePaymentLink(ctx *gin.Context) {
	loanID := ctx.Param("id")
	if loanID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	paymentData, err := c.paymentService.GeneratePaymentLink(ctx, loanID)
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

// HandlePaymentWebhook godoc
// @Summary Handle payment webhook
// @Description Process payment webhook from payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentData body map[string]interface{} true "Payment webhook data"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /webhook/payment [post]
func (c *PaymentController) HandlePaymentWebhook(ctx *gin.Context) {
	var paymentData map[string]interface{}
	if err := ctx.ShouldBindJSON(&paymentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid webhook payload",
			"details": err.Error(),
		})
		return
	}

	if err := c.paymentService.HandlePaymentWebhook(ctx, paymentData); err != nil {
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
