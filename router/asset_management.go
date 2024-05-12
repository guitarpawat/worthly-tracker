package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.openly.dev/pointy"
	"net/http"
	"strconv"
	"strings"
	"worthly-tracker/db"
	"worthly-tracker/model"
	"worthly-tracker/ports"
)

func assetsManagementRouter(api *echo.Group) {
	as := assetManagementService{
		assetRepo:     db.GetAssetRepo(),
		assetTypeRepo: db.GetAssetTypeRepo(),
		dbConn:        db.GetDB(),
	}

	api.GET("/asset_types", as.getAssetTypes)
	api.POST("/asset_types", as.updateAssetType)
	api.PUT("/asset_types", as.addAssetType)
	api.DELETE("/asset_types", as.deleteAssetType)
	api.GET("/assets", as.getAssets)
	api.POST("/assets", as.updateAsset)
	api.PUT("/assets", as.addAsset)
	api.DELETE("/assets", as.deleteAsset)
}

type assetManagementService struct {
	assetRepo     ports.AssetRepo
	assetTypeRepo ports.AssetTypeRepo
	dbConn        ports.Connection
}

func (a assetManagementService) getAssetTypes(c echo.Context) error {
	isActiveParam := c.QueryParam("is_active")
	var isActive *bool
	if strings.ToLower(isActiveParam) == "true" {
		isActive = pointy.Bool(true)
	} else if strings.ToLower(isActiveParam) == "false" {
		isActive = pointy.Bool(false)
	} else if isActiveParam != "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid isActive param: %s", isActiveParam))
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	assetTypeDetail, err := a.assetTypeRepo.Get(isActive, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot fetch asset type: %w", err))
	}

	if len(assetTypeDetail) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("no asset type in database"))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return c.JSON(http.StatusOK, assetTypeDetail)
}

func (a assetManagementService) updateAssetType(c echo.Context) error {
	var body model.AssetTypeDetail
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	if err = a.upsertAssetType(body, true); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (a assetManagementService) addAssetType(c echo.Context) error {
	var body model.AssetTypeDetail
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	if err = a.upsertAssetType(body, false); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (a assetManagementService) upsertAssetType(body model.AssetTypeDetail, requireId bool) error {
	if requireId && body.Id == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("asset type id is required"))
	} else if !requireId && body.Id != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("asset type id must be null"))
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	if err = a.assetTypeRepo.Upsert(body, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot update asset type: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return nil
}

func (a assetManagementService) deleteAssetType(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse id: %w", err))
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	if err = a.assetTypeRepo.Delete(id, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot delete asset type: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return c.JSON(http.StatusOK, nil)
}

func (a assetManagementService) getAssets(c echo.Context) error {
	isActiveParam := c.QueryParam("is_active")
	var isActive *bool
	if strings.ToLower(isActiveParam) == "true" {
		isActive = pointy.Bool(true)
	} else if strings.ToLower(isActiveParam) == "false" {
		isActive = pointy.Bool(false)
	} else if isActiveParam != "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid isActive param: %s", isActiveParam))
	}

	typeIdParam := c.QueryParam("type_id")
	var typeId *int
	if typeIdParam != "" {
		id, err := strconv.Atoi(typeIdParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse id: %w", err))
		}
		typeId = &id
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	assetDetail, err := a.assetRepo.Get(isActive, typeId, tx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot fetch asset: %w", err))
	}

	if len(assetDetail) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("no asset in database"))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return c.JSON(http.StatusOK, assetDetail)
}

func (a assetManagementService) updateAsset(c echo.Context) error {
	var body model.AssetDetail
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	if err = a.upsertAsset(body, true); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (a assetManagementService) addAsset(c echo.Context) error {
	var body model.AssetDetail
	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot bind request body: %w", err))
	}

	if err = a.upsertAsset(body, false); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (a assetManagementService) upsertAsset(body model.AssetDetail, requireId bool) error {
	if requireId && body.Id == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("asset id is required"))
	} else if !requireId && body.Id != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("asset id must be null"))
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	if err = a.assetRepo.Upsert(body, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot update asset: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return nil
}

func (a assetManagementService) deleteAsset(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("cannot parse id: %w", err))
	}

	tx, err := a.dbConn.BeginTx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot create db transaction: %w", err))
	}
	defer tx.Rollback()

	if err = a.assetRepo.Delete(id, tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot delete asset: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot commit db transaction: %w", err))
	}

	return c.JSON(http.StatusOK, nil)
}
