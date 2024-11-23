package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/pkg/response"
)

// DoProduce handles the request to do kafka produce function
// @Summary Run DoProduce function
// @Success 200 string string
// @Router /v1/kafka [get]
func (h *httpHandler) DoKafkaProduce(c echo.Context) error {
	ctx := context.Background()
	err := h.d.Service.DoKafkaProduce(ctx)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error - [DoKafkaProduce] unable to produce data to Kafka: %v", err), "ERR0")
	}
	return response.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
