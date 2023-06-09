package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"worthly-tracker/db"
	"worthly-tracker/model"
)

func recordsRouter(api *echo.Group) {
	rs := recordService{db.GetRecordRepo(), db.GetDB()}
	api.GET("/draft", rs.getRecordDraft)
	api.POST("/", rs.postRecord)
	api.GET("/:date", rs.getRecordByDate)
	api.GET("/", rs.getRecordByDate)
}

type recordService struct {
	recordRepo db.RecordRepo
	conn       db.Connection
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("date is invalid format (YYYY-MM-DD): %v", err))
	}

	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot create db transaction: %v", err))
	}
	defer tx.Rollback()

	if date == nil {
		if date, err = r.recordRepo.GetLatestDate(tx); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot get latest record date: %v", err))
		}
		if date == nil {
			return echo.NewHTTPError(http.StatusNotFound, "no any record in the system")
		}
	}

	dateList, err := r.recordRepo.GetDate(*date, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot get date list: %v", err))
	}

	records, err := r.recordRepo.GetRecordByDate(*date, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot get record: %v", err))
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot create db transaction: %v", err))
	}
	defer tx.Rollback()

	res, err := r.recordRepo.GetRecordDraft(tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot get record draft: %v", err))
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("cannot bind request body: %v", err))
	}

	if body.Date == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing date")
	}

	if body.Assets == nil || len(body.Assets) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "nothing to upsert")
	}

	tx, err := r.conn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot create db transaction: %v", err))
	}
	defer tx.Rollback()

	for _, v := range body.Assets {
		err = r.recordRepo.UpsertRecord(v, *body.Date, tx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot upsert asset record: %v", err))
		}
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot commit db transaction: %v", err))
	}

	return c.JSON(http.StatusOK, nil)
}
