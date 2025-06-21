package models

type ErrorResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message   string      `json:"message" example:"Something went wrong while processing your request."` // Error message to be shown for user
	ErrorText string      `json:"error,omitempty" example:"Nil pointer reference"`                       // Error message for debugging
	Result    interface{} `json:"result,omitempty"`                                                      // Custom data for needed for specific case
}

type LoanResponse struct {
	ID                   string  `json:"id"`
	Amount               float64 `json:"amount"`
	RepaymentCadenceDays int     `json:"repayment_cadence_days"`
	RepaymentRepetition  int     `json:"repayment_repetition"`
	InterestPercentage   float64 `json:"interest_percentage"`
	InterestAmount       float64 `json:"interest_amount"`
	Status               string  `json:"status"`
	TotalOutstanding     float64 `json:"total_outstanding"`
}

type PaymentLinkResponse struct {
	ID                   string  `json:"id"`
	TotalRepaymentAmount float64 `json:"total_repayment_amount"`
	PaymentLink          string  `json:"payment_link"`
}

type BorrowerResponse struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
	IsDelinquent bool   `json:"is_delinquent"`
}
