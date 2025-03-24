package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// BaseResponse is the standard response format for the API
type BaseResponse struct {
	Status     string      `json:"status"`     // "success" or "error"
	Message    string      `json:"message"`    // A human-readable message
	Data       interface{} `json:"data"`       // The actual response data, can be null for error responses
	Error      interface{} `json:"error"`      // Error details, null for success responses
	Pagination *Pagination `json:"pagination"` // Pagination info, null for non-paginated responses
	Timestamp  time.Time   `json:"timestamp"`  // When the response was generated
	RequestID  string      `json:"request_id"` // Unique identifier for the request
}

// Pagination contains pagination metadata
type Pagination struct {
	Total   int64 `json:"total"`    // Total number of records
	Page    int   `json:"page"`     // Current page number
	PerPage int   `json:"per_page"` // Number of records per page
	Pages   int   `json:"pages"`    // Total number of pages
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code    string      `json:"code"`              // Machine-readable error code
	Message string      `json:"message"`           // Human-readable error message
	Details interface{} `json:"details,omitempty"` // Additional error details
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(ctx *fiber.Ctx, data interface{}, message string) *BaseResponse {
	return &BaseResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		Error:      nil,
		Pagination: nil,
		Timestamp:  time.Now(),
		RequestID:  ctx.GetRespHeader("X-Request-ID", ""),
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(ctx *fiber.Ctx, err ErrorDetail, message string) *BaseResponse {
	return &BaseResponse{
		Status:     "error",
		Message:    message,
		Data:       nil,
		Error:      err,
		Pagination: nil,
		Timestamp:  time.Now(),
		RequestID:  ctx.GetRespHeader("X-Request-ID", ""),
	}
}

// NewPaginatedResponse creates a new paginated success response
func NewPaginatedResponse(ctx *fiber.Ctx, data interface{}, pagination *Pagination, message string) *BaseResponse {
	return &BaseResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		Error:      nil,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  ctx.GetRespHeader("X-Request-ID", ""),
	}
}

// Success is a helper function to send a success response
func Success(ctx *fiber.Ctx, data interface{}, message string) error {
	return ctx.Status(fiber.StatusOK).JSON(NewSuccessResponse(ctx, data, message))
}

// Created is a helper function to send a 201 Created response
func Created(ctx *fiber.Ctx, data interface{}, message string) error {
	return ctx.Status(fiber.StatusCreated).JSON(NewSuccessResponse(ctx, data, message))
}

// NoContent is a helper function to send a 204 No Content response
func NoContent(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNoContent)
}

// Error is a helper function to send an error response
func Error(ctx *fiber.Ctx, statusCode int, errorCode string, errorMessage string, details interface{}) error {
	err := ErrorDetail{
		Code:    errorCode,
		Message: errorMessage,
		Details: details,
	}

	return ctx.Status(statusCode).JSON(NewErrorResponse(ctx, err, errorMessage))
}

// BadRequest is a helper function to send a 400 Bad Request response
func BadRequest(ctx *fiber.Ctx, errorMessage string, details interface{}) error {
	return Error(ctx, fiber.StatusBadRequest, "BAD_REQUEST", errorMessage, details)
}

// Unauthorized is a helper function to send a 401 Unauthorized response
func Unauthorized(ctx *fiber.Ctx, errorMessage string) error {
	return Error(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", errorMessage, nil)
}

// Forbidden is a helper function to send a 403 Forbidden response
func Forbidden(ctx *fiber.Ctx, errorMessage string) error {
	return Error(ctx, fiber.StatusForbidden, "FORBIDDEN", errorMessage, nil)
}

// NotFound is a helper function to send a 404 Not Found response
func NotFound(ctx *fiber.Ctx, errorMessage string) error {
	return Error(ctx, fiber.StatusNotFound, "NOT_FOUND", errorMessage, nil)
}

// InternalServerError is a helper function to send a 500 Internal Server Error response
func InternalServerError(ctx *fiber.Ctx, errorMessage string) error {
	return Error(ctx, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", errorMessage, nil)
}

// WithPagination is a helper function to send a paginated success response
func WithPagination(ctx *fiber.Ctx, data interface{}, total int64, page, perPage int, message string) error {
	pages := int(total) / perPage
	if int(total)%perPage > 0 {
		pages++
	}

	pagination := &Pagination{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Pages:   pages,
	}

	return ctx.Status(fiber.StatusOK).JSON(NewPaginatedResponse(ctx, data, pagination, message))
}
