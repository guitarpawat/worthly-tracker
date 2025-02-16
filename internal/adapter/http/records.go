package http

import (
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/service"
	"github.com/guitarpawat/worthly-tracker/internal/view/component"
	"github.com/guitarpawat/worthly-tracker/internal/view/page"
	"github.com/labstack/echo/v4"
	"github.com/rickb777/date/v2"
	"net/http"
	"time"
)

type RecordHandler struct {
	service *service.Records
}

func NewRecordHandler(service *service.Records) *RecordHandler {
	return &RecordHandler{
		service: service,
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
		return render(ctx, http.StatusOK, component.GetRecordTable(resp))
	} else {
		return render(ctx, http.StatusOK, page.GetRecord(resp))
	}
}

func (h *RecordHandler) GetRecordsForDraft(ctx echo.Context) error {
	return nil
}

func (h *RecordHandler) PutRecord(ctx echo.Context) error {
	return nil
}

func (h *RecordHandler) PatchRecord(ctx echo.Context) error {
	return nil
}
