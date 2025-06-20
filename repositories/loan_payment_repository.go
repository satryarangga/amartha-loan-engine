package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanPaymentRepository interface {
	CommonRepository[models.LoanPayment]
}
