package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
)

type BorrowerRepository interface {
	CommonRepository[models.Borrower]
}
