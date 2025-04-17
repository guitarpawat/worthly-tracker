package db

import (
	"database/sql"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}

type RecordSuite struct {
	suite.Suite
	conn Conn
	tx   tx
	repo *RecordsRepository
}

func (s *RecordSuite) SetupSuite() {
	log := logs.New(logs.TestConfig)
	conn, err := NewInMemorySqlite(log)
	if err != nil {
		s.Require().NoError(err)
	}
	s.conn = conn
}

func (s *RecordSuite) SetupTest() {
	repo, tx, err := NewRecordsRepository(s.conn).beginTx()
	if err != nil {
		s.Require().NoError(err)
	}
	s.repo = repo
	s.tx = tx
}

func (s *RecordSuite) TearDownTest() {
	s.Require().NoError(s.tx.Rollback())
}

func (s *RecordSuite) TestGetDate_NoRecord() {
	now := date.Today()
	actual, err := s.repo.GetDate(now)
	s.Require().NoError(err)

	expect := &DateList{
		Current: now,
		Prev:    make([]date.Date, 0, 12),
		Next:    make([]date.Date, 0, 12),
	}

	s.Require().Equal(expect, actual)
}

func (s *RecordSuite) TestGetDate_WithRecord() {
	s.Require().NoError(s.mockRecords())
	cases := []struct {
		current  date.Date
		prevFrom date.Date
		prevTo   date.Date
		nextFrom date.Date
		nextTo   date.Date
	}{
		{
			current:  date.MustParseISO("2022-12-31"),
			nextFrom: date.MustParseISO("2023-01-01"),
			nextTo:   date.MustParseISO("2023-01-12"),
		},
		{
			current:  date.MustParseISO("2023-01-01"),
			nextFrom: date.MustParseISO("2023-01-02"),
			nextTo:   date.MustParseISO("2023-01-13"),
		},
		{
			current:  date.MustParseISO("2023-01-04"),
			prevFrom: date.MustParseISO("2023-01-01"),
			prevTo:   date.MustParseISO("2023-01-03"),
			nextFrom: date.MustParseISO("2023-01-05"),
			nextTo:   date.MustParseISO("2023-01-16"),
		},
		{
			current:  date.MustParseISO("2023-01-15"),
			prevFrom: date.MustParseISO("2023-01-03"),
			prevTo:   date.MustParseISO("2023-01-14"),
			nextFrom: date.MustParseISO("2023-01-16"),
			nextTo:   date.MustParseISO("2023-01-27"),
		},
		{
			current:  date.MustParseISO("2023-01-28"),
			prevFrom: date.MustParseISO("2023-01-16"),
			prevTo:   date.MustParseISO("2023-01-27"),
			nextFrom: date.MustParseISO("2023-01-29"),
			nextTo:   date.MustParseISO("2023-01-31"),
		},
		{
			current:  date.MustParseISO("2023-01-31"),
			prevFrom: date.MustParseISO("2023-01-19"),
			prevTo:   date.MustParseISO("2023-01-30"),
		},
		{
			current:  date.MustParseISO("2023-02-01"),
			prevFrom: date.MustParseISO("2023-01-20"),
			prevTo:   date.MustParseISO("2023-01-31"),
		},
	}

	for _, c := range cases {
		s.Run(c.current.String(), func() {
			actual, err := s.repo.GetDate(c.current)
			s.Require().NoError(err)
			s.Require().Equal(c.current, actual.Current)

			if c.prevFrom != date.Zero && c.prevTo != date.Zero {
				cur := c.prevTo
				ptr := 0
				for cur >= c.prevFrom {
					s.Require().Equal(cur, actual.Prev[ptr], cur.String()+"!="+actual.Prev[ptr].String())
					ptr++
					cur = cur.AddDate(0, 0, -1)
				}
				s.Require().Equal(ptr, len(actual.Prev))
			}

			if c.nextFrom != date.Zero && c.nextTo != date.Zero {
				cur := c.nextFrom
				ptr := 0
				for cur <= c.nextTo {
					s.Require().Equal(cur, actual.Next[ptr], cur.String()+"!="+actual.Next[ptr].String())
					ptr++
					cur = cur.AddDate(0, 0, 1)
				}
				s.Require().Equal(ptr, len(actual.Next))
			}
		})
	}
}

