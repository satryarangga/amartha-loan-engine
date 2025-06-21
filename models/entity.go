package models

import (
	"time"

	"github.com/lib/pq"
)

type Borrower struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	PhoneNumber string    `gorm:"not null;unique" json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Loan struct {
	ID                   string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BorrowerID           string     `gorm:"type:uuid;not null" json:"borrower_id"`
	Amount               float64    `gorm:"not null" json:"amount"`
	RepaymentCadenceDays int        `gorm:"not null" json:"repayment_cadence_days"`
	RepaymentRepetition  int        `gorm:"not null" json:"repayment_repetition"`
	InterestPercentage   float64    `gorm:"not null" json:"interest_percentage"`
	InterestAmount       float64    `gorm:"not null" json:"interest_amount"`
	Status               LoanStatus `gorm:"not null;default:'active'" json:"status"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	Borrower      Borrower       `gorm:"foreignKey:BorrowerID" json:"-"`
	LoanSchedules []LoanSchedule `gorm:"foreignKey:LoanID" json:"-"`
	LoanPayments  []LoanPayment  `gorm:"foreignKey:LoanID" json:"-"`
}

type LoanSchedule struct {
	ID             string             `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LoanID         string             `gorm:"type:uuid;not null" json:"loan_id"`
	DueDate        time.Time          `gorm:"not null" json:"due_date"`
	BasicAmount    float64            `gorm:"not null" json:"basic_amount"`
	InterestAmount float64            `gorm:"not null" json:"interest_amount"`
	TotalPayment   float64            `gorm:"not null" json:"total_payment"`
	Status         LoanScheduleStatus `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`

	Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

type LoanPayment struct {
	ID              string            `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LoanID          string            `gorm:"type:uuid;not null" json:"loan_id"`
	LoanScheduleIDs pq.StringArray    `gorm:"type:uuid[]" json:"loan_schedule_ids"`
	TotalPayment    float64           `gorm:"not null" json:"total_payment"`
	PaymentMethod   string            `gorm:"not null" json:"payment_method"`
	Status          LoanPaymentStatus `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}
