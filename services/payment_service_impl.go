package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/repositories"
)

type PaymentServiceImpl struct {
	loanPaymentRepo  repositories.LoanPaymentRepository
	loanScheduleRepo repositories.LoanScheduleRepository
}

func NewPaymentService(loanPaymentRepo repositories.LoanPaymentRepository, loanScheduleRepo repositories.LoanScheduleRepository) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		loanPaymentRepo:  loanPaymentRepo,
		loanScheduleRepo: loanScheduleRepo,
	}
}

func (s *PaymentServiceImpl) GeneratePaymentLink(ctx context.Context, loanID string) (map[string]interface{}, error) {
	return nil, nil
}

func (s *PaymentServiceImpl) HandlePaymentWebhook(ctx context.Context, paymentData map[string]interface{}) error {
	return nil
}
