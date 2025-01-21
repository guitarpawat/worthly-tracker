package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/entity"
	"github.com/rickb777/date/v2"
	"gorm.io/gorm"
)

type ValueOffsetRepository struct {
	db *gorm.DB
}

var _ TxRepository[*ValueOffsetRepository] = (*ValueOffsetRepository)(nil)

func (r *ValueOffsetRepository) BeginTx(ctx context.Context) (*ValueOffsetRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &ValueOffsetRepository{db: tx}, &Tx{tx: tx}
}

func (r *ValueOffsetRepository) FindByEffectiveDateAndAssetId(ctx context.Context, effectiveDate *date.Date, assetId *int) (res []entity.ValueOffset, err error) {
	where := make(map[string]any)
	if effectiveDate != nil {
		where["EffectiveDate"] = effectiveDate
	}
	if assetId != nil {
		where["AssetId"] = assetId
	}

	err = r.db.WithContext(ctx).Preload("Asset").Order("EffectiveDate, Asset.Sequence, AssetId").Where(where).Select(&res).Error

	return res, err
}

func (r *ValueOffsetRepository) Upsert(ctx context.Context, offset entity.ValueOffset) error {
	return r.db.WithContext(ctx).Save(offset).Error
}

func (r *ValueOffsetRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(id).Error
}
