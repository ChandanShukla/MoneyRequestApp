package response

type ClientSuccessResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
