package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type BorrowerService interface {
	GetBorrowerByID(ctx context.Context, id string) (*models.Borrower, error)
	CreateBorrower(ctx context.Context, borrower *models.BorrowerRequest) (*models.Borrower, error)
}
