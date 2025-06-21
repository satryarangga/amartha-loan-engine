package helpers

import (
	"testing"

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
