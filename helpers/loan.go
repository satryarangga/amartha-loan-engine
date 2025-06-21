package helpers

import (
	"time"

	"github.com/satryarangga/amartha-loan-engine/models"
)

// this use O(n) time complexity which is fine for this case since the number of loan schedules is limited
func CalculateTotalOutstanding(loan *models.Loan) float64 {
	var totalOutstanding float64
	for _, schedule := range loan.LoanSchedules {
		if schedule.Status == "pending" {
			totalOutstanding += schedule.TotalPayment
		}
	}
	return totalOutstanding
}

func GetTotalRepaymentAmount(loan *models.Loan) float64 {
	return loan.Amount + loan.InterestAmount
}

func IsBorrowerDelinquent(loanSchedules []models.LoanSchedule) bool {
	if len(loanSchedules) == 0 {
		return false
	}

	now := time.Now()
	overdueCount := 0
	const maxOverdueThreshold = 2

	for _, schedule := range loanSchedules {
		if overdueCount >= maxOverdueThreshold {
			return true
		}

		if schedule.Status == models.LoanScheduleStatusPending && schedule.DueDate.Before(now) {
			overdueCount++
		}
	}

	return overdueCount >= maxOverdueThreshold
}
