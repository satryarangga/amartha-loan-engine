package helpers

import "github.com/satryarangga/amartha-loan-engine/models"

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
