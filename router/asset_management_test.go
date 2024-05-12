package router

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/db"
	"worthly-tracker/logs"
	"worthly-tracker/mocks"
)

func TestAssetManagementSuite(t *testing.T) {
	suite.Run(t, new(AssetManagementSuite))
}

type AssetManagementSuite struct {
	suite.Suite

	service       assetManagementService
	assetTypeRepo *mocks.MockAssetTypeRepo
	assetRepo     *mocks.MockAssetRepo
	dbMock        sqlmock.Sqlmock
}

func (a *AssetManagementSuite) SetupSuite() {
	logs.Init()
	config.InitTest()
}

func (a *AssetManagementSuite) SetupTest() {
	a.assetTypeRepo = mocks.NewMockAssetTypeRepo(a.T())
	a.assetRepo = mocks.NewMockAssetRepo(a.T())
	a.dbMock = db.InitMock()
}

func (a *AssetManagementSuite) TestGetAssetTypes_400_InvalidIsActiveParam() {

}

func (a *AssetManagementSuite) TestGetAssetTypes_404_NoAssetTypeInDatabase() {

}

func (a *AssetManagementSuite) TestGetAssetTypes_200_NoIsActiveParam() {

}

func (a *AssetManagementSuite) TestGetAssetTypes_200_WithIsActiveParam() {

}

func (a *AssetManagementSuite) TestUpdateAssetType_400_MalformedBody() {

}

func (a *AssetManagementSuite) TestUpdateAssetType_400_NoId() {

}

func (a *AssetManagementSuite) TestUpdateAssetType_200_Success() {

}

func (a *AssetManagementSuite) TestAddAssetType_400_MalformedBody() {

}

func (a *AssetManagementSuite) TestAddAssetType_400_NoId() {

}

func (a *AssetManagementSuite) TestAddAssetType_200_Success() {

}

func (a *AssetManagementSuite) TestDeleteAssetType_400_IdIsNotInteger() {

}

func (a *AssetManagementSuite) TestDeleteAssetType_200_Success() {

}

func (a *AssetManagementSuite) TestGetAssets_400_InvalidIsActiveParam() {

}

func (a *AssetManagementSuite) TestGetAssets_400_InvalidTypeIdParam() {

}

func (a *AssetManagementSuite) TestGetAssets_404_NoAssetInDatabase() {

}

func (a *AssetManagementSuite) TestGetAssets_200_NoIsActiveParam() {

}

func (a *AssetManagementSuite) TestGetAssets_200_NoTypeIdParam() {

}

func (a *AssetManagementSuite) TestGetAssets_200_WithAllParams() {

}

func (a *AssetManagementSuite) TestUpdateAsset_400_MalformedBody() {

}

func (a *AssetManagementSuite) TestUpdateAsset_400_NoId() {

}

func (a *AssetManagementSuite) TestUpdateAsset_200_Success() {

}

func (a *AssetManagementSuite) TestAddAsset_400_MalformedBody() {

}

func (a *AssetManagementSuite) TestAddAsset_400_NoId() {

}

func (a *AssetManagementSuite) TestAddAsset_200_Success() {

}

func (a *AssetManagementSuite) TestDeleteAsset_400_IdIsNotInteger() {

}

func (a *AssetManagementSuite) TestDeleteAsset_200_Success() {

}
