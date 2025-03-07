package dtos

type ErrorDTO struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}
