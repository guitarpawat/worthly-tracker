//go:build test || integration

package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/openlyinc/pointy"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"worthly-tracker/config"
	"worthly-tracker/logs"
	"worthly-tracker/model"
)

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}

type RecordSuite struct {
	suite.Suite
	tx   *sqlx.Tx
	repo RecordRepo
}

func (s *RecordSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
	Init()
	s.repo = GetRecordRepo()
}

func (s *RecordSuite) SetupTest() {
	var err error
	s.tx, err = GetDB().BeginTx()
	s.Require().NoError(err)
}

func (s *RecordSuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *RecordSuite) TestGetDate_NoRecord() {
	now := model.Date(time.Now())
	actual, err := s.repo.GetDate(now, s.tx)
	s.Require().NoError(err)

	expect := &model.DateList{
		Current: &now,
		Prev:    make([]model.Date, 0, 12),
		Next:    make([]model.Date, 0, 12),
	}

	s.Require().Equal(expect, actual)
}

func (s *RecordSuite) TestGetDate_WithRecord() {
	mockRecords(s.tx)
	cases := []struct {
		current  model.Date
		prevFrom model.Date
		prevTo   model.Date
		nextFrom model.Date
		nextTo   model.Date
	}{
		{
			current:  model.MustNewDate("2022-12-31"),
			nextFrom: model.MustNewDate("2023-01-01"),
			nextTo:   model.MustNewDate("2023-01-12"),
		},
		{
			current:  model.MustNewDate("2023-01-01"),
			nextFrom: model.MustNewDate("2023-01-02"),
			nextTo:   model.MustNewDate("2023-01-13"),
		},
		{
			current:  model.MustNewDate("2023-01-04"),
			prevFrom: model.MustNewDate("2023-01-01"),
			prevTo:   model.MustNewDate("2023-01-03"),
			nextFrom: model.MustNewDate("2023-01-05"),
			nextTo:   model.MustNewDate("2023-01-16"),
		},
		{
			current:  model.MustNewDate("2023-01-15"),
			prevFrom: model.MustNewDate("2023-01-03"),
			prevTo:   model.MustNewDate("2023-01-14"),
			nextFrom: model.MustNewDate("2023-01-16"),
			nextTo:   model.MustNewDate("2023-01-27"),
		},
		{
			current:  model.MustNewDate("2023-01-28"),
			prevFrom: model.MustNewDate("2023-01-16"),
			prevTo:   model.MustNewDate("2023-01-27"),
			nextFrom: model.MustNewDate("2023-01-29"),
			nextTo:   model.MustNewDate("2023-01-31"),
		},
		{
			current:  model.MustNewDate("2023-01-31"),
			prevFrom: model.MustNewDate("2023-01-19"),
			prevTo:   model.MustNewDate("2023-01-30"),
		},
		{
			current:  model.MustNewDate("2023-02-01"),
			prevFrom: model.MustNewDate("2023-01-20"),
			prevTo:   model.MustNewDate("2023-01-31"),
		},
	}

	for _, c := range cases {
		s.Run(c.current.String(), func() {
			actual, err := s.repo.GetDate(c.current, s.tx)
			s.Require().NoError(err)
			s.Require().Equal(c.current, *actual.Current)

			if !time.Time(c.prevFrom).IsZero() && !time.Time(c.prevTo).IsZero() {
				cur := c.prevTo
				ptr := 0
				for time.Time(cur).After(time.Time(c.prevFrom)) || time.Time(cur).Equal(time.Time(c.prevFrom)) {
					s.Require().Equal(cur, actual.Prev[ptr], cur.String()+"!="+actual.Prev[ptr].String())
					ptr++
					cur = model.Date(time.Time(cur).AddDate(0, 0, -1))
				}
				s.Require().Equal(ptr, len(actual.Prev))
			}

			if !time.Time(c.nextFrom).IsZero() && !time.Time(c.nextTo).IsZero() {
				cur := c.nextFrom
				ptr := 0
				for time.Time(cur).Before(time.Time(c.nextTo)) || time.Time(cur).Equal(time.Time(c.nextTo)) {
					s.Require().Equal(cur, actual.Next[ptr], cur.String()+"!="+actual.Next[ptr].String())
					ptr++
					cur = model.Date(time.Time(cur).AddDate(0, 0, 1))
				}
				s.Require().Equal(ptr, len(actual.Next))
			}
		})
	}
}

func (s *RecordSuite) TestGetLatestDate_NoRecord() {
	date, err := s.repo.GetLatestDate(s.tx)
	s.Require().NoError(err)
	s.Require().Nil(date)
}

