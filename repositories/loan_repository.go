package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanRepository interface {
	CommonRepository[models.Loan]
}
