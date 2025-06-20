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

func (s *BorrowerServiceImpl) GetBorrowers(ctx context.Context) ([]models.Borrower, error) {
	return s.borrowerRepo.FindAll(ctx, models.FindAllParam{})
}

func (s *BorrowerServiceImpl) GetBorrowerByID(ctx context.Context, id string) (*models.Borrower, error) {
	if id == "" {
		return nil, errors.New("borrower ID is required")
	}
	return s.borrowerRepo.FindByID(ctx, id)
}

func (s *BorrowerServiceImpl) CreateBorrower(ctx context.Context, borrower *models.Borrower) error {
	if borrower.FirstName == "" || borrower.LastName == "" || borrower.PhoneNumber == "" {
		return errors.New("first name, last name, and phone number are required")
	}
	_, err := s.borrowerRepo.Insert(ctx, nil, borrower)
	return err
}

func (s *BorrowerServiceImpl) UpdateBorrower(ctx context.Context, borrower *models.Borrower) error {
	if borrower.ID.String() == "" {
		return errors.New("borrower ID is required")
	}
	return s.borrowerRepo.Update(ctx, nil, borrower)
}
