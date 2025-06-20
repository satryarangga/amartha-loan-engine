package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type PaymentService interface {
	GeneratePaymentLink(ctx context.Context, loanID string) (map[string]interface{}, error)
	HandlePaymentWebhook(ctx context.Context, paymentData map[string]interface{}) error
	GetPaymentHistory(ctx context.Context, loanID string) ([]models.LoanPayment, error)
}
