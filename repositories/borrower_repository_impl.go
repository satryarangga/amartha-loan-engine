package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type BorrowerRepositoryImpl struct {
	DB *gorm.DB
	CommonRepository[models.Borrower]
}

func NewBorrowerRepository(db *gorm.DB) *BorrowerRepositoryImpl {
	return &BorrowerRepositoryImpl{
		DB:               db,
		CommonRepository: NewCommonRepository[models.Borrower](db),
	}
}

func (r *BorrowerRepositoryImpl) FindOneByPhoneNumber(ctx context.Context, phoneNumber string) (models.Borrower, error) {
	var borrower models.Borrower
	err := r.DB.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&borrower).Error
	return borrower, err
}
