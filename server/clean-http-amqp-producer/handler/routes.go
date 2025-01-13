package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *httpHandler) initRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)

	v1 := e.Group("/v1")
	v1.GET("/example", h.DoExample)
	v1.GET("/test", h.DoWiremock)
	v1.GET("/db", h.DoDBTest)
	v1.GET("/cache", h.DoSetGetCache)
	v1.POST("/amqp", h.DoProduce)
}
