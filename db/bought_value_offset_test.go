//go:build test || integration

package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/logs"
	"worthly-tracker/ports"
)

func TestBoughtValueOffsetTestSuite(t *testing.T) {
	suite.Run(t, new(BoughtValueOffsetTestSuite))
}

type BoughtValueOffsetTestSuite struct {
	suite.Suite
	tx   *sqlx.Tx
	repo ports.AssetTypeRepo
}

func (s *BoughtValueOffsetTestSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
	Init()
	s.repo = GetAssetTypeRepo()
}

func (s *BoughtValueOffsetTestSuite) SetupTest() {
	var err error
	s.tx, err = GetDB().BeginTx()
	s.Require().NoError(err)
}

func (s *BoughtValueOffsetTestSuite) TearDownTest() {
	err := s.tx.Rollback()
	s.Require().NoError(err)
}

func (s *BoughtValueOffsetTestSuite) TestGet_Empty() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestGet_Success() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByAssetId_Empty() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByAssetId_Success() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByDate_Empty() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestGetAllByDateSuccess() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestUpsert_Insert() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestUpsert_Update() {
	// TODO: Implement this test
}

func (s *BoughtValueOffsetTestSuite) TestDelete() {
	// TODO: Implement this test
}