func (s *RecordSuite) TestGetLatestDate_WithRecord() {
	s.Require().NoError(mockRecords(s.tx))
	actual, err := s.repo.GetLatestDate(s.tx)
	s.Require().NoError(err)

	expect, err := model.NewDate("2023-01-31")
	s.Require().NoError(err)
	s.Require().Equal(expect, *actual)
}

func (s *RecordSuite) TestUpsertRecord_Insert() {
	s.Require().NoError(mockRecords(s.tx))
	record := model.AssetRecord{
		Id:               nil,
		AssetId:          pointy.Int(2),
		Name:             pointy.String("TFFIF"),
		Broker:           pointy.String("SCBS"),
		DefaultIncrement: nil,
		BoughtValue:      pointy.Pointer(decimal.NewFromInt(50)),
		CurrentValue:     pointy.Pointer(decimal.NewFromInt(60)),
		RealizedValue:    pointy.Pointer(decimal.NewFromInt(10)),
		Note:             nil,
	}

	date := model.MustNewDate("2023-08-23")
	s.Require().NoError(s.repo.UpsertRecord(record, date, s.tx))

	var id int
	s.Require().NoError(s.tx.Get(&id, "SELECT id FROM records WHERE date = '2023-08-23'"))
	s.Require().Equal(63, id)
}

func (s *RecordSuite) TestUpsertRecord_Update() {
	s.Require().NoError(mockRecords(s.tx))
	record := model.AssetRecord{
		Id:               pointy.Int(1),
		AssetId:          pointy.Int(2),
		Name:             pointy.String("TFFIF"),
		Broker:           pointy.String("SCBS"),
		DefaultIncrement: nil,
		BoughtValue:      pointy.Pointer(decimal.NewFromInt(50)),
		CurrentValue:     pointy.Pointer(decimal.NewFromInt(60)),
		RealizedValue:    pointy.Pointer(decimal.NewFromInt(10)),
		Note:             nil,
	}

	date := model.MustNewDate("2023-08-23")
	s.Require().NoError(s.repo.UpsertRecord(record, date, s.tx))

	var id int
	s.Require().NoError(s.tx.Get(&id, "SELECT id FROM records WHERE date = '2023-08-23'"))
	s.Require().Equal(1, id)
}

func (s *RecordSuite) TestDeleteRecordById() {
	s.Require().NoError(mockRecords(s.tx))
	s.Require().NoError(s.repo.DeleteRecordById(5, s.tx))
	row := s.tx.QueryRowx("SELECT * FROM records WHERE id = 5")

	var id *int = nil
	err := row.Scan(&id)
	s.Require().ErrorIs(err, sql.ErrNoRows)
}

func (s *RecordSuite) TestDeleteRecordByDate() {
	s.Require().NoError(mockRecords(s.tx))
	rows, err := s.repo.DeleteRecordByDate(model.MustNewDate("2023-01-02"), s.tx)
	s.Require().NoError(err)
	s.Require().EqualValues(2, rows)

	row := s.tx.QueryRowx("SELECT * FROM records WHERE id = 3")
	var id *int = nil
	err = row.Scan(&id)
	s.Require().ErrorIs(err, sql.ErrNoRows)

	row = s.tx.QueryRowx("SELECT * FROM records WHERE id = 4")
	err = row.Scan(&id)
	s.Require().ErrorIs(err, sql.ErrNoRows)
}

func (s *RecordSuite) TestGetRecordByDate_NotFound() {
	s.Require().NoError(mockRecords(s.tx))
	records, err := s.repo.GetRecordByDate(model.MustNewDate("2023-05-03"), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(records))
}

