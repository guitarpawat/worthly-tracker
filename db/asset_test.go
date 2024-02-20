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

func TestAssetSuite(t *testing.T) {
	suite.Run(t, new(AssetSuite))
}

type AssetSuite struct {
	suite.Suite
	tx   *sqlx.Tx
	repo ports.AssetRepo
}

func (s *AssetSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
	Init()
	s.repo = GetAssetRepo()
}

func (s *AssetSuite) SetupTest() {
	var err error
	s.tx, err = GetDB().BeginTx()
	s.Require().NoError(err)
}

func (s *AssetSuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *AssetSuite) TestGet_NotFound() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, pointy.Int(5), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(resp))
}

func (s *AssetSuite) TestGet_ByIsActive() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(pointy.Bool(true), nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(5, len(resp))
	s.Require().Equal(pointy.String("d2"), resp[0].Name)
	s.Require().Equal(pointy.Bool(true), resp[0].IsActive)
	s.Require().Equal(pointy.String("d1"), resp[1].Name)
	s.Require().Equal(pointy.Bool(true), resp[1].IsActive)
	s.Require().Equal(pointy.String("a1"), resp[2].Name)
	s.Require().Equal(pointy.Bool(true), resp[2].IsActive)
	s.Require().Equal(pointy.String("a2"), resp[3].Name)
	s.Require().Equal(pointy.Bool(true), resp[3].IsActive)
	s.Require().Equal(pointy.String("b1"), resp[4].Name)
	s.Require().Equal(pointy.Bool(true), resp[4].IsActive)
}

func (s *AssetSuite) TestGet_ByTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, pointy.Int(3), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(resp))
	s.Require().Equal(pointy.String("c1"), resp[0].Name)
	s.Require().Equal(pointy.Bool(false), resp[0].IsActive)
	s.Require().Equal(pointy.String("c2"), resp[1].Name)
	s.Require().Equal(pointy.Bool(false), resp[1].IsActive)
}

func (s *AssetSuite) TestGet_ByIsActiveAndTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(pointy.Bool(false), pointy.Int(2), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(resp))
	s.Require().Equal(pointy.String("b2"), resp[0].Name)
	s.Require().Equal(pointy.Bool(false), resp[0].IsActive)
}

func (s *AssetSuite) TestGet_ByAll() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.Get(nil, nil, s.tx)
	s.Require().NoError(err)
	s.Require().Equal(8, len(resp))
	s.Require().Equal(pointy.String("c1"), resp[0].Name)
	s.Require().Equal(pointy.String("c2"), resp[1].Name)
	s.Require().Equal(pointy.String("d2"), resp[2].Name)
	s.Require().Equal(pointy.String("d1"), resp[3].Name)
	s.Require().Equal(pointy.String("a1"), resp[4].Name)
	s.Require().Equal(pointy.String("a2"), resp[5].Name)
	s.Require().Equal(pointy.String("b1"), resp[6].Name)
	s.Require().Equal(pointy.String("b2"), resp[7].Name)
}

func (s *AssetSuite) TestGetNames_NotFound() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(nil, pointy.Int(5), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(resp))
}

func (s *AssetSuite) TestGetNames_ByIsActive() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(pointy.Bool(true), nil, s.tx)
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
	resp, err := s.repo.GetNames(nil, pointy.Int(3), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(resp))
	s.Require().Equal("c1", resp[0].Name)
	s.Require().Equal(5, resp[0].Id)
	s.Require().Equal("c2", resp[1].Name)
	s.Require().Equal(6, resp[1].Id)
}

func (s *AssetSuite) TestGetNames_ByIsActiveAndTypeId() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(pointy.Bool(false), pointy.Int(2), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(resp))
	s.Require().Equal("b2", resp[0].Name)
	s.Require().Equal(4, resp[0].Id)
}

func (s *AssetSuite) TestGetNames_ByAll() {
	s.Require().NoError(s.mockAssets())
	resp, err := s.repo.GetNames(nil, nil, s.tx)
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
		Id:               nil,
		Name:             pointy.String("Monika"),
		Broker:           pointy.String("ddlc"),
		TypeId:           pointy.Int(5),
		TypeName:         nil,
		DefaultIncrement: pointy.Pointer(decimal.NewFromInt(1000)),
		Sequence:         nil,
		IsActive:         pointy.Bool(true),
	}
	s.Require().NoError(s.repo.Upsert(assetDetail, s.tx))

	res, err := s.repo.Get(pointy.Bool(true), pointy.Int(5), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(9, *res[0].Id)
	s.Require().Equal("Monika", *res[0].Name)
	s.Require().Equal("ddlc", *res[0].Broker)
	s.Require().Equal(5, *res[0].TypeId)
	s.Require().Equal("Just Monika", *res[0].TypeName)
	s.Require().True(decimal.NewFromInt(1000).Equal(*res[0].DefaultIncrement))
	s.Require().Equal(0, *res[0].Sequence)
	s.Require().True(*res[0].IsActive)
}

func (s *AssetSuite) TestUpsert_Update() {
	s.Require().NoError(s.mockAssets())
	var assetDetail = model.AssetDetail{
		Id:               pointy.Int(3),
		Name:             pointy.String("Monika"),
		Broker:           pointy.String("ddlc"),
		TypeId:           pointy.Int(5),
		TypeName:         nil,
		DefaultIncrement: pointy.Pointer(decimal.NewFromInt(1000)),
		Sequence:         nil,
		IsActive:         pointy.Bool(true),
	}
	s.Require().NoError(s.repo.Upsert(assetDetail, s.tx))

	res, err := s.repo.Get(pointy.Bool(true), pointy.Int(5), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(3, *res[0].Id)
	s.Require().Equal("Monika", *res[0].Name)
	s.Require().Equal("ddlc", *res[0].Broker)
	s.Require().Equal(5, *res[0].TypeId)
	s.Require().Equal("Just Monika", *res[0].TypeName)
	s.Require().True(decimal.NewFromInt(1000).Equal(*res[0].DefaultIncrement))
	s.Require().Equal(0, *res[0].Sequence)
	s.Require().True(*res[0].IsActive)

	res, err = s.repo.Get(pointy.Bool(true), pointy.Int(2), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(res))
}

func (s *AssetSuite) TestUpsert_Delete() {
	s.Require().NoError(s.mockAssets())
	s.Require().NoError(s.repo.Delete(3, s.tx))

	res, err := s.repo.Get(pointy.Bool(true), pointy.Int(2), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(res))
}

func (s *AssetSuite) TestUpdateSequence() {
	s.Require().NoError(s.mockAssets())
	var sequence = model.SequenceDetail{
		Id:       3,
		Sequence: 99,
	}
	s.Require().NoError(s.repo.UpdateSequence(sequence, s.tx))

	res, err := s.repo.Get(pointy.Bool(true), pointy.Int(2), s.tx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(res))
	s.Require().Equal(3, *res[0].Id)
	s.Require().Equal(99, *res[0].Sequence)
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
