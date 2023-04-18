package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"worthly-tracker/config"
	"worthly-tracker/db"
	"worthly-tracker/logs"
	"worthly-tracker/router"
)

func init() {
	logs.Init()
	config.Init()
	db.Init()
}

func main() {
	e := echo.New()
	router.RegisterMiddleware(e)
	router.RegisterRoutes(e)

	err := e.Start(fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")))
	logs.Log().Panicf("Server error: %v\n", err)
}
