package http

import (
	"errors"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/ports"
	"github.com/guitarpawat/worthly-tracker/internal/view"
	"github.com/labstack/echo/v4"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type RecordHandlerTestSuite struct {
	suite.Suite
	mockService *ports.MockRecordsService
	handler     *RecordHandler
	echo        *echo.Echo
}

func TestRecordHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RecordHandlerTestSuite))
}

func (s *RecordHandlerTestSuite) SetupTest() {
	s.mockService = new(ports.MockRecordsService)
	s.handler = NewRecordHandler(s.mockService)
	s.echo = echo.New()
}

func (s *RecordHandlerTestSuite) TearDownTest() {
	s.mockService.AssertExpectations(s.T())
}

func (s *RecordHandlerTestSuite) TestGetRecordsByDate_Error_InvalidDateFormat() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	q := url.Values{}
	q.Add("date", "invalid-date")
	c.SetPath("/records")
	c.QueryParams().Add("date", "invalid-date")

	// Execute
	err := s.handler.GetRecordsByDate(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "cannot parse input date")
}

func (s *RecordHandlerTestSuite) TestGetRecordsByDate_Error_ServiceError() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records")
	c.QueryParams().Add("date", d.String())

	s.mockService.On("GetByDate", d).Return(view.GetRecordTable{}, errors.New("service error"))

	// Execute
	err := s.handler.GetRecordsByDate(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusInternalServerError, httpErr.Code)
	s.Equal("service error", httpErr.Message.(error).Error())
}

func (s *RecordHandlerTestSuite) TestGetRecordsByDate_Success_NoDate_FullPage() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	response := view.GetRecordTable{
		Dates: view.DateResult{
			CurrentDate: date.Today(),
			PastDate:    []date.Date{},
			FutureDate:  []date.Date{},
		},
		AssetTypes: []view.RecordTableAssetType{
			{
				Id:     1,
				Name:   "Stock",
				IsCash: false,
				Records: []view.RecordTableData{
					{
						Id:               1,
						AssetId:          1,
						Name:             "test",
						Broker:           "test",
						BoughtValue:      decimal.NewFromInt(1000),
						CurrentValue:     decimal.NewFromInt(1500),
						RealizedValue:    decimal.NewFromInt(100),
						DefaultIncrement: decimal.Zero,
						IsCash:           false,
						Note:             "",
					},
				},
			},
		},
	}

	s.mockService.On("GetByDate", date.Zero).Return(response, nil)

	// Execute
	err := s.handler.GetRecordsByDate(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	// Should render the full page
	s.Contains(rec.Body.String(), "<!doctype html>")
}

func (s *RecordHandlerTestSuite) TestGetRecordsByDate_Success_WithDate_HtmxRequest() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constant.HeaderKeyHtmxRequest, constant.HeaderValueHtmxRequest)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records")
	c.QueryParams().Add("date", d.String())

	response := view.GetRecordTable{
		Dates: view.DateResult{
			CurrentDate: d,
			PastDate:    []date.Date{},
			FutureDate:  []date.Date{},
		},
		AssetTypes: []view.RecordTableAssetType{
			{
				Id:     1,
				Name:   "Test Asset Type",
				IsCash: true,
				Records: []view.RecordTableData{
					{
						Id:               1,
						AssetId:          1,
						Name:             "test",
						Broker:           "test",
						BoughtValue:      decimal.NewFromInt(1000),
						CurrentValue:     decimal.NewFromInt(1500),
						RealizedValue:    decimal.NewFromInt(100),
						DefaultIncrement: decimal.Zero,
						IsCash:           false,
						Note:             "",
					},
				},
			},
		},
	}

	s.mockService.On("GetByDate", d).Return(response, nil)

	// Execute
	err := s.handler.GetRecordsByDate(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	// Should only render the table component, not the full page
	s.NotContains(rec.Body.String(), "<!doctype html>")
}

func (s *RecordHandlerTestSuite) TestGetRecordsForDraft_Error() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records/draft")

	s.mockService.On("GetDraft").Return(view.EditRecordTable{}, errors.New("draft error"))

	// Execute
	err := s.handler.GetRecordsForDraft(c)

	// Verify
	s.Error(err)
	s.Equal("draft error", err.Error())
}

func (s *RecordHandlerTestSuite) TestGetRecordsForDraft_Success() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records/draft")

	response := view.EditRecordTable{
		Date:       date.Today(),
		Action:     constant.EditActionNew,
		AssetTypes: []view.RecordTableAssetType{},
	}

	s.mockService.On("GetDraft").Return(response, nil)

	// Execute
	err := s.handler.GetRecordsForDraft(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), "Add Record")
}