func (s *RecordSuite) TestGetRecordByDate_Found() {
	s.Require().NoError(mockRecords(s.tx))
	s.Require().NoError(mockRecordsWithMultipleType(s.tx))

	actual, err := s.repo.GetRecordByDate(model.MustNewDate("2023-01-03"), s.tx)
	s.Require().NoError(err)

	expect := []model.AssetTypeRecord{
		{
			Id:          pointy.Int(2),
			Name:        pointy.String("MF"),
			IsCash:      pointy.Bool(false),
			IsLiability: pointy.Bool(false),
			Assets: []model.AssetRecord{
				{
					Id:               pointy.Int(64),
					AssetId:          pointy.Int(4),
					Name:             pointy.String("b2"),
					Broker:           pointy.String("finno"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(164)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(264)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(364)),
					Note:             pointy.String("test"),
				},
				{
					Id:               pointy.Int(63),
					AssetId:          pointy.Int(3),
					Name:             pointy.String("b1"),
					Broker:           pointy.String("finno"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(2000)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(163)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(263)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(363)),
					Note:             pointy.String("test2"),
				},
			},
		},
		{
			Id:          pointy.Int(1),
			Name:        pointy.String("Stocks"),
			IsCash:      pointy.Bool(true),
			IsLiability: pointy.Bool(true),
			Assets: []model.AssetRecord{
				{
					Id:               pointy.Int(5),
					AssetId:          pointy.Int(1),
					Name:             pointy.String("a1"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(105)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(205)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(305)),
					Note:             nil,
				},
				{
					Id:               pointy.Int(6),
					AssetId:          pointy.Int(2),
					Name:             pointy.String("a2"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(106)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(206)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(306)),
					Note:             nil,
				},
			},
		},
	}

	s.Require().Equal(len(expect), len(actual))
	s.Require().Equal(expect[0].Id, actual[0].Id)
	s.Require().Equal(expect[0].Name, actual[0].Name)
	s.Require().Equal(expect[0].IsLiability, actual[0].IsLiability)
	s.Require().Equal(expect[0].IsCash, actual[0].IsCash)
	s.Require().Equal(len(expect[0].Assets), len(actual[0].Assets))

	s.Require().Equal(expect[0].Assets[0].Id, actual[0].Assets[0].Id)
	s.Require().Equal(expect[0].Assets[0].AssetId, actual[0].Assets[0].AssetId)
	s.Require().Equal(expect[0].Assets[0].Name, actual[0].Assets[0].Name)
	s.Require().Equal(expect[0].Assets[0].Broker, actual[0].Assets[0].Broker)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Equal(*actual[0].Assets[0].DefaultIncrement))
	s.Require().True(expect[0].Assets[0].BoughtValue.Equal(*actual[0].Assets[0].BoughtValue))
	s.Require().True(expect[0].Assets[0].CurrentValue.Equal(*actual[0].Assets[0].CurrentValue))
	s.Require().True(expect[0].Assets[0].RealizedValue.Equal(*actual[0].Assets[0].RealizedValue))
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Equal(*actual[0].Assets[1].DefaultIncrement))
	s.Require().True(expect[0].Assets[1].BoughtValue.Equal(*actual[0].Assets[1].BoughtValue))
	s.Require().True(expect[0].Assets[1].CurrentValue.Equal(*actual[0].Assets[1].CurrentValue))
	s.Require().True(expect[0].Assets[1].RealizedValue.Equal(*actual[0].Assets[1].RealizedValue))
	s.Require().Equal(expect[0].Assets[1].Note, actual[0].Assets[1].Note)

	s.Require().Equal(expect[0].Id, actual[0].Id)
	s.Require().Equal(expect[0].Name, actual[0].Name)
	s.Require().Equal(expect[0].IsLiability, actual[0].IsLiability)
	s.Require().Equal(expect[0].IsCash, actual[0].IsCash)
	s.Require().Equal(len(expect[1].Assets), len(actual[1].Assets))

	s.Require().Equal(expect[1].Assets[0].Id, actual[1].Assets[0].Id)
	s.Require().Equal(expect[1].Assets[0].AssetId, actual[1].Assets[0].AssetId)
	s.Require().Equal(expect[1].Assets[0].Name, actual[1].Assets[0].Name)
	s.Require().Equal(expect[1].Assets[0].Broker, actual[1].Assets[0].Broker)
	s.Require().True(expect[1].Assets[0].DefaultIncrement.Equal(*actual[1].Assets[0].DefaultIncrement))
	s.Require().True(expect[1].Assets[0].BoughtValue.Equal(*actual[1].Assets[0].BoughtValue))
	s.Require().True(expect[1].Assets[0].CurrentValue.Equal(*actual[1].Assets[0].CurrentValue))
	s.Require().True(expect[1].Assets[0].RealizedValue.Equal(*actual[1].Assets[0].RealizedValue))
	s.Require().Equal(expect[1].Assets[0].Note, actual[1].Assets[0].Note)

	s.Require().Equal(expect[1].Assets[1].Id, actual[1].Assets[1].Id)
	s.Require().Equal(expect[1].Assets[1].AssetId, actual[1].Assets[1].AssetId)
	s.Require().Equal(expect[1].Assets[1].Name, actual[1].Assets[1].Name)
	s.Require().Equal(expect[1].Assets[1].Broker, actual[1].Assets[1].Broker)
	s.Require().True(expect[1].Assets[1].DefaultIncrement.Equal(*actual[1].Assets[1].DefaultIncrement))
	s.Require().True(expect[1].Assets[1].BoughtValue.Equal(*actual[1].Assets[1].BoughtValue))
	s.Require().True(expect[1].Assets[1].CurrentValue.Equal(*actual[1].Assets[1].CurrentValue))
	s.Require().True(expect[1].Assets[1].RealizedValue.Equal(*actual[1].Assets[1].RealizedValue))
	s.Require().Equal(expect[1].Assets[1].Note, actual[1].Assets[1].Note)
}

