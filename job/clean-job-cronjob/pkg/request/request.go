// Package request provides the request functions for the whole project
package request

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		Context   echo.Context
		validator *validator.Validate
	}
)

// ContextWrapper wraps the echo.Context to provide the request functions
func ContextWrapper(ctx echo.Context) contextWrapperService {
	return &contextWrapper{
		Context:   ctx,
		validator: validator.New(),
	}
}

// Bind binds the data to the request context
func (c *contextWrapper) Bind(data any) error {
	if err := c.Context.Bind(data); err != nil {
		slog.Error("error - [request.Bind] bind data failed", slog.Any("error", err))
		return err
	}

	if err := c.validator.Struct(data); err != nil {
		slog.Error("error - [request.Bind] validate data failed", slog.Any("error", err))
		return err
	}

	return nil
}
