package repositories

import (
	"context"

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

func (r *LoanRepositoryImpl) FindOneByBorrowerID(ctx context.Context, borrowerID string) (models.Loan, error) {
	var loan models.Loan
	err := r.DB.WithContext(ctx).
		Where("borrower_id = ? and status = ?", borrowerID, models.LoanStatusActive).
		Preload("LoanSchedules").
		First(&loan).Error
	return loan, err
}
