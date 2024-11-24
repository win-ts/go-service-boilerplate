package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/pkg/request"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/pkg/response"
)

// DoProduce handles the request to produce message to AMQP
// @Summary Run DoProduce function
// @Accept json
// @Produce json
// @Param body body dto.Event true "Event"
// @Success 200 string string
// @Router /v1/amqp [post]
func (h *httpHandler) DoProduce(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(dto.Event)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error - [DoProduce] unable to bind request: %v", err), "ERR0")
	}

	if err := h.d.Service.DoProduce(ctx, *req); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error - [DoProduce] unable to retrieve data: %v", err), "ERR0")
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]interface{}{"success": true})
}
