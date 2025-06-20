package services

import (
	"context"
	"errors"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"
)

type BorrowerServiceImpl struct {
	borrowerRepo repositories.BorrowerRepository
}

func NewBorrowerService(borrowerRepo repositories.BorrowerRepository) *BorrowerServiceImpl {
	return &BorrowerServiceImpl{
		borrowerRepo: borrowerRepo,
	}
}

func (s *BorrowerServiceImpl) GetBorrowerByID(ctx context.Context, id string) (*models.Borrower, error) {
	if id == "" {
		return nil, errors.New("borrower ID is required")
	}
	return s.borrowerRepo.FindByID(ctx, id, []string{})
}

func (s *BorrowerServiceImpl) CreateBorrower(ctx context.Context, borrower *models.Borrower) error {
	_, err := s.borrowerRepo.Insert(ctx, nil, borrower)
	return err
}
