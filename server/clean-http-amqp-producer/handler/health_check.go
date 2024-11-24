package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/pkg/response"
)

// HealthCheck checks the health of the service
func (h *httpHandler) HealthCheck(c echo.Context) error {
	startTime := time.Now()
	duration := time.Since(startTime)

	if duration > 1000 {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error - [HealthCheck] service is unhealthy: %v", duration), "ERR0")
	}

	return response.SuccessResponse(c, http.StatusOK, "[HealthCheck] service is healthy")
}
