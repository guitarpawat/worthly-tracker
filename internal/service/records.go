package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/guitarpawat/worthly-tracker/internal/view"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
	"strconv"
)

type Records struct {
	repo *db.RecordsRepository
}

func NewRecords(records *db.RecordsRepository) *Records {
	return &Records{
		repo: records,
	}
}

type CreateRecordRequest struct {
	Date    date.Date                   `validate:"required"`
	Records []CreateRecordRequestRecord `validate:"required"`
}

type CreateRecordRequestRecord struct {
	AssetId       string `validate:"number"`
	BoughtValue   decimal.Decimal
	CurrentValue  decimal.Decimal
	RealizedValue decimal.Decimal
	Note          string
}

type UpdateRecordRequest struct {
	Date    date.Date                   `validate:"required"`
	Records []UpdateRecordRequestRecord `validate:"required"`
}

type UpdateRecordRequestRecord struct {
	Id            string `validate:"number"`
	AssetId       string `validate:"number"`
	BoughtValue   decimal.Decimal
	CurrentValue  decimal.Decimal
	RealizedValue decimal.Decimal
	Note          string
}

func (r *Records) GetByDate(d date.Date) (view.GetRecordTable, error) {
	var err error
	if date.Zero == d {
		d, err = r.repo.GetLatestDate()
		if err != nil {
			return view.GetRecordTable{}, fmt.Errorf("cannot get latest date: %w", err)
		}
	}

	records, err := r.repo.GetRecordByDate(d)
	if err != nil {
		return view.GetRecordTable{}, fmt.Errorf("cannot get records: %w", err)
	}

	dateList, err := r.repo.GetDate(d)
	if err != nil {
		return view.GetRecordTable{}, fmt.Errorf("cannot get date list: %w", err)
	}

	var assetTypes []view.RecordTableAssetType
	for _, assetType := range records {
		assetTypes = append(assetTypes, view.RecordTableAssetType{
			Id:      int(assetType.Id.Int64),
			Name:    assetType.Name,
			Records: make([]view.RecordTableData, 0),
		})

		for _, record := range assetType.Assets {
			assetTypes[len(assetTypes)-1].Records = append(assetTypes[len(assetTypes)-1].Records, view.RecordTableData{
				Id:               int(record.Id.Int64),
				AssetId:          record.AssetId,
				Name:             record.Name,
				Broker:           record.Broker,
				BoughtValue:      nullDecimalToDecimal(record.BoughtValue),
				CurrentValue:     nullDecimalToDecimal(record.CurrentValue),
				RealizedValue:    nullDecimalToDecimal(record.RealizedValue),
				DefaultIncrement: nullDecimalToDecimal(record.DefaultIncrement),
				IsCash:           assetType.IsCash,
				Note:             record.Note.String,
			})
		}
	}

	var dates view.DateResult
	dates.CurrentDate = dateList.Current
	dates.FutureDate = dateList.Next
	dates.PastDate = dateList.Prev

	return view.GetRecordTable{
		Dates:      dates,
		AssetTypes: assetTypes,
	}, nil
}

func (r *Records) GetByDateForEdit(d date.Date) (view.EditRecordTable, error) {
	var err error
	if date.Zero == d {
		d, err = r.repo.GetLatestDate()
		if err != nil {
			return view.EditRecordTable{}, fmt.Errorf("cannot get latest date: %w", err)
		}
	}

	records, err := r.repo.GetRecordByDate(d)
	if err != nil {
		return view.EditRecordTable{}, fmt.Errorf("cannot get records: %w", err)
	}

	var assetTypes []view.RecordTableAssetType
	for _, assetType := range records {
		assetTypes = append(assetTypes, view.RecordTableAssetType{
			Id:      int(assetType.Id.Int64),
			Name:    assetType.Name,
			Records: make([]view.RecordTableData, 0),
		})

		for _, record := range assetType.Assets {
			assetTypes[len(assetTypes)-1].Records = append(assetTypes[len(assetTypes)-1].Records, view.RecordTableData{
				Id:               int(record.Id.Int64),
				AssetId:          record.AssetId,
				Name:             record.Name,
				Broker:           record.Broker,
				BoughtValue:      nullDecimalToDecimal(record.BoughtValue),
				CurrentValue:     nullDecimalToDecimal(record.CurrentValue),
				RealizedValue:    nullDecimalToDecimal(record.RealizedValue),
				DefaultIncrement: nullDecimalToDecimal(record.DefaultIncrement),
				IsCash:           assetType.IsCash,
				Note:             record.Note.String,
			})
		}
	}

	return view.EditRecordTable{
		AssetTypes: assetTypes,
		Date:       d,
		Action:     constant.EditActionUpdate,
	}, nil
}

