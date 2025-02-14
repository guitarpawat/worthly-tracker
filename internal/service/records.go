package service

import (
	"context"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/rickb777/date/v2"
	"github.com/samber/lo"
)

type Records struct {
	repo *db.RecordsRepository
}

func NewRecords(records *db.RecordsRepository) *Records {
	return &Records{
		repo: records,
	}
}

func (r *Records) GetByDate(ctx context.Context, d date.Date) (model.GetRecordResponse, error) {
	var err error
	if date.Zero == d {
		d, err = r.repo.GetLatestDate(ctx)
		if err != nil {
			return model.GetRecordResponse{}, fmt.Errorf("cannot get latest date: %w", err)
		}
	}

	records, err := r.repo.FindByDate(ctx, d)
	if err != nil {
		return model.GetRecordResponse{}, fmt.Errorf("cannot get records: %w", err)
	}

	dateList, err := r.repo.FindPastAndFutureDate(ctx, d)
	if err != nil {
		return model.GetRecordResponse{}, fmt.Errorf("cannot get date list: %w", err)
	}

	mapByAssetType := lo.GroupBy(records, func(item model.Record) int {
		return item.Asset.AssetType.Id
	})

	var assetTypes = make([]model.GetRecordAssetType, 0, len(mapByAssetType))

	for _, record := range mapByAssetType {
		var assetType model.GetRecordAssetType
		assetType.Name = record[0].Asset.AssetType.Name
		assetType.Records = record
		assetTypes = append(assetTypes, assetType)
	}

	return model.GetRecordResponse{
		Dates:      dateList,
		AssetTypes: assetTypes,
	}, nil
}
