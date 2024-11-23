package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/pkg/response"
)

// DoWiremock handles the request to do wiremock function
// @Summary Run DoWiremock function
// @Success 200 string string
// @Router /v1/test [get]
func (h *httpHandler) DoWiremock(c echo.Context) error {
	ctx := context.Background()
	getTest, err := h.d.Service.DoWiremock(ctx)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error - [DoExample] unable to retrieve data: %v", err), "ERR0")
	}
	return response.SuccessResponse(c, http.StatusOK, getTest)
}
