package repositories

import (
	"amartha/models"

	"gorm.io/gorm"
)

type BorrowerRepository struct {
	db *gorm.DB
}

func NewBorrowerRepository(db *gorm.DB) *BorrowerRepository {
	return &BorrowerRepository{db: db}
}

func (r *BorrowerRepository) Create(borrower *models.Borrower) error {
	return r.db.Create(borrower).Error
}

func (r *BorrowerRepository) FindAll() ([]models.Borrower, error) {
	var borrowers []models.Borrower
	err := r.db.Find(&borrowers).Error
	return borrowers, err
}

func (r *BorrowerRepository) FindByID(id string) (*models.Borrower, error) {
	var borrower models.Borrower
	err := r.db.Where("id = ?", id).First(&borrower).Error
	if err != nil {
		return nil, err
	}
	return &borrower, nil
}

func (r *BorrowerRepository) Update(borrower *models.Borrower) error {
	return r.db.Save(borrower).Error
}

func (r *BorrowerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Borrower{}).Error
}
