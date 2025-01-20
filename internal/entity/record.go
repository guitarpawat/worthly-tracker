package entity

import (
	"database/sql"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type Record struct {
	Id            int             `json:"id" gorm:"primaryKey"`
	AssetId       int             `json:"assetId" gorm:"index:idx_record_asset_id_date"`
	Asset         Asset           `json:"asset" gorm:"foreignKey:AssetId;references:Id"`
	Date          date.Date       `json:"date" gorm:"type:date;index:idx_record_asset_id_date;index:idx_record_asset_date"`
	BoughtValue   decimal.Decimal `json:"boughtValue" gorm:"type:decimal(14,2);ddefault:0"`
	CurrentValue  decimal.Decimal `json:"currentValue" gorm:"type:decimal(14,2);ddefault:0"`
	RealizedValue decimal.Decimal `json:"realizedValue" gorm:"type:decimal(14,2);ddefault:0"`
	Note          sql.NullString  `json:"note" gorm:""`
}
