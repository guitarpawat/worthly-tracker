//go:build test || unit

package router

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/openlyinc/pointy"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/db"
	"worthly-tracker/logs"
	"worthly-tracker/mocks"
	"worthly-tracker/model"
)

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}

type RecordSuite struct {
	suite.Suite
	repo    *mocks.MockRecordRepo
	service recordService
	dbMock  sqlmock.Sqlmock
}

func (s *RecordSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
}

func (s *RecordSuite) SetupTest() {
	s.dbMock = db.InitMock()
	s.repo = mocks.NewMockRecordRepo(s.T())
	s.service = recordService{
		recordRepo: s.repo,
		conn:       db.GetDB(),
	}
}

func (s *RecordSuite) TestGetRecordByDate_400_InvalidDateFormat() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")
	c.SetParamNames("date")
	c.SetParamValues("31-01-2023")

	err := s.service.getRecordByDate(c)
	s.Require().ErrorContains(err, "date is invalid format (YYYY-MM-DD):")
	s.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)

	s.True(s.repo.AssertExpectations(s.T()))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "GetLatestDate"))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "GetRecordByDate"))
}

func (s *RecordSuite) TestGetRecordByDate_404_DateNotFound() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")
	c.SetParamNames("date")
	c.SetParamValues("2023-01-31")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	s.repo.On("GetDate", model.MustNewDate("2023-01-31"), mock.Anything).Return(
		&model.DateList{
			Current: pointy.Pointer(model.MustNewDate("2023-01-31")),
			Prev:    nil,
			Next:    nil,
		}, nil)

	s.repo.On("GetRecordByDate", model.MustNewDate("2023-01-31"), mock.Anything).Return(
		make([]model.AssetTypeRecord, 0), nil)

	err := s.service.getRecordByDate(c)
	s.Require().ErrorContains(err, "no record found on requested date")
	s.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)

	s.Require().NoError(s.dbMock.ExpectationsWereMet())
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "GetLatestDate"))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "GetRecordByDate"))
}

func (s *RecordSuite) TestGetRecordByDate_404_NoRecordsInSystem() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	s.repo.On("GetLatestDate", mock.Anything).Return(nil, nil)

	err := s.service.getRecordByDate(c)
	s.Require().ErrorContains(err, "no any record in the system")
	s.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)

	s.NoError(s.dbMock.ExpectationsWereMet())
	s.True(s.repo.AssertExpectations(s.T()))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "GetRecordByDate"))
}

func (s *RecordSuite) TestGetRecordByDate_200_SuccessWithDate() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	date := model.MustNewDate("2023-01-31")
	s.repo.On("GetLatestDate", mock.Anything).Return(pointy.Pointer(date), nil)

	dateList := mockDateList()
	s.repo.On("GetDate", date, mock.Anything).Return(pointy.Pointer(dateList), nil)

	records := mockRecord()
	s.repo.On("GetRecordByDate", date, mock.Anything).Return([]model.AssetTypeRecord{records}, nil)

	err := s.service.getRecordByDate(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	expResp, err := json.Marshal(
		getRecordByDateResponse{
			Date: pointy.Pointer(dateList),
			Types: []model.AssetTypeRecord{
				records,
			},
		},
	)
	s.Require().NoError(err)
	s.Require().Equal(strings.TrimSpace(string(expResp)), strings.TrimSpace(w.Body.String()))
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())
}

func (s *RecordSuite) TestGetRecordByDate_200_SuccessWithoutDate() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")
	c.SetParamNames("date")
	c.SetParamValues("2023-01-31")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	date := model.MustNewDate("2023-01-31")
	dateList := mockDateList()
	s.repo.On("GetDate", date, mock.Anything).Return(pointy.Pointer(dateList), nil)

	records := mockRecord()
	s.repo.On("GetRecordByDate", date, mock.Anything).Return([]model.AssetTypeRecord{records}, nil)

	err := s.service.getRecordByDate(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	expResp, err := json.Marshal(
		getRecordByDateResponse{
			Date: pointy.Pointer(dateList),
			Types: []model.AssetTypeRecord{
				records,
			},
		},
	)
	s.Require().NoError(err)
	s.Require().Equal(strings.TrimSpace(string(expResp)), strings.TrimSpace(w.Body.String()))
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())
}

func (s *RecordSuite) TestGetRecordDraft_200_Success() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/draft")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	record := mockRecord()
	s.repo.On("GetRecordDraft", mock.Anything).Return([]model.AssetTypeRecord{record}, nil)

	err := s.service.getRecordDraft(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	expResp, err := json.Marshal([]model.AssetTypeRecord{record})
	s.Require().NoError(err)
	s.Require().Equal(strings.TrimSpace(string(expResp)), strings.TrimSpace(w.Body.String()))
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())
}

func (s *RecordSuite) TestPostRecord_400_NoRecord() {
	body := []string{
		`{"date": "2022-01-31"}`,
		`{"date": "2022-01-31", "assets": null}`,
		`{"date": "2022-01-31", "assets": []}`,
	}

	for _, v := range body {
		s.Run(v, func() {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(v))
			r.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c := echo.New().NewContext(r, w)
			c.SetPath("/api/records/")

			err := s.service.postRecord(c)
			s.Require().ErrorContains(err, "nothing to upsert")
			s.Require().True(s.repo.AssertExpectations(s.T()))
			s.Require().True(s.repo.AssertNotCalled(s.T(), "UpsertRecord"))
			s.Require().NoError(s.dbMock.ExpectationsWereMet())
		})
	}

}

func (s *RecordSuite) TestPostRecord_400_NoDate() {
	body := `{"assets": [{"id": 1,"boughtValue": "500.00"}]}`

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/")

	err := s.service.postRecord(c)
	s.Require().ErrorContains(err, "missing date")
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().True(s.repo.AssertNotCalled(s.T(), "UpsertRecord"))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())
}

func (s *RecordSuite) TestPostRecord_200_Success() {
	body := `{"date": "2022-01-31","assets": [{"id": 1,"boughtValue": "500.00"}]}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	s.repo.On("UpsertRecord", mock.MatchedBy(func(a model.AssetRecord) bool {
		return *a.Id == 1 && a.BoughtValue.Equal(decimal.NewFromInt(500))
	}), model.MustNewDate("2022-01-31"), mock.Anything).Return(nil)

	err := s.service.postRecord(c)
	s.Require().NoError(err)
	s.Require().True(s.repo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())
}

func mockDateList() model.DateList {
	return model.DateList{
		Current: pointy.Pointer(model.MustNewDate("2022-02-01")),
		Prev: []model.Date{
			model.MustNewDate("2022-01-01"),
		},
		Next: nil,
	}
}

func mockRecord() model.AssetTypeRecord {
	return model.AssetTypeRecord{
		Id:          pointy.Int(1),
		Name:        pointy.String("test"),
		IsCash:      pointy.Bool(true),
		IsLiability: pointy.Bool(false),
		Assets: []model.AssetRecord{
			{
				Id:               pointy.Int(1),
				AssetId:          pointy.Int(1),
				Name:             pointy.String("test"),
				Broker:           pointy.String("test"),
				DefaultIncrement: pointy.Pointer(decimal.NewFromInt(50)),
				BoughtValue:      pointy.Pointer(decimal.NewFromInt(50)),
				CurrentValue:     pointy.Pointer(decimal.NewFromInt(50)),
				RealizedValue:    pointy.Pointer(decimal.NewFromInt(50)),
				Note:             nil,
			},
		},
	}
}
