package handler

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"cloud-sprint/internal/api/request"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/pkg/util"
)

// UserHandler handles user-related requests
type UserHandler struct {
	store db.Querier
}

// NewUserHandler creates a new user handler
func NewUserHandler(store db.Querier) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// GetCurrentUser gets the current authenticated user
// @Summary Get current user
// @Description Get the current authenticated user's information
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.store.GetUserByID(c.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
	}

	userResponse := response.NewUserResponse(user)
	return response.Success(c, userResponse, "User retrieved successfully")
}

// GetUser gets a user by ID
// @Summary Get a user by ID
// @Description Get a user's information by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.store.GetUserByID(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
	}

	userResponse := response.NewUserResponse(user)
	return response.Success(c, userResponse, "User retrieved successfully")
}

// ListUsers lists users with pagination
// @Summary List users
// @Description List users with pagination
// @Tags users
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {array} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	users, err := h.store.ListUsers(c.Context(), db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list users")
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.NewUserResponse(user))
	}

	total := int64(len(userResponses))
	return response.WithPagination(c, userResponses, total, offset/limit+1, limit, "Users retrieved successfully")
}

// UpdateUser updates a user
// @Summary Update a user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.UpdateUserRequest true "Update user request"
// @Security BearerAuth
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Get the user ID from the URL
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get the current user ID
	currentUserID, ok := c.Locals("userID").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Only allow users to update their own profile
	if currentUserID != int64(id) {
		return fiber.NewError(fiber.StatusForbidden, "Cannot update other users")
	}

	// Parse request body
	var req request.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create update params
	params := db.UpdateUserParams{
		ID: int64(id),
	}

	// Set optional fields
	if req.Username != nil {
		params.Username = sql.NullString{
			String: *req.Username,
			Valid:  true,
		}
	}

	if req.Email != nil {
		params.Email = sql.NullString{
			String: *req.Email,
			Valid:  true,
		}
	}

	if req.FullName != nil {
		params.FullName = sql.NullString{
			String: *req.FullName,
			Valid:  true,
		}
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
		}

		params.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	// Update user
	user, err := h.store.UpdateUser(c.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}

	userResponse := response.NewUserResponse(user)
	return response.Success(c, userResponse, "User updated successfully")
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Get the user ID from the URL
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get the current user ID
	currentUserID, ok := c.Locals("userID").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Only allow users to delete their own profile
	if currentUserID != int64(id) {
		return fiber.NewError(fiber.StatusForbidden, "Cannot delete other users")
	}

	// Delete user
	err = h.store.DeleteUser(c.Context(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}

	return response.NoContent(c)
}
