package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type LoanRepositoryImpl struct {
	DB *gorm.DB
	CommonRepository[models.Loan]
}

func NewLoanRepository(db *gorm.DB) *LoanRepositoryImpl {
	return &LoanRepositoryImpl{
		DB:               db,
		CommonRepository: NewCommonRepository[models.Loan](db),
	}
}
