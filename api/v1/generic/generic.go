package generic

// ErrorResponse is the response body for custom errors. Can be used anywhere where failure response body isnt available or required.
type ErrorResponse struct {
	Message string `json:"message" yaml:"message"`
	Error   string `json:"error" yaml:"error"`
} //@name GenericErrorResponse

// Response is the response body for custom messages. Can be used anywhere where success response body isnt available or required.
type Response struct {
	Message string `json:"message" yaml:"message"`
} //@name GenericResponse
