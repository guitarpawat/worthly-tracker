package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"gorm.io/gorm"
)

type AssetsRepository struct {
	db *gorm.DB
}

var _ TxRepository[*AssetsRepository] = (*AssetsRepository)(nil)

func (r *AssetsRepository) BeginTx(ctx context.Context) (*AssetsRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &AssetsRepository{db: tx}, &Tx{tx: tx}
}

func (r *AssetsRepository) FindByIsActiveAndTypeId(ctx context.Context, isActive *bool, assetTypeId *int) (res []model.Asset, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}
	if assetTypeId != nil {
		where["AssetTypeId"] = *assetTypeId
	}

	err = r.db.WithContext(ctx).Preload("AssetType").Order("AssetType.Sequence, Sequence").Where(where).Find(&res).Error

	return res, err
}

func (r *AssetsRepository) FindNames(ctx context.Context, isActive *bool, assetTypeId *int) (res string, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}
	if assetTypeId != nil {
		where["AssetTypeId"] = *assetTypeId
	}

	err = r.db.WithContext(ctx).Order("Name").Model(new(model.Asset)).Where(where).Pluck("name", &res).Error
	return res, err
}

func (r *AssetsRepository) Upsert(ctx context.Context, asset model.Asset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *AssetsRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(new(model.Asset), id).Error
}