func (s *RecordSuite) TestGetLatestDate_NoRecord() {
	d, err := s.repo.GetLatestDate()
	s.Require().NoError(err)
	s.Require().Zero(d)
}

func (s *RecordSuite) TestGetLatestDate_WithRecord() {
	s.Require().NoError(s.mockRecords())
	actual, err := s.repo.GetLatestDate()
	s.Require().NoError(err)

	expect := date.MustParseISO("2023-01-31")
	s.Require().Equal(expect, actual)
}

func (s *RecordSuite) TestUpsertRecord_Insert() {
	s.Require().NoError(s.mockRecords())
	record := AssetRecord{
		Id:               sql.NullInt64{Valid: false},
		AssetId:          2,
		Name:             "TFFIF",
		Broker:           "SCBS",
		DefaultIncrement: decimal.NullDecimal{Valid: false},
		BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(50)),
		CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(60)),
		RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(10)),
		Note:             sql.NullString{Valid: false},
	}

	d := date.MustParseISO("2023-08-23")
	s.Require().NoError(s.repo.UpsertRecord(record, d))

	var id int
	s.Require().NoError(s.tx.Get(&id, "SELECT id FROM records WHERE date = '2023-08-23'"))
	s.Require().Equal(63, id)
}

func (s *RecordSuite) TestUpsertRecord_Update() {
	s.Require().NoError(s.mockRecords())
	record := AssetRecord{
		Id:               sql.NullInt64{Int64: 1, Valid: true},
		AssetId:          2,
		Name:             "TFFIF",
		Broker:           "SCBS",
		DefaultIncrement: decimal.NullDecimal{Valid: false},
		BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(50)),
		CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(60)),
		RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(10)),
		Note:             sql.NullString{Valid: false},
	}

	d := date.MustParseISO("2023-08-23")
	s.Require().NoError(s.repo.UpsertRecord(record, d))

	var id int
	s.Require().NoError(s.tx.Get(&id, "SELECT id FROM records WHERE date = '2023-08-23'"))
	s.Require().Equal(1, id)
}

func (s *RecordSuite) TestDeleteRecordById() {
	s.Require().NoError(s.mockRecords())
	s.Require().NoError(s.repo.DeleteRecordById(5))
	row := s.tx.QueryRowx("SELECT * FROM records WHERE id = 5")

	var id *int = nil
	err := row.Scan(&id)
	s.Require().ErrorIs(err, sql.ErrNoRows)
}

func (s *RecordSuite) TestDeleteRecordByDate() {
	s.Require().NoError(s.mockRecords())
	rows, err := s.repo.DeleteRecordByDate(date.MustParseISO("2023-01-02"))
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
	s.Require().NoError(s.mockRecords())
	records, err := s.repo.GetRecordByDate(date.MustParseISO("2023-05-03"))
	s.Require().NoError(err)
	s.Require().Equal(0, len(records))
}

