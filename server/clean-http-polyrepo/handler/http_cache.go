package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/response"
)

// DoSetGetCache handles the request to do cache set and get function
// @Summary Run DoSetGetCache function
// @Success 200 string string
// @Router /v1/cache [get]
func (h *httpHandler) DoSetGetCache(c echo.Context) error {
	ctx := context.Background()
	test, err := h.d.Service.DoSetGetCache(ctx)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error - [DoSetGetCache] unable to retrieve data: %v", err), "ERR0")
	}
	return response.SuccessResponse(c, http.StatusOK, test)
}