func (s *RecordSuite) TestGetRecordDraft_EmptyRecord() {
	s.Require().NoError(mockAssets(s.tx))
	s.Require().NoError(mockInactiveAssets(s.tx))

	actual, err := s.repo.GetRecordDraft(s.tx)
	s.Require().NoError(err)

	expect := []model.AssetTypeRecord{
		{
			Id:          pointy.Int(1),
			Name:        pointy.String("Stocks"),
			IsCash:      pointy.Bool(true),
			IsLiability: pointy.Bool(true),
			Assets: []model.AssetRecord{
				{
					Id:               nil,
					AssetId:          pointy.Int(1),
					Name:             pointy.String("a1"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      nil,
					CurrentValue:     nil,
					RealizedValue:    nil,
					Note:             nil,
				},
				{
					Id:               nil,
					AssetId:          pointy.Int(2),
					Name:             pointy.String("a2"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      nil,
					CurrentValue:     nil,
					RealizedValue:    nil,
					Note:             nil,
				},
			},
		},
	}

	s.Require().Equal(len(expect), len(actual))
	s.Require().Equal(expect[0].Id, actual[0].Id)
	s.Require().Equal(expect[0].Name, actual[0].Name)
	s.Require().Equal(expect[0].IsLiability, actual[0].IsLiability)
	s.Require().Equal(expect[0].IsCash, actual[0].IsCash)
	s.Require().Equal(len(expect[0].Assets), len(actual[0].Assets))

	s.Require().Equal(expect[0].Assets[0].Id, actual[0].Assets[0].Id)
	s.Require().Equal(expect[0].Assets[0].AssetId, actual[0].Assets[0].AssetId)
	s.Require().Equal(expect[0].Assets[0].Name, actual[0].Assets[0].Name)
	s.Require().Equal(expect[0].Assets[0].Broker, actual[0].Assets[0].Broker)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Equal(*actual[0].Assets[0].DefaultIncrement))
	s.Require().Nil(actual[0].Assets[0].BoughtValue)
	s.Require().Nil(actual[0].Assets[0].CurrentValue)
	s.Require().Nil(actual[0].Assets[0].RealizedValue)
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Equal(*actual[0].Assets[1].DefaultIncrement))
	s.Require().Nil(actual[0].Assets[1].BoughtValue)
	s.Require().Nil(actual[0].Assets[1].CurrentValue)
	s.Require().Nil(actual[0].Assets[1].RealizedValue)
	s.Require().Equal(expect[0].Assets[1].Note, actual[0].Assets[1].Note)
}

func (s *RecordSuite) TestGetRecordDraft_PartialHitRecord() {
	s.Require().NoError(mockAssets(s.tx))
	s.Require().NoError(mockInactiveAssets(s.tx))

	actual, err := s.repo.GetRecordDraft(s.tx)
	s.Require().NoError(err)

	expect := []model.AssetTypeRecord{
		{
			Id:          pointy.Int(1),
			Name:        pointy.String("Stocks"),
			IsCash:      pointy.Bool(true),
			IsLiability: pointy.Bool(true),
			Assets: []model.AssetRecord{
				{
					Id:               nil,
					AssetId:          pointy.Int(1),
					Name:             pointy.String("a1"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(101)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(201)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(301)),
					Note:             nil,
				},
				{
					Id:               nil,
					AssetId:          pointy.Int(2),
					Name:             pointy.String("a2"),
					Broker:           pointy.String("scbs"),
					DefaultIncrement: pointy.Pointer(decimal.NewFromInt(0)),
					BoughtValue:      pointy.Pointer(decimal.NewFromInt(102)),
					CurrentValue:     pointy.Pointer(decimal.NewFromInt(202)),
					RealizedValue:    pointy.Pointer(decimal.NewFromInt(302)),
					Note:             nil,
				},
			},
		},
	}

	s.Require().Equal(len(expect), len(actual))
	s.Require().Equal(expect[0].Id, actual[0].Id)
	s.Require().Equal(expect[0].Name, actual[0].Name)
	s.Require().Equal(expect[0].IsLiability, actual[0].IsLiability)
	s.Require().Equal(expect[0].IsCash, actual[0].IsCash)
	s.Require().Equal(len(expect[0].Assets), len(actual[0].Assets))

	s.Require().Equal(expect[0].Assets[0].Id, actual[0].Assets[0].Id)
	s.Require().Equal(expect[0].Assets[0].AssetId, actual[0].Assets[0].AssetId)
	s.Require().Equal(expect[0].Assets[0].Name, actual[0].Assets[0].Name)
	s.Require().Equal(expect[0].Assets[0].Broker, actual[0].Assets[0].Broker)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Equal(*actual[0].Assets[0].DefaultIncrement))
	s.Require().Nil(actual[0].Assets[0].BoughtValue)
	s.Require().Nil(actual[0].Assets[0].CurrentValue)
	s.Require().Nil(actual[0].Assets[0].RealizedValue)
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Equal(*actual[0].Assets[1].DefaultIncrement))
	s.Require().Nil(actual[0].Assets[1].BoughtValue)
	s.Require().Nil(actual[0].Assets[1].CurrentValue)
	s.Require().Nil(actual[0].Assets[1].RealizedValue)
	s.Require().Equal(expect[0].Assets[1].Note, actual[0].Assets[1].Note)
}

