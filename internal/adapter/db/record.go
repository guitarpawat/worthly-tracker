package db

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type RecordsRepository struct {
	conn Conn
	db   DB
}

var _ TxRepository[*RecordsRepository] = (*RecordsRepository)(nil)

func NewRecordsRepository(conn Conn) *RecordsRepository {
	return &RecordsRepository{conn: conn, db: conn}
}

type DateList struct {
	Current date.Date   `json:"current,omitempty" format:"date"` // Selected date
	Prev    []date.Date `json:"prev,omitempty" format:"date"`    // Prev 12 days from selected date
	Next    []date.Date `json:"next,omitempty" format:"date"`    // Next 12 days from selected date
}

type AssetTypeRecord struct {
	Id          sql.NullInt64 `json:"id,omitempty" example:"1"`
	Name        string        `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      bool          `json:"isCash,omitempty" example:"false"`
	IsLiability bool          `json:"isLiability,omitempty" example:"false"`
	Assets      []AssetRecord `json:"assets,omitempty"`
}

type AssetRecord struct {
	Id               sql.NullInt64       `json:"id,omitempty" example:"1"`
	AssetId          int                 `json:"assetId,omitempty" example:"1"`
	Name             string              `json:"name,omitempty" example:"BTP"`
	Broker           string              `json:"broker,omitempty" example:"SCBAM"`
	DefaultIncrement decimal.NullDecimal `json:"defaultIncrement,omitempty" example:"0.00"`
	BoughtValue      decimal.NullDecimal `json:"boughtValue,omitempty" example:"100.00"`
	CurrentValue     decimal.NullDecimal `json:"currentValue,omitempty" example:"101.50"`
	RealizedValue    decimal.NullDecimal `json:"realizedValue,omitempty" example:"0.00"`
	Note             sql.NullString      `json:"note,omitempty" example:"Something worth mention"`
}

type joinedRecord struct {
	Id               sql.NullInt64       `db:"id"`
	AssetId          int                 `db:"asset_id"`
	Name             string              `db:"name"`
	Broker           string              `db:"broker"`
	DefaultIncrement decimal.NullDecimal `db:"default_increment"`
	BoughtValue      decimal.NullDecimal `db:"bought_value"`
	CurrentValue     decimal.NullDecimal `db:"current_value"`
	RealizedValue    decimal.NullDecimal `db:"realized_value"`
	Note             sql.NullString      `db:"note"`
	TypeId           int                 `db:"type_id"`
	TypeName         string              `db:"type_name"`
	IsCash           bool                `db:"is_cash"`
	IsActive         bool                `db:"is_active"`
	IsLiability      bool                `db:"is_liability"`
}

func (r *RecordsRepository) BeginTx() (*RecordsRepository, TX, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &RecordsRepository{conn: r.conn, db: tx}, tx, nil
}

func (r *RecordsRepository) beginTx() (*RecordsRepository, tx, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return nil, nil, err
	}
	return &RecordsRepository{conn: r.conn, db: tx}, tx, nil
}

func (r *RecordsRepository) GetDate(current date.Date) (*DateList, error) {
	prevRows, err := r.db.Query("SELECT DISTINCT date FROM records WHERE date < ? ORDER BY date DESC LIMIT 12", current.String())
	if err != nil {
		return nil, err
	}
	defer prevRows.Close()

	var prev = make([]date.Date, 0, 12)
	for prevRows.Next() {
		var d date.Date
		if err = prevRows.Scan(&d); err != nil {
			return nil, err
		}
		prev = append(prev, d)
	}

	nextRows, err := r.db.Query("SELECT DISTINCT date FROM records WHERE date > ? ORDER BY date LIMIT 12", current.String())
	if err != nil {
		return nil, err
	}
	defer nextRows.Close()

	var next = make([]date.Date, 0, 12)
	for nextRows.Next() {
		var d date.Date
		if err = nextRows.Scan(&d); err != nil {
			return nil, err
		}
		next = append(next, d)
	}

	return &DateList{
		Current: current,
		Prev:    prev,
		Next:    next,
	}, nil
}

func (r *RecordsRepository) GetLatestDate() (date.Date, error) {
	row := r.db.QueryRowx("SELECT date FROM records ORDER BY date DESC LIMIT 1")

	var d date.Date
	err := row.Scan(&d)
	if errors.Is(err, sql.ErrNoRows) {
		return date.Zero, nil
	}
	if err != nil {
		return date.Zero, err
	}
	return d, nil
}

func (r *RecordsRepository) mapJoinedRecord(records []joinedRecord) []AssetTypeRecord {
	// assume that the records is order by type_id and asset_id
	types := make([]AssetTypeRecord, 0)
	for _, v := range records {
		if len(types) == 0 || int(types[len(types)-1].Id.Int64) != v.TypeId {
			newType := AssetTypeRecord{
				Id:          sql.NullInt64{Int64: int64(v.TypeId), Valid: true},
				Name:        v.TypeName,
				IsCash:      v.IsCash,
				IsLiability: v.IsLiability,
				Assets:      make([]AssetRecord, 0),
			}
			types = append(types, newType)
		}

		types[len(types)-1].Assets = append(types[len(types)-1].Assets, AssetRecord{
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

func (r *RecordsRepository) GetRecordByDate(d date.Date) ([]AssetTypeRecord, error) {
	rows, err := r.db.Queryx(`SELECT r.id, r.asset_id, a.name, a.broker, a.default_increment, 
r.bought_value, r.current_value, r.realized_value, r.note, a.type_id, t.name AS type_name, t.is_cash, t.is_liability 
FROM records r LEFT JOIN assets a on a.id = r.asset_id LEFT JOIN asset_types t on t.id = a.type_id WHERE r.date = $1 ORDER BY t.sequence, a.sequence`, d)
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

func (r *RecordsRepository) GetRecordDraft() ([]AssetTypeRecord, error) {
	rows, err := r.db.Queryx(`SELECT r.id, a.id as asset_id, a.name, a.broker, a.default_increment,r.bought_value, r.current_value, r.realized_value, 
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

func (r *RecordsRepository) UpsertRecord(record AssetRecord, d date.Date) (err error) {
	if !record.Id.Valid {
		_, err = r.db.Exec("INSERT INTO records(asset_id, date, bought_value, current_value, realized_value, note) VALUES ($1, $2, $3, $4, $5, $6)",
			record.AssetId, d, record.BoughtValue, record.CurrentValue, record.RealizedValue, record.Note)
	} else {
		_, err = r.db.Exec("UPDATE records SET asset_id = $1, date = $2, bought_value = $3, current_value = $4, realized_value = $5, note = $6 WHERE id = $7",
			record.AssetId, d, record.BoughtValue, record.CurrentValue, record.RealizedValue, record.Note, record.Id)
	}
	return
}

func (r *RecordsRepository) DeleteRecordById(id int) error {
	_, err := r.db.Exec("DELETE FROM records WHERE id = $1", id)
	return err
}

func (r *RecordsRepository) DeleteRecordByDate(d date.Date) (int64, error) {
	exec, err := r.db.Exec("DELETE FROM records WHERE date = $1", d)
	if err != nil {
		return -1, err
	}
	return exec.RowsAffected()
}
