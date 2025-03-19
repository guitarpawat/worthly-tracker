package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/service"
	"github.com/guitarpawat/worthly-tracker/internal/view/component"
	"github.com/guitarpawat/worthly-tracker/internal/view/page"
	"github.com/labstack/echo/v4"
	"github.com/rickb777/date/v2"
	"net/http"
	"time"
)

type RecordHandler struct {
	service   *service.Records
	validator *validator.Validate
}

func NewRecordHandler(service *service.Records) *RecordHandler {
	return &RecordHandler{
		service:   service,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *RecordHandler) GetRecordsByDate(ctx echo.Context) error {
	var d date.Date
	var err error
	dateParam := ctx.QueryParam("date")
	if dateParam == "" {
		d = date.Zero
	} else {
		d, err = date.Parse(time.DateOnly, dateParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse input date: %w", err))
		}
	}
	resp, err := h.service.GetByDate(ctx.Request().Context(), d)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if ctx.Request().Header.Get(constant.HeaderKeyHtmxRequest) == constant.HeaderValueHtmxRequest {
		return htmxRender(ctx, http.StatusOK, component.GetRecordTable(resp))
	} else {
		return htmxRender(ctx, http.StatusOK, page.GetRecord(resp))
	}
}

func (h *RecordHandler) GetRecordsForDraft(ctx echo.Context) error {
	draft, err := h.service.GetDraft(ctx.Request().Context())
	if err != nil {
		return err
	}

	return htmxRender(ctx, http.StatusOK, page.EditRecord(draft))
}

func (h *RecordHandler) GetRecordsForEdit(ctx echo.Context) error {
	var d date.Date
	var err error
	dateParam := ctx.QueryParam("date")
	if dateParam == "" {
		d = date.Zero
	} else {
		d, err = date.Parse(time.DateOnly, dateParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse input date: %w", err))
		}
	}

	draft, err := h.service.GetByDateForEdit(ctx.Request().Context(), d)
	if err != nil {
		return err
	}

	return htmxRender(ctx, http.StatusOK, page.EditRecord(draft))
}

func (h *RecordHandler) CreateRecord(ctx echo.Context) error {
	var body model.CreateRecordRequest
	err := ctx.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	err = h.validator.Struct(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("fail while validate request body: %w", err))
	}

	err = h.service.CreateRecords(ctx.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return htmxRedirect(ctx, "/")
}

func (h *RecordHandler) UpdateRecord(ctx echo.Context) error {
	var body model.UpdateRecordRequest
	err := ctx.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	err = h.validator.Struct(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("fail while validate request body: %w", err))
	}

	err = h.service.UpdateRecords(ctx.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return htmxRedirect(ctx, "/")
}

func (h *RecordHandler) DeleteRecord(ctx echo.Context) error {
	var d date.Date
	var err error
	dateParam := ctx.QueryParam("date")
	if dateParam == "" {
		d = date.Zero
	} else {
		d, err = date.Parse(time.DateOnly, dateParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse input date: %w", err))
		}
	}

	err = h.service.DeleteRecords(ctx.Request().Context(), d)
	if err != nil {
		return err
	}

	return htmxRedirect(ctx, "/")
}
