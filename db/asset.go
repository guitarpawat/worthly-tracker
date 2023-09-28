package db

import (
	"github.com/jmoiron/sqlx"
	"worthly-tracker/model"
)

var assetRepo = &SqliteAssetRepo{}

func GetAssetRepo() *SqliteAssetRepo {
	return assetRepo
}

type SqliteAssetRepo struct {
}

func (s *SqliteAssetRepo) Get(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.AssetDetail, error) {

	rows, err := tx.Queryx(`SELECT a.id, a.name, a.broker, a.type_id, a.default_increment, a.sequence, a.is_active, t.name AS type_name 
FROM assets a LEFT JOIN asset_types t on a.type_id = t.id WHERE (a.is_active = ? OR ? is null) AND (a.type_id = ? OR ? is null) ORDER BY t.sequence, a.sequence`, isActive, isActive, typeId, typeId)

	if err != nil {
		return nil, err
	}

	details := make([]model.AssetDetail, 0)
	for rows.Next() {
		detail := model.AssetDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (s *SqliteAssetRepo) GetNames(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.NameDetail, error) {
	rows, err := tx.Queryx("SELECT id, name FROM assets WHERE (is_active = ? OR ? IS NULL) AND (type_id = ? OR ? IS NULL) ORDER BY name",
		isActive, isActive, typeId, typeId)

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

func (s *SqliteAssetRepo) Upsert(asset model.AssetDetail, tx *sqlx.Tx) (err error) {
	if asset.Id == nil {
		_, err = tx.Exec("INSERT INTO assets(name, broker, type_id, default_increment, sequence, is_active) VALUES (?,?,?,?,?,?)",
			asset.Name, asset.Broker, asset.TypeId, asset.DefaultIncrement, 0, asset.IsActive)
	} else {
		_, err = tx.Exec("UPDATE assets SET name = ?, broker = ?, type_id = ?, default_increment = ?, sequence = 0, is_active = ? WHERE id = ?",
			asset.Name, asset.Broker, asset.TypeId, asset.DefaultIncrement, asset.IsActive, asset.Id)
	}

	return
}

func (s *SqliteAssetRepo) Delete(id int, tx *sqlx.Tx) (err error) {
	_, err = tx.Exec("DELETE FROM assets WHERE id = ?", id)
	return
}

func (s *SqliteAssetRepo) UpdateSequence(sequence model.SequenceDetail, tx *sqlx.Tx) (err error) {
	_, err = tx.Exec("UPDATE assets SET sequence = ? WHERE id = ?", sequence.Sequence, sequence.Id)
	return
}
