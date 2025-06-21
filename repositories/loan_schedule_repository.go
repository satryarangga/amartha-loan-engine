package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type LoanScheduleRepository interface {
	CommonRepository[models.LoanSchedule]

	FindDueRepaymentSchedules(ctx context.Context, loanID string) ([]models.LoanSchedule, error)

	UpdateStatusByIDs(ctx context.Context, tx *gorm.DB, ids []string, status models.LoanScheduleStatus) error
}
