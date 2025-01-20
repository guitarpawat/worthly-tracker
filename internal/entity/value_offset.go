package entity

import (
	"database/sql"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type ValueOffset struct {
	Id            int             `json:"id" gorm:"primaryKey"`
	EffectiveDate date.Date       `json:"effectiveDate" gorm:"type:date;"`
	AssetId       int             `json:"assetId" gorm:"index:idx_value_offset_name"`
	Asset         Asset           `json:"asset" gorm:"foreignKey:AssetId;references:Id"`
	OffsetPrice   decimal.Decimal `json:"offsetPrice" gorm:"type:decimal(14,2);dindex:idx_value_offset_price"`
	Note          sql.NullString  `json:"note" gorm:""`
}
