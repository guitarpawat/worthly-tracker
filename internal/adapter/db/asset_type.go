package db

import (
	"database/sql"
	"errors"
)

type AssetTypesRepository struct {
	conn Conn
	db   DB
}

var _ TxRepository[*AssetTypesRepository] = (*AssetTypesRepository)(nil)

func NewAssetTypesRepository(conn Conn) *AssetTypesRepository {
	return &AssetTypesRepository{conn: conn, db: conn}
}

type AssetTypeDetail struct {
	Id          sql.NullInt64 `json:"id,omitempty" example:"1"`
	Name        string        `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      bool          `json:"isCash,omitempty" example:"false" db:"is_cash"`
	IsLiability bool          `json:"isLiability,omitempty" example:"false" db:"is_liability"`
	Sequence    sql.NullInt64 `json:"sequence,omitempty" example:"1"`
	IsActive    bool          `json:"isActive,omitempty" example:"true" db:"is_active"`
}

type AssetTypeNameDetail struct {
	Id   int    `json:"id" example:"1"`
	Name string `json:"name" example:"BTP"`
}

type AssetTypeSequenceDetail struct {
	Id       int `json:"id" example:"1"`
	Sequence int `json:"sequence" example:"1"`
}

func (r *AssetTypesRepository) BeginTx() (*AssetTypesRepository, TX, error) {
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

func (r *AssetTypesRepository) Get(isActive *bool) ([]AssetTypeDetail, error) {
	rows, err := r.db.Queryx("SELECT * FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY sequence", isActive, isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	details := make([]AssetTypeDetail, 0)
	for rows.Next() {
		detail := AssetTypeDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *AssetTypesRepository) GetNames(isActive *bool) ([]AssetTypeNameDetail, error) {
	rows, err := r.db.Queryx("SELECT id, name FROM asset_types WHERE (is_active = ? OR ? IS NULL) ORDER BY name", isActive, isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	details := make([]AssetTypeNameDetail, 0)
	for rows.Next() {
		detail := AssetTypeNameDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *AssetTypesRepository) Upsert(assetType AssetTypeDetail) (err error) {
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

func (r *AssetTypesRepository) UpdateSequence(sequences AssetTypeSequenceDetail) (err error) {
	_, err = r.db.Exec("UPDATE asset_types SET sequence = ? WHERE id = ?", sequences.Sequence, sequences.Id)
	return
}
