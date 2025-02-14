package db

import (
	"context"
	"database/sql"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/pkg/errors"
	"github.com/rickb777/date/v2"
	"gorm.io/gorm"
)

type RecordsRepository struct {
	db *gorm.DB
}

var _ TxRepository[*RecordsRepository] = (*RecordsRepository)(nil)

func (r *RecordsRepository) BeginTx(ctx context.Context) (*RecordsRepository, *Tx) {
	tx := r.db.WithContext(ctx).Begin()
	return &RecordsRepository{db: tx}, &Tx{tx: tx}
}

func (r *RecordsRepository) FindByDate(ctx context.Context, current date.Date) (res []model.Record, err error) {
	err = r.db.WithContext(ctx).Joins("Asset").Joins("Asset.AssetType").Order("Asset__AssetType.Sequence, Asset.Sequence, records.Id").
		Where(&model.Record{Date: current}).Find(&res).Error

	return res, err
}

func (r *RecordsRepository) FindForDraft(ctx context.Context) (res []model.Record, err error) {
	latest, err := r.GetLatestDate(ctx)
	if err != nil {
		return nil, err
	}

	where := map[string]any{
		"Date":                     latest,
		"Asset.IsActive":           true,
		"Asset.AssetType.IsActive": true,
	}

	err = r.db.WithContext(ctx).Joins("Asset").Joins("Asset.AssetType").Order("Asset__AssetType.Sequence, Asset.Sequence, records.Id").
		Where(where).Find(&res).Error

	return res, err
}

func (r *RecordsRepository) Upsert(ctx context.Context, record model.Record) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *RecordsRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(id).Error
}

func (r *RecordsRepository) DeleteByDate(ctx context.Context, date date.Date) error {
	return r.db.WithContext(ctx).Delete(&model.Record{Date: date}).Error
}

func (r *RecordsRepository) GetLatestDate(ctx context.Context) (res date.Date, err error) {
	err = r.db.Table("records").WithContext(ctx).Order("Date").Limit(1).Select("Date").Scan(&res).Error
	if err == nil && res == date.Zero {
		return res, sql.ErrNoRows
	}
	return res, err
}

func (r *RecordsRepository) FindPastAndFutureDate(ctx context.Context, current date.Date) (model.DateResult, error) {
	var exists bool
	err := r.db.WithContext(ctx).Table("records").Select("count(*) > 0").Where("Date = ?", current).Find(&exists).Error
	if err != nil {
		return model.DateResult{}, err
	}
	if !exists {
		return model.DateResult{}, errors.Errorf("date %s not found", current)
	}

	var past []date.Date

	err = r.db.WithContext(ctx).Table("records").Order("Date").Limit(12).Where("Date < ?", current).Pluck("Date", &past).Error
	if err != nil {
		return model.DateResult{}, err
	}

	var future []date.Date
	err = r.db.WithContext(ctx).Table("records").Order("Date").Limit(12).Where("Date > ?", current).Pluck("Date", &future).Error
	if err != nil {
		return model.DateResult{}, err
	}

	return model.DateResult{
		PastDate:    past,
		CurrentDate: current,
		FutureDate:  future,
	}, nil
}
