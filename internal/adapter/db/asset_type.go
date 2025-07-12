package db

import (
	"database/sql"
	"errors"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/ports"
)

type AssetTypesRepository struct {
	conn ports.Conn
	db   ports.DB
}

var _ ports.AssetTypesRepository = (*AssetTypesRepository)(nil)

func NewAssetTypesRepository(conn ports.Conn) *AssetTypesRepository {
	return &AssetTypesRepository{conn: conn, db: conn}
}

func (r *AssetTypesRepository) BeginTx() (ports.AssetTypesRepository, ports.TX, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &AssetTypesRepository{conn: r.conn, db: tx}, tx, nil
}

func (r *AssetTypesRepository) beginTx() (*AssetTypesRepository, tx, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &AssetTypesRepository{conn: r.conn, db: tx}, tx, nil
}

func (r *AssetTypesRepository) Get(isActive *bool) ([]model.AssetTypeDetail, error) {
	rows, err := r.db.Queryx("SELECT * FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY sequence", isActive, isActive)

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

func (r *AssetTypesRepository) GetNames(isActive *bool) ([]model.AssetTypeNameDetail, error) {
	rows, err := r.db.Queryx("SELECT id, name FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY name", isActive, isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	details := make([]model.AssetTypeNameDetail, 0)
	for rows.Next() {
		detail := model.AssetTypeNameDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *AssetTypesRepository) Upsert(assetType model.AssetTypeDetail) (err error) {
	if !assetType.Id.Valid {
		_, err = r.db.Exec("INSERT INTO asset_types (name, is_cash, is_liability, sequence, is_active) VALUES (?, ?, ?, ?, ?)",
			assetType.Name, assetType.IsCash, assetType.IsLiability, 0, assetType.IsActive)
	} else {
		_, err = r.db.Exec("UPDATE asset_types SET name = ?, is_cash = ?, is_liability = ?, is_active = ?, sequence = ? WHERE id = ?",
			assetType.Name, assetType.IsCash, assetType.IsLiability, assetType.IsActive, assetType.Sequence, assetType.Id)
	}

	return
}

func (r *AssetTypesRepository) Delete(id int) (err error) {
	_, err = r.db.Exec("DELETE FROM asset_types WHERE id = ?", id)
	return
}

func (r *AssetTypesRepository) UpdateSequence(sequences model.AssetTypeSequenceDetail) (err error) {
	_, err = r.db.Exec("UPDATE asset_types SET sequence = ? WHERE id = ?", sequences.Sequence, sequences.Id)
	return
}
