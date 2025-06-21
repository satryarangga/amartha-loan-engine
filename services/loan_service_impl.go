package services

import (
	"context"
	"errors"
	"time"

	"github.com/satryarangga/amartha-loan-engine/helpers"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"
	"gorm.io/gorm"
)

type LoanServiceImpl struct {
	loanRepo         repositories.LoanRepository
	loanScheduleRepo repositories.LoanScheduleRepository
	borrowerRepo     repositories.BorrowerRepository
}

func NewLoanService(loanRepo repositories.LoanRepository, loanScheduleRepo repositories.LoanScheduleRepository, borrowerRepo repositories.BorrowerRepository) *LoanServiceImpl {
	return &LoanServiceImpl{
		loanRepo:         loanRepo,
		loanScheduleRepo: loanScheduleRepo,
		borrowerRepo:     borrowerRepo,
	}
}

func (s *LoanServiceImpl) GetLoanByID(ctx context.Context, id string) (*models.LoanResponse, error) {
	if id == "" {
		return nil, errors.New("loan ID is required")
	}
	loan, err := s.loanRepo.FindByID(ctx, id, []string{"LoanSchedules"})
	if err != nil {
		return nil, err
	}

	loanResponse := models.LoanResponse{
		ID:                   loan.ID,
		Amount:               loan.Amount,
		RepaymentCadenceDays: loan.RepaymentCadenceDays,
		RepaymentRepetition:  loan.RepaymentRepetition,
		InterestPercentage:   loan.InterestPercentage,
		InterestAmount:       loan.InterestAmount,
		Status:               string(loan.Status),
		TotalOutstanding:     helpers.CalculateTotalOutstanding(loan),
	}

	return &loanResponse, nil
}

func (s *LoanServiceImpl) CreateLoan(ctx context.Context, req *models.LoanRequest) error {
	borrower, err := s.borrowerRepo.FindByID(ctx, req.BorrowerID, []string{})
	if err != nil {
		return err
	}

	if borrower == nil {
		return errors.New("borrower not found")
	}

	interestAmount := req.Amount * req.InterestPercentage / 100

	loan := models.Loan{
		BorrowerID:           borrower.ID,
		Amount:               req.Amount,
		RepaymentCadenceDays: req.RepaymentCadenceDays,
		RepaymentRepetition:  req.RepaymentRepetition,
		InterestPercentage:   req.InterestPercentage,
		InterestAmount:       interestAmount,
		Status:               "active",
	}

	err = s.loanRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		loanID, err := s.loanRepo.Insert(ctx, tx, &loan)
		if err != nil {
			return err
		}

		loanSchedules := make([]models.LoanSchedule, req.RepaymentRepetition)
		basicRepaymentAmount := loan.Amount / float64(req.RepaymentRepetition)
		interestRepaymentAmount := loan.InterestAmount / float64(req.RepaymentRepetition)
		totalRepaymentAmount := basicRepaymentAmount + interestRepaymentAmount
		for i := 1; i <= req.RepaymentRepetition; i++ {
			loanSchedules[i-1] = models.LoanSchedule{
				LoanID:         loanID,
				DueDate:        time.Now().AddDate(0, 0, req.RepaymentCadenceDays*i),
				BasicAmount:    basicRepaymentAmount,
				InterestAmount: interestRepaymentAmount,
				TotalPayment:   totalRepaymentAmount,
				Status:         "pending",
			}
		}

		for _, loanSchedule := range loanSchedules {
			_, err := s.loanScheduleRepo.Insert(ctx, tx, &loanSchedule)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
