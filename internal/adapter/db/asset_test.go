package db

import (
	"database/sql"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/ports"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAssetSuite(t *testing.T) {
	suite.Run(t, new(AssetSuite))
}

type AssetSuite struct {
	suite.Suite
	conn ports.Conn
	tx   tx
	repo *AssetsRepository
}

func (s *AssetSuite) SetupSuite() {
	log := logs.New(logs.TestConfig)
	conn, err := NewInMemorySqlite(log)
	if err != nil {
		s.Require().NoError(err)
	}
	s.conn = conn
}

func (s *AssetSuite) SetupTest() {
	repo, tx, err := NewAssetsRepository(s.conn).beginTx()
	if err != nil {
		s.Require().NoError(err)
	}
	s.repo = repo
	s.tx = tx
}

func (s *AssetSuite) TearDownTest() {
	s.Require().NoError(s.tx.Rollback())
}

func (s *AssetSuite) TestGet_NotFound() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, lo.ToPtr(5))
	s.Require().NoError(err)
	s.Require().Equal(0, len(resp))
}

func (s *AssetSuite) TestGet_ByIsActive() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(lo.ToPtr(true), nil)
	s.Require().NoError(err)
	s.Require().Equal(5, len(resp))
	s.Require().Equal("d2", resp[0].Name)
	s.Require().Equal(true, resp[0].IsActive)
	s.Require().Equal("d1", resp[1].Name)
	s.Require().Equal(true, resp[1].IsActive)
	s.Require().Equal("a1", resp[2].Name)
	s.Require().Equal(true, resp[2].IsActive)
	s.Require().Equal("a2", resp[3].Name)
	s.Require().Equal(true, resp[3].IsActive)
	s.Require().Equal("b1", resp[4].Name)
	s.Require().Equal(true, resp[4].IsActive)
}

func (s *AssetSuite) TestGet_ByTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, lo.ToPtr(3))
	s.Require().NoError(err)
	s.Require().Equal(2, len(resp))
	s.Require().Equal("c1", resp[0].Name)
	s.Require().Equal(false, resp[0].IsActive)
	s.Require().Equal("c2", resp[1].Name)
	s.Require().Equal(false, resp[1].IsActive)
}

func (s *AssetSuite) TestGet_ByIsActiveAndTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(lo.ToPtr(false), lo.ToPtr(2))
	s.Require().NoError(err)
	s.Require().Equal(1, len(resp))
	s.Require().Equal("b2", resp[0].Name)
	s.Require().Equal(false, resp[0].IsActive)
}

func (s *AssetSuite) TestGet_ByAll() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, nil)
	s.Require().NoError(err)
	s.Require().Equal(8, len(resp))
	s.Require().Equal("c1", resp[0].Name)
	s.Require().Equal("c2", resp[1].Name)
	s.Require().Equal("d2", resp[2].Name)
	s.Require().Equal("d1", resp[3].Name)
	s.Require().Equal("a1", resp[4].Name)
	s.Require().Equal("a2", resp[5].Name)
	s.Require().Equal("b1", resp[6].Name)
	s.Require().Equal("b2", resp[7].Name)
}

func (s *AssetSuite) TestGetNames_NotFound() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(nil, lo.ToPtr(5))
	s.Require().NoError(err)
	s.Require().Equal(0, len(resp))
}

func (s *AssetSuite) TestGetNames_ByIsActive() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(lo.ToPtr(true), nil)
	s.Require().NoError(err)
	s.Require().Equal(5, len(resp))
	s.Require().Equal(1, resp[0].Id)
	s.Require().Equal("a1", resp[0].Name)
	s.Require().Equal(2, resp[1].Id)
	s.Require().Equal("a2", resp[1].Name)
	s.Require().Equal(3, resp[2].Id)
	s.Require().Equal("b1", resp[2].Name)
	s.Require().Equal(7, resp[3].Id)
	s.Require().Equal("d1", resp[3].Name)
	s.Require().Equal(8, resp[4].Id)
	s.Require().Equal("d2", resp[4].Name)
}

func (s *AssetSuite) TestGetNames_ByTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(nil, lo.ToPtr(3))
	s.Require().NoError(err)
	s.Require().Equal(2, len(resp))
	s.Require().Equal("c1", resp[0].Name)
	s.Require().Equal(5, resp[0].Id)
	s.Require().Equal("c2", resp[1].Name)
	s.Require().Equal(6, resp[1].Id)
}

func (s *AssetSuite) TestGetNames_ByIsActiveAndTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(lo.ToPtr(false), lo.ToPtr(2))
	s.Require().NoError(err)
	s.Require().Equal(1, len(resp))
	s.Require().Equal("b2", resp[0].Name)
	s.Require().Equal(4, resp[0].Id)
}

