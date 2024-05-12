package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "worthly-tracker/docs"
)

func RegisterRoutes(e *echo.Echo) {
	staticRouter(e.Group(""))

	api := e.Group("/api")
	recordsRouter(api.Group("/records"))
	assetsManagementRouter(api.Group("/assets_management"))
	configsRouter(api.Group("/configs"))

	swagger := e.Group("/swagger")
	swagger.Any("*", echoSwagger.WrapHandler)
	swagger.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	swagger.GET("//index.html", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
}
