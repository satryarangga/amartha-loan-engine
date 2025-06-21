package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/satryarangga/amartha-loan-engine/mock"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewBorrowerService(t *testing.T) {
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.borrowerRepo)
	assert.Equal(t, mockLoanRepo, service.loanRepo)
}

func TestBorrowerServiceImpl_GetBorrowerByID_Success(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()
	borrowerID := "test-borrower-id"
	expectedBorrower := &models.Borrower{
		ID:          borrowerID,
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	expectedLoan := models.Loan{
		ID:         "test-loan-id",
		BorrowerID: borrowerID,
		Amount:     100000,
		Status:     models.LoanStatusActive,
		LoanSchedules: []models.LoanSchedule{
			{
				ID:           "schedule-1",
				LoanID:       "test-loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPending,
				DueDate:      time.Now().AddDate(0, 0, -5), // 5 days overdue
			},
			{
				ID:           "schedule-2",
				LoanID:       "test-loan-id",
				TotalPayment: 110000,
				Status:       models.LoanScheduleStatusPending,
				DueDate:      time.Now().AddDate(0, 0, -3), // 3 days overdue
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("FindByID", ctx, borrowerID, []string{}).Return(expectedBorrower, nil)
	mockLoanRepo.On("FindOneByBorrowerID", ctx, borrowerID).Return(expectedLoan, nil)
	// Act
	result, err := service.GetBorrowerByID(ctx, borrowerID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, borrowerID, result.ID)
	assert.Equal(t, expectedBorrower.FirstName, result.FirstName)
	assert.Equal(t, expectedBorrower.LastName, result.LastName)
	assert.Equal(t, expectedBorrower.PhoneNumber, result.PhoneNumber)
	assert.True(t, result.IsDelinquent) // Should be delinquent with 2 overdue payments
	mockRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

func TestBorrowerServiceImpl_GetBorrowerByID_EmptyID(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()

	// Act
	result, err := service.GetBorrowerByID(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "borrower ID is required", err.Error())
}

func TestBorrowerServiceImpl_GetBorrowerByID_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()
	borrowerID := "test-borrower-id"
	expectedError := errors.New("database error")

	mockRepo.On("FindByID", ctx, borrowerID, []string{}).Return(nil, expectedError)

	// Act
	result, err := service.GetBorrowerByID(ctx, borrowerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestBorrowerServiceImpl_GetBorrowerByID_LoanRepositoryError(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()
	borrowerID := "test-borrower-id"
	expectedBorrower := &models.Borrower{
		ID:          borrowerID,
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	expectedError := errors.New("loan database error")

	mockRepo.On("FindByID", ctx, borrowerID, []string{}).Return(expectedBorrower, nil)
	mockLoanRepo.On("FindOneByBorrowerID", ctx, borrowerID).Return(models.Loan{}, expectedError)

	// Act
	result, err := service.GetBorrowerByID(ctx, borrowerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

func TestBorrowerServiceImpl_CreateBorrower_Success(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()
	borrower := &models.Borrower{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}
	expectedID := "new-borrower-id"

	mockRepo.On("Insert", ctx, (*gorm.DB)(nil), borrower).Return(expectedID, nil)

	// Act
	err := service.CreateBorrower(ctx, borrower)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBorrowerServiceImpl_CreateBorrower_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := mock.NewBorrowerRepository(t)
	mockLoanRepo := mock.NewLoanRepository(t)
	service := NewBorrowerService(mockRepo, mockLoanRepo)

	ctx := context.Background()
	borrower := &models.Borrower{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}
	expectedError := errors.New("database error")

	mockRepo.On("Insert", ctx, (*gorm.DB)(nil), borrower).Return("", expectedError)

	// Act
	err := service.CreateBorrower(ctx, borrower)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
