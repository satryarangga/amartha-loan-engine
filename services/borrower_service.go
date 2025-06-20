package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type BorrowerService interface {
	GetBorrowers(ctx context.Context, param models.FindAllParam) ([]models.Borrower, error)
	GetBorrowerByID(ctx context.Context, id string) (*models.Borrower, error)
	CreateBorrower(ctx context.Context, borrower *models.BorrowerRequest) (*models.Borrower, error)
	UpdateBorrower(ctx context.Context, borrower *models.BorrowerRequest) (*models.Borrower, error)
	DeleteBorrower(ctx context.Context, id string) error
}
