package repositories

import (
	"context"
	"time"

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

func (r *LoanScheduleRepositoryImpl) FindDueRepaymentSchedules(ctx context.Context, loanID string) ([]models.LoanSchedule, error) {
	var loanSchedules []models.LoanSchedule
	err := r.DB.WithContext(ctx).Where("loan_id = ? and status = ? and due_date <= ?", loanID, models.LoanScheduleStatusPending, time.Now().AddDate(0, 0, 3)).Find(&loanSchedules).Error
	return loanSchedules, err
}

func (r *LoanScheduleRepositoryImpl) UpdateStatusByIDs(ctx context.Context, tx *gorm.DB, ids []string, status models.LoanScheduleStatus) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&models.LoanSchedule{}).Where("id IN (?)", ids).Update("status", status).Error
}
