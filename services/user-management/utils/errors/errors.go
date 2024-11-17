package errors

import (
	"errors"
)

var (
	ErrorBadRequest   = New("Bad Request")
	ErrorUnauthorized = New("Unauthorized")
	ErrorNotFound     = New("Not Found")
	ErrorInternal     = New("Internal Server Error")
	ErrorConflict     = New("Conflict")
	ErrorForbidden    = New("Forbidden")
)

type Error struct {
	appError error
	svcError error
}

func NewError(svcError error, appError error) Error {
	return Error{
		svcError: svcError,
		appError: appError,
	}
}

func (e Error) AppError() error {
	return e.appError
}

func (e Error) SvcError() error {
	return e.svcError
}

func (e Error) Error() string {
	return errors.Join(e.appError, e.svcError).Error()
}

func New(err string) error {
	return errors.New(err)
}
