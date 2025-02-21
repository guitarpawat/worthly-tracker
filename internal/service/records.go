package service

import (
	"context"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/rickb777/date/v2"
)

type Records struct {
	repo *db.RecordsRepository
}

func NewRecords(records *db.RecordsRepository) *Records {
	return &Records{
		repo: records,
	}
}

func (r *Records) GetByDate(ctx context.Context, d date.Date) (model.GetRecordTableView, error) {
	var err error
	if date.Zero == d {
		d, err = r.repo.GetLatestDate(ctx)
		if err != nil {
			return model.GetRecordTableView{}, fmt.Errorf("cannot get latest date: %w", err)
		}
	}

	records, err := r.repo.FindByDate(ctx, d)
	if err != nil {
		return model.GetRecordTableView{}, fmt.Errorf("cannot get records: %w", err)
	}

	dateList, err := r.repo.FindPastAndFutureDate(ctx, d)
	if err != nil {
		return model.GetRecordTableView{}, fmt.Errorf("cannot get date list: %w", err)
	}

	var assetTypes []model.RecordTableViewAssetType

	for idx, record := range records {
		if idx == 0 || assetTypes[len(assetTypes)-1].Id != record.Asset.AssetType.Id {
			assetTypes = append(assetTypes, model.RecordTableViewAssetType{
				Id:      record.Asset.AssetType.Id,
				Name:    record.Asset.AssetType.Name,
				Records: make([]model.RecordTableViewData, 0),
			})
		}
		assetTypes[len(assetTypes)-1].Records = append(assetTypes[len(assetTypes)-1].Records, model.RecordTableViewData{
			Id:               record.Id,
			AssetId:          record.AssetId,
			Name:             record.Asset.Name,
			Broker:           record.Asset.Broker,
			BoughtValue:      record.BoughtValue,
			CurrentValue:     record.CurrentValue,
			RealizedValue:    record.RealizedValue,
			DefaultIncrement: record.Asset.DefaultIncrement,
			IsCash:           record.Asset.AssetType.IsCash,
			Note:             record.Note.String,
		})
	}

	return model.GetRecordTableView{
		Dates:      dateList,
		AssetTypes: assetTypes,
	}, nil
}

func (r *Records) GetDraft(ctx context.Context) (model.EditRecordTableView, error) {
	records, err := r.repo.FindForDraft(ctx)
	if err != nil {
		return model.EditRecordTableView{}, fmt.Errorf("cannot get record draft: %w", err)
	}

	var assetTypes []model.RecordTableViewAssetType
	for idx, record := range records {
		if idx == 0 || assetTypes[len(assetTypes)-1].Id != record.Asset.AssetType.Id {
			assetTypes = append(assetTypes, model.RecordTableViewAssetType{
				Id:      record.Asset.AssetType.Id,
				Name:    record.Asset.AssetType.Name,
				Records: make([]model.RecordTableViewData, 0),
			})
		}
		assetTypes[len(assetTypes)-1].Records = append(assetTypes[len(assetTypes)-1].Records, model.RecordTableViewData{
			Id:               record.Id,
			AssetId:          record.AssetId,
			Name:             record.Asset.Name,
			Broker:           record.Asset.Broker,
			BoughtValue:      record.BoughtValue,
			CurrentValue:     record.CurrentValue,
			RealizedValue:    record.RealizedValue,
			DefaultIncrement: record.Asset.DefaultIncrement,
			IsCash:           record.Asset.AssetType.IsCash,
			Note:             record.Note.String,
		})
	}

	return model.EditRecordTableView{
		AssetTypes: assetTypes,
		Date:       date.Today(),
		Action:     constant.EditActionNew,
	}, nil
}