func (s *RecordHandlerTestSuite) TestGetRecordsForEdit_Error_InvalidDateFormat() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records/edit")
	c.QueryParams().Add("date", "invalid-date")

	// Execute
	err := s.handler.GetRecordsForEdit(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "cannot parse input date")
}

func (s *RecordHandlerTestSuite) TestGetRecordsForEdit_Error_ServiceError() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records/edit")
	c.QueryParams().Add("date", d.String())

	s.mockService.On("GetByDateForEdit", d).Return(view.EditRecordTable{}, errors.New("service error"))

	// Execute
	err := s.handler.GetRecordsForEdit(c)

	// Verify
	s.Error(err)
	s.Equal("service error", err.Error())
}

func (s *RecordHandlerTestSuite) TestGetRecordsForEdit_Success() {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records/edit")
	c.QueryParams().Add("date", d.String())

	response := view.EditRecordTable{
		Date:       d,
		Action:     constant.EditActionUpdate,
		AssetTypes: []view.RecordTableAssetType{},
	}

	s.mockService.On("GetByDateForEdit", d).Return(response, nil)

	// Execute
	err := s.handler.GetRecordsForEdit(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), "Edit Record")
}

func (s *RecordHandlerTestSuite) TestCreateRecord_Error_InvalidRequestBody() {
	// Setup
	requestBody := "invalid-json"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	// Execute
	err := s.handler.CreateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "cannot bind request body")
}

func (s *RecordHandlerTestSuite) TestCreateRecord_Error_ValidationFailed() {
	// Setup
	requestBody := `{"records":[{"assetId":"","boughtValue":"1000","currentValue":"1500","realizedValue":"500"}]}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	// Execute
	err := s.handler.CreateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "fail while validate request body")
}

func (s *RecordHandlerTestSuite) TestCreateRecord_Error_ServiceError() {
	// Setup
	d := date.MustParse(time.DateOnly, "2023-01-01")
	requestBody := `{"date":"2023-01-01","records":[{"assetId":"1","boughtValue":"1000","currentValue":"1500","realizedValue":"500"}]}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	expectedRequest := model.CreateRecordRequest{
		Date: d,
		Records: []model.CreateRecordRequestRecord{
			{
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "",
			},
		},
	}

	s.mockService.On("CreateRecords", mock.MatchedBy(func(req model.CreateRecordRequest) bool {
		return assert.Equal(s.T(), expectedRequest.Date, req.Date) &&
			assert.Equal(s.T(), len(expectedRequest.Records), len(req.Records)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].AssetId, req.Records[0].AssetId) &&
			assert.True(s.T(), expectedRequest.Records[0].BoughtValue.Equal(req.Records[0].BoughtValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].CurrentValue.Equal(req.Records[0].CurrentValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].RealizedValue.Equal(req.Records[0].RealizedValue)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Note, req.Records[0].Note)
	})).Return(errors.New("service error"))

	// Execute
	err := s.handler.CreateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusInternalServerError, httpErr.Code)
	s.Equal("service error", httpErr.Message.(error).Error())
}

func (s *RecordHandlerTestSuite) TestCreateRecord_Success() {
	// Setup
	d := date.MustParse(time.DateOnly, "2023-01-01")
	requestBody := `{"date":"2023-01-01","records":[{"assetId":"1","boughtValue":"1000","currentValue":"1500","realizedValue":"500","note":"Test note"}]}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	expectedRequest := model.CreateRecordRequest{
		Date: d,
		Records: []model.CreateRecordRequestRecord{
			{
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	s.mockService.On("CreateRecords", mock.MatchedBy(func(req model.CreateRecordRequest) bool {
		return assert.Equal(s.T(), expectedRequest.Date, req.Date) &&
			assert.Equal(s.T(), len(expectedRequest.Records), len(req.Records)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].AssetId, req.Records[0].AssetId) &&
			assert.True(s.T(), expectedRequest.Records[0].BoughtValue.Equal(req.Records[0].BoughtValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].CurrentValue.Equal(req.Records[0].CurrentValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].RealizedValue.Equal(req.Records[0].RealizedValue)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Note, req.Records[0].Note)
	})).Return(nil)

	// Execute
	err := s.handler.CreateRecord(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("/", rec.Header().Get("HX-Redirect"))
}

func (s *RecordHandlerTestSuite) TestUpdateRecord_Error_InvalidRequestBody() {
	// Setup
	requestBody := "invalid-json"
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	// Execute
	err := s.handler.UpdateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "cannot bind request body")
}

func (s *RecordHandlerTestSuite) TestUpdateRecord_Error_ValidationFailed() {
	// Setup
	requestBody := `{"records":[{"id":"1","assetId":"","boughtValue":"1000","currentValue":"1500","realizedValue":"500"}]}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	// Execute
	err := s.handler.UpdateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "fail while validate request body")
}

