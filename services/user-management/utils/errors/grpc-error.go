package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
)

type GRPCError struct {
	Code    codes.Code `json:"-"`
	Status  int        `json:"status"`
	Message string     `json:"message"`
}

func (e GRPCError) String() string {
	output, _ := json.Marshal(e)
	return string(output)
}

func GRPCFromError(err error) GRPCError {
	var grpcError GRPCError
	var svcError Error

	if errors.As(err, &svcError) {
		grpcError.Message = svcError.AppError().Error()
		svcError := svcError.SvcError()
		switch svcError {
		case ErrorBadRequest:
			grpcError.Code = codes.InvalidArgument
			grpcError.Status = http.StatusBadRequest
		case ErrorUnauthorized:
			grpcError.Code = codes.Unauthenticated
			grpcError.Status = http.StatusUnauthorized
		case ErrorForbidden:
			grpcError.Code = codes.PermissionDenied
			grpcError.Status = http.StatusForbidden
		case ErrorNotFound:
			grpcError.Code = codes.NotFound
			grpcError.Status = http.StatusNotFound
		case ErrorConflict:
			grpcError.Code = codes.AlreadyExists
			grpcError.Status = http.StatusConflict
		case ErrorInternal:
			grpcError.Code = codes.Internal
			grpcError.Status = http.StatusInternalServerError
		}
	} else {
		grpcError.Code = codes.Unknown
		grpcError.Status = http.StatusInternalServerError
		grpcError.Message = err.Error()
	}

	// Remove the grpc.code from the message if it exists
	messageParts := strings.Split(grpcError.Message, "desc = ")
	if len(messageParts) > 1 {
		grpcError.Message = messageParts[1]
	}
	// Remove any newline and text after it
	messageParts = strings.Split(grpcError.Message, "\n")
	if len(messageParts) > 1 {
		grpcError.Message = messageParts[0]
	}

	return GRPCError{
		Code:    grpcError.Code,
		Status:  grpcError.Status,
		Message: grpcError.Message,
	}
}