func (r *Records) GetDraft() (view.EditRecordTable, error) {
	records, err := r.repo.GetRecordDraft()
	if err != nil {
		return view.EditRecordTable{}, fmt.Errorf("cannot get record draft: %w", err)
	}

	var assetTypes []view.RecordTableAssetType
	for _, assetType := range records {
		assetTypes = append(assetTypes, view.RecordTableAssetType{
			Id:      int(assetType.Id.Int64),
			Name:    assetType.Name,
			Records: make([]view.RecordTableData, 0),
		})

		for _, record := range assetType.Assets {
			assetTypes[len(assetTypes)-1].Records = append(assetTypes[len(assetTypes)-1].Records, view.RecordTableData{
				Id:               int(record.Id.Int64),
				AssetId:          record.AssetId,
				Name:             record.Name,
				Broker:           record.Broker,
				BoughtValue:      nullDecimalToDecimal(record.BoughtValue),
				CurrentValue:     nullDecimalToDecimal(record.CurrentValue),
				RealizedValue:    nullDecimalToDecimal(record.RealizedValue),
				DefaultIncrement: nullDecimalToDecimal(record.DefaultIncrement),
				IsCash:           assetType.IsCash,
				Note:             record.Note.String,
			})
		}
	}

	return view.EditRecordTable{
		AssetTypes: assetTypes,
		Date:       date.Today(),
		Action:     constant.EditActionNew,
	}, nil
}

func (r *Records) CreateRecords(req CreateRecordRequest) error {
	repo, tx, err := r.repo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var assetRecords []db.AssetRecord
	for _, record := range req.Records {
		assetId, _ := strconv.Atoi(record.AssetId) // Already validate request
		assetRecords = append(assetRecords, db.AssetRecord{
			Id:            sql.NullInt64{Valid: false},
			AssetId:       assetId,
			BoughtValue:   decimal.NewNullDecimal(record.BoughtValue),
			CurrentValue:  decimal.NewNullDecimal(record.CurrentValue),
			RealizedValue: decimal.NewNullDecimal(record.RealizedValue),
			Note:          stringToNullString(record.Note),
		})
	}

	for _, record := range assetRecords {
		err := repo.UpsertRecord(record, req.Date)
		if err != nil {
			return fmt.Errorf("cannot insert record: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (r *Records) UpdateRecords(ctx context.Context, req UpdateRecordRequest) error {
	repo, tx, err := r.repo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var assetRecords []db.AssetRecord
	for _, record := range req.Records {
		recordId, _ := strconv.Atoi(record.Id)     // Already validate request
		assetId, _ := strconv.Atoi(record.AssetId) // Already validate request
		assetRecords = append(assetRecords, db.AssetRecord{
			Id:            sql.NullInt64{Int64: int64(recordId), Valid: true},
			AssetId:       assetId,
			BoughtValue:   decimal.NewNullDecimal(record.BoughtValue),
			CurrentValue:  decimal.NewNullDecimal(record.CurrentValue),
			RealizedValue: decimal.NewNullDecimal(record.RealizedValue),
			Note:          stringToNullString(record.Note),
		})
	}

	for _, record := range assetRecords {
		err := repo.UpsertRecord(record, req.Date)
		if err != nil {
			return fmt.Errorf("cannot insert record: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (r *Records) DeleteRecords(d date.Date) error {
	_, err := r.repo.DeleteRecordByDate(d)
	return err
}

func nullDecimalToDecimal(d decimal.NullDecimal) decimal.Decimal {
	if d.Valid {
		return d.Decimal
	}
	return decimal.Zero
}

func stringToNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
