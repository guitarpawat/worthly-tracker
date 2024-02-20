//go:generate mockery
//go:generate swag fmt
//go:generate swag init
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

//	@title			Worthly Tracker
//	@version		0.1
//	@host			localhost:8080
//	@schemes		http
//	@contact.name	Pawat Nakpiphatkul
//	@contact.url	https://github.com/guitarpawat/worthly-tracker/issues
func main() {
	e := echo.New()
	router.RegisterMiddleware(e)
	router.RegisterRoutes(e)

	err := e.Start(fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")))
	logs.Log().Panicf("Server error: %v\n", err)
}
