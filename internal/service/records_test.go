package service

import (
	"database/sql"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/ports"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RecordsTestSuite struct {
	suite.Suite
	mockRepo *ports.MockRecordsRepository
	service  *Records
}

func TestRecordsTestSuite(t *testing.T) {
	suite.Run(t, new(RecordsTestSuite))
}

func (s *RecordsTestSuite) SetupTest() {
	s.mockRepo = ports.NewMockRecordsRepository(s.T())
	s.service = NewRecords(s.mockRepo)
}

func (s *RecordsTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func (s *RecordsTestSuite) TestGetByDate_Error_ZeroValueDate_CannotGetLatestDate() {
	// Setup
	d := date.Zero
	s.mockRepo.EXPECT().GetLatestDate().Return(date.Zero, fmt.Errorf("failed to get latest date"))

	// Execute
	_, err := s.service.GetByDate(d)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get latest date")
}

func (s *RecordsTestSuite) TestGetByDate_Error_FailToGetRecordByDate() {
	// Setup
	d := date.Today()
	s.mockRepo.EXPECT().GetRecordByDate(d).Return(nil, fmt.Errorf("failed to get records"))

	// Execute
	_, err := s.service.GetByDate(d)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get records")
}

func (s *RecordsTestSuite) TestGetByDate_Error_FailToGetDate() {
	// Setup
	d := date.Today()
	var emptyRecords []model.AssetTypeRecord
	s.mockRepo.EXPECT().GetRecordByDate(d).Return(emptyRecords, nil)
	s.mockRepo.EXPECT().GetDate(d).Return(nil, fmt.Errorf("failed to get date list"))

	// Execute
	_, err := s.service.GetByDate(d)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get date list")
}

func (s *RecordsTestSuite) TestGetByDate_Success_EmptyRecords() {
	// Setup
	d := date.Today()
	yesterday := d.AddDate(0, 0, -1)
	twoDaysAgo := d.AddDate(0, 0, -2)
	tomorrow := d.AddDate(0, 0, 1)
	twoDaysLater := d.AddDate(0, 0, 2)

	var emptyRecords []model.AssetTypeRecord
	dateList := &model.DateList{
		Current: d,
		Prev:    []date.Date{yesterday, twoDaysAgo},
		Next:    []date.Date{tomorrow, twoDaysLater},
	}

	s.mockRepo.EXPECT().GetRecordByDate(d).Return(emptyRecords, nil)
	s.mockRepo.EXPECT().GetDate(d).Return(dateList, nil)

	// Execute
	result, err := s.service.GetByDate(d)

	// Verify
	s.NoError(err)
	s.Equal(d, result.Dates.CurrentDate)
	s.Equal([]date.Date{yesterday, twoDaysAgo}, result.Dates.PastDate)
	s.Equal([]date.Date{tomorrow, twoDaysLater}, result.Dates.FutureDate)
	s.Empty(result.AssetTypes)
}

func (s *RecordsTestSuite) TestGetByDate_Success_EmptyDates() {
	// Setup
	d := date.Today()
	assetType := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType}
	dateList := &model.DateList{
		Current: d,
		Prev:    []date.Date{},
		Next:    []date.Date{},
	}

	s.mockRepo.EXPECT().GetRecordByDate(d).Return(records, nil)
	s.mockRepo.EXPECT().GetDate(d).Return(dateList, nil)

	// Execute
	result, err := s.service.GetByDate(d)

	// Verify
	s.NoError(err)
	s.Equal(d, result.Dates.CurrentDate)
	s.Empty(result.Dates.PastDate)
	s.Empty(result.Dates.FutureDate)
	s.Len(result.AssetTypes, 1)
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetByDate_Success_HasRecordsAndDates() {
	// Setup
	d := date.Today()
	yesterday := d.AddDate(0, 0, -1)
	tomorrow := d.AddDate(0, 0, 1)

	assetType1 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	assetType2 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 2, Valid: true},
		Name:        "Bonds",
		IsCash:      false,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 201, Valid: true},
				AssetId:          2,
				Name:             "Bond1",
				Broker:           "Broker2",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(10000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(10500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{Valid: false},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType1, assetType2}
	dateList := &model.DateList{
		Current: d,
		Prev:    []date.Date{yesterday},
		Next:    []date.Date{tomorrow},
	}

	s.mockRepo.EXPECT().GetRecordByDate(d).Return(records, nil)
	s.mockRepo.EXPECT().GetDate(d).Return(dateList, nil)

	// Execute
	result, err := s.service.GetByDate(d)

	// Verify
	s.NoError(err)
	s.Equal(d, result.Dates.CurrentDate)
	s.Equal([]date.Date{yesterday}, result.Dates.PastDate)
	s.Equal([]date.Date{tomorrow}, result.Dates.FutureDate)
	s.Len(result.AssetTypes, 2)

	// Verify first asset type
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)

	// Verify second asset type
	s.Equal(2, result.AssetTypes[1].Id)
	s.Equal("Bonds", result.AssetTypes[1].Name)
	s.Len(result.AssetTypes[1].Records, 1)
	s.Equal(201, result.AssetTypes[1].Records[0].Id)
	s.Equal(2, result.AssetTypes[1].Records[0].AssetId)
	s.Equal("Bond1", result.AssetTypes[1].Records[0].Name)
	s.Equal("Broker2", result.AssetTypes[1].Records[0].Broker)
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[1].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(10000).Equal(result.AssetTypes[1].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(10500).Equal(result.AssetTypes[1].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[1].Records[0].RealizedValue))
	s.Empty(result.AssetTypes[1].Records[0].Note)
	s.False(result.AssetTypes[1].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetByDate_Success_ZeroValueDate_GetWithLatestDate() {
	// Setup
	d := date.Zero
	latestDate := date.Today()
	yesterday := latestDate.AddDate(0, 0, -1)
	tomorrow := latestDate.AddDate(0, 0, 1)

	assetType := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType}
	dateList := &model.DateList{
		Current: latestDate,
		Prev:    []date.Date{yesterday},
		Next:    []date.Date{tomorrow},
	}

	s.mockRepo.EXPECT().GetLatestDate().Return(latestDate, nil)
	s.mockRepo.EXPECT().GetRecordByDate(latestDate).Return(records, nil)
	s.mockRepo.EXPECT().GetDate(latestDate).Return(dateList, nil)

	// Execute
	result, err := s.service.GetByDate(d)

	// Verify
	s.NoError(err)
	s.Equal(latestDate, result.Dates.CurrentDate)
	s.Equal([]date.Date{yesterday}, result.Dates.PastDate)
	s.Equal([]date.Date{tomorrow}, result.Dates.FutureDate)
	s.Len(result.AssetTypes, 1)
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetByDateForEdit_Error_ZeroValueDate_CannotGetLatestDate() {
	// Setup
	d := date.Zero
	s.mockRepo.EXPECT().GetLatestDate().Return(date.Zero, fmt.Errorf("failed to get latest date"))

	// Execute
	_, err := s.service.GetByDateForEdit(d)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get latest date")
}

func (s *RecordsTestSuite) TestGetByDateForEdit_Error_FailToGetRecordByDate() {
	// Setup
	d := date.Today()
	s.mockRepo.EXPECT().GetRecordByDate(d).Return(nil, fmt.Errorf("failed to get records"))

	// Execute
	_, err := s.service.GetByDateForEdit(d)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get records")
}

func (s *RecordsTestSuite) TestGetByDateForEdit_Success_WithProvidedDate() {
	// Setup
	d := date.Today()
	assetType := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType}
	s.mockRepo.EXPECT().GetRecordByDate(d).Return(records, nil)

	// Execute
	result, err := s.service.GetByDateForEdit(d)

	// Verify
	s.NoError(err)
	s.Equal(d, result.Date)
	s.Equal(constant.EditActionUpdate, result.Action)
	s.Len(result.AssetTypes, 1)
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetByDateForEdit_Success_ZeroValueDate_GetWithLatestDate() {
	// Setup
	d := date.Zero
	latestDate := date.Today()
	assetType := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType}
	s.mockRepo.EXPECT().GetLatestDate().Return(latestDate, nil)
	s.mockRepo.EXPECT().GetRecordByDate(latestDate).Return(records, nil)

	// Execute
	result, err := s.service.GetByDateForEdit(d)

	// Verify
	s.NoError(err)
	s.Equal(latestDate, result.Date)
	s.Equal(constant.EditActionUpdate, result.Action)
	s.Len(result.AssetTypes, 1)
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetByDateForEdit_Success_MultipleAssetTypes() {
	// Setup
	d := date.Today()
	assetType1 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	assetType2 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 2, Valid: true},
		Name:        "Bonds",
		IsCash:      false,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 201, Valid: true},
				AssetId:          2,
				Name:             "Bond1",
				Broker:           "Broker2",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(10000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(10500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{Valid: false},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType1, assetType2}
	s.mockRepo.On("GetRecordByDate", d).Return(records, nil)

	// Execute
	result, err := s.service.GetByDateForEdit(d)

	// Verify
	s.NoError(err)
	s.Equal(d, result.Date)
	s.Equal(constant.EditActionUpdate, result.Action)
	s.Len(result.AssetTypes, 2)

	// Verify first asset type
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)

	// Verify second asset type
	s.Equal(2, result.AssetTypes[1].Id)
	s.Equal("Bonds", result.AssetTypes[1].Name)
	s.Len(result.AssetTypes[1].Records, 1)
	s.Equal(201, result.AssetTypes[1].Records[0].Id)
	s.Equal(2, result.AssetTypes[1].Records[0].AssetId)
	s.Equal("Bond1", result.AssetTypes[1].Records[0].Name)
	s.Equal("Broker2", result.AssetTypes[1].Records[0].Broker)
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[1].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(10000).Equal(result.AssetTypes[1].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(10500).Equal(result.AssetTypes[1].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[1].Records[0].RealizedValue))
	s.Empty(result.AssetTypes[1].Records[0].Note)
	s.False(result.AssetTypes[1].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetDraft_Error_FailToGetRecordDraft() {
	// Setup
	s.mockRepo.EXPECT().GetRecordDraft().Return(nil, fmt.Errorf("failed to get record draft"))

	// Execute
	_, err := s.service.GetDraft()

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot get record draft")
}

func (s *RecordsTestSuite) TestGetDraft_Success_EmptyRecords() {
	// Setup
	emptyRecords := []model.AssetTypeRecord{}
	s.mockRepo.EXPECT().GetRecordDraft().Return(emptyRecords, nil)

	// Execute
	result, err := s.service.GetDraft()

	// Verify
	s.NoError(err)
	s.Equal(date.Today(), result.Date)
	s.Equal(constant.EditActionNew, result.Action)
	s.Empty(result.AssetTypes)
}

func (s *RecordsTestSuite) TestGetDraft_Success_WithRecords() {
	// Setup
	assetType := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType}
	s.mockRepo.EXPECT().GetRecordDraft().Return(records, nil)

	// Execute
	result, err := s.service.GetDraft()

	// Verify
	s.NoError(err)
	s.Equal(date.Today(), result.Date)
	s.Equal(constant.EditActionNew, result.Action)
	s.Len(result.AssetTypes, 1)
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestGetDraft_Success_MultipleAssetTypes() {
	// Setup
	assetType1 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 1, Valid: true},
		Name:        "Stocks",
		IsCash:      true,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 101, Valid: true},
				AssetId:          1,
				Name:             "AAPL",
				Broker:           "Broker1",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(100)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(1500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{String: "Test note", Valid: true},
			},
		},
	}

	assetType2 := model.AssetTypeRecord{
		Id:          sql.NullInt64{Int64: 2, Valid: true},
		Name:        "Bonds",
		IsCash:      false,
		IsLiability: false,
		Assets: []model.AssetRecord{
			{
				Id:               sql.NullInt64{Int64: 201, Valid: true},
				AssetId:          2,
				Name:             "Bond1",
				Broker:           "Broker2",
				DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(1000)),
				BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(10000)),
				CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(10500)),
				RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(500)),
				Note:             sql.NullString{Valid: false},
			},
		},
	}

	records := []model.AssetTypeRecord{assetType1, assetType2}
	s.mockRepo.EXPECT().GetRecordDraft().Return(records, nil)

	// Execute
	result, err := s.service.GetDraft()

	// Verify
	s.NoError(err)
	s.Equal(date.Today(), result.Date)
	s.Equal(constant.EditActionNew, result.Action)
	s.Len(result.AssetTypes, 2)

	// Verify first asset type
	s.Equal(1, result.AssetTypes[0].Id)
	s.Equal("Stocks", result.AssetTypes[0].Name)
	s.Len(result.AssetTypes[0].Records, 1)
	s.Equal(101, result.AssetTypes[0].Records[0].Id)
	s.Equal(1, result.AssetTypes[0].Records[0].AssetId)
	s.Equal("AAPL", result.AssetTypes[0].Records[0].Name)
	s.Equal("Broker1", result.AssetTypes[0].Records[0].Broker)
	s.True(decimal.NewFromInt(100).Equal(result.AssetTypes[0].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[0].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(1500).Equal(result.AssetTypes[0].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[0].Records[0].RealizedValue))
	s.Equal("Test note", result.AssetTypes[0].Records[0].Note)
	s.True(result.AssetTypes[0].Records[0].IsCash)

	// Verify second asset type
	s.Equal(2, result.AssetTypes[1].Id)
	s.Equal("Bonds", result.AssetTypes[1].Name)
	s.Len(result.AssetTypes[1].Records, 1)
	s.Equal(201, result.AssetTypes[1].Records[0].Id)
	s.Equal(2, result.AssetTypes[1].Records[0].AssetId)
	s.Equal("Bond1", result.AssetTypes[1].Records[0].Name)
	s.Equal("Broker2", result.AssetTypes[1].Records[0].Broker)
	s.True(decimal.NewFromInt(1000).Equal(result.AssetTypes[1].Records[0].DefaultIncrement))
	s.True(decimal.NewFromInt(10000).Equal(result.AssetTypes[1].Records[0].BoughtValue))
	s.True(decimal.NewFromInt(10500).Equal(result.AssetTypes[1].Records[0].CurrentValue))
	s.True(decimal.NewFromInt(500).Equal(result.AssetTypes[1].Records[0].RealizedValue))
	s.Empty(result.AssetTypes[1].Records[0].Note)
	s.False(result.AssetTypes[1].Records[0].IsCash)
}

