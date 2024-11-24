// Package response provides the response functions for the whole project
package response

import "github.com/labstack/echo/v4"

// ErrorResponsePayload represents the message of response payload
type ErrorResponsePayload struct {
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message"`
}

// ErrorResponse returns an error response wrapped with ErrorResponsePayload
func ErrorResponse(c echo.Context, statusCode int, message string, errorCode string) error {
	return c.JSON(statusCode, &ErrorResponsePayload{ErrorCode: errorCode, Message: message})
}

// SuccessResponse returns a success response with the data
func SuccessResponse(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, data)
}
