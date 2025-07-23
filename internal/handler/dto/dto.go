package dto

import "time"

type GeneralResponse struct {
	Data   any    `json:"data"`
	Status Status `json:"status"`
}

type Status struct {
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewGeneralResponse(data any, status, message, err string, statusCode int) GeneralResponse {
	return GeneralResponse{
		Data: data,
		Status: Status{
			Status:     status,
			StatusCode: statusCode,
			Message:    message,
			Error:      err,
			TimeStamp:  time.Now().UTC(),
		},
	}
}
