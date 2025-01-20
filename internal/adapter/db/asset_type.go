package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/entity"
	"gorm.io/gorm"
)

type AssetTypeRepository struct {
	db *gorm.DB
}

var _ TxRepository[*AssetTypeRepository] = (*AssetTypeRepository)(nil)

func (r *AssetTypeRepository) BeginTx(ctx context.Context) (*AssetTypeRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &AssetTypeRepository{db: tx}, &Tx{tx: tx}
}

func (r *AssetTypeRepository) FindByIsActive(ctx context.Context, isActive *bool, preload bool) (res []entity.AssetType, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = isActive
	}

	if preload {
		err = r.db.WithContext(ctx).Preload("Assets").Order("Sequence").Where(where).Find(&res).Error
	} else {
		err = r.db.WithContext(ctx).Order("Sequence").Where(where).Find(&res).Error
	}

	return res, err
}

func (r *AssetTypeRepository) FindNames(ctx context.Context, isActive *bool) (res string, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}

	err = r.db.WithContext(ctx).Order("Name").Model(new(entity.AssetType)).Where(where).Pluck("name", &res).Error
	return res, err
}

func (r *AssetTypeRepository) Upsert(ctx context.Context, asset entity.AssetType) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *AssetTypeRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(new(entity.AssetType), id).Error
}