func mockAssets(tx *sqlx.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Stocks', true, true, 1, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a1', 'scbs', 1, 0, 1, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a2', 'scbs', 1, 0, 2, true)")

	return
}

func mockInactiveAssets(tx *sqlx.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Bonds', true, true, 1, false)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b1', 'scbs', 2, 0, 1, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b2', 'scbs', 2, 0, 1, true)")

	tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Cash', true, true, 1, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c1', 'scbs', 3, 0, 1, false)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c2', 'scbs', 3, 0, 1, false)")

	return
}

func mockPartialRecords(tx *sqlx.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	err = mockAssets(tx)
	if err != nil {
		return
	}

	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-01', 101, 201, 301, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-01', 102, 202, 302, null)")

	return
}

func mockRecords(tx *sqlx.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	err = mockAssets(tx)
	if err != nil {
		return
	}

	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-01', 101, 201, 301, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-01', 102, 202, 302, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-02', 103, 203, 303, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-02', 104, 204, 304, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-03', 105, 205, 305, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-03', 106, 206, 306, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-04', 107, 207, 307, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-04', 108, 208, 308, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-05', 109, 209, 309, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-05', 110, 210, 310, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-06', 111, 211, 311, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-06', 112, 212, 312, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-07', 113, 213, 313, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-07', 114, 214, 314, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-08', 115, 215, 315, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-08', 116, 216, 316, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-09', 117, 217, 317, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-09', 118, 218, 318, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-10', 119, 219, 319, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-10', 120, 220, 320, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-11', 121, 221, 321, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-11', 122, 222, 322, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-12', 123, 223, 323, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-12', 124, 224, 324, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-13', 125, 225, 325, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-13', 126, 226, 326, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-14', 127, 227, 327, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-14', 128, 228, 328, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-15', 129, 229, 329, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-15', 130, 230, 330, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-16', 131, 231, 331, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-16', 132, 232, 332, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-17', 133, 233, 333, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-17', 134, 234, 334, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-18', 135, 235, 335, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-18', 136, 236, 336, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-19', 137, 237, 337, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-19', 138, 238, 338, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-20', 139, 239, 339, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-20', 140, 240, 340, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-21', 141, 241, 341, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-21', 142, 242, 342, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-22', 143, 243, 343, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-22', 144, 244, 344, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-23', 145, 245, 345, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-23', 146, 246, 346, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-24', 147, 247, 347, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-24', 148, 248, 348, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-25', 149, 249, 349, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-25', 150, 250, 350, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-26', 151, 251, 351, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-26', 152, 252, 352, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-27', 153, 253, 353, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-27', 154, 254, 354, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-28', 155, 255, 355, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-28', 156, 256, 356, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-29', 157, 257, 357, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-29', 158, 258, 358, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-30', 159, 259, 359, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-30', 160, 260, 360, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-31', 161, 261, 361, null)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-31', 162, 262, 362, null)")

	return
}

func mockRecordsWithMultipleType(tx *sqlx.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('MF', false, false, 0, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b1', 'finno', 2, 2000, 2, true)")
	tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b2', 'finno', 2, 0, 1, true)")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (3, '2023-01-03', 163, 263, 363, 'test2')")
	tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (4, '2023-01-03', 164, 264, 364, 'test')")
	return
}
