package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanService interface {
	GetLoanByID(ctx context.Context, id string) (*models.Loan, error)
	CreateLoan(ctx context.Context, loan *models.LoanRequest) error
}
