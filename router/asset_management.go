package router

import (
	"github.com/labstack/echo/v4"
	"worthly-tracker/db"
	"worthly-tracker/ports"
)

func assetsManagementRouter(api *echo.Group) {
	as := assetManagementService{
		assetRepo:     db.GetAssetRepo(),
		assetTypeRepo: db.GetAssetTypeRepo(),
		offsetRepo:    db.GetBoughtValueOffsetRepo(),
		dbConn:        db.GetDB(),
	}

	api.GET("/asset_types", as.getAssetTypes)
	api.POST("/asset_types", as.updateAssetTypes)
	api.PUT("/asset_types", as.addAssetTypes)
	api.DELETE("/asset_types", as.deleteAssetTypes)
	api.GET("/assets", as.getAssets)
	api.POST("/assets", as.updateAssets)
	api.PUT("/assets", as.addAssets)
	api.DELETE("/assets", as.deleteAssets)
}

type assetManagementService struct {
	assetRepo     ports.AssetRepo
	assetTypeRepo ports.AssetTypeRepo
	offsetRepo    ports.BoughtValueOffsetRepo
	dbConn        ports.Connection
}

func (a assetManagementService) getAssetTypes(c echo.Context) error {
	return nil
}

func (a assetManagementService) updateAssetTypes(c echo.Context) error {
	return nil
}

func (a assetManagementService) addAssetTypes(c echo.Context) error {
	return nil
}

func (a assetManagementService) deleteAssetTypes(c echo.Context) error {
	return nil
}

func (a assetManagementService) getAssets(c echo.Context) error {
	return nil
}

func (a assetManagementService) updateAssets(c echo.Context) error {
	return nil
}

func (a assetManagementService) addAssets(c echo.Context) error {
	return nil
}

func (a assetManagementService) deleteAssets(c echo.Context) error {
	return nil
}
