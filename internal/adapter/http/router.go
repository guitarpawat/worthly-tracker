package http

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/guitarpawat/worthly-tracker/resource"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type RouterConfig struct {
	Port int
}

type Router struct {
	e              *echo.Echo
	cfg            RouterConfig
	recordsHandler *RecordHandler
}

func NewRouter(cfg RouterConfig, recordsHandler *RecordHandler) *Router {
	e := echo.New()
	r := &Router{
		e:              e,
		cfg:            cfg,
		recordsHandler: recordsHandler,
	}

	r.init()
	return r
}

func (r *Router) init() {
	r.registerMiddleWare()
	r.registerStaticContent()
	r.registerRecordsRoutes()
}

func (r *Router) Start() error {
	return r.e.Start(fmt.Sprintf("0.0.0.0:%d", r.cfg.Port))
}

func (r *Router) registerMiddleWare() {
	r.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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

	r.e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logs.Log().Errorf("RECOVER: error: %v\nStack: %s", err, string(stack))
			return err
		},
	}))
}

func (r *Router) registerStaticContent() {
	r.e.StaticFS("/", echo.MustSubFS(resource.StaticContent, "static"))
}

func (r *Router) registerRecordsRoutes() {
	r.e.GET("/", r.recordsHandler.GetRecordsByDate)
	r.e.GET("/records", r.recordsHandler.GetRecordsByDate)
	r.e.GET("/records/partial/summary_table", r.recordsHandler.GetPartialRecordsByDate)
}

func render(ctx echo.Context, status int, t templ.Component) error {
	err := t.Render(ctx.Request().Context(), ctx.Response().Writer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Errorf("failed to render response template: %w", err).Error(),
		})
	}

	ctx.Response().Writer.WriteHeader(status)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return nil
}
