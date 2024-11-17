package middleware

import (
	"strings"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/errors"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a middleware that verifies the token from the Authorization header.
func AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := VerifyClaimsFromAuthHeader(ctx)
		if err != nil {
			apiErr := errors.FromError(err)
			return ctx.Status(apiErr.Status).JSON(apiErr)
		}

		// Pass the token to the context
		ctx.Locals("token", token)

		// Proceed to the next handler
		return ctx.Next()
	}
}

// VerifyClaimsFromAuthHeader verifies the token from the Authorization header.
func VerifyClaimsFromAuthHeader(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", errors.NewError(errors.ErrorBadRequest, errors.New("authorization header is missing"))
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
		return "", errors.NewError(errors.ErrorBadRequest, errors.New("invalid authorization header"))
	}
	tokenString := fields[1]

	return tokenString, nil
}
