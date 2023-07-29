package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"worthly-tracker/logs"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				logs.Log().Errorf("REQUEST: uri: %v, error: %v\n", v.URI, v.Error)
			} else {
				logs.Log().Infof("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())
}
