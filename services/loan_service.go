package services

import (
	"context"

	"github.com/satryarangga/amartha-loan-engine/models"
)

type LoanService interface {
	GetLoans(ctx context.Context) ([]models.Loan, error)
	GetLoanByID(ctx context.Context, id string) (*models.Loan, error)
	CreateLoan(ctx context.Context, loan *models.Loan) error
	UpdateLoan(ctx context.Context, loan *models.Loan) error
	DeleteLoan(ctx context.Context, id string) error
}