func (s *RecordSuite) TestGetRecordByDate_Found() {
	s.Require().NoError(s.mockRecords())
	s.Require().NoError(s.mockRecordsWithMultipleType())

	actual, err := s.repo.GetRecordByDate(date.MustParseISO("2023-01-03"))
	s.Require().NoError(err)

	expect := []AssetTypeRecord{
		{
			Id:          sql.NullInt64{Int64: 2, Valid: true},
			Name:        "MF",
			IsCash:      false,
			IsLiability: false,
			Assets: []AssetRecord{
				{
					Id:               sql.NullInt64{Int64: 64, Valid: true},
					AssetId:          4,
					Name:             "b2",
					Broker:           "finno",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(164)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(264)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(364)),
					Note:             sql.NullString{String: "test", Valid: true},
				},
				{
					Id:               sql.NullInt64{Int64: 63, Valid: true},
					AssetId:          3,
					Name:             "b1",
					Broker:           "finno",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(2000)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(163)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(263)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(363)),
					Note:             sql.NullString{String: "test2", Valid: true},
				},
			},
		},
		{
			Id:          sql.NullInt64{Int64: 1, Valid: true},
			Name:        "Stocks",
			IsCash:      true,
			IsLiability: true,
			Assets: []AssetRecord{
				{
					Id:               sql.NullInt64{Int64: 5, Valid: true},
					AssetId:          1,
					Name:             "a1",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(105)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(205)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(305)),
					Note:             sql.NullString{Valid: false},
				},
				{
					Id:               sql.NullInt64{Int64: 6, Valid: true},
					AssetId:          2,
					Name:             "a2",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(106)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(206)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(306)),
					Note:             sql.NullString{Valid: false},
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
	s.Require().True(actual[0].Assets[0].DefaultIncrement.Valid)
	s.Require().True(actual[0].Assets[0].BoughtValue.Valid)
	s.Require().True(actual[0].Assets[0].CurrentValue.Valid)
	s.Require().True(actual[0].Assets[0].RealizedValue.Valid)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Decimal.Equal(actual[0].Assets[0].DefaultIncrement.Decimal))
	s.Require().True(expect[0].Assets[0].BoughtValue.Decimal.Equal(actual[0].Assets[0].BoughtValue.Decimal))
	s.Require().True(expect[0].Assets[0].CurrentValue.Decimal.Equal(actual[0].Assets[0].CurrentValue.Decimal))
	s.Require().True(expect[0].Assets[0].RealizedValue.Decimal.Equal(actual[0].Assets[0].RealizedValue.Decimal))
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(actual[0].Assets[1].DefaultIncrement.Valid)
	s.Require().True(actual[0].Assets[1].BoughtValue.Valid)
	s.Require().True(actual[0].Assets[1].CurrentValue.Valid)
	s.Require().True(actual[0].Assets[1].RealizedValue.Valid)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Decimal.Equal(actual[0].Assets[1].DefaultIncrement.Decimal))
	s.Require().True(expect[0].Assets[1].BoughtValue.Decimal.Equal(actual[0].Assets[1].BoughtValue.Decimal))
	s.Require().True(expect[0].Assets[1].CurrentValue.Decimal.Equal(actual[0].Assets[1].CurrentValue.Decimal))
	s.Require().True(expect[0].Assets[1].RealizedValue.Decimal.Equal(actual[0].Assets[1].RealizedValue.Decimal))
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
	s.Require().True(expect[1].Assets[0].DefaultIncrement.Decimal.Equal(actual[1].Assets[0].DefaultIncrement.Decimal))
	s.Require().True(expect[1].Assets[0].BoughtValue.Decimal.Equal(actual[1].Assets[0].BoughtValue.Decimal))
	s.Require().True(expect[1].Assets[0].CurrentValue.Decimal.Equal(actual[1].Assets[0].CurrentValue.Decimal))
	s.Require().True(expect[1].Assets[0].RealizedValue.Decimal.Equal(actual[1].Assets[0].RealizedValue.Decimal))
	s.Require().Equal(expect[1].Assets[0].Note, actual[1].Assets[0].Note)

	s.Require().Equal(expect[1].Assets[1].Id, actual[1].Assets[1].Id)
	s.Require().Equal(expect[1].Assets[1].AssetId, actual[1].Assets[1].AssetId)
	s.Require().Equal(expect[1].Assets[1].Name, actual[1].Assets[1].Name)
	s.Require().Equal(expect[1].Assets[1].Broker, actual[1].Assets[1].Broker)
	s.Require().True(expect[1].Assets[1].DefaultIncrement.Decimal.Equal(actual[1].Assets[1].DefaultIncrement.Decimal))
	s.Require().True(expect[1].Assets[1].BoughtValue.Decimal.Equal(actual[1].Assets[1].BoughtValue.Decimal))
	s.Require().True(expect[1].Assets[1].CurrentValue.Decimal.Equal(actual[1].Assets[1].CurrentValue.Decimal))
	s.Require().True(expect[1].Assets[1].RealizedValue.Decimal.Equal(actual[1].Assets[1].RealizedValue.Decimal))
	s.Require().Equal(expect[1].Assets[1].Note, actual[1].Assets[1].Note)
}

