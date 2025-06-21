package services

import (
	"context"
	"errors"
	"testing"

	"github.com/satryarangga/amartha-loan-engine/mock"
	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewPaymentService(t *testing.T) {
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)

	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockLoanRepo, service.loanRepo)
	assert.Equal(t, mockLoanPaymentRepo, service.loanPaymentRepo)
	assert.Equal(t, mockLoanScheduleRepo, service.loanScheduleRepo)
	assert.Equal(t, mockBorrowerRepo, service.borrowerRepo)
}

func TestPaymentServiceImpl_GeneratePaymentLink_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentLinkRequest{
		BorrowerID:    "borrower-id",
		PaymentMethod: "bank_transfer",
	}

	borrower := &models.Borrower{
		ID:          "borrower-id",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}

	loan := models.Loan{
		ID:     "loan-id",
		Amount: 1000000,
	}

	loanSchedules := []models.LoanSchedule{
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
			Status:       models.LoanScheduleStatusPending,
		},
	}

	paymentID := "payment-id"

	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(borrower, nil)
	mockLoanRepo.On("FindOneByBorrowerID", ctx, "borrower-id").Return(loan, nil)
	mockLoanScheduleRepo.On("FindDueRepaymentSchedules", ctx, "loan-id").Return(loanSchedules, nil)
	mockLoanPaymentRepo.On("Insert", ctx, (*gorm.DB)(nil), testifymock.AnythingOfType("*models.LoanPayment")).Return(paymentID, nil)

	// Act
	result, err := service.GeneratePaymentLink(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, paymentID, result.ID)
	assert.Equal(t, 220000.0, result.TotalRepaymentAmount) // 110000 * 2
	assert.Contains(t, result.PaymentLink, paymentID)
	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
	mockLoanScheduleRepo.AssertExpectations(t)
	mockLoanPaymentRepo.AssertExpectations(t)
}

func TestPaymentServiceImpl_GeneratePaymentLink_BorrowerNotFound(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentLinkRequest{
		BorrowerID:    "borrower-id",
		PaymentMethod: "bank_transfer",
	}

	expectedError := errors.New("borrower not found")
	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(nil, expectedError)

	// Act
	result, err := service.GeneratePaymentLink(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockBorrowerRepo.AssertExpectations(t)
}

func TestPaymentServiceImpl_GeneratePaymentLink_LoanNotFound(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentLinkRequest{
		BorrowerID:    "borrower-id",
		PaymentMethod: "bank_transfer",
	}

	borrower := &models.Borrower{
		ID:          "borrower-id",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}

	expectedError := errors.New("loan not found")
	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(borrower, nil)
	mockLoanRepo.On("FindOneByBorrowerID", ctx, "borrower-id").Return(models.Loan{}, expectedError)

	// Act
	result, err := service.GeneratePaymentLink(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
}

func TestPaymentServiceImpl_GeneratePaymentLink_NoSchedules(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentLinkRequest{
		BorrowerID:    "borrower-id",
		PaymentMethod: "bank_transfer",
	}

	borrower := &models.Borrower{
		ID:          "borrower-id",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
	}

	loan := models.Loan{
		ID:     "loan-id",
		Amount: 1000000,
	}

	mockBorrowerRepo.On("FindByID", ctx, "borrower-id", []string{}).Return(borrower, nil)
	mockLoanRepo.On("FindOneByBorrowerID", ctx, "borrower-id").Return(loan, nil)
	mockLoanScheduleRepo.On("FindDueRepaymentSchedules", ctx, "loan-id").Return([]models.LoanSchedule{}, nil)

	// Act
	result, err := service.GeneratePaymentLink(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no loan schedules found", err.Error())
	mockBorrowerRepo.AssertExpectations(t)
	mockLoanRepo.AssertExpectations(t)
	mockLoanScheduleRepo.AssertExpectations(t)
}

func TestPaymentServiceImpl_HandlePaymentWebhook_Success(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentWebhookRequest{
		ExternalID:    "payment-id",
		PaymentStatus: "paid",
	}

	mockLoanRepo.On("WithTransaction", ctx, testifymock.AnythingOfType("repositories.TransactionFunc")).Return(nil)

	// Act
	err := service.HandlePaymentWebhook(ctx, request)

	// Assert
	assert.NoError(t, err)
	mockLoanRepo.AssertExpectations(t)
}

func TestPaymentServiceImpl_HandlePaymentWebhook_NotPaidStatus(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentWebhookRequest{
		ExternalID:    "payment-id",
		PaymentStatus: "failed",
	}

	// Act
	err := service.HandlePaymentWebhook(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payment status from PG is not paid", err.Error())
}

func TestPaymentServiceImpl_HandlePaymentWebhook_PaymentNotFound(t *testing.T) {
	// Arrange
	mockLoanRepo := mock.NewLoanRepository(t)
	mockLoanPaymentRepo := mock.NewLoanPaymentRepository(t)
	mockLoanScheduleRepo := mock.NewLoanScheduleRepository(t)
	mockBorrowerRepo := mock.NewBorrowerRepository(t)
	service := NewPaymentService(mockLoanRepo, mockLoanPaymentRepo, mockLoanScheduleRepo, mockBorrowerRepo)

	ctx := context.Background()
	request := models.PaymentWebhookRequest{
		ExternalID:    "payment-id",
		PaymentStatus: "paid",
	}

	expectedError := errors.New("payment not found")
	mockLoanRepo.On("WithTransaction", ctx, testifymock.AnythingOfType("repositories.TransactionFunc")).Return(expectedError)

	// Act
	err := service.HandlePaymentWebhook(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockLoanRepo.AssertExpectations(t)
}
