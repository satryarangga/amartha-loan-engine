package helpers

import (
	"testing"
	"time"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalOutstanding_WithPendingSchedules(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:     "loan-id",
		Amount: 1000000,
		LoanSchedules: []models.LoanSchedule{
			{
				ID:           "schedule-1",
				LoanID:       "loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPending,
			},
			{
				ID:           "schedule-2",
				LoanID:       "loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPaid,
			},
			{
				ID:           "schedule-3",
				LoanID:       "loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPending,
			},
		},
	}

	// Act
	result := CalculateTotalOutstanding(loan)

	// Assert
	expected := 220000.0 // 110000 + 110000 (only pending schedules)
	assert.Equal(t, expected, result)
}

func TestCalculateTotalOutstanding_WithNoPendingSchedules(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:     "loan-id",
		Amount: 1000000,
		LoanSchedules: []models.LoanSchedule{
			{
				ID:           "schedule-1",
				LoanID:       "loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPaid,
			},
			{
				ID:           "schedule-2",
				LoanID:       "loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPaid,
			},
		},
	}

	// Act
	result := CalculateTotalOutstanding(loan)

	// Assert
	expected := 0.0 // No pending schedules
	assert.Equal(t, expected, result)
}

func TestCalculateTotalOutstanding_WithEmptySchedules(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:            "loan-id",
		Amount:        1000000,
		LoanSchedules: []models.LoanSchedule{},
	}

	// Act
	result := CalculateTotalOutstanding(loan)

	// Assert
	expected := 0.0
	assert.Equal(t, expected, result)
}

func TestCalculateTotalOutstanding_WithNilSchedules(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:     "loan-id",
		Amount: 1000000,
	}

	// Act
	result := CalculateTotalOutstanding(loan)

	// Assert
	expected := 0.0
	assert.Equal(t, expected, result)
}

func TestGetTotalRepaymentAmount(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:             "loan-id",
		Amount:         1000000,
		InterestAmount: 100000,
	}

	// Act
	result := GetTotalRepaymentAmount(loan)

	// Assert
	expected := 1100000.0 // 1000000 + 100000
	assert.Equal(t, expected, result)
}

func TestGetTotalRepaymentAmount_WithZeroInterest(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:             "loan-id",
		Amount:         1000000,
		InterestAmount: 0,
	}

	// Act
	result := GetTotalRepaymentAmount(loan)

	// Assert
	expected := 1000000.0 // 1000000 + 0
	assert.Equal(t, expected, result)
}

func TestGetTotalRepaymentAmount_WithZeroAmount(t *testing.T) {
	// Arrange
	loan := &models.Loan{
		ID:             "loan-id",
		Amount:         0,
		InterestAmount: 100000,
	}

	// Act
	result := GetTotalRepaymentAmount(loan)

	// Assert
	expected := 100000.0 // 0 + 100000
	assert.Equal(t, expected, result)
}

func TestIsBorrowerDelinquent_WithDelinquentBorrower(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -5), // 5 days overdue
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -3), // 3 days overdue
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -1), // 1 day overdue
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.True(t, result, "Borrower should be delinquent with 3 overdue payments")
}

func TestIsBorrowerDelinquent_WithNonDelinquentBorrower(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -1), // 1 day overdue
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -3), // 3 days overdue but paid
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, 5), // 5 days in future
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.False(t, result, "Borrower should not be delinquent with only 1 overdue payment")
}

func TestIsBorrowerDelinquent_WithExactlyTwoOverduePayments(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -5), // 5 days overdue
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -3), // 3 days overdue
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, 5), // 5 days in future
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.True(t, result, "Borrower should be delinquent with exactly 2 overdue payments")
}

func TestIsBorrowerDelinquent_WithNoOverduePayments(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, 5), // 5 days in future
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -3), // 3 days overdue but paid
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, 10), // 10 days in future
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.False(t, result, "Borrower should not be delinquent with no overdue payments")
}

func TestIsBorrowerDelinquent_WithAllPaidSchedules(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -5), // 5 days overdue but paid
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -3), // 3 days overdue but paid
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -1), // 1 day overdue but paid
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.False(t, result, "Borrower should not be delinquent with all paid schedules")
}

func TestIsBorrowerDelinquent_WithEmptySchedules(t *testing.T) {
	// Arrange
	loanSchedules := []models.LoanSchedule{}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.False(t, result, "Borrower should not be delinquent with empty schedules")
}

func TestIsBorrowerDelinquent_WithMixedStatuses(t *testing.T) {
	// Arrange
	now := time.Now()
	loanSchedules := []models.LoanSchedule{
		{
			ID:           "schedule-1",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -10), // 10 days overdue
		},
		{
			ID:           "schedule-2",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPaid,
			DueDate:      now.AddDate(0, 0, -5), // 5 days overdue but paid
		},
		{
			ID:           "schedule-3",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, -2), // 2 days overdue
		},
		{
			ID:           "schedule-4",
			LoanID:       "loan-id",
			TotalPayment: 110000,
			Status:       models.LoanScheduleStatusPending,
			DueDate:      now.AddDate(0, 0, 5), // 5 days in future
		},
	}

	// Act
	result := IsBorrowerDelinquent(loanSchedules)

	// Assert
	assert.True(t, result, "Borrower should be delinquent with 2 overdue payments (2 pending + 1 paid)")
}
