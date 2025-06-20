package services

import (
	"errors"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/repositories"
)

type BorrowerService struct {
	borrowerRepo *repositories.BorrowerRepository
}

func NewBorrowerService(borrowerRepo *repositories.BorrowerRepository) *BorrowerService {
	return &BorrowerService{
		borrowerRepo: borrowerRepo,
	}
}

func (s *BorrowerService) GetBorrowers() ([]models.Borrower, error) {
	return s.borrowerRepo.FindAll()
}

func (s *BorrowerService) GetBorrowerByID(id string) (*models.Borrower, error) {
	if id == "" {
		return nil, errors.New("borrower ID is required")
	}
	return s.borrowerRepo.FindByID(id)
}

func (s *BorrowerService) CreateBorrower(borrower *models.Borrower) error {
	if borrower.FirstName == "" || borrower.LastName == "" || borrower.PhoneNumber == "" {
		return errors.New("first name, last name, and phone number are required")
	}
	return s.borrowerRepo.Create(borrower)
}

func (s *BorrowerService) UpdateBorrower(borrower *models.Borrower) error {
	if borrower.ID.String() == "" {
		return errors.New("borrower ID is required")
	}
	return s.borrowerRepo.Update(borrower)
}

func (s *BorrowerService) DeleteBorrower(id string) error {
	if id == "" {
		return errors.New("borrower ID is required")
	}
	return s.borrowerRepo.Delete(id)
}
