package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IController interface {
	RegisterRouter(api fiber.Router)
}

type EndpointHandler (func(c *fiber.Ctx) error)

func _BetterValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return strings.Split(fe.Error(), "Error:")[1] // default error
}

type APIResult[T any] struct {
	Data             *T                    `json:"data,omitempty"`
	Code             int                   `json:"code"`
	Error            string                `json:"error,omitempty"`
	ValidationErrors []_APIValidationError `json:"validationErrors,omitempty"`
}
type _APIValidationError struct {
	Param   string
	Message string
}

func OKResult[T any](v T, statusCode int) APIResult[T] {
	return APIResult[T]{
		Code: statusCode,
		Data: &v,
	}
}

// ErrResult constructs a Result with the given error set.
func ErrResult(err error, code int) APIResult[interface{}] {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]_APIValidationError, len(ve))
		for i, fe := range ve {
			out[i] = _APIValidationError{fe.Field(), _BetterValidationMessage(fe)}
		}
		return APIResult[any]{
			Data:             nil,
			Code:             code,
			Error:            "validation error",
			ValidationErrors: out,
		}
	}

	return APIResult[any]{
		Data:  nil,
		Code:  code,
		Error: err.Error(),
	}
}
