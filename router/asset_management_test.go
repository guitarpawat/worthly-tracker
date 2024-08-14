//go:build test || unit

package router

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sqlx/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.openly.dev/pointy"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"worthly-tracker/config"
	"worthly-tracker/db"
	"worthly-tracker/logs"
	"worthly-tracker/mocks"
	"worthly-tracker/model"
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
	a.service = assetManagementService{
		assetTypeRepo: a.assetTypeRepo,
		assetRepo:     a.assetRepo,
		dbConn:        db.GetDB(),
	}
}

func (a *AssetManagementSuite) TestGetAssetTypes_400_InvalidIsActiveParam() {
	q := make(url.Values)
	q.Set("is_active", "invalid")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	err := a.service.getAssetTypes(c)
	a.Require().ErrorContains(err, "invalid isActive param: invalid")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestGetAssetTypes_404_NoAssetTypeInDatabase() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectRollback()
	a.assetTypeRepo.EXPECT().Get((*bool)(nil), mock.Anything).Return(nil, nil)

	err := a.service.getAssetTypes(c)
	a.Require().ErrorContains(err, "no asset type in database")
	a.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestGetAssetTypes_200_NoIsActiveParam() {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()

	res := []model.AssetTypeDetail{
		{
			Id:          pointy.Int(1),
			Name:        pointy.String("Stock"),
			IsCash:      pointy.Bool(false),
			IsLiability: pointy.Bool(false),
			Sequence:    pointy.Int(1),
			IsActive:    pointy.Bool(true),
		},
	}
	a.assetTypeRepo.EXPECT().Get((*bool)(nil), mock.Anything).Return(res, nil)

	err := a.service.getAssetTypes(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
	var bodyParsed []model.AssetTypeDetail
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	a.Require().NoError(err)
	a.Require().Equal(res, bodyParsed)
}

func (a *AssetManagementSuite) TestGetAssetTypes_200_WithIsActiveParam() {
	q := make(url.Values)
	q.Set("is_active", "true")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()

	res := []model.AssetTypeDetail{
		{
			Id:          pointy.Int(1),
			Name:        pointy.String("Stock"),
			IsCash:      pointy.Bool(false),
			IsLiability: pointy.Bool(false),
			Sequence:    pointy.Int(1),
			IsActive:    pointy.Bool(true),
		},
	}
	a.assetTypeRepo.EXPECT().Get(pointy.Bool(true), mock.Anything).Return(res, nil)

	err := a.service.getAssetTypes(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
	var bodyParsed []model.AssetTypeDetail
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	a.Require().NoError(err)
	a.Require().Equal(res, bodyParsed)
}

func (a *AssetManagementSuite) TestUpdateAssetType_400_MalformedBody() {
	body := `{"id": "1", "name": "Stock", "is_cash": false, "is_liability": false, "sequence": "1", "is_active": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	err := a.service.updateAssetType(c)
	a.Require().ErrorContains(err, "cannot bind request body")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestUpdateAssetType_400_NoId() {
	body := `{"name": "Stock", "is_cash": false, "is_liability": false, "sequence": 1, "is_active": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	err := a.service.updateAssetType(c)
	a.Require().ErrorContains(err, "asset type id is required")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestUpdateAssetType_200_Success() {
	body := `{"id": 1, "name": "Stock", "isCash": false, "isLiability": false, "sequence": 1, "isActive": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetTypeRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Run(func(assetType model.AssetTypeDetail, _ *sqlx.Tx) {
		a.Require().Equal(1, *assetType.Id)
		a.Require().Equal("Stock", *assetType.Name)
		a.Require().Equal(false, *assetType.IsCash)
		a.Require().Equal(false, *assetType.IsLiability)
		a.Require().Equal(1, *assetType.Sequence)
		a.Require().Equal(true, *assetType.IsActive)
	}).Return(nil)

	err := a.service.updateAssetType(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}

func (a *AssetManagementSuite) TestAddAssetType_400_MalformedBody() {
	body := `{"name": "Stock", "is_cash": false, "is_liability": false, "sequence": "1", "is_active": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	err := a.service.addAssetType(c)
	a.Require().ErrorContains(err, "cannot bind request body")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestAddAssetType_400_HasId() {
	body := `{"id": 1, "name": "Stock", "is_cash": false, "is_liability": false, "sequence": 1, "is_active": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	err := a.service.addAssetType(c)
	a.Require().ErrorContains(err, "asset type id must be null")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestAddAssetType_200_Success() {
	body := `{"name": "Stock", "isCash": false, "isLiability": false, "sequence": 1, "isActive": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetTypeRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Run(func(assetType model.AssetTypeDetail, _ *sqlx.Tx) {
		a.Require().Nil(assetType.Id)
		a.Require().Equal("Stock", *assetType.Name)
		a.Require().Equal(false, *assetType.IsCash)
		a.Require().Equal(false, *assetType.IsLiability)
		a.Require().Equal(1, *assetType.Sequence)
		a.Require().Equal(true, *assetType.IsActive)
	}).Return(nil)

	err := a.service.addAssetType(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}

func (a *AssetManagementSuite) TestDeleteAssetType_400_IdIsNotInteger() {
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := a.service.deleteAssetType(c)
	a.Require().ErrorContains(err, "cannot parse id")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestDeleteAssetType_200_Success() {
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset_types/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetTypeRepo.EXPECT().Delete(1, mock.Anything).Return(nil)

	err := a.service.deleteAssetType(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}

func (a *AssetManagementSuite) TestGetAssets_400_InvalidIsActiveParam() {
	q := make(url.Values)
	q.Set("is_active", "invalid")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	err := a.service.getAssets(c)
	a.Require().ErrorContains(err, "invalid isActive param: invalid")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestGetAssets_400_InvalidTypeIdParam() {
	q := make(url.Values)
	q.Set("type_id", "abc")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	err := a.service.getAssets(c)
	a.Require().ErrorContains(err, "cannot parse id")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestGetAssets_404_NoAssetInDatabase() {
	q := make(url.Values)
	q.Set("is_active", "true")
	q.Set("type_id", "1")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectRollback()
	a.assetRepo.EXPECT().Get(pointy.Bool(true), pointy.Int(1), mock.Anything).Return(nil, nil)

	err := a.service.getAssets(c)
	a.Require().ErrorContains(err, "no asset in database")
	a.Require().Equal(http.StatusNotFound, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestGetAssets_200_NoIsActiveParam() {
	q := make(url.Values)
	q.Set("type_id", "1")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()

	res := []model.AssetDetail{
		{
			Id:               pointy.Int(1),
			Name:             pointy.String("ACN"),
			Broker:           pointy.String("SCBAM"),
			TypeId:           pointy.Int(1),
			TypeName:         pointy.String("Stock"),
			DefaultIncrement: nil,
			Sequence:         pointy.Int(1),
			IsActive:         pointy.Bool(true),
		},
	}
	a.assetRepo.EXPECT().Get((*bool)(nil), pointy.Int(1), mock.Anything).Return(res, nil)

	err := a.service.getAssets(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
	var bodyParsed []model.AssetDetail
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	a.Require().NoError(err)
	a.Require().Equal(res, bodyParsed)
}

func (a *AssetManagementSuite) TestGetAssets_200_NoTypeIdParam() {
	q := make(url.Values)
	q.Set("is_active", "true")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()

	res := []model.AssetDetail{
		{
			Id:               pointy.Int(1),
			Name:             pointy.String("ACN"),
			Broker:           pointy.String("SCBAM"),
			TypeId:           pointy.Int(1),
			TypeName:         pointy.String("Stock"),
			DefaultIncrement: nil,
			Sequence:         pointy.Int(1),
			IsActive:         pointy.Bool(true),
		},
	}
	a.assetRepo.EXPECT().Get(pointy.Bool(true), (*int)(nil), mock.Anything).Return(res, nil)

	err := a.service.getAssets(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
	var bodyParsed []model.AssetDetail
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	a.Require().NoError(err)
	a.Require().Equal(res, bodyParsed)
}

func (a *AssetManagementSuite) TestGetAssets_200_WithAllParams() {
	q := make(url.Values)
	q.Set("is_active", "true")
	q.Set("type_id", "1")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()

	res := []model.AssetDetail{
		{
			Id:               pointy.Int(1),
			Name:             pointy.String("ACN"),
			Broker:           pointy.String("SCBAM"),
			TypeId:           pointy.Int(1),
			TypeName:         pointy.String("Stock"),
			DefaultIncrement: nil,
			Sequence:         pointy.Int(1),
			IsActive:         pointy.Bool(true),
		},
	}
	a.assetRepo.EXPECT().Get(pointy.Bool(true), pointy.Int(1), mock.Anything).Return(res, nil)

	err := a.service.getAssets(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
	var bodyParsed []model.AssetDetail
	err = json.Unmarshal(w.Body.Bytes(), &bodyParsed)
	a.Require().NoError(err)
	a.Require().Equal(res, bodyParsed)
}

func (a *AssetManagementSuite) TestUpdateAsset_400_MalformedBody() {
	body := `{"id": "1", "name": "ACN", "broker": "SCBAM", "typeId": "1", "typeName": "Stock", "defaultIncrement": 0, "sequence": "1", "is_active": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	err := a.service.updateAsset(c)
	a.Require().ErrorContains(err, "cannot bind request body")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestUpdateAsset_400_NoId() {
	body := `{"name": "ACN", "broker": "SCBAM", "typeId": 1, "typeName": "Stock", "defaultIncrement": "0", "sequence": 1, "is_active": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	err := a.service.updateAsset(c)
	a.Require().ErrorContains(err, "asset id is required")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestUpdateAsset_200_Success() {
	body := `{"id": 1, "name": "ACN", "broker": "SCBAM", "typeId": 1, "typeName": "Stock", "defaultIncrement": null, "sequence": 1, "isActive": true}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Run(func(asset model.AssetDetail, _ *sqlx.Tx) {
		a.Require().Equal(1, *asset.Id)
		a.Require().Equal("ACN", *asset.Name)
		a.Require().Equal("SCBAM", *asset.Broker)
		a.Require().Equal(1, *asset.TypeId)
		a.Require().Equal("Stock", *asset.TypeName)
		a.Require().Nil(asset.DefaultIncrement)
		a.Require().Equal(1, *asset.Sequence)
		a.Require().Equal(true, *asset.IsActive)
	}).Return(nil)

	err := a.service.updateAsset(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}

func (a *AssetManagementSuite) TestAddAsset_400_MalformedBody() {
	body := `{"name": "ACN", "broker": "SCBAM", "typeId": "1", "typeName": "Stock", "defaultIncrement": 0, "sequence": "1", "is_active": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/asset")

	err := a.service.addAsset(c)
	a.Require().ErrorContains(err, "cannot bind request body")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestAddAsset_400_WithId() {
	body := `{"id": 1, "name": "ACN", "broker": "SCBAM", "typeId": 1, "typeName": "Stock", "defaultIncrement": "0", "sequence": 1, "is_active": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	err := a.service.addAsset(c)
	a.Require().ErrorContains(err, "asset id must be null")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestAddAsset_200_Success() {
	body := `{"name": "ACN", "broker": "SCBAM", "typeId": 1, "typeName": "Stock", "defaultIncrement": null, "sequence": 1, "isActive": true}`
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Run(func(asset model.AssetDetail, _ *sqlx.Tx) {
		a.Require().Nil(asset.Id)
		a.Require().Equal("ACN", *asset.Name)
		a.Require().Equal("SCBAM", *asset.Broker)
		a.Require().Equal(1, *asset.TypeId)
		a.Require().Equal("Stock", *asset.TypeName)
		a.Require().Nil(asset.DefaultIncrement)
		a.Require().Equal(1, *asset.Sequence)
		a.Require().Equal(true, *asset.IsActive)
	}).Return(nil)

	err := a.service.addAsset(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}

func (a *AssetManagementSuite) TestDeleteAsset_400_IdIsNotInteger() {
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := a.service.deleteAsset(c)
	a.Require().ErrorContains(err, "cannot parse id")
	a.Require().Equal(http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func (a *AssetManagementSuite) TestDeleteAsset_200_Success() {
	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	c.SetPath("/api/assets_management/assets/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	a.dbMock.ExpectBegin()
	a.dbMock.ExpectCommit()
	a.assetRepo.EXPECT().Delete(1, mock.Anything).Return(nil)

	err := a.service.deleteAsset(c)
	a.Require().NoError(err)
	a.Require().Equal(http.StatusOK, w.Code)
}
