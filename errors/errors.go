package errors

import (
	"net/http"
)

type AppError struct {
	Message string
	Code    int
}

func (Ae AppError) Error() string {
	return Ae.Message
}

func BadRequest(err error) error {
	return &AppError{
		Message: err.Error(),
		Code:    http.StatusBadRequest,
	}
}

func NotFound(err error) error {
	return &AppError{
		Message: err.Error(),
		Code:    http.StatusNotFound,
	}
}
