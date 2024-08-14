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
	api.DELETE("/asset_types/:id", as.deleteAssetType)
	api.GET("/assets", as.getAssets)
	api.POST("/assets", as.updateAsset)
	api.PUT("/assets", as.addAsset)
	api.DELETE("/assets/:id", as.deleteAsset)
}

type assetManagementService struct {
	assetRepo     ports.AssetRepo
	assetTypeRepo ports.AssetTypeRepo
	dbConn        ports.Connection
}

//	@Summary		Get asset types
//	@Tags			asset_management
//	@Description	Get asset types by filtered is_active, or get all asset types if not specified
//	@Param			is_active	query	string	false	"Filter asset types by is_active"
//	@Produce		json
//	@Success		200	{array}		model.AssetTypeDetail
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		404	{object}	nil	"No any asset types found"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/asset_types [get]
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

//	@Summary		Update asset type
//	@Tags			asset_management
//	@Description	Update asset type data to database
//	@Accept			json
//	@Param			request	body	model.AssetTypeDetail	true	"Asset type data for update"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to update asset type"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/asset_types [post]
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

//	@Summary		Add new asset type
//	@Tags			asset_management
//	@Description	Insert new asset type to database
//	@Accept			json
//	@Param			request	body	model.AssetTypeDetail	true	"Asset type data for insert"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to insert asset type"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/asset_types [put]
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

//	@Summary		Delete asset type
//	@Tags			asset_management
//	@Description	Delete asset type from database
//	@Param			id	path	int	true	"Asset type id"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to delete asset type"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/asset_types/{id} [delete]
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

//	@Summary		Get assets
//	@Tags			asset_management
//	@Description	Get asset data from database filtered by is_active and type_id, or get all if not specified
//	@Param			is_active	query	string	false	"Filter asset types by is_active"
//	@Param			type_id		query	string	false	"Filter asset types by type_id"
//	@Produce		json
//	@Success		200	{object}	model.AssetDetail	"Success to get assets"
//	@Failure		400	{object}	nil					"Input validation failed"
//	@Failure		404	{object}	nil					"Asset not found"
//	@Failure		500	{object}	nil					"Generic server error"
//	@Router			/api/asset [get]
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

//	@Summary		Update asset
//	@Tags			asset_management
//	@Description	Update asset data to database
//	@Accept			json
//	@Param			request	body	model.AssetDetail	true	"Asset data for update"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to update asset"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/assets [post]
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

//	@Summary		Add new asset
//	@Tags			asset_management
//	@Description	Insert new asset to database
//	@Accept			json
//	@Param			request	body	model.AssetDetail	true	"Asset data for insert"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to insert asset"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/assets [put]
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

//	@Summary		Delete asset
//	@Tags			asset_management
//	@Description	Delete asset from database
//	@Param			id	path	int	true	"Asset id"
//	@Produce		json
//	@Success		200	{object}	nil	"Success to delete asset"
//	@Failure		400	{object}	nil	"Input validation failed"
//	@Failure		500	{object}	nil	"Generic server error"
//	@Router			/api/assets/{id} [delete]
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
