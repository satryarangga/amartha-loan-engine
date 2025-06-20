package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"

	"gorm.io/gorm"
)

type LoanPaymentRepository struct {
	db *gorm.DB
}

func NewLoanPaymentRepository(db *gorm.DB) *LoanPaymentRepository {
	return &LoanPaymentRepository{db: db}
}

func (r *LoanPaymentRepository) Create(payment *models.LoanPayment) error {
	return r.db.Create(payment).Error
}

func (r *LoanPaymentRepository) FindByLoanID(loanID string) ([]models.LoanPayment, error) {
	var payments []models.LoanPayment
	err := r.db.Where("loan_id = ?", loanID).Find(&payments).Error
	return payments, err
}

func (r *LoanPaymentRepository) FindByID(id string) (*models.LoanPayment, error) {
	var payment models.LoanPayment
	err := r.db.Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *LoanPaymentRepository) Update(payment *models.LoanPayment) error {
	return r.db.Save(payment).Error
}
