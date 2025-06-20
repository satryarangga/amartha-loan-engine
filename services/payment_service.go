package services

import (
	"context"
)

type PaymentService interface {
	GeneratePaymentLink(ctx context.Context, loanID string) (map[string]interface{}, error)
	HandlePaymentWebhook(ctx context.Context, paymentData map[string]interface{}) error
}
