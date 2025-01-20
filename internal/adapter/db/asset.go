package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/entity"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

var _ TxRepository[*AssetRepository] = (*AssetRepository)(nil)

func (r *AssetRepository) BeginTx(ctx context.Context) (*AssetRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &AssetRepository{db: tx}, &Tx{tx: tx}
}

func (r *AssetRepository) FindByIsActiveAndTypeId(ctx context.Context, isActive *bool, assetTypeId *int, preload bool) (res []entity.Asset, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}
	if assetTypeId != nil {
		where["AssetTypeId"] = *assetTypeId
	}

	if preload {
		err = r.db.WithContext(ctx).Preload("AssetType").Order("AssetType.Sequence, Sequence").Where(where).Find(&res).Error
	} else {
		err = r.db.WithContext(ctx).Preload("AssetType.Sequence").Order("AssetType.Sequence, Sequence").Where(where).Find(&res).Error
	}
	return res, err
}

func (r *AssetRepository) FindNames(ctx context.Context, isActive *bool, assetTypeId *int) (res string, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}
	if assetTypeId != nil {
		where["AssetTypeId"] = *assetTypeId
	}

	err = r.db.WithContext(ctx).Order("Name").Model(new(entity.Asset)).Where(where).Pluck("name", &res).Error
	return res, err
}

func (r *AssetRepository) Upsert(ctx context.Context, asset entity.Asset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *AssetRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(new(entity.Asset), id).Error
}
