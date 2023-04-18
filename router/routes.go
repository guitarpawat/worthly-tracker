package router

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	staticRouter(e.Group(""))

	api := e.Group("/api")
	recordsRouter(api.Group("/records"))
}
