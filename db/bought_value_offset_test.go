//go:build test || integration

package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/logs"
	"worthly-tracker/model"
	"worthly-tracker/ports"
)

func TestBoughtValueOffsetTestSuite(t *testing.T) {
	suite.Run(t, new(BoughtValueOffsetTestSuite))
}

type BoughtValueOffsetTestSuite struct {
	suite.Suite
	tx   *sqlx.Tx
	repo ports.BoughtValueOffsetRepo
}

func (s *BoughtValueOffsetTestSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
	Init()
	s.repo = GetBoughtValueOffsetRepo()
}

func (s *BoughtValueOffsetTestSuite) SetupTest() {
	var err error
	s.tx, err = GetDB().BeginTx()
	s.Require().NoError(err)

	s.Require().NoError(s.mockBoughtValueOffsets())
}

func (s *BoughtValueOffsetTestSuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *BoughtValueOffsetTestSuite) TestGet() {
	testData := []struct {
		name          string
		date          string
		id            int
		assetId       int
		effectiveDate string
		offsetPrice   string
		note          *string
	}{
		{
			"sequential dates",
			"2023-08-23",
			2,
			1,
			"2023-08-23",
			"2.0",
			pointy.String("test2"),
		},
		{
			"non-sequential dates 1",
			"2023-08-23",
			4,
			2,
			"2023-08-23",
			"4.0",
			pointy.String("test4"),
		},
		{
			"non-sequential dates 2",
			"2023-08-23",
			9,
			3,
			"2023-08-23",
			"9.0",
			pointy.String("test9"),
		},
		{
			"has date before",
			"2023-08-23",
			10,
			4,
			"2023-08-22",
			"9.4",
			pointy.String("test10"),
		},
		{
			"no date before",
			"2023-08-23",
			0,
			5,
			"2023-08-23",
			"0",
			nil,
		},
		{
			"invalid asset id",
			"2023-08-23",
			0,
			6,
			"2023-08-23",
			"0",
			nil,
		},
	}

	for _, test := range testData {
		s.Run(test.name, func() {
			currentDate, err := model.NewDate(test.date)
			s.Require().NoError(err)

			effectiveDate, err := model.NewDate(test.effectiveDate)
			s.Require().NoError(err)

			res, err := s.repo.Get(currentDate, test.assetId, s.tx)
			s.Require().NoError(err)
			s.Require().NotNil(res)
			s.Require().Equal(test.id, *res.Id)
			s.Require().Equal(test.assetId, res.AssetId)
			s.Require().Equal(effectiveDate, res.EffectiveDate)
			s.Require().True(res.OffsetPrice.Equal(decimal.RequireFromString(test.offsetPrice)))
			s.Require().Equal(test.note, res.Note)
		})
	}
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByAssetId_Empty() {
	res, err := s.repo.GetAllByAssetId(6, s.tx)
	s.Require().NoError(err)
	s.Require().Empty(res)
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByAssetId_Success() {
	res, err := s.repo.GetAllByAssetId(1, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(3, len(res))

	s.Require().Equal(1, *res[0].Id)
	s.Require().Equal(1, res[0].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-22"), res[0].EffectiveDate)
	s.Require().True(res[0].OffsetPrice.Equal(decimal.NewFromInt(1)))
	s.Require().Equal("test1", *res[0].Note)

	s.Require().Equal(2, *res[1].Id)
	s.Require().Equal(1, res[1].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-23"), res[1].EffectiveDate)
	s.Require().True(res[1].OffsetPrice.Equal(decimal.NewFromInt(2)))
	s.Require().Equal("test2", *res[1].Note)

	s.Require().Equal(3, *res[2].Id)
	s.Require().Equal(1, res[2].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-24"), res[2].EffectiveDate)
	s.Require().True(res[2].OffsetPrice.Equal(decimal.NewFromInt(3)))
	s.Require().Equal("test3", *res[2].Note)
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByDate_Empty() {
	res, err := s.repo.GetAllByDate(model.MustNewDate("1990-08-23"), s.tx)
	s.Require().NoError(err)
	s.Require().Empty(res)
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByDateSuccess() {
	res, err := s.repo.GetAllByDate(model.MustNewDate("2023-08-23"), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))

	s.Require().Equal(2, *res[0].Id)
	s.Require().Equal(1, res[0].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-23"), res[0].EffectiveDate)
	s.Require().True(res[0].OffsetPrice.Equal(decimal.NewFromInt(2)))
	s.Require().Equal("test2", *res[0].Note)

	s.Require().Equal(4, *res[1].Id)
	s.Require().Equal(2, res[1].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-23"), res[1].EffectiveDate)
	s.Require().True(res[1].OffsetPrice.Equal(decimal.NewFromInt(4)))
	s.Require().Equal("test4", *res[1].Note)

	s.Require().Equal(9, *res[2].Id)
	s.Require().Equal(3, res[2].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-23"), res[2].EffectiveDate)
	s.Require().True(res[2].OffsetPrice.Equal(decimal.NewFromInt(9)))
	s.Require().Equal("test9", *res[2].Note)

	s.Require().Equal(10, *res[3].Id)
	s.Require().Equal(4, res[3].AssetId)
	s.Require().Equal(model.MustNewDate("2023-08-22"), res[3].EffectiveDate)
	s.Require().True(res[3].OffsetPrice.Equal(decimal.RequireFromString("9.40")))
	s.Require().Equal("test10", *res[3].Note)
}

func (s *BoughtValueOffsetTestSuite) TestUpsert_Insert() {
	data := model.OffsetDetail{
		AssetId:       5,
		EffectiveDate: model.MustNewDate("2023-08-23"),
		OffsetPrice:   decimal.NewFromInt(12),
		Note:          pointy.String("test12"),
	}

	err := s.repo.Upsert(data, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(model.MustNewDate("2023-08-23"), 5, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(12, *res.Id)
	s.Require().Equal(data.AssetId, res.AssetId)
	s.Require().Equal(data.EffectiveDate, res.EffectiveDate)
	s.Require().True(data.OffsetPrice.Equal(res.OffsetPrice))
	s.Require().Equal(*data.Note, *res.Note)
}

func (s *BoughtValueOffsetTestSuite) TestUpsert_Update() {
	data := model.OffsetDetail{
		Id:            pointy.Int(1),
		AssetId:       5,
		EffectiveDate: model.MustNewDate("2023-08-23"),
		OffsetPrice:   decimal.NewFromInt(12),
		Note:          pointy.String("test12"),
	}

	err := s.repo.Upsert(data, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(model.MustNewDate("2023-08-23"), 5, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(*data.Id, *res.Id)
	s.Require().Equal(data.AssetId, res.AssetId)
	s.Require().Equal(data.EffectiveDate, res.EffectiveDate)
	s.Require().True(data.OffsetPrice.Equal(res.OffsetPrice))
	s.Require().Equal(*data.Note, *res.Note)
}

func (s *BoughtValueOffsetTestSuite) TestDelete() {
	err := s.repo.Delete(1, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(model.MustNewDate("2023-08-22"), 1, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, *res.Id)
	s.Require().True(res.OffsetPrice.Equal(decimal.Zero))
}

func (s *BoughtValueOffsetTestSuite) mockBoughtValueOffsets() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Stocks', true, true, 2, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a1', 'scbs', 1, 0, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('a2', 'scbs', 1, 0, 2, true)")

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Bonds', true, true, 3, false)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b1', 'scbs', 2, 0, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('b2', 'scbs', 2, 0, 1, false)")

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Cash', true, true, 1, false)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c1', 'scbs', 3, 0, 1, false)")

	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (1, '2023-08-22', '1.00', 'test1')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (1, '2023-08-23', '2.00', 'test2')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (1, '2023-08-24', '3.00', 'test3')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (2, '2023-08-23', '4.00', 'test4')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (2, '2023-08-22', '5.00', 'test5')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (2, '2023-08-24', '6.00', 'test6')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (3, '2023-08-22', '7.00', 'test7')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (3, '2023-08-24', '8.00', 'test8')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (3, '2023-08-23', '9.00', 'test9')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (4, '2023-08-22', '9.40', 'test10')")
	s.tx.MustExec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (5, '2023-08-24', '9.50', 'test11')")

	return
}
