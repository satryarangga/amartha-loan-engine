package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/satryarangga/amartha-loan-engine/helpers"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	loanRepo         repositories.LoanRepository
	loanPaymentRepo  repositories.LoanPaymentRepository
	loanScheduleRepo repositories.LoanScheduleRepository
	borrowerRepo     repositories.BorrowerRepository
}

func NewPaymentService(
	loanRepo repositories.LoanRepository,
	loanPaymentRepo repositories.LoanPaymentRepository,
	loanScheduleRepo repositories.LoanScheduleRepository,
	borrowerRepo repositories.BorrowerRepository,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		loanRepo:         loanRepo,
		loanPaymentRepo:  loanPaymentRepo,
		loanScheduleRepo: loanScheduleRepo,
		borrowerRepo:     borrowerRepo,
	}
}

func (s *PaymentServiceImpl) GeneratePaymentLink(ctx context.Context, request models.PaymentLinkRequest) (*models.PaymentLinkResponse, error) {
	// 1. Find by borrower ID
	borrower, err := s.borrowerRepo.FindByID(ctx, request.BorrowerID, []string{})
	if err != nil {
		return nil, err
	}

	// 2. Get loan by borrower ID along with loan schedules
	loan, err := s.loanRepo.FindOneByBorrowerID(ctx, borrower.ID)
	if err != nil {
		return nil, err
	}

	//3 Get all loan schedules that are pending and due date is less than 3 days from now
	loanSchedules, err := s.loanScheduleRepo.FindDueRepaymentSchedules(ctx, loan.ID)
	if err != nil {
		return nil, err
	}

	if len(loanSchedules) == 0 {
		return nil, errors.New("no loan schedules found")
	}

	// 3. Show total outstanding that needs to be paid
	var totalRepaymentAmount float64
	loanScheduleIDs := []string{}
	for _, loanSchedule := range loanSchedules {
		totalRepaymentAmount += loanSchedule.TotalPayment
		loanScheduleIDs = append(loanScheduleIDs, loanSchedule.ID)
	}

	// 4 Create new loan_payments table with status pending and payment_method
	loanPayment := models.LoanPayment{
		LoanID:          loan.ID,
		LoanScheduleIDs: loanScheduleIDs,
		TotalPayment:    totalRepaymentAmount,
		PaymentMethod:   request.PaymentMethod,
	}
	loanPaymentID, err := s.loanPaymentRepo.Insert(ctx, nil, &loanPayment)
	if err != nil {
		return nil, err
	}

	// 5. Generate payment link (Assume hitting Payment Gateway API and the Payment Gateway API returns payment link)
	paymentLink := fmt.Sprintf("https://example.com/payment-link?external_id=%s", loanPaymentID)

	return &models.PaymentLinkResponse{
		ID:                   loanPaymentID,
		TotalRepaymentAmount: totalRepaymentAmount,
		PaymentLink:          paymentLink,
	}, nil
}

func (s *PaymentServiceImpl) HandlePaymentWebhook(ctx context.Context, request models.PaymentWebhookRequest) error {
	if request.PaymentStatus != "paid" {
		return errors.New("payment status from PG is not paid")
	}

	err := s.loanRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// 1.Find loan payment with ID
		loanPayment, err := s.loanPaymentRepo.FindByID(ctx, request.ExternalID, []string{})
		if err != nil {
			return err
		}

		// 2. Update Status on Loan Payment
		loanPayment.Status = models.LoanPaymentStatusPaid
		err = s.loanPaymentRepo.Update(ctx, nil, loanPayment)
		if err != nil {
			return err
		}

		// 3. Find Loan Detail
		loan, err := s.loanRepo.FindByID(ctx, loanPayment.LoanID, []string{"LoanSchedules"})
		if err != nil {
			return err
		}

		// 4. Calculate total paid repayment amount
		var totalPaidRepaymentAmount float64
		for _, loanSchedule := range loan.LoanSchedules {
			if loanSchedule.Status == models.LoanScheduleStatusPaid {
				totalPaidRepaymentAmount += loanSchedule.TotalPayment
			}
		}

		// 5. Update Status of Loan Schedules
		err = s.loanScheduleRepo.UpdateStatusByIDs(ctx, tx, loanPayment.LoanScheduleIDs, models.LoanScheduleStatusPaid)
		if err != nil {
			return err
		}

		// 6. Update Status of Loan if no more outstanding repayment amount
		if totalPaidRepaymentAmount+loanPayment.TotalPayment >= helpers.GetTotalRepaymentAmount(loan) {
			loan.Status = models.LoanStatusPaid
			err = s.loanRepo.Update(ctx, tx, loan)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
