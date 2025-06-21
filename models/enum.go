package models

type LoanStatus string

const (
	LoanStatusActive LoanStatus = "active"
	LoanStatusPaid   LoanStatus = "paid"
)

type LoanScheduleStatus string

const (
	LoanScheduleStatusPending LoanScheduleStatus = "pending"
	LoanScheduleStatusPaid    LoanScheduleStatus = "paid"
)

type LoanPaymentStatus string

const (
	LoanPaymentStatusPending LoanPaymentStatus = "pending"
	LoanPaymentStatusPaid    LoanPaymentStatus = "paid"
)
