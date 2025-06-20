package models

type ErrorResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message   string      `json:"message" example:"Something went wrong while processing your request."` // Error message to be shown for user
	ErrorText string      `json:"error,omitempty" example:"Nil pointer reference"`                       // Error message for debugging
	Result    interface{} `json:"result,omitempty"`                                                      // Custom data for needed for specific case
}
