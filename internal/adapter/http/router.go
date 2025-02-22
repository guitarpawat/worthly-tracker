package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/view/component"
	"github.com/guitarpawat/worthly-tracker/internal/view/page"
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

func (r *Router) Close() error {
	return r.e.Shutdown(context.Background())
}

func (r *Router) registerMiddleWare() {
	r.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				logs.Log().Errorf("REQUEST: uri: %v %v, error: %v\n", c.Request().Method, v.URI, v.Error)
			} else {
				logs.Log().Infof("REQUEST: uri: %v %v, status: %v\n", c.Request().Method, v.URI, v.Status)
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

	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		var finalErr error
		var finalHttpCode int
		httpErr := new(echo.HTTPError)
		if errors.As(err, &httpErr) {
			if httpErr.Internal == nil {
				finalErr = fmt.Errorf("%v", httpErr.Message)
			} else {
				finalErr = httpErr.Internal
			}
			finalHttpCode = httpErr.Code
		} else {
			finalErr = err
			finalHttpCode = http.StatusInternalServerError
		}

		if c.Request().Header.Get(constant.HeaderKeyHtmxRequest) == constant.HeaderValueHtmxRequest {
			c.Response().Header().Set(constant.HeaderKeyHtmxRetarget, "#page-header")
			c.Response().Header().Set(constant.HeaderKeyHtmxReswap, "afterend")
			err = render(c, finalHttpCode, component.Error(finalHttpCode, true, finalErr))
		} else {
			err = render(c, finalHttpCode, page.Error(finalHttpCode, finalErr))
		}

		if err != nil {
			logs.Log().Errorf("cannot render error page: error: %v", err)
		}
	}
}

func (r *Router) registerStaticContent() {
	r.e.StaticFS("/", echo.MustSubFS(resource.StaticContent, "static"))
}

func (r *Router) registerRecordsRoutes() {
	r.e.GET("/", r.recordsHandler.GetRecordsByDate)
	r.e.GET("/records", r.recordsHandler.GetRecordsByDate)
	r.e.GET("/records/new", r.recordsHandler.GetRecordsForDraft)
	r.e.GET("/records/edit", r.recordsHandler.GetRecordsForEdit)
	r.e.POST("/records", r.recordsHandler.CreateRecord)
	r.e.PUT("/records", r.recordsHandler.UpdateRecord)
}

func render(ctx echo.Context, status int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(status, buf.String())
}
