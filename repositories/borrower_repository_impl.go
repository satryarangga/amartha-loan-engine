package repositories

import (
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
