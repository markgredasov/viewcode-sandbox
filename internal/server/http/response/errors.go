package response

type ErrorResponse struct {
	Error   error  `json:"error"`
	Message string `json:"message" example:"INTERNAL_ERROR"`
}
