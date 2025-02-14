package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"gorm.io/gorm"
)

type AssetTypesRepository struct {
	db *gorm.DB
}

var _ TxRepository[*AssetTypesRepository] = (*AssetTypesRepository)(nil)

func (r *AssetTypesRepository) BeginTx(ctx context.Context) (*AssetTypesRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &AssetTypesRepository{db: tx}, &Tx{tx: tx}
}

func (r *AssetTypesRepository) FindByIsActive(ctx context.Context, isActive *bool) (res []model.AssetType, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = isActive
	}

	err = r.db.WithContext(ctx).Preload("Assets").Order("Sequence").Where(where).Find(&res).Error

	return res, err
}

func (r *AssetTypesRepository) FindNames(ctx context.Context, isActive *bool) (res string, err error) {
	where := make(map[string]any)
	if isActive != nil {
		where["IsActive"] = *isActive
	}

	err = r.db.WithContext(ctx).Order("Name").Model(new(model.AssetType)).Where(where).Pluck("name", &res).Error
	return res, err
}

func (r *AssetTypesRepository) Upsert(ctx context.Context, asset model.AssetType) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *AssetTypesRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(new(model.AssetType), id).Error
}
