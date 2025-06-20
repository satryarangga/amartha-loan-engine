package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"

	"github.com/google/uuid"
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
	if loanID == "" {
		return nil, errors.New("loan ID is required")
	}

	// Get pending schedules for the loan
	schedules, err := s.getPendingSchedules(ctx, loanID)
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

func (s *PaymentServiceImpl) HandlePaymentWebhook(ctx context.Context, paymentData map[string]interface{}) error {
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
	schedules, err := s.getPendingSchedules(ctx, loanID)
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
	_, err = s.loanPaymentRepo.Insert(ctx, nil, payment)
	if err != nil {
		return err
	}

	// Update schedule statuses to paid
	for _, schedule := range schedules {
		schedule.Status = "paid"
		if err := s.loanScheduleRepo.Update(ctx, nil, &schedule); err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentServiceImpl) GetPaymentHistory(ctx context.Context, loanID string) ([]models.LoanPayment, error) {
	if loanID == "" {
		return nil, errors.New("loan ID is required")
	}

	// This would need a custom method in the repository to filter by loan ID
	// For now, we'll get all payments and filter in memory
	allPayments, err := s.loanPaymentRepo.FindAll(ctx, models.FindAllParam{})
	if err != nil {
		return nil, err
	}

	// Filter payments by loan ID
	var filteredPayments []models.LoanPayment
	loanUUID := uuid.MustParse(loanID)
	for _, payment := range allPayments {
		if payment.LoanID == loanUUID {
			filteredPayments = append(filteredPayments, payment)
		}
	}

	return filteredPayments, nil
}

func (s *PaymentServiceImpl) getPendingSchedules(ctx context.Context, loanID string) ([]models.LoanSchedule, error) {
	// This would need a custom method in the repository to filter by loan ID and status
	// For now, we'll get all schedules and filter in memory
	allSchedules, err := s.loanScheduleRepo.FindAll(ctx, models.FindAllParam{})
	if err != nil {
		return nil, err
	}

	// Filter schedules by loan ID and pending status
	var pendingSchedules []models.LoanSchedule
	loanUUID := uuid.MustParse(loanID)
	for _, schedule := range allSchedules {
		if schedule.LoanID == loanUUID && schedule.Status == "pending" {
			pendingSchedules = append(pendingSchedules, schedule)
		}
	}

	return pendingSchedules, nil
}