func (s *AssetSuite) TestGetNames_ByAll() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(nil, nil)
	s.Require().NoError(err)
	s.Require().Equal(8, len(resp))
	s.Require().Equal("a1", resp[0].Name)
	s.Require().Equal("a2", resp[1].Name)
	s.Require().Equal("b1", resp[2].Name)
	s.Require().Equal("b2", resp[3].Name)
	s.Require().Equal("c1", resp[4].Name)
	s.Require().Equal("c2", resp[5].Name)
	s.Require().Equal("d1", resp[6].Name)
	s.Require().Equal("d2", resp[7].Name)
}

func (s *AssetSuite) TestUpsert_Insert() {
	s.Require().NoError(s.mockAssets())
	var assetDetail = model.AssetDetail{
		Id:               sql.NullInt64{Valid: false},
		Name:             "Monika",
		Broker:           "ddlc",
		TypeId:           5,
		TypeName:         "",
		DefaultIncrement: decimal.NullDecimal{Decimal: decimal.NewFromInt(1000), Valid: true},
		Sequence:         sql.NullInt64{Valid: false},
		IsActive:         true,
	}
	s.Require().NoError(s.repo.Upsert(assetDetail))

	res, err := s.repo.Get(lo.ToPtr(true), lo.ToPtr(5))
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(sql.NullInt64{Int64: 9, Valid: true}, res[0].Id)
	s.Require().Equal("Monika", res[0].Name)
	s.Require().Equal("ddlc", res[0].Broker)
	s.Require().Equal(5, res[0].TypeId)
	s.Require().Equal("Just Monika", res[0].TypeName)
	s.Require().True(res[0].DefaultIncrement.Valid)
	s.Require().True(decimal.NewFromInt(1000).Equal(res[0].DefaultIncrement.Decimal))
	s.Require().Equal(sql.NullInt64{Int64: 0, Valid: true}, res[0].Sequence)
	s.Require().True(res[0].IsActive)
}

func (s *AssetSuite) TestUpsert_Update() {
	s.Require().NoError(s.mockAssets())
	var assetDetail = model.AssetDetail{
		Id:               sql.NullInt64{Int64: 3, Valid: true},
		Name:             "Monika",
		Broker:           "ddlc",
		TypeId:           5,
		TypeName:         "",
		DefaultIncrement: decimal.NullDecimal{Decimal: decimal.NewFromInt(1000), Valid: true},
		Sequence:         sql.NullInt64{Valid: false},
		IsActive:         true,
	}
	s.Require().NoError(s.repo.Upsert(assetDetail))

	res, err := s.repo.Get(lo.ToPtr(true), lo.ToPtr(5))
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(sql.NullInt64{Int64: 3, Valid: true}, res[0].Id)
	s.Require().Equal("Monika", res[0].Name)
	s.Require().Equal("ddlc", res[0].Broker)
	s.Require().Equal(5, res[0].TypeId)
	s.Require().Equal("Just Monika", res[0].TypeName)
	s.Require().True(res[0].DefaultIncrement.Valid)
	s.Require().True(decimal.NewFromInt(1000).Equal(res[0].DefaultIncrement.Decimal))
	s.Require().Equal(sql.NullInt64{Int64: 0, Valid: true}, res[0].Sequence)
	s.Require().True(res[0].IsActive)

	res, err = s.repo.Get(lo.ToPtr(true), lo.ToPtr(2))
	s.Require().NoError(err)
	s.Require().Equal(0, len(res))
}

func (s *AssetSuite) TestUpsert_Delete() {
	s.Require().NoError(s.mockAssets())
	s.Require().NoError(s.repo.Delete(3))

	res, err := s.repo.Get(lo.ToPtr(true), lo.ToPtr(2))
	s.Require().NoError(err)
	s.Require().Equal(0, len(res))
}

func (s *AssetSuite) TestUpdateSequence() {
	s.Require().NoError(s.mockAssets())
	var sequence = model.AssetSequenceDetail{
		Id:       3,
		Sequence: 99,
	}
	s.Require().NoError(s.repo.UpdateSequence(sequence))

	res, err := s.repo.Get(lo.ToPtr(true), lo.ToPtr(2))
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(sql.NullInt64{Int64: 3, Valid: true}, res[0].Id)
	s.Require().Equal(sql.NullInt64{Int64: 99, Valid: true}, res[0].Sequence)
}

func (s *AssetSuite) mockAssets() (err error) {
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
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('c2', 'scbs', 3, 0, 1, false)")

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Mutual Funds', true, true, 1, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('d1', 'scbs', 4, 0, 2, true)")
	s.tx.MustExec("INSERT INTO assets(name, broker ,type_id, default_increment, sequence, is_active) VALUES ('d2', 'scbs', 4, 0, 1, true)")

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Just Monika', true, true, 1, true)")

	return
}
