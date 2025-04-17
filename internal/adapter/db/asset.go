package db

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type AssetsRepository struct {
	conn Conn
	db   DB
}

var _ TxRepository[*AssetsRepository] = (*AssetsRepository)(nil)

func NewAssetsRepository(conn Conn) *AssetsRepository {
	return &AssetsRepository{conn: conn, db: conn}
}

type AssetDetail struct {
	Id               sql.NullInt64       `json:"id,omitempty" example:"1"`
	Name             string              `json:"name,omitempty" example:"BTP"`
	Broker           string              `json:"broker,omitempty" example:"SCBAM"`
	TypeId           int                 `json:"typeId,omitempty" db:"type_id" example:"1"`
	TypeName         string              `json:"typeName,omitempty" db:"type_name" example:"Mutual Fund"`
	DefaultIncrement decimal.NullDecimal `json:"defaultIncrement,omitempty" db:"default_increment" example:"1000.00"`
	Sequence         sql.NullInt64       `json:"sequence,omitempty" example:"1"`
	IsActive         bool                `json:"isActive,omitempty" db:"is_active" example:"true"`
}

func (r *AssetsRepository) BeginTx() (*AssetsRepository, TX, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &AssetsRepository{conn: r.conn, db: tx}, tx, nil
}

func (r *AssetsRepository) beginTx() (*AssetsRepository, tx, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &AssetsRepository{conn: r.conn, db: tx}, tx, nil
}

type AssetNameDetail struct {
	Id   int    `json:"id" example:"1"`
	Name string `json:"name" example:"BTP"`
}

type AssetSequenceDetail struct {
	Id       int `json:"id" example:"1"`
	Sequence int `json:"sequence" example:"1"`
}

func (r *AssetsRepository) Get(isActive *bool, typeId *int) ([]AssetDetail, error) {
	rows, err := r.db.Queryx(`SELECT a.id, a.name, a.broker, a.type_id, a.default_increment, a.sequence, a.is_active, t.name AS type_name 
FROM assets a LEFT JOIN asset_types t on a.type_id = t.id WHERE (a.is_active = ? OR ? is null) AND (a.type_id = ? OR ? is null) ORDER BY t.sequence, a.sequence`, isActive, isActive, typeId, typeId)

	if err != nil {
		return nil, err
	}

	details := make([]AssetDetail, 0)
	for rows.Next() {
		detail := AssetDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *AssetsRepository) GetNames(isActive *bool, typeId *int) ([]AssetNameDetail, error) {
	rows, err := r.db.Queryx("SELECT id, name FROM assets WHERE (is_active = ? OR ? IS NULL) AND (type_id = ? OR ? IS NULL) ORDER BY name",
		isActive, isActive, typeId, typeId)

	if err != nil {
		return nil, err
	}

	details := make([]AssetNameDetail, 0)
	for rows.Next() {
		detail := AssetNameDetail{}
		if err = rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *AssetsRepository) Upsert(asset AssetDetail) (err error) {
	if !asset.Id.Valid {
		_, err = r.db.Exec("INSERT INTO assets(name, broker, type_id, default_increment, sequence, is_active) VALUES (?,?,?,?,?,?)",
			asset.Name, asset.Broker, asset.TypeId, asset.DefaultIncrement, 0, asset.IsActive)
	} else {
		_, err = r.db.Exec("UPDATE assets SET name = ?, broker = ?, type_id = ?, default_increment = ?, sequence = 0, is_active = ? WHERE id = ?",
			asset.Name, asset.Broker, asset.TypeId, asset.DefaultIncrement, asset.IsActive, asset.Id)
	}

	return
}

func (r *AssetsRepository) Delete(id int) (err error) {
	_, err = r.db.Exec("DELETE FROM assets WHERE id = ?", id)
	return
}

func (r *AssetsRepository) UpdateSequence(sequence AssetSequenceDetail) (err error) {
	_, err = r.db.Exec("UPDATE assets SET sequence = ? WHERE id = ?", sequence.Sequence, sequence.Id)
	return
}
