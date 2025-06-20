package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type LoanPaymentRepositoryImpl struct {
	DB *gorm.DB
	CommonRepository[models.LoanPayment]
}

func NewLoanPaymentRepository(db *gorm.DB) *LoanPaymentRepositoryImpl {
	return &LoanPaymentRepositoryImpl{
		DB:               db,
		CommonRepository: NewCommonRepository[models.LoanPayment](db),
	}
}