func (s *RecordsTestSuite) TestCreateRecords_Error_BeginTxFails() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
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

	s.mockRepo.EXPECT().BeginTx().Return(nil, nil, fmt.Errorf("database connection error"))

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.Error(err)
	s.Equal("database connection error", err.Error())
}

func (s *RecordsTestSuite) TestCreateRecords_Error_UpsertRecordFails() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
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

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(fmt.Errorf("database error"))

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot insert record")
}

func (s *RecordsTestSuite) TestCreateRecords_Error_CommitFails() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
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

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(fmt.Errorf("commit failed"))

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot commit transaction")
}

func (s *RecordsTestSuite) TestCreateRecords_Success_SingleRecord() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
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

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestCreateRecords_Success_MultipleRecords() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
		Records: []model.CreateRecordRequestRecord{
			{
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note 1",
			},
			{
				AssetId:       "2",
				BoughtValue:   decimal.NewFromInt(2000),
				CurrentValue:  decimal.NewFromInt(2500),
				RealizedValue: decimal.NewFromInt(1000),
				Note:          "Test note 2",
			},
			{
				AssetId:       "3",
				BoughtValue:   decimal.NewFromInt(3000),
				CurrentValue:  decimal.NewFromInt(3500),
				RealizedValue: decimal.NewFromInt(1500),
				Note:          "", // Empty note should create a null string
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)

	expectedRecord1 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note 1", Valid: true},
	}

	expectedRecord2 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       2,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(2000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(2500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		Note:          sql.NullString{String: "Test note 2", Valid: true},
	}

	expectedRecord3 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       3,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(3000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(3500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		Note:          sql.NullString{Valid: false},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord1, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord2, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord3, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestCreateRecords_Success_InvalidAssetIdIgnored() {
	// Setup
	testDate := date.Today()
	req := model.CreateRecordRequest{
		Date: testDate,
		Records: []model.CreateRecordRequestRecord{
			{
				AssetId:       "invalid", // This should be ignored as validation should happen before reaching this method
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)

	// We expect 0 since the strconv.Atoi will fail on "invalid" and return 0
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       0,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestCreateRecords_Success_LargeDataSet() {
	// Setup
	testDate := date.Today()

	// Create a large data set of records
	var records []model.CreateRecordRequestRecord
	for i := 1; i <= 50; i++ {
		records = append(records, model.CreateRecordRequestRecord{
			AssetId:       fmt.Sprintf("%d", i),
			BoughtValue:   decimal.NewFromInt(int64(i * 100)),
			CurrentValue:  decimal.NewFromInt(int64(i * 150)),
			RealizedValue: decimal.NewFromInt(int64(i * 50)),
			Note:          fmt.Sprintf("Asset %d note", i),
		})
	}

	req := model.CreateRecordRequest{
		Date:    testDate,
		Records: records,
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)

	// Setup expectations for each record
	for i := 1; i <= 50; i++ {
		expectedRecord := model.AssetRecord{
			Id:            sql.NullInt64{Valid: false},
			AssetId:       i,
			BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(int64(i * 100))),
			CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(int64(i * 150))),
			RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(int64(i * 50))),
			Note:          sql.NullString{String: fmt.Sprintf("Asset %d note", i), Valid: true},
		}
		s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	}

	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestCreateRecords_Success_SpecialDecimalValues() {
	// Setup
	testDate := date.Today()

	// Create records with special decimal values
	req := model.CreateRecordRequest{
		Date: testDate,
		Records: []model.CreateRecordRequestRecord{
			{
				AssetId:       "1",
				BoughtValue:   decimal.NewFromFloat(123.45), // Test decimal places
				CurrentValue:  decimal.NewFromFloat(678.90),
				RealizedValue: decimal.NewFromFloat(555.45),
				Note:          "Decimal places test",
			},
			{
				AssetId:       "2",
				BoughtValue:   decimal.NewFromFloat(0.0001), // Test very small value
				CurrentValue:  decimal.NewFromFloat(0.0002),
				RealizedValue: decimal.NewFromFloat(0.0001),
				Note:          "Small decimal test",
			},
			{
				AssetId:       "3",
				BoughtValue:   decimal.NewFromFloat(9999999.99), // Test large value
				CurrentValue:  decimal.NewFromFloat(9999999.99),
				RealizedValue: decimal.Zero, // Test zero
				Note:          "Large value test",
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)

	// Define expected records with appropriate decimal conversions
	expectedRecord1 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromFloat(123.45)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromFloat(678.90)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromFloat(555.45)),
		Note:          sql.NullString{String: "Decimal places test", Valid: true},
	}

	expectedRecord2 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       2,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromFloat(0.0001)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromFloat(0.0002)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromFloat(0.0001)),
		Note:          sql.NullString{String: "Small decimal test", Valid: true},
	}

	expectedRecord3 := model.AssetRecord{
		Id:            sql.NullInt64{Valid: false},
		AssetId:       3,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromFloat(9999999.99)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromFloat(9999999.99)),
		RealizedValue: decimal.NewNullDecimal(decimal.Zero),
		Note:          sql.NullString{String: "Large value test", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord1, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord2, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord3, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.CreateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestUpdateRecords_Error_BeginTxFails() {
	// Setup
	testDate := date.Today()
	req := model.UpdateRecordRequest{
		Date: testDate,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "101",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	s.mockRepo.EXPECT().BeginTx().Return(nil, nil, fmt.Errorf("database connection error"))

	// Execute
	err := s.service.UpdateRecords(req)

	// Verify
	s.Error(err)
	s.Equal("database connection error", err.Error())
}

func (s *RecordsTestSuite) TestUpdateRecords_Error_UpsertRecordFails() {
	// Setup
	testDate := date.Today()
	req := model.UpdateRecordRequest{
		Date: testDate,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "101",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 101, Valid: true},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(fmt.Errorf("database error"))

	// Execute
	err := s.service.UpdateRecords(req)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot insert record")
}

func (s *RecordsTestSuite) TestUpdateRecords_Error_CommitFails() {
	// Setup
	testDate := date.Today()
	req := model.UpdateRecordRequest{
		Date: testDate,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "101",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 101, Valid: true},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(fmt.Errorf("commit failed"))

	// Execute
	err := s.service.UpdateRecords(req)

	// Verify
	s.Error(err)
	s.Contains(err.Error(), "cannot commit transaction")
}

func (s *RecordsTestSuite) TestUpdateRecords_Success_SingleRecord() {
	// Setup
	testDate := date.Today()
	req := model.UpdateRecordRequest{
		Date: testDate,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "101",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note",
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	expectedRecord := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 101, Valid: true},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note", Valid: true},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.UpdateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestUpdateRecords_Success_MultipleRecords() {
	// Setup
	testDate := date.Today()
	req := model.UpdateRecordRequest{
		Date: testDate,
		Records: []model.UpdateRecordRequestRecord{
			{
				Id:            "101",
				AssetId:       "1",
				BoughtValue:   decimal.NewFromInt(1000),
				CurrentValue:  decimal.NewFromInt(1500),
				RealizedValue: decimal.NewFromInt(500),
				Note:          "Test note 1",
			},
			{
				Id:            "102",
				AssetId:       "2",
				BoughtValue:   decimal.NewFromInt(2000),
				CurrentValue:  decimal.NewFromInt(2500),
				RealizedValue: decimal.NewFromInt(1000),
				Note:          "Test note 2",
			},
			{
				Id:            "103",
				AssetId:       "3",
				BoughtValue:   decimal.NewFromInt(3000),
				CurrentValue:  decimal.NewFromInt(3500),
				RealizedValue: decimal.NewFromInt(1500),
				Note:          "", // Empty note should create a null string
			},
		},
	}

	mockTx := ports.NewMockTX(s.T())

	s.mockRepo.EXPECT().BeginTx().Return(s.mockRepo, mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)

	expectedRecord1 := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 101, Valid: true},
		AssetId:       1,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(500)),
		Note:          sql.NullString{String: "Test note 1", Valid: true},
	}

	expectedRecord2 := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 102, Valid: true},
		AssetId:       2,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(2000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(2500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(1000)),
		Note:          sql.NullString{String: "Test note 2", Valid: true},
	}

	expectedRecord3 := model.AssetRecord{
		Id:            sql.NullInt64{Int64: 103, Valid: true},
		AssetId:       3,
		BoughtValue:   decimal.NewNullDecimal(decimal.NewFromInt(3000)),
		CurrentValue:  decimal.NewNullDecimal(decimal.NewFromInt(3500)),
		RealizedValue: decimal.NewNullDecimal(decimal.NewFromInt(1500)),
		Note:          sql.NullString{Valid: false},
	}

	s.mockRepo.EXPECT().UpsertRecord(expectedRecord1, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord2, testDate).Return(nil)
	s.mockRepo.EXPECT().UpsertRecord(expectedRecord3, testDate).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Execute
	err := s.service.UpdateRecords(req)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestDeleteRecords_Error() {
	// Setup
	testDate := date.Today()
	s.mockRepo.EXPECT().DeleteRecordByDate(testDate).Return(int64(0), fmt.Errorf("failed to delete records"))

	// Execute
	err := s.service.DeleteRecords(testDate)

	// Verify
	s.Error(err)
	s.Equal("failed to delete records", err.Error())
}

func (s *RecordsTestSuite) TestDeleteRecords_Success_NoRecordsDeleted() {
	// Setup
	testDate := date.Today()
	s.mockRepo.EXPECT().DeleteRecordByDate(testDate).Return(int64(0), nil)

	// Execute
	err := s.service.DeleteRecords(testDate)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestDeleteRecords_Success_SingleRecordDeleted() {
	// Setup
	testDate := date.Today()
	s.mockRepo.EXPECT().DeleteRecordByDate(testDate).Return(int64(1), nil)

	// Execute
	err := s.service.DeleteRecords(testDate)

	// Verify
	s.NoError(err)
}

func (s *RecordsTestSuite) TestDeleteRecords_Success_MultipleRecordsDeleted() {
	// Setup
	testDate := date.Today()
	s.mockRepo.EXPECT().DeleteRecordByDate(testDate).Return(int64(5), nil)

	// Execute
	err := s.service.DeleteRecords(testDate)

	// Verify
	s.NoError(err)
}
