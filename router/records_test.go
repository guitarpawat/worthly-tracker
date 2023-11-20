//go:build test || unit

package router

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"
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
	recordRepo *mocks.MockRecordRepo
	offsetRepo *mocks.MockBoughtValueOffsetRepo
	service    recordService
	dbMock     sqlmock.Sqlmock
}

func (s *RecordSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
}

func (s *RecordSuite) SetupTest() {
	s.dbMock = db.InitMock()
	s.recordRepo = mocks.NewMockRecordRepo(s.T())
	s.offsetRepo = mocks.NewMockBoughtValueOffsetRepo(s.T())
	s.service = recordService{
		recordRepo: s.recordRepo,
		offsetRepo: s.offsetRepo,
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

	s.True(s.recordRepo.AssertExpectations(s.T()))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "GetLatestDate"))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "GetRecordByDate"))
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

	s.recordRepo.On("GetDate", model.MustNewDate("2023-01-31"), mock.Anything).Return(
		&model.DateList{
			Current: pointy.Pointer(model.MustNewDate("2023-01-31")),
			Prev:    nil,
			Next:    nil,
		}, nil)

	s.recordRepo.On("GetRecordByDate", model.MustNewDate("2023-01-31"), mock.Anything).Return(
		make([]model.AssetTypeRecord, 0), nil)

	err := s.service.getRecordByDate(c)
	s.Require().ErrorContains(err, "no record found on requested date")
	s.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)

	s.Require().NoError(s.dbMock.ExpectationsWereMet())
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "GetLatestDate"))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "GetRecordByDate"))
}

func (s *RecordSuite) TestGetRecordByDate_404_NoRecordsInSystem() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	s.recordRepo.On("GetLatestDate", mock.Anything).Return(nil, nil)

	err := s.service.getRecordByDate(c)
	s.Require().ErrorContains(err, "no any record in the system")
	s.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)

	s.NoError(s.dbMock.ExpectationsWereMet())
	s.True(s.recordRepo.AssertExpectations(s.T()))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "GetRecordByDate"))
}

func (s *RecordSuite) TestGetRecordByDate_200_SuccessWithDate() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/records/:date")

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	date := model.MustNewDate("2023-01-31")
	s.recordRepo.On("GetLatestDate", mock.Anything).Return(pointy.Pointer(date), nil)

	dateList := mockDateList()
	s.recordRepo.On("GetDate", date, mock.Anything).Return(pointy.Pointer(dateList), nil)

	records := mockRecord()
	s.recordRepo.On("GetRecordByDate", date, mock.Anything).Return([]model.AssetTypeRecord{records}, nil)

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
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
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
	s.recordRepo.On("GetDate", date, mock.Anything).Return(pointy.Pointer(dateList), nil)

	records := mockRecord()
	s.recordRepo.On("GetRecordByDate", date, mock.Anything).Return([]model.AssetTypeRecord{records}, nil)

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
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
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
	s.recordRepo.On("GetRecordDraft", mock.Anything).Return([]model.AssetTypeRecord{record}, nil)

	err := s.service.getRecordDraft(c)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, w.Code)

	expResp, err := json.Marshal([]model.AssetTypeRecord{record})
	s.Require().NoError(err)
	s.Require().Equal(strings.TrimSpace(string(expResp)), strings.TrimSpace(w.Body.String()))
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
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
			s.Require().True(s.recordRepo.AssertExpectations(s.T()))
			s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "UpsertRecord"))
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
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
	s.Require().True(s.recordRepo.AssertNotCalled(s.T(), "UpsertRecord"))
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

	s.recordRepo.On("UpsertRecord", mock.MatchedBy(func(a model.AssetRecord) bool {
		return *a.Id == 1 && a.BoughtValue.Equal(decimal.NewFromInt(500))
	}), model.MustNewDate("2022-01-31"), mock.Anything).Return(nil)

	err := s.service.postRecord(c)
	s.Require().NoError(err)
	s.Require().True(s.recordRepo.AssertExpectations(s.T()))
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

func (s *RecordSuite) TestGetOffsetByDate_400_InvalidDateFormat() {
	date := []string{
		"31-01-2023",
		"abc",
		"AAAA-AA-AA",
		"1234",
	}

	for _, d := range date {
		s.Run(d, func() {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			c := echo.New().NewContext(r, w)
			c.SetPath("/offset/:date")
			c.SetParamNames("date")
			c.SetParamValues(d)

			err := s.service.getOffsetByDate(c)
			s.Require().ErrorContains(err, "date is invalid format (YYYY-MM-DD):")
			s.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)

			s.True(s.offsetRepo.AssertExpectations(s.T()))
			s.Require().True(s.offsetRepo.AssertNotCalled(s.T(), "GetAllByDate"))
		})
	}
}

func (s *RecordSuite) TestGetOffsetByDate_400_NoDateParam() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/offset/:date")

	err := s.service.getOffsetByDate(c)
	s.Require().ErrorContains(err, "no date specified")
	s.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)

	s.True(s.offsetRepo.AssertExpectations(s.T()))
	s.Require().True(s.offsetRepo.AssertNotCalled(s.T(), "GetAllByDate"))
}

func (s *RecordSuite) TestGetOffsetByDate_200_EmptyRecord() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	date := model.MustNewDate("2023-01-31")
	c.SetPath("/offset/:date")
	c.SetParamNames("date")
	c.SetParamValues(date.String())

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	s.offsetRepo.EXPECT().GetAllByDate(date, mock.Anything).Return([]model.OffsetDetail{}, nil)

	err := s.service.getOffsetByDate(c)
	s.Require().NoError(err)
	s.Require().True(s.offsetRepo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())

	s.Require().Equal(http.StatusOK, w.Code)
	bodyParsed := make([]model.OffsetDetail, 0)
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	s.Require().NoError(err)
	s.Require().Empty(bodyParsed)
}

func (s *RecordSuite) TestGetOffsetByDate_200_WithResponse() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	date := model.MustNewDate("2023-01-31")
	c.SetPath("/offset/:date")
	c.SetParamNames("date")
	c.SetParamValues(date.String())

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	s.offsetRepo.EXPECT().GetAllByDate(date, mock.Anything).Return([]model.OffsetDetail{
		{
			Id:            pointy.Int(1),
			AssetId:       2,
			EffectiveDate: date,
			OffsetPrice:   decimal.NewFromInt(3),
			Note:          pointy.String("4"),
		},
	}, nil)

	err := s.service.getOffsetByDate(c)
	s.Require().NoError(err)
	s.Require().True(s.offsetRepo.AssertExpectations(s.T()))
	s.Require().NoError(s.dbMock.ExpectationsWereMet())

	s.Require().Equal(http.StatusOK, w.Code)
	bodyParsed := make([]model.OffsetDetail, 0)
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	s.Require().NoError(err)
	s.Require().Equal(1, len(bodyParsed))
	s.Require().Equal(1, *bodyParsed[0].Id)
	s.Require().Equal(2, bodyParsed[0].AssetId)
	s.Require().Equal(date, bodyParsed[0].EffectiveDate)
	s.Require().Equal(decimal.NewFromInt(3), bodyParsed[0].OffsetPrice)
	s.Require().Equal("4", *bodyParsed[0].Note)
}
