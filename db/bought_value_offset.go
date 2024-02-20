package db

import (
	"database/sql"
	"errors"
	"github.com/go-sqlx/sqlx"
	"github.com/shopspring/decimal"
	"go.openly.dev/pointy"
	"worthly-tracker/model"
)

var boughtValueOffsetRepo = &SqliteOffsetRepo{}

func GetBoughtValueOffsetRepo() *SqliteOffsetRepo {
	return boughtValueOffsetRepo
}

type SqliteOffsetRepo struct {
}

func (s *SqliteOffsetRepo) Get(date model.Date, assetId int, tx *sqlx.Tx) (model.OffsetDetail, error) {
	rows := tx.QueryRowx("SELECT id, asset_id, effective_date, offset_price, note FROM bought_value_offsets WHERE asset_id = ? AND effective_date <= ? ORDER BY effective_date DESC LIMIT 1", assetId, date.String())

	offset := new(model.OffsetDetail)
	err := rows.Scan(&offset.Id, &offset.AssetId, &offset.EffectiveDate, &offset.OffsetPrice, &offset.Note)
	if errors.Is(err, sql.ErrNoRows) {
		return model.OffsetDetail{
			Id:            pointy.Int(0),
			AssetId:       assetId,
			EffectiveDate: date,
			OffsetPrice:   decimal.Zero,
			Note:          nil,
		}, nil
	}
	if err != nil {
		return model.OffsetDetail{}, err
	}

	return *offset, nil
}

func (s *SqliteOffsetRepo) GetAllByAssetId(assetId int, tx *sqlx.Tx) ([]model.OffsetDetail, error) {
	rows, err := tx.Queryx("SELECT id, asset_id, effective_date, offset_price, note FROM bought_value_offsets WHERE asset_id = ?", assetId)
	if err != nil {
		return nil, err
	}

	var offsets []model.OffsetDetail
	for rows.Next() {
		offset := model.OffsetDetail{}
		err := rows.Scan(&offset.Id, &offset.AssetId, &offset.EffectiveDate, &offset.OffsetPrice, &offset.Note)
		if err != nil {
			return nil, err
		}
		offsets = append(offsets, offset)
	}

	return offsets, nil
}

func (s *SqliteOffsetRepo) GetAllByDate(date model.Date, tx *sqlx.Tx) ([]model.OffsetDetail, error) {
	rows, err := tx.Queryx(`SELECT id, asset_id, MAX(effective_date), offset_price, note FROM bought_value_offsets WHERE effective_date <= ? GROUP BY asset_id`, date.String())
	if err != nil {
		return nil, err
	}

	var offsets []model.OffsetDetail
	for rows.Next() {
		offset := model.OffsetDetail{}
		var effectiveDate string
		err := rows.Scan(&offset.Id, &offset.AssetId, &effectiveDate, &offset.OffsetPrice, &offset.Note)
		if err != nil {
			return nil, err
		}

		offset.EffectiveDate, err = model.NewDate(effectiveDate)
		if err != nil {
			return nil, err
		}

		offsets = append(offsets, offset)
	}

	return offsets, nil
}

func (s *SqliteOffsetRepo) Upsert(data model.OffsetDetail, tx *sqlx.Tx) (err error) {
	if data.Id == nil {
		_, err = tx.Exec("INSERT INTO bought_value_offsets(asset_id, effective_date, offset_price, note) VALUES (?,?,?,?)",
			data.AssetId, data.EffectiveDate.String(), data.OffsetPrice, data.Note)
	} else {
		_, err = tx.Exec("UPDATE bought_value_offsets SET asset_id = ?, effective_date = ?, offset_price = ?, note = ? WHERE id = ?",
			data.AssetId, data.EffectiveDate.String(), data.OffsetPrice, data.Note, data.Id)
	}

	return
}

func (s *SqliteOffsetRepo) Delete(id int, tx *sqlx.Tx) (err error) {
	_, err = tx.Exec("DELETE FROM bought_value_offsets WHERE id = ?", id)
	return
}
