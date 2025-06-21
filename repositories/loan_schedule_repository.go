package repositories

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanScheduleRepository interface {
	CommonRepository[models.LoanSchedule]

	FindDueRepaymentSchedules(ctx context.Context, loanID string) ([]models.LoanSchedule, error)
}