func (s *RecordSuite) TestGetRecordDraft_EmptyRecord() {
	s.Require().NoError(s.mockAssets())
	s.Require().NoError(s.mockInactiveAssets())

	actual, err := s.repo.GetRecordDraft()
	s.Require().NoError(err)

	expect := []AssetTypeRecord{
		{
			Id:          sql.NullInt64{Int64: 1, Valid: true},
			Name:        "Stocks",
			IsCash:      true,
			IsLiability: true,
			Assets: []AssetRecord{
				{
					Id:               sql.NullInt64{Valid: false},
					AssetId:          1,
					Name:             "a1",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NullDecimal{Valid: false},
					CurrentValue:     decimal.NullDecimal{Valid: false},
					RealizedValue:    decimal.NullDecimal{Valid: false},
					Note:             sql.NullString{Valid: false},
				},
				{
					Id:               sql.NullInt64{Valid: false},
					AssetId:          2,
					Name:             "a2",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NullDecimal{Valid: false},
					CurrentValue:     decimal.NullDecimal{Valid: false},
					RealizedValue:    decimal.NullDecimal{Valid: false},
					Note:             sql.NullString{Valid: false},
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
	s.Require().True(actual[0].Assets[0].DefaultIncrement.Valid)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Decimal.Equal(actual[0].Assets[0].DefaultIncrement.Decimal))
	s.Require().False(actual[0].Assets[0].BoughtValue.Valid)
	s.Require().False(actual[0].Assets[0].CurrentValue.Valid)
	s.Require().False(actual[0].Assets[0].RealizedValue.Valid)
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(actual[0].Assets[1].DefaultIncrement.Valid)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Decimal.Equal(actual[0].Assets[1].DefaultIncrement.Decimal))
	s.Require().False(actual[0].Assets[1].BoughtValue.Valid)
	s.Require().False(actual[0].Assets[1].CurrentValue.Valid)
	s.Require().False(actual[0].Assets[1].RealizedValue.Valid)
	s.Require().Equal(expect[0].Assets[1].Note, actual[0].Assets[1].Note)
}

func (s *RecordSuite) TestGetRecordDraft_PartialHitRecord() {
	s.Require().NoError(s.mockAssets())
	s.Require().NoError(s.mockInactiveAssets())

	actual, err := s.repo.GetRecordDraft()
	s.Require().NoError(err)

	expect := []AssetTypeRecord{
		{
			Id:          sql.NullInt64{Int64: 1, Valid: true},
			Name:        "Stocks",
			IsCash:      true,
			IsLiability: true,
			Assets: []AssetRecord{
				{
					Id:               sql.NullInt64{Valid: false},
					AssetId:          1,
					Name:             "a1",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(101)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(201)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(301)),
					Note:             sql.NullString{Valid: false},
				},
				{
					Id:               sql.NullInt64{Valid: false},
					AssetId:          2,
					Name:             "a2",
					Broker:           "scbs",
					DefaultIncrement: decimal.NewNullDecimal(decimal.NewFromInt(0)),
					BoughtValue:      decimal.NewNullDecimal(decimal.NewFromInt(102)),
					CurrentValue:     decimal.NewNullDecimal(decimal.NewFromInt(202)),
					RealizedValue:    decimal.NewNullDecimal(decimal.NewFromInt(302)),
					Note:             sql.NullString{Valid: false},
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
	s.Require().True(actual[0].Assets[0].DefaultIncrement.Valid)
	s.Require().True(expect[0].Assets[0].DefaultIncrement.Decimal.Equal(actual[0].Assets[0].DefaultIncrement.Decimal))
	s.Require().False(actual[0].Assets[0].BoughtValue.Valid)
	s.Require().False(actual[0].Assets[0].CurrentValue.Valid)
	s.Require().False(actual[0].Assets[0].RealizedValue.Valid)
	s.Require().Equal(expect[0].Assets[0].Note, actual[0].Assets[0].Note)

	s.Require().Equal(expect[0].Assets[1].Id, actual[0].Assets[1].Id)
	s.Require().Equal(expect[0].Assets[1].AssetId, actual[0].Assets[1].AssetId)
	s.Require().Equal(expect[0].Assets[1].Name, actual[0].Assets[1].Name)
	s.Require().Equal(expect[0].Assets[1].Broker, actual[0].Assets[1].Broker)
	s.Require().True(actual[0].Assets[1].DefaultIncrement.Valid)
	s.Require().True(expect[0].Assets[1].DefaultIncrement.Decimal.Equal(actual[0].Assets[1].DefaultIncrement.Decimal))
	s.Require().False(actual[0].Assets[1].BoughtValue.Valid)
	s.Require().False(actual[0].Assets[1].CurrentValue.Valid)
	s.Require().False(actual[0].Assets[1].RealizedValue.Valid)
	s.Require().Equal(expect[0].Assets[1].Note, actual[0].Assets[1].Note)
}

func (s *RecordSuite) mockAssets() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Stocks', true, true, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a1', 'scbs', 1, 0, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a2', 'scbs', 1, 0, 2, true)")

	return
}

func (s *RecordSuite) mockInactiveAssets() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Bonds', true, true, 1, false)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b1', 'scbs', 2, 0, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b2', 'scbs', 2, 0, 1, true)")

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Cash', true, true, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c1', 'scbs', 3, 0, 1, false)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c2', 'scbs', 3, 0, 1, false)")

	return
}

func (s *RecordSuite) mockPartialRecords() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	err = s.mockAssets()
	if err != nil {
		return
	}

	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-01', 101, 201, 301, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-01', 102, 202, 302, null)")

	return
}

