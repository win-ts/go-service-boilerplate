package service

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
)

// Auth authenticates the request of http handler middleware
func (s *service) Auth(c echo.Context, token string) (echo.Context, error) {
	ctx := c.Request().Context()

	if err := s.authMiddlewareRepository.VerifyToken(ctx, strings.TrimPrefix(token, "Bearer ")); err != nil {
		slog.Error("error verifying token",
			slog.Any("error", err),
		)
		return nil, err
	}

	return c, nil
}
