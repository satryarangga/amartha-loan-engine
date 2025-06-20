package repositories

import (
	"amartha/models"

	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

func (r *LoanRepository) FindAll() ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.Preload("Borrower").Preload("LoanSchedules").Find(&loans).Error
	return loans, err
}

func (r *LoanRepository) FindByID(id string) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.Preload("Borrower").Preload("LoanSchedules").Where("id = ?", id).First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *LoanRepository) FindByBorrowerID(borrowerID string) ([]models.Loan, error) {
	var loans []models.Loan
	err := r.db.Preload("Borrower").Preload("LoanSchedules").Where("borrower_id = ?", borrowerID).Find(&loans).Error
	return loans, err
}

func (r *LoanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

func (r *LoanRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Loan{}).Error
}
