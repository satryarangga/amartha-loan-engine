package helpers

import "github.com/satryarangga/amartha-loan-engine/models"

func CalculateTotalOutstanding(loan *models.Loan) float64 {
	var totalOutstanding float64
	for _, schedule := range loan.LoanSchedules {
		if schedule.Status == "pending" {
			totalOutstanding += schedule.TotalPayment
		}
	}
	return totalOutstanding
}
