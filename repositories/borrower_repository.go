package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type BorrowerRepository interface {
	CommonRepository[models.Borrower]

	FindOneByPhoneNumber(ctx context.Context, phoneNumber string) (models.Borrower, error)
}
