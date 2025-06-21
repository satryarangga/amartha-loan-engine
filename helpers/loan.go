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
	maxOverdueRepayment := 2
	for _, schedule := range loanSchedules {
		if schedule.Status == models.LoanScheduleStatusPending && schedule.DueDate.Before(time.Now()) {
			maxOverdueRepayment--
		}
	}
	return maxOverdueRepayment <= 0
}
