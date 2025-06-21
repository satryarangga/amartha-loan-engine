package models

type BorrowerRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoanRequest struct {
	BorrowerID           string  `json:"borrower_id" binding:"required" description:"Borrower ID"`
	Amount               float64 `json:"amount" binding:"required" description:"Loan amount"`
	RepaymentCadenceDays int     `json:"repayment_cadence_days" binding:"required" description:"Repayment cadence days (If weekly then 7)"`
	RepaymentRepetition  int     `json:"repayment_repetition" binding:"required" description:"How many times the loan will be repaid"`
	InterestPercentage   float64 `json:"interest_percentage" binding:"required" description:"Interest percentage"`
}

type PaymentLinkRequest struct {
	BorrowerID    string `json:"borrower_id" binding:"required" description:"Borrower ID"`
	PaymentMethod string `json:"payment_method" binding:"required" description:"Payment method"`
}

type PaymentWebhookRequest struct {
	LoanPaymentID string `json:"loan_payment_id" binding:"required" description:"Loan payment ID"`
	PaymentStatus string `json:"payment_status" binding:"required" description:"Payment status"`
}