func (s *RecordHandlerTestSuite) TestUpdateRecord_Error_ServiceError() {
	// Setup
	d := date.MustParse(time.DateOnly, "2023-01-01")
	requestBody := `{"date":"2023-01-01","records":[{"id":"1","assetId":"1","boughtValue":"1000","currentValue":"1500","realizedValue":"500"}]}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	expectedRequest := model.UpdateRecordRequest{
		Date: d,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "1",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "",
			},
		},
	}

	s.mockService.On("UpdateRecords", mock.MatchedBy(func(req model.UpdateRecordRequest) bool {
		return assert.Equal(s.T(), expectedRequest.Date, req.Date) &&
			assert.Equal(s.T(), len(expectedRequest.Records), len(req.Records)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Id, req.Records[0].Id) &&
			assert.Equal(s.T(), expectedRequest.Records[0].AssetId, req.Records[0].AssetId) &&
			assert.True(s.T(), expectedRequest.Records[0].BoughtValue.Equal(req.Records[0].BoughtValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].CurrentValue.Equal(req.Records[0].CurrentValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].RealizedValue.Equal(req.Records[0].RealizedValue)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Note, req.Records[0].Note)
	})).Return(errors.New("service error"))

	// Execute
	err := s.handler.UpdateRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusInternalServerError, httpErr.Code)
	s.Equal("service error", httpErr.Message.(error).Error())
}

func (s *RecordHandlerTestSuite) TestUpdateRecord_Success() {
	// Setup
	d := date.MustParse(time.DateOnly, "2023-01-01")
	requestBody := `{"date":"2023-01-01","records":[{"id":"1","assetId":"1","boughtValue":"1000","currentValue":"1500","realizedValue":"500","note":"Test note"}]}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")

	expectedRequest := model.UpdateRecordRequest{
		Date: d,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "1",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	s.mockService.On("UpdateRecords", mock.MatchedBy(func(req model.UpdateRecordRequest) bool {
		return assert.Equal(s.T(), expectedRequest.Date, req.Date) &&
			assert.Equal(s.T(), len(expectedRequest.Records), len(req.Records)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Id, req.Records[0].Id) &&
			assert.Equal(s.T(), expectedRequest.Records[0].AssetId, req.Records[0].AssetId) &&
			assert.True(s.T(), expectedRequest.Records[0].BoughtValue.Equal(req.Records[0].BoughtValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].CurrentValue.Equal(req.Records[0].CurrentValue)) &&
			assert.True(s.T(), expectedRequest.Records[0].RealizedValue.Equal(req.Records[0].RealizedValue)) &&
			assert.Equal(s.T(), expectedRequest.Records[0].Note, req.Records[0].Note)
	})).Return(nil)

	// Execute
	err := s.handler.UpdateRecord(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("/", rec.Header().Get("HX-Redirect"))
}

func (s *RecordHandlerTestSuite) TestDeleteRecord_Error_InvalidDateFormat() {
	// Setup
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/records")
	c.QueryParams().Add("date", "invalid-date")

	// Execute
	err := s.handler.DeleteRecord(c)

	// Verify
	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	s.Require().True(ok)
	s.Equal(http.StatusBadRequest, httpErr.Code)
	s.Contains(httpErr.Message.(error).Error(), "cannot parse input date")
}

func (s *RecordHandlerTestSuite) TestDeleteRecord_Error_ServiceError() {
	// Setup
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records")
	c.QueryParams().Add("date", d.String())

	s.mockService.On("DeleteRecords", d).Return(errors.New("service error"))

	// Execute
	err := s.handler.DeleteRecord(c)

	// Verify
	s.Error(err)
	s.Equal("service error", err.Error())
}

func (s *RecordHandlerTestSuite) TestDeleteRecord_Success() {
	// Setup
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	d := date.Today()
	c.SetPath("/records")
	c.QueryParams().Add("date", d.String())

	s.mockService.On("DeleteRecords", d).Return(nil)

	// Execute
	err := s.handler.DeleteRecord(c)

	// Verify
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("/", rec.Header().Get("HX-Redirect"))
}
