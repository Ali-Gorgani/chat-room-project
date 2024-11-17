package handler

import (
	"log"
	"strconv"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/errors"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(usecase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: *usecase,
	}
}

// FindUserByID godoc
// @Summary      Get user by ID
// @Description  Get details of a specific user by their ID
// @Tags         Users
// @Security BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [get]
func (h *UserHandler) FindUserByID(ctx *fiber.Ctx) error {
	strId := ctx.Params("id")
	log.Printf("Received id: %s", strId) // Log the id value
	id, err := strconv.Atoi(strId)
	if err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	user := domain.User{
		ID: id,
	}
	foundUser, err := h.usecase.FindUserByID(ctx.Context(), user)
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainUserToUserResponse(foundUser)

	return ctx.JSON(res)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body  CreateUserRequest  true  "User to create"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	var user CreateUserRequest
	if err := ctx.BodyParser(&user); err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	createdUser, err := h.usecase.CreateUser(ctx.Context(), CreateUserRequestToDomainUser(user))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainUserToUserResponse(createdUser)

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  Update the details of a specific user
// @Tags         Users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path    int                true  "User ID"
// @Param        user  body    UpdateUserRequest  true  "Updated user details"
// @Success      200   {object}  UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      403   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	strId := ctx.Params("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	var user UpdateUserRequest
	if err := ctx.BodyParser(&user); err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	updatedUser, err := h.usecase.UpdateUser(ctx.Context(), UpdateUserRequestToDomainUser(user, id))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainUserToUserResponse(updatedUser)

	return ctx.JSON(res)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         Users
// @Security BearerAuth
// @Param        id  path  int  true  "User ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(ctx *fiber.Ctx) error {
	strId := ctx.Params("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		apiErr := errors.FromError(errors.NewError(errors.ErrorBadRequest, err))
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	user := domain.User{
		ID: id,
	}
	err = h.usecase.DeleteUser(ctx.Context(), user)
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
