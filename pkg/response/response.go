// Package response provides the response functions for the whole project
package response

import "github.com/labstack/echo/v4"

// MsgResponse represents the message of response payload
type MsgResponse struct {
	Message string `json:"message"`
}

// ErrResponse returns an error response wrapped with MsgResponse
func ErrResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, &MsgResponse{Message: message})
}

// SuccessResponse returns a success response with the data
func SuccessResponse(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, data)
}
