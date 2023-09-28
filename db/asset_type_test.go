//go:build test || integration

package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/suite"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/logs"
	"worthly-tracker/model"
	"worthly-tracker/ports"
)

func TestAssetTypeSuite(t *testing.T) {
	suite.Run(t, new(AssetSuite))
}

type AssetTypeSuite struct {
	suite.Suite
	tx   *sqlx.Tx
	repo ports.AssetTypeRepo
}

func (s *AssetTypeSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
	Init()
	s.repo = GetAssetTypeRepo()
}

func (s *AssetTypeSuite) SetupTest() {
	var err error
	s.tx, err = GetDB().BeginTx()
	s.Require().NoError(err)
}

func (s *AssetTypeSuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *AssetTypeSuite) TestGet_Empty() {
	res, err := s.repo.Get(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Nil(res)
}

func (s *AssetTypeSuite) TestGet_IsActive() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.Get(pointy.Bool(true), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(res))
	s.Require().Equal("Stocks", *res[0].Name)
	s.Require().Equal(true, *res[0].IsActive)
	s.Require().Equal("Mutual Funds", *res[1].Name)
	s.Require().Equal(true, *res[1].IsActive)
}

func (s *AssetTypeSuite) TestGet_All() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.Get(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal("Stocks", *res[0].Name)
	s.Require().Equal(true, *res[0].IsActive)
	s.Require().Equal("Bonds", *res[1].Name)
	s.Require().Equal(false, *res[1].IsActive)
	s.Require().Equal("Cash", *res[2].Name)
	s.Require().Equal(false, *res[2].IsActive)
	s.Require().Equal("Mutual Funds", *res[3].Name)
	s.Require().Equal(true, *res[3].IsActive)
}

func (s *AssetTypeSuite) TestGetNames_Empty() {
	res, err := s.repo.GetNames(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Nil(res)
}

func (s *AssetTypeSuite) TestGetNames_IsActive() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.GetNames(pointy.Bool(false), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(res))
	s.Require().Equal("Bonds", res[0].Name)
	s.Require().Equal(1, res[0].Name)
	s.Require().Equal("Cash", res[1].Name)
	s.Require().Equal(2, res[1].Name)
}

func (s *AssetTypeSuite) TestGetNames_All() {
	s.Require().NoError(s.mockAssetType())
	res, err := s.repo.GetNames(pointy.Bool(false), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal("Stocks", res[0].Name)
	s.Require().Equal(0, res[0].Id)
	s.Require().Equal("Bonds", res[1].Name)
	s.Require().Equal(1, res[1].Id)
	s.Require().Equal("Cash", res[2].Name)
	s.Require().Equal(2, res[2].Id)
	s.Require().Equal("Mutual Funds", res[3].Name)
	s.Require().Equal(3, res[3].Id)
}

func (s *AssetTypeSuite) TestUpsert_Insert() {
	s.Require().NoError(s.mockAssetType())
	req := model.AssetTypeDetail{
		Id:          nil,
		Name:        pointy.String("Test"),
		IsCash:      pointy.Bool(false),
		IsLiability: pointy.Bool(false),
		Sequence:    pointy.Int(0),
		IsActive:    pointy.Bool(true),
	}

	err := s.repo.Upsert(req, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(pointy.Bool(true), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(3, len(res))
	s.Require().Equal("Test", *res[2].Name)
	s.Require().Equal(5, *res[2].Id)
	s.Require().Equal(true, *res[2].IsActive)
	s.Require().Equal(false, *res[2].IsCash)
	s.Require().Equal(false, *res[2].IsLiability)
	s.Require().Equal(0, *res[2].Sequence)

	res, err = s.repo.Get(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(5, len(res))
}

func (s *AssetTypeSuite) TestUpsert_Update() {
	s.Require().NoError(s.mockAssetType())
	req := model.AssetTypeDetail{
		Id:          pointy.Int(2),
		Name:        pointy.String("Test"),
		IsCash:      pointy.Bool(false),
		IsLiability: pointy.Bool(false),
		Sequence:    pointy.Int(0),
		IsActive:    pointy.Bool(true),
	}

	err := s.repo.Upsert(req, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(pointy.Bool(true), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(3, len(res))
	s.Require().Equal("Test", *res[2].Name)
	s.Require().Equal(2, *res[2].Id)
	s.Require().Equal(true, *res[2].IsActive)
	s.Require().Equal(false, *res[2].IsCash)
	s.Require().Equal(false, *res[2].IsLiability)
	s.Require().Equal(0, *res[2].Sequence)

	res, err = s.repo.Get(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
}

func (s *AssetTypeSuite) TestDelete() {
	s.Require().NoError(s.mockAssetType())
	err := s.repo.Delete(1, s.tx)
	s.Require().NoError(err)
}

func (s *AssetTypeSuite) TestUpdateSequence() {
	s.Require().NoError(s.mockAssetType())
	req := model.SequenceDetail{
		Id:       1,
		Sequence: 99,
	}

	err := s.repo.UpdateSequence(req, s.tx)
	s.Require().NoError(err)

	res, err := s.repo.Get(nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(4, len(res))
	s.Require().Equal(1, *res[0].Id)
	s.Require().Equal(99, *res[0].Sequence)
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
