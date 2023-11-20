package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"worthly-tracker/db"
	"worthly-tracker/model"
	"worthly-tracker/ports"
)

func recordsRouter(api *echo.Group) {
	rs := recordService{db.GetRecordRepo(), db.GetBoughtValueOffsetRepo(), db.GetDB()}
	api.GET("/draft", rs.getRecordDraft)
	api.POST("/", rs.postRecord)
	api.GET("/:date", rs.getRecordByDate)
	api.GET("/", rs.getRecordByDate)
	api.GET("/offset/:date", rs.getOffsetByDate)
}

type recordService struct {
	recordRepo ports.RecordRepo
	offsetRepo ports.BoughtValueOffsetRepo
	conn       ports.Connection
}

type getRecordByDateResponse struct {
	Date  *model.DateList         `json:"date"`
	Types []model.AssetTypeRecord `json:"types"`
}

func (r recordService) getRecordByDate(c echo.Context) error {
	dateParam := c.Param("date")
	var date *model.Date = nil
	var dateObj model.Date
	var err error
	if dateParam != "" {
		dateObj, err = model.NewDate(dateParam)
		date = &dateObj
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("date is invalid format (YYYY-MM-DD): %w", err))
	}

	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	if date == nil {
		if date, err = r.recordRepo.GetLatestDate(tx); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot get latest record date: %w", err))
		}
		if date == nil {
			return echo.NewHTTPError(http.StatusNotFound, "no any record in the system")
		}
	}

	dateList, err := r.recordRepo.GetDate(*date, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot get date list: %w", err))
	}

	records, err := r.recordRepo.GetRecordByDate(*date, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot get record: %w", err))
	}
	if records == nil || len(records) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no record found on requested date")
	}

	tx.Commit()

	return c.JSON(http.StatusOK, getRecordByDateResponse{
		Date:  dateList,
		Types: records,
	})
}

func (r recordService) getRecordDraft(c echo.Context) error {
	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	res, err := r.recordRepo.GetRecordDraft(tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot get record draft: %w", err))
	}

	tx.Commit()

	return c.JSON(http.StatusOK, res)
}

func (r recordService) getOffsetByDate(c echo.Context) error {
	dateParam := c.Param("date")
	if dateParam == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no date specified"))
	}
	var date model.Date
	var err error
	date, err = model.NewDate(dateParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("date is invalid format (YYYY-MM-DD): %w", err))
	}

	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	res, err := r.offsetRepo.GetAllByDate(date, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot get offset: %w", err))
	}

	tx.Commit()

	return c.JSON(http.StatusOK, res)
}

type postRecordRequest struct {
	Assets []model.AssetRecord `json:"assets"`
	Date   *model.Date         `json:"date"`
}

func (r recordService) postRecord(c echo.Context) error {
	var body postRecordRequest
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	if body.Date == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing date")
	}

	if body.Assets == nil || len(body.Assets) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "nothing to upsert")
	}

	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	for _, v := range body.Assets {
		err = r.recordRepo.UpsertRecord(v, *body.Date, tx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot upsert asset record: %w", err))
		}
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return c.JSON(http.StatusOK, nil)
}
