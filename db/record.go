package db

import (
	"database/sql"
	"errors"
	"github.com/go-sqlx/sqlx"
	"github.com/shopspring/decimal"
	"worthly-tracker/model"
)

var recordRepo = &SqliteRecordRepo{}

func GetRecordRepo() *SqliteRecordRepo {
	return recordRepo
}

type SqliteRecordRepo struct {
}

func (r *SqliteRecordRepo) GetDate(current model.Date, tx *sqlx.Tx) (*model.DateList, error) {
	prevRows, err := tx.Query("SELECT DISTINCT date FROM records WHERE date < ? ORDER BY date DESC LIMIT 12", current.String())
	if err != nil {
		return nil, err
	}
	defer prevRows.Close()

	var prev = make([]model.Date, 0, 12)
	for prevRows.Next() {
		var date model.Date
		if err = prevRows.Scan(&date); err != nil {
			return nil, err
		}
		prev = append(prev, date)
	}

	nextRows, err := tx.Query("SELECT DISTINCT date FROM records WHERE date > ? ORDER BY date LIMIT 12", current.String())
	if err != nil {
		return nil, err
	}
	defer nextRows.Close()

	var next = make([]model.Date, 0, 12)
	for nextRows.Next() {
		var date model.Date
		if err = nextRows.Scan(&date); err != nil {
			return nil, err
		}
		next = append(next, date)
	}

	return &model.DateList{
		Current: &current,
		Prev:    prev,
		Next:    next,
	}, nil
}

func (r *SqliteRecordRepo) GetLatestDate(tx *sqlx.Tx) (*model.Date, error) {
	row := tx.QueryRow("SELECT date FROM records ORDER BY date DESC LIMIT 1")

	var date *model.Date
	err := row.Scan(&date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return date, nil
}

type joinedRecord struct {
	Id               *int             `db:"id"`
	AssetId          *int             `db:"asset_id"`
	Name             *string          `db:"name"`
	Broker           *string          `db:"broker"`
	DefaultIncrement *decimal.Decimal `db:"default_increment"`
	BoughtValue      *decimal.Decimal `db:"bought_value"`
	CurrentValue     *decimal.Decimal `db:"current_value"`
	RealizedValue    *decimal.Decimal `db:"realized_value"`
	Note             *string          `db:"note"`
	TypeId           *int             `db:"type_id"`
	TypeName         *string          `db:"type_name"`
	IsCash           *bool            `db:"is_cash"`
	IsActive         *bool            `db:"is_active"`
	IsLiability      *bool            `db:"is_liability"`
}

func (r *SqliteRecordRepo) mapJoinedRecord(records []joinedRecord) []model.AssetTypeRecord {
	// assume that the records is order by type_id and asset_id
	types := make([]model.AssetTypeRecord, 0)
	for _, v := range records {
		if len(types) == 0 || *types[len(types)-1].Id != *v.TypeId {
			newType := model.AssetTypeRecord{
				Id:          v.TypeId,
				Name:        v.TypeName,
				IsCash:      v.IsCash,
				IsLiability: v.IsLiability,
				Assets:      make([]model.AssetRecord, 0),
			}
			types = append(types, newType)
		}

		types[len(types)-1].Assets = append(types[len(types)-1].Assets, model.AssetRecord{
			Id:               v.Id,
			AssetId:          v.AssetId,
			Name:             v.Name,
			Broker:           v.Broker,
			DefaultIncrement: v.DefaultIncrement,
			BoughtValue:      v.BoughtValue,
			CurrentValue:     v.CurrentValue,
			RealizedValue:    v.RealizedValue,
			Note:             v.Note,
		})
	}

	return types
}

func (r *SqliteRecordRepo) GetRecordByDate(date model.Date, tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	rows, err := tx.Queryx(`SELECT r.id, r.asset_id, a.name, a.broker, a.default_increment, 
r.bought_value, r.current_value, r.realized_value, r.note, a.type_id, t.name AS type_name, t.is_cash, t.is_liability 
FROM records r LEFT JOIN assets a on a.id = r.asset_id LEFT JOIN asset_types t on t.id = a.type_id WHERE r.date = $1 ORDER BY t.sequence, a.sequence`, date.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]joinedRecord, 0)
	for rows.Next() {
		record := joinedRecord{}
		if err = rows.StructScan(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return r.mapJoinedRecord(records), nil
}

func (r *SqliteRecordRepo) GetRecordDraft(tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	rows, err := tx.Queryx(`SELECT r.id, a.id as asset_id, a.name, a.broker, a.default_increment,r.bought_value, r.current_value, r.realized_value, 
r.note, a.type_id, t.name AS type_name, t.is_cash, t.is_liability FROM asset_types t LEFT JOIN assets a on t.id = a.type_id LEFT JOIN records r ON 
a.id = r.asset_id AND r.date = (SELECT date FROM records ORDER BY date DESC LIMIT 1) WHERE a.is_active AND t.is_active = true ORDER BY t.sequence, a.sequence`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]joinedRecord, 0)
	for rows.Next() {
		record := joinedRecord{}
		if err = rows.StructScan(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return r.mapJoinedRecord(records), nil
}

func (r *SqliteRecordRepo) UpsertRecord(record model.AssetRecord, date model.Date, tx *sqlx.Tx) (err error) {
	if record.Id == nil {
		_, err = tx.Exec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES ($1, $2, $3, $4, $5, $6)",
			record.AssetId, date.String(), record.BoughtValue, record.CurrentValue, record.RealizedValue, record.Note)
	} else {
		_, err = tx.Exec("UPDATE records SET asset_id = $1, date = $2, bought_value = $3, current_value = $4, realized_value = $5, note = $6 WHERE id = $7",
			record.AssetId, date.String(), record.BoughtValue, record.CurrentValue, record.RealizedValue, record.Note, record.Id)
	}
	return
}

func (r *SqliteRecordRepo) DeleteRecordById(id int, tx *sqlx.Tx) error {
	_, err := tx.Exec("DELETE FROM records WHERE id = $1", id)
	return err
}

func (r *SqliteRecordRepo) DeleteRecordByDate(date model.Date, tx *sqlx.Tx) (int64, error) {
	exec, err := tx.Exec("DELETE FROM records WHERE date = $1", date.String())
	if err != nil {
		return -1, err
	}
	return exec.RowsAffected()
}
