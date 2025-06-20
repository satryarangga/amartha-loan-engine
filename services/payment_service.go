package services

import (
	"amartha/models"
	"amartha/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	loanPaymentRepo  *repositories.LoanPaymentRepository
	loanScheduleRepo *repositories.LoanScheduleRepository
}

func NewPaymentService(loanPaymentRepo *repositories.LoanPaymentRepository, loanScheduleRepo *repositories.LoanScheduleRepository) *PaymentService {
	return &PaymentService{
		loanPaymentRepo:  loanPaymentRepo,
		loanScheduleRepo: loanScheduleRepo,
	}
}

func (s *PaymentService) GeneratePaymentLink(loanID string) (map[string]interface{}, error) {
	if loanID == "" {
		return nil, errors.New("loan ID is required")
	}

	// Get pending schedules for the loan
	schedules, err := s.loanScheduleRepo.FindPendingByLoanID(loanID)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, errors.New("no pending payments found for this loan")
	}

	// Calculate total amount due
	var totalAmount float64
	var scheduleIDs []uuid.UUID
	for _, schedule := range schedules {
		totalAmount += schedule.TotalPayment
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}

	// Generate payment link (in a real scenario, this would integrate with a payment gateway)
	paymentLink := fmt.Sprintf("https://payment-gateway.com/pay?loan_id=%s&amount=%.2f", loanID, totalAmount)

	return map[string]interface{}{
		"payment_link": paymentLink,
		"total_amount": totalAmount,
		"schedule_ids": scheduleIDs,
		"expires_at":   time.Now().Add(24 * time.Hour), // Link expires in 24 hours
	}, nil
}

func (s *PaymentService) HandlePaymentWebhook(paymentData map[string]interface{}) error {
	loanID, ok := paymentData["loan_id"].(string)
	if !ok {
		return errors.New("loan_id is required in payment data")
	}

	amount, ok := paymentData["amount"].(float64)
	if !ok {
		return errors.New("amount is required in payment data")
	}

	paymentMethod, ok := paymentData["payment_method"].(string)
	if !ok {
		paymentMethod = "unknown"
	}

	// Get pending schedules for the loan
	schedules, err := s.loanScheduleRepo.FindPendingByLoanID(loanID)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return errors.New("no pending payments found for this loan")
	}

	// Create payment record
	payment := &models.LoanPayment{
		ID:            uuid.New(),
		LoanID:        uuid.MustParse(loanID),
		TotalPayment:  amount,
		PaymentMethod: paymentMethod,
	}

	// Collect schedule IDs
	var scheduleIDs []uuid.UUID
	for _, schedule := range schedules {
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}
	payment.LoanScheduleIDs = scheduleIDs

	// Save payment record
	if err := s.loanPaymentRepo.Create(payment); err != nil {
		return err
	}

	// Update schedule statuses to paid
	for _, schedule := range schedules {
		if err := s.loanScheduleRepo.UpdateStatus(schedule.ID.String(), "paid"); err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentService) GetPaymentHistory(loanID string) ([]models.LoanPayment, error) {
	if loanID == "" {
		return nil, errors.New("loan ID is required")
	}
	return s.loanPaymentRepo.FindByLoanID(loanID)
}
