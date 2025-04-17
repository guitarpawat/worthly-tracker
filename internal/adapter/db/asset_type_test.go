package db

import (
	"database/sql"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAssetTypeSuite(t *testing.T) {
	suite.Run(t, new(AssetTypeSuite))
}

type AssetTypeSuite struct {
	suite.Suite
	conn Conn
	tx   tx
	repo *AssetTypesRepository
}

func (s *AssetTypeSuite) SetupSuite() {
	log := logs.New(logs.TestConfig)
	conn, err := NewInMemorySqlite(log)
	if err != nil {
		s.Require().NoError(err)
	}
	s.conn = conn
}

func (s *AssetTypeSuite) SetupTest() {
	repo, tx, err := NewAssetTypesRepository(s.conn).beginTx()
	if err != nil {
		s.Require().NoError(err)
	}
	s.repo = repo
	s.tx = tx
}

func (s *AssetTypeSuite) TearDownTest() {
	s.Require().NoError(s.tx.Rollback())
}

func (s *AssetTypeSuite) TestGet_Empty() {
	res, err := s.repo.Get(nil)
	s.Require().NoError(err)
	s.Require().Empty(res)
}

func (s *AssetTypeSuite) TestGet_IsActive() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.Get(lo.ToPtr(true))
	s.Require().NoError(err)
	s.Require().Equal(2, len(res))
	s.Require().Equal("Stocks", res[0].Name)
	s.Require().Equal(true, res[0].IsActive)
	s.Require().Equal("Mutual Funds", res[1].Name)
	s.Require().Equal(true, res[1].IsActive)
}

func (s *AssetTypeSuite) TestGet_All() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.Get(nil)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal("Cash", res[0].Name)
	s.Require().Equal(false, res[0].IsActive)
	s.Require().Equal("Stocks", res[1].Name)
	s.Require().Equal(true, res[1].IsActive)
	s.Require().Equal("Bonds", res[2].Name)
	s.Require().Equal(false, res[2].IsActive)
	s.Require().Equal("Mutual Funds", res[3].Name)
	s.Require().Equal(true, res[3].IsActive)
}

func (s *AssetTypeSuite) TestGetNames_Empty() {
	res, err := s.repo.GetNames(nil)
	s.Require().NoError(err)
	s.Require().Empty(res)
}

func (s *AssetTypeSuite) TestGetNames_IsActive() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.GetNames(lo.ToPtr(false))
	s.Require().NoError(err)
	s.Require().Equal(2, len(res))
	s.Require().Equal("Bonds", res[0].Name)
	s.Require().Equal(2, res[0].Id)
	s.Require().Equal("Cash", res[1].Name)
	s.Require().Equal(3, res[1].Id)
}

func (s *AssetTypeSuite) TestGetNames_All() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.GetNames(nil)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal("Bonds", res[0].Name)
	s.Require().Equal(2, res[0].Id)
	s.Require().Equal("Cash", res[1].Name)
	s.Require().Equal(3, res[1].Id)
	s.Require().Equal("Mutual Funds", res[2].Name)
	s.Require().Equal(4, res[2].Id)
	s.Require().Equal("Stocks", res[3].Name)
	s.Require().Equal(1, res[3].Id)
}

func (s *AssetTypeSuite) TestUpsert_Insert() {
	s.Require().NoError(s.mockAssetType())
	req := AssetTypeDetail{
		Id:          sql.NullInt64{Valid: false},
		Name:        "Test",
		IsCash:      false,
		IsLiability: false,
		Sequence:    sql.NullInt64{Int64: 0, Valid: true},
		IsActive:    true,
	}

	err := s.repo.Upsert(req)
	s.Require().NoError(err)

	res, err := s.repo.Get(lo.ToPtr(true))
	s.Require().NoError(err)
	s.Require().Equal(3, len(res))
	s.Require().Equal("Test", res[0].Name)
	s.Require().Equal(sql.NullInt64{Int64: 5, Valid: true}, res[0].Id)
	s.Require().Equal(true, res[0].IsActive)
	s.Require().Equal(false, res[0].IsCash)
	s.Require().Equal(false, res[0].IsLiability)
	s.Require().Equal(sql.NullInt64{Int64: 0, Valid: true}, res[0].Sequence)

	res, err = s.repo.Get(nil)
	s.Require().NoError(err)
	s.Require().Equal(5, len(res))
}

func (s *AssetTypeSuite) TestUpsert_Update() {
	s.Require().NoError(s.mockAssetType())
	req := AssetTypeDetail{
		Id:          sql.NullInt64{Int64: 2, Valid: true},
		Name:        "Test",
		IsCash:      false,
		IsLiability: false,
		Sequence:    sql.NullInt64{Int64: 0, Valid: true},
		IsActive:    true,
	}

	err := s.repo.Upsert(req)
	s.Require().NoError(err)

	res, err := s.repo.Get(lo.ToPtr(true))
	s.Require().NoError(err)
	s.Require().Equal(3, len(res))
	s.Require().Equal("Test", res[0].Name)
	s.Require().Equal(sql.NullInt64{Int64: 2, Valid: true}, res[0].Id)
	s.Require().Equal(true, res[0].IsActive)
	s.Require().Equal(false, res[0].IsCash)
	s.Require().Equal(false, res[0].IsLiability)
	s.Require().Equal(sql.NullInt64{Int64: 0, Valid: true}, res[0].Sequence)

	res, err = s.repo.Get(nil)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
}

func (s *AssetTypeSuite) TestDelete() {
	s.Require().NoError(s.mockAssetType())
	err := s.repo.Delete(1)
	s.Require().NoError(err)
}

func (s *AssetTypeSuite) TestUpdateSequence() {
	s.Require().NoError(s.mockAssetType())
	req := AssetTypeSequenceDetail{
		Id:       1,
		Sequence: 99,
	}

	err := s.repo.UpdateSequence(req)
	s.Require().NoError(err)

	res, err := s.repo.Get(nil)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal(sql.NullInt64{Int64: 1, Valid: true}, res[3].Id)
	s.Require().Equal(sql.NullInt64{Int64: 99, Valid: true}, res[3].Sequence)
}

func (s *AssetTypeSuite) mockAssetType() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Stocks', true, true, 2, true)")
	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Bonds', true, true, 3, false)")
	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Cash', true, true, 1, false)")
	s.tx.MustExec("INSERT INTO asset_types(name, is_cash, is_liability, sequence, is_active) VALUES ('Mutual Funds', true, true, 4, true)")
	return
}
