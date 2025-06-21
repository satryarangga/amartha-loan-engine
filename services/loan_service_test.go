package services

import (
	"context"
	"errors"
	"testing"

	"github.com/satryarangga/amartha-loan-engine/mock"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestNewLoanService(t *testing.T) {
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)

	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockLoanRepo, service.loanRepo)
	assert.Equal(t, mockLoanScheduleRepo, service.loanScheduleRepo)
	assert.Equal(t, mockBorrowerRepo, service.borrowerRepo)
}

func TestLoanServiceImpl_GetLoanByID_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	loanID := "test-loan-id"
	expectedLoan := &models.Loan{
		ID:                   loanID,
		BorrowerID:           "borrower-id",
		Amount:               1000000,
		RepaymentCadenceDays: 7,
		RepaymentRepetition:  12,
		InterestPercentage:   10,
		InterestAmount:       100000,
		Status:               models.LoanStatusActive,
		LoanSchedules: []models.LoanSchedule{
			{
				ID:           "schedule-1",
				LoanID:       loanID,
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPending,
			},
			{
				ID:           "schedule-2",
				LoanID:       loanID,
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPaid,
			},
		},
	}

	mockLoanRepo.On("FindByID", ctx, loanID, []string{"LoanSchedules"}).Return(expectedLoan, nil)

	// Act
	result, err := service.GetLoanByID(ctx, loanID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, loanID, result.ID)
	assert.Equal(t, expectedLoan.Amount, result.Amount)
	assert.Equal(t, expectedLoan.RepaymentCadenceDays, result.RepaymentCadenceDays)
	assert.Equal(t, expectedLoan.RepaymentRepetition, result.RepaymentRepetition)
	assert.Equal(t, expectedLoan.InterestPercentage, result.InterestPercentage)
	assert.Equal(t, expectedLoan.InterestAmount, result.InterestAmount)
	assert.Equal(t, string(expectedLoan.Status), result.Status)
	assert.Equal(t, 110000.0, result.TotalOutstanding) // Only pending schedule
	mockLoanRepo.AssertExpectations(t)
}

func TestLoanServiceImpl_GetLoanByID_EmptyID(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()

	// Act
	result, err := service.GetLoanByID(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan ID is required", err.Error())
}

func TestLoanServiceImpl_GetLoanByID_RepositoryError(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	loanID := "test-loan-id"
	expectedError := errors.New("database error")

	mockLoanRepo.On("FindByID", ctx, loanID, []string{"LoanSchedules"}).Return(nil, expectedError)

	// Act
	result, err := service.GetLoanByID(ctx, loanID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockLoanRepo.AssertExpectations(t)
}

func TestLoanServiceImpl_CreateLoan_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := &models.LoanRequest{
		BorrowerID:           "borrower-id",
		Amount:               1000000,
		RepaymentCadenceDays: 7,
		RepaymentRepetition:  12,
		InterestPercentage:   10,
	}

	borrower := &models.Borrower{
		ID:          "borrower-id",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}

	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(borrower, nil)
	mockLoanRepo.On("WithTransaction", ctx, testifymock.AnythingOfType("repositories.TransactionFunc")).Return(nil)

	// Act
	err := service.CreateLoan(ctx, request)

	// Assert
	assert.NoError(t, err)
	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

func TestLoanServiceImpl_CreateLoan_BorrowerNotFound(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := &models.LoanRequest{
		BorrowerID:           "borrower-id",
		Amount:               1000000,
		RepaymentCadenceDays: 7,
		RepaymentRepetition:  12,
		InterestPercentage:   10,
	}

	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(nil, nil)

	// Act
	err := service.CreateLoan(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "borrower not found", err.Error())
	mockBorrowerRepo.AssertExpectations(t)
}

func TestLoanServiceImpl_CreateLoan_BorrowerRepositoryError(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewLoanService(mockLoanRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := &models.LoanRequest{
		BorrowerID:           "borrower-id",
		Amount:               1000000,
		RepaymentCadenceDays: 7,
		RepaymentRepetition:  12,
		InterestPercentage:   10,
	}

	expectedError := errors.New("database error")
	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(nil, expectedError)

	// Act
	err := service.CreateLoan(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockBorrowerRepo.AssertExpectations(t)
}
