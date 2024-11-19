package errors

import (
	"errors"
	"net/http"
)

type APIError struct {
	Status  int
	Message string
}

func HTTPFromError(err error) APIError {
	var apiError APIError
	var svcError Error

	if errors.As(err, &svcError) {
		apiError.Message = svcError.AppError().Error()
		svcError := svcError.SvcError()
		switch svcError {
		case ErrorBadRequest:
			apiError.Status = http.StatusBadRequest
		case ErrorUnauthorized:
			apiError.Status = http.StatusUnauthorized
		case ErrorForbidden:
			apiError.Status = http.StatusForbidden
		case ErrorNotFound:
			apiError.Status = http.StatusNotFound
		case ErrorConflict:
			apiError.Status = http.StatusConflict
		case ErrorInternal:
			apiError.Status = http.StatusInternalServerError
		}
	}
	return apiError
}
