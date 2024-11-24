package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DoDBTest handles the request to do db test query function
// @Summary Run DoDBTest function
// @Success 200 string string
// @Router /v1/db[get]
func (h *httpHandler) DoDBTest(c echo.Context) error {
	res, err := h.d.Service.DoDBTest()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
