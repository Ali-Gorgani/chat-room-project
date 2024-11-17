package errors

import (
	"encoding/json"
	"strings"

	"google.golang.org/grpc/codes"
)

// UnifiedError represents a unified structure for HTTP and gRPC errors.
type UnifiedError struct {
	Status  int        `json:"status"`
	Code    codes.Code `json:"-"`
	Message string     `json:"message"`
	Type    string     `json:"type"` // "HTTP" or "gRPC"
}

func (e UnifiedError) String() string {
	output, _ := json.Marshal(e)
	return string(output)
}

// FromError maps an error to UnifiedError, leveraging HTTPFromError and GRPCFromError.
func FromError(err error) UnifiedError {
	if err == nil {
		return UnifiedError{}
	}

	if strings.Contains(err.Error(), "rpc") {
		// If it's a gRPC error, use GRPCFromError
		grpcError := GRPCFromError(err)
		return UnifiedError{
			Status:  grpcError.Status,
			Code:    grpcError.Code,
			Message: grpcError.Message,
			Type:    "gRPC",
		}
	}

	// For HTTP or other errors, use HTTPFromError
	httpError := HTTPFromError(err)
	return UnifiedError{
		Status:  httpError.Status,
		Code:    codes.Code(httpError.Status), // Approximate mapping for HTTP status
		Message: httpError.Message,
		Type:    "HTTP",
	}
}
