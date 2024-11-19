package handler

import (
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/errors"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase usecase.AuthUseCase
}

func NewAuthHandler(usecase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		usecase: *usecase,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login request body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var loginRequest LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	loginResponse, err := h.usecase.Login(ctx.Context(), LoginRequestToDomainAuth(loginRequest))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainAuthToLoginResponse(loginResponse)

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  res.RefreshTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// Logout godoc
// @Summary Logout user
// @Description Logout by deleting the refresh token cookie and entry in database
// @Tags auth
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /logout [post]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	refreshTokenCookie := ctx.Cookies("refresh_token")
	if refreshTokenCookie == "" {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, fmt.Errorf("refresh token cookie is missing")))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	auth := domain.Auth{
		RefreshToken: refreshTokenCookie,
	}

	err := h.usecase.Logout(ctx.Context(), auth)
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	ctx.ClearCookie("refresh_token")

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh the access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshTokenRequest body RefreshTokenRequest true "Refresh token request body"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /refresh-token [post]
func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	var refreshTokenRequest RefreshTokenRequest
	if err := ctx.BodyParser(&refreshTokenRequest); err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	refreshTokenResponse, err := h.usecase.RefreshToken(ctx.Context(), RefreshTokenRequestToDomainAuth(refreshTokenRequest))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainAuthToRefreshTokenResponse(refreshTokenResponse)

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// RevokeToken godoc
// @Summary Revoke a refresh token
// @Description Revoke a refresh token to invalidate future access tokens
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param revokeTokenRequest body RevokeTokenRequest true "Revoke token request body"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /revoke-token [post]
func (h *AuthHandler) RevokeToken(ctx *fiber.Ctx) error {
	var revokeTokenRequest RevokeTokenRequest
	if err := ctx.BodyParser(&revokeTokenRequest); err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	err := h.usecase.RevokeToken(ctx.Context(), RevokeTokenRequestToDomainAuth(revokeTokenRequest))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}
