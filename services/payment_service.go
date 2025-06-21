package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type PaymentService interface {
	GeneratePaymentLink(ctx context.Context, paymentLinkRequest models.PaymentLinkRequest) (map[string]interface{}, error)
	HandlePaymentWebhook(ctx context.Context, paymentWebhookRequest models.PaymentWebhookRequest) error
}
