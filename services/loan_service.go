package services

import (
	"errors"
	"time"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"

	"github.com/google/uuid"
)

type LoanService struct {
	loanRepo         *repositories.LoanRepository
	loanScheduleRepo *repositories.LoanScheduleRepository
}

func NewLoanService(loanRepo *repositories.LoanRepository, loanScheduleRepo *repositories.LoanScheduleRepository) *LoanService {
	return &LoanService{
		loanRepo:         loanRepo,
		loanScheduleRepo: loanScheduleRepo,
	}
}

func (s *LoanService) GetLoans() ([]models.Loan, error) {
	return s.loanRepo.FindAll()
}

func (s *LoanService) GetLoanByID(id string) (*models.Loan, error) {
	if id == "" {
		return nil, errors.New("loan ID is required")
	}
	return s.loanRepo.FindByID(id)
}

func (s *LoanService) CreateLoan(loan *models.Loan) error {
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
	if err := s.loanRepo.Create(loan); err != nil {
		return err
	}

	// Generate loan schedules
	return s.generateLoanSchedules(loan)
}

func (s *LoanService) generateLoanSchedules(loan *models.Loan) error {
	// Calculate basic amount per schedule (assuming equal distribution)
	basicAmountPerSchedule := loan.Amount / float64(loan.RepaymentCadenceDays)
	interestAmountPerSchedule := loan.InterestAmount / float64(loan.RepaymentCadenceDays)

	var schedules []models.LoanSchedule
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

		schedules = append(schedules, schedule)
	}

	return s.loanScheduleRepo.CreateBatch(schedules)
}

func (s *LoanService) UpdateLoan(loan *models.Loan) error {
	if loan.ID.String() == "" {
		return errors.New("loan ID is required")
	}
	return s.loanRepo.Update(loan)
}

func (s *LoanService) DeleteLoan(id string) error {
	if id == "" {
		return errors.New("loan ID is required")
	}
	return s.loanRepo.Delete(id)
}
