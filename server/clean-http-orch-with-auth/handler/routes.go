package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/pkg/response"
)

func (h *httpHandler) initRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)

	v1 := e.Group("/v1", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return response.ErrorResponse(c, http.StatusUnauthorized, "error - missing token", "ERR1")
			}

			newCtx, err := h.d.Service.Auth(c, token)
			if err != nil {
				return response.ErrorResponse(c, http.StatusUnauthorized, "error - invalid token", "ERR1")
			}

			return next(newCtx)
		}
	})
	v1.GET("/example", h.DoExample)
	v1.GET("/test", h.DoWiremock)
}
