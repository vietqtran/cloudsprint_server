package response

import (
	"cloud-sprint/internal/constants"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Code       constants.HttpStatusCode `json:"code"`
	Data       interface{}              `json:"data"`
	Pagination *Pagination              `json:"pagination,omitempty"`
	Timestamp  time.Time                `json:"timestamp"`
	RequestID  string                   `json:"request_id,omitempty"`
	Trace      error                    `json:"-"`
}

type Pagination struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Pages   int   `json:"pages"`
}

func (r *BaseResponse) Send(c *fiber.Ctx) error {
	statusCode := int(r.Code)
	return c.Status(statusCode).JSON(r)
}

func NewSuccessResponse(c *fiber.Ctx, code constants.HttpStatusCode, data interface{}, message string) *BaseResponse {
	return &BaseResponse{
		Status:     "success",
		Message:    message,
		Code:       code,
		Data:       data,
		Pagination: nil,
		Timestamp:  time.Now(),
		RequestID:  c.GetRespHeader("X-Request-ID", ""),
	}
}

func NewErrorResponse(c *fiber.Ctx, code constants.HttpStatusCode, message string, trace error) *BaseResponse {
	return &BaseResponse{
		Status:    "error",
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now(),
		RequestID: c.GetRespHeader("X-Request-ID", ""),
		Trace:     trace,
	}
}

func NewPaginatedResponse(c *fiber.Ctx, code constants.HttpStatusCode, data interface{}, pagination *Pagination, message string) *BaseResponse {
	return &BaseResponse{
		Status:     "success",
		Message:    message,
		Code:       code,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  c.GetRespHeader("X-Request-ID", ""),
	}
}

func WithPagination(c *fiber.Ctx, data interface{}, total int64, page, perPage int, message string) error {
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

	return NewPaginatedResponse(c, constants.StatusOK, data, pagination, message).Send(c)
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return NewSuccessResponse(c, constants.StatusOK, data, message).Send(c)
}

func Created(c *fiber.Ctx, data interface{}, message string) error {
	return NewSuccessResponse(c, constants.StatusCreated, data, message).Send(c)
}

func BadRequest(c *fiber.Ctx, message string, err error) error {
	return NewErrorResponse(c, constants.StatusBadRequest, message, err).Send(c)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return NewErrorResponse(c, constants.StatusUnauthorized, message, nil).Send(c)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return NewErrorResponse(c, constants.StatusForbidden, message, nil).Send(c)
}

func NotFound(c *fiber.Ctx, message string) error {
	return NewErrorResponse(c, constants.StatusNotFound, message, nil).Send(c)
}

func InternalServerError(c *fiber.Ctx, message string, trace error) error {
	return NewErrorResponse(c, constants.StatusInternalServerError, message, trace).Send(c)
}
