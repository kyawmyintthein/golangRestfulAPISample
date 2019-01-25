package model


type SuccessResponse struct{
	Success  bool `json:"success" example:"true"`
	Data     interface{}   `json:"data,omitempty" `
}

type ErrorResponse struct{
	Success  bool `json:"success" example:"false"`
	Error    HttpError `json:"error,omitempty"`
}

type HttpError struct {
	Code    uint32 `json:"code" example:"40001"`
	Message string `json:"message" example:"status bad request"`
}