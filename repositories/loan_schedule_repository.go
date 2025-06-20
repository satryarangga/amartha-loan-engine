package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanScheduleRepository interface {
	CommonRepository[models.LoanSchedule]
}
