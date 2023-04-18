package router

import (
	"github.com/labstack/echo/v4"
	"worthly-tracker/resource"
)

func staticRouter(e *echo.Group) {
	e.StaticFS("/css", resource.MustLoadDirFS("static/css"))
	e.StaticFS("/js", resource.MustLoadDirFS("static/js"))
	e.StaticFS("/img", resource.MustLoadDirFS("static/img"))
	e.StaticFS("/favicon.ico", resource.MustLoadDirFS("static/favicon.ico"))

	e.StaticFS("/", resource.MustLoadDirFS("static/index.html"))
	e.StaticFS("/add", resource.MustLoadDirFS("static/postRecord.html"))
	e.StaticFS("/edit", resource.MustLoadDirFS("static/postRecord.html"))
}
