package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type LoanScheduleRepositoryImpl struct {
	DB *gorm.DB
	CommonRepository[models.LoanSchedule]
}

func NewLoanScheduleRepository(db *gorm.DB) *LoanScheduleRepositoryImpl {
	return &LoanScheduleRepositoryImpl{
		DB:               db,
		CommonRepository: NewCommonRepository[models.LoanSchedule](db),
	}
}
