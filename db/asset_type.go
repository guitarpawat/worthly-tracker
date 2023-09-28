package db

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"worthly-tracker/model"
)

var assetTypeRepo = &SqliteAssetTypeRepo{}

func GetAssetTypeRepo() *SqliteAssetTypeRepo {
	return assetTypeRepo
}

type SqliteAssetTypeRepo struct {
}

func (s *SqliteAssetTypeRepo) Get(isActive *bool, tx *sqlx.Tx) ([]model.AssetTypeDetail, error) {
	rows, err := tx.Queryx("SELECT * FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY sequence", isActive, isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	details := make([]model.AssetTypeDetail, 0)
	for rows.Next() {
		detail := model.AssetTypeDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (s *SqliteAssetTypeRepo) GetNames(isActive *bool, tx *sqlx.Tx) ([]model.NameDetail, error) {
	rows, err := tx.Queryx("SELECT id, name FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY name", isActive, isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	details := make([]model.NameDetail, 0)
	for rows.Next() {
		detail := model.NameDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (s *SqliteAssetTypeRepo) Upsert(assetType model.AssetTypeDetail, tx *sqlx.Tx) (err error) {
	if assetType.Id == nil {
		_, err = tx.Exec("INSERT INTO asset_types (name, is_cash, is_liability, sequence, is_active) VALUES (?, ?, ?, ?, ?)",
			assetType.Name, assetType.IsCash, assetType.IsLiability, 0, assetType.IsActive)
	} else {
		_, err = tx.Exec("UPDATE asset_types SET name = ?, is_cash = ?, is_liability = ?, is_active = ? WHERE id = ?",
			assetType.Name, assetType.IsCash, assetType.IsLiability, assetType.IsActive, assetType.Id)
	}

	return
}

func (s *SqliteAssetTypeRepo) Delete(id int, tx *sqlx.Tx) (err error) {
	_, err = tx.Exec("DELETE FROM asset_types WHERE id = ?", id)
	return
}

func (s *SqliteAssetTypeRepo) UpdateSequence(sequences model.SequenceDetail, tx *sqlx.Tx) (err error) {
	_, err = tx.Exec("UPDATE asset_types SET sequence = ? WHERE id = ?", sequences.Sequence, sequences.Id)
	return
}