func (s *RecordSuite) mockRecords() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	err = s.mockAssets()
	if err != nil {
		return
	}

	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-01', 101, 201, 301, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-01', 102, 202, 302, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-02', 103, 203, 303, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-02', 104, 204, 304, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-03', 105, 205, 305, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-03', 106, 206, 306, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-04', 107, 207, 307, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-04', 108, 208, 308, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-05', 109, 209, 309, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-05', 110, 210, 310, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-06', 111, 211, 311, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-06', 112, 212, 312, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-07', 113, 213, 313, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-07', 114, 214, 314, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-08', 115, 215, 315, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-08', 116, 216, 316, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-09', 117, 217, 317, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-09', 118, 218, 318, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-10', 119, 219, 319, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-10', 120, 220, 320, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-11', 121, 221, 321, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-11', 122, 222, 322, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-12', 123, 223, 323, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-12', 124, 224, 324, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-13', 125, 225, 325, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-13', 126, 226, 326, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-14', 127, 227, 327, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-14', 128, 228, 328, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-15', 129, 229, 329, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-15', 130, 230, 330, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-16', 131, 231, 331, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-16', 132, 232, 332, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-17', 133, 233, 333, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-17', 134, 234, 334, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-18', 135, 235, 335, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-18', 136, 236, 336, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-19', 137, 237, 337, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-19', 138, 238, 338, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-20', 139, 239, 339, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-20', 140, 240, 340, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-21', 141, 241, 341, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-21', 142, 242, 342, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-22', 143, 243, 343, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-22', 144, 244, 344, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-23', 145, 245, 345, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-23', 146, 246, 346, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-24', 147, 247, 347, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-24', 148, 248, 348, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-25', 149, 249, 349, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-25', 150, 250, 350, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-26', 151, 251, 351, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-26', 152, 252, 352, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-27', 153, 253, 353, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-27', 154, 254, 354, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-28', 155, 255, 355, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-28', 156, 256, 356, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-29', 157, 257, 357, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-29', 158, 258, 358, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-30', 159, 259, 359, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-30', 160, 260, 360, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (1, '2023-01-31', 161, 261, 361, null)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (2, '2023-01-31', 162, 262, 362, null)")

	return
}

func (s *RecordSuite) mockRecordsWithMultipleType() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('MF', false, false, 0, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b1', 'finno', 2, 2000, 2, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b2', 'finno', 2, 0, 1, true)")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (3, '2023-01-03', 163, 263, 363, 'test2')")
	s.tx.MustExec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES (4, '2023-01-03', 164, 264, 364, 'test')")
	return
}
