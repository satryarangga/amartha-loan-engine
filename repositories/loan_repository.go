package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanRepository interface {
	CommonRepository[models.Loan]

	FindOneByBorrowerID(ctx context.Context, borrowerID string) (models.Loan, error)
}
