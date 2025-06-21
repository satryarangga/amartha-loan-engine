package services

import (
	"context"
	"errors"

	"github.com/satryarangga/amartha-loan-engine/helpers"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"
)

type BorrowerServiceImpl struct {
	borrowerRepo repositories.BorrowerRepository
	loanRepo     repositories.LoanRepository
}

func NewBorrowerService(
	borrowerRepo repositories.BorrowerRepository,
	loanRepo repositories.LoanRepository,
) *BorrowerServiceImpl {
	return &BorrowerServiceImpl{
		borrowerRepo: borrowerRepo,
		loanRepo:     loanRepo,
	}
}

func (s *BorrowerServiceImpl) GetBorrowerByID(ctx context.Context, id string) (*models.BorrowerResponse, error) {
	if id == "" {
		return nil, errors.New("borrower ID is required")
	}
	borrower, err := s.borrowerRepo.FindByID(ctx, id, []string{})
	if err != nil {
		return nil, err
	}

	loan, err := s.loanRepo.FindOneByBorrowerID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.BorrowerResponse{
		ID:           borrower.ID,
		FirstName:    borrower.FirstName,
		LastName:     borrower.LastName,
		PhoneNumber:  borrower.PhoneNumber,
		IsDelinquent: helpers.IsBorrowerDelinquent(loan.LoanSchedules),
	}, nil
}

func (s *BorrowerServiceImpl) CreateBorrower(ctx context.Context, borrower *models.Borrower) error {
	_, err := s.borrowerRepo.Insert(ctx, nil, borrower)
	return err
}
