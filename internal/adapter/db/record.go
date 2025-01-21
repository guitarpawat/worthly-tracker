package db

import (
	"context"
	"github.com/guitarpawat/worthly-tracker/internal/entity"
	"github.com/rickb777/date/v2"
	"gorm.io/gorm"
)

type RecordRepository struct {
	db *gorm.DB
}

var _ TxRepository[*RecordRepository] = (*RecordRepository)(nil)

func (r *RecordRepository) BeginTx(ctx context.Context) (*RecordRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &RecordRepository{db: tx}, &Tx{tx: tx}
}

func (r *RecordRepository) FindByDate(ctx context.Context, current date.Date) (res []entity.Record, err error) {
	err = r.db.WithContext(ctx).Preload("Asset", "Asset.AssetType").Order("Asset.AssetType.Sequence, Asset.Sequence, Id").
		Where(&entity.Record{Date: current}).Find(&res).Error

	return res, err
}

func (r *RecordRepository) FindForDraft(ctx context.Context) (res []entity.Record, err error) {
	latest, err := r.GetLatestDate(ctx)
	if err != nil {
		return nil, err
	}

	where := map[string]any{
		"Date":                     latest,
		"Asset.IsActive":           true,
		"Asset.AssetType.IsActive": true,
	}

	err = r.db.WithContext(ctx).Preload("Asset", "Asset.AssetType").Order("Asset.AssetType.Sequence, Asset.Sequence, Id").
		Where(where).Find(&res).Error

	return res, err
}

func (r *RecordRepository) Upsert(ctx context.Context, record entity.Record) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *RecordRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(id).Error
}

func (r *RecordRepository) DeleteByDate(ctx context.Context, date date.Date) error {
	return r.db.WithContext(ctx).Delete(&entity.Record{Date: date}).Error
}

func (r *RecordRepository) GetLatestDate(ctx context.Context) (res date.Date, err error) {
	err = r.db.WithContext(ctx).Order("Date").Limit(1).Select("Date").Scan(&res).Error
	return res, err
}

func (r *RecordRepository) FindPastAndFutureDate(ctx context.Context, current date.Date) (res []date.Date, err error) {
	var past []date.Date

	err = r.db.WithContext(ctx).Order("Date").Limit(12).Where("Date < ?", current).Pluck("Date", &past).Error
	if err != nil {
		return nil, err
	}

	var future []date.Date
	err = r.db.WithContext(ctx).Order("Date").Limit(12).Where("Date > ?", current).Pluck("Date", &future).Error
	if err != nil {
		return nil, err
	}

	return append(append(past, current), future...), nil
}
