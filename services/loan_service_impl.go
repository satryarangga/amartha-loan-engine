package services

import (
	"context"
	"errors"
	"time"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"

	"github.com/google/uuid"
)

type LoanServiceImpl struct {
	loanRepo         repositories.LoanRepository
	loanScheduleRepo repositories.LoanScheduleRepository
}

func NewLoanService(loanRepo repositories.LoanRepository, loanScheduleRepo repositories.LoanScheduleRepository) *LoanServiceImpl {
	return &LoanServiceImpl{
		loanRepo:         loanRepo,
		loanScheduleRepo: loanScheduleRepo,
	}
}

func (s *LoanServiceImpl) GetLoans(ctx context.Context) ([]models.Loan, error) {
	return s.loanRepo.FindAll(ctx, models.FindAllParam{})
}

func (s *LoanServiceImpl) GetLoanByID(ctx context.Context, id string) (*models.Loan, error) {
	if id == "" {
		return nil, errors.New("loan ID is required")
	}
	return s.loanRepo.FindByID(ctx, id)
}

func (s *LoanServiceImpl) CreateLoan(ctx context.Context, loan *models.Loan) error {
	// Validate required fields
	if loan.BorrowerID.String() == "" {
		return errors.New("borrower ID is required")
	}
	if loan.Amount <= 0 {
		return errors.New("loan amount must be greater than 0")
	}
	if loan.RepaymentCadenceDays <= 0 {
		return errors.New("repayment cadence days must be greater than 0")
	}
	if loan.InterestPercentage < 0 {
		return errors.New("interest percentage cannot be negative")
	}

	// Calculate interest amount
	loan.InterestAmount = (loan.Amount * loan.InterestPercentage) / 100

	// Set default status
	if loan.Status == "" {
		loan.Status = "active"
	}

	// Create the loan
	_, err := s.loanRepo.Insert(ctx, nil, loan)
	if err != nil {
		return err
	}

	// Generate loan schedules
	return s.generateLoanSchedules(ctx, loan)
}

func (s *LoanServiceImpl) generateLoanSchedules(ctx context.Context, loan *models.Loan) error {
	// Calculate basic amount per schedule (assuming equal distribution)
	basicAmountPerSchedule := loan.Amount / float64(loan.RepaymentCadenceDays)
	interestAmountPerSchedule := loan.InterestAmount / float64(loan.RepaymentCadenceDays)

	startDate := time.Now()

	for i := 0; i < loan.RepaymentCadenceDays; i++ {
		dueDate := startDate.AddDate(0, 0, i+1) // Start from tomorrow

		schedule := models.LoanSchedule{
			ID:             uuid.New(),
			LoanID:         loan.ID,
			DueDate:        dueDate,
			BasicAmount:    basicAmountPerSchedule,
			InterestAmount: interestAmountPerSchedule,
			TotalPayment:   basicAmountPerSchedule + interestAmountPerSchedule,
			Status:         "pending",
		}

		// Insert each schedule individually since CreateBatch is not in CommonRepository
		_, err := s.loanScheduleRepo.Insert(ctx, nil, &schedule)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *LoanServiceImpl) UpdateLoan(ctx context.Context, loan *models.Loan) error {
	if loan.ID.String() == "" {
		return errors.New("loan ID is required")
	}
	return s.loanRepo.Update(ctx, nil, loan)
}
