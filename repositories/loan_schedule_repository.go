package repositories

import (
	"github.com/satryarangga/amartha-loan-engine/models"

	"gorm.io/gorm"
)

type LoanScheduleRepository struct {
	db *gorm.DB
}

func NewLoanScheduleRepository(db *gorm.DB) *LoanScheduleRepository {
	return &LoanScheduleRepository{db: db}
}

func (r *LoanScheduleRepository) Create(schedule *models.LoanSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *LoanScheduleRepository) CreateBatch(schedules []models.LoanSchedule) error {
	return r.db.Create(&schedules).Error
}

func (r *LoanScheduleRepository) FindByLoanID(loanID string) ([]models.LoanSchedule, error) {
	var schedules []models.LoanSchedule
	err := r.db.Where("loan_id = ?", loanID).Find(&schedules).Error
	return schedules, err
}

func (r *LoanScheduleRepository) FindByID(id string) (*models.LoanSchedule, error) {
	var schedule models.LoanSchedule
	err := r.db.Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *LoanScheduleRepository) Update(schedule *models.LoanSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *LoanScheduleRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.LoanSchedule{}).Where("id = ?", id).Update("status", status).Error
}

func (r *LoanScheduleRepository) FindPendingByLoanID(loanID string) ([]models.LoanSchedule, error) {
	var schedules []models.LoanSchedule
	err := r.db.Where("loan_id = ? AND status = ?", loanID, "pending").Find(&schedules).Error
	return schedules, err
}
