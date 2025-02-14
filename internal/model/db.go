package model

import (
	"database/sql"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type Asset struct {
	Id               int             `json:"id" gorm:"primaryKey"`
	Name             string          `json:"name" gorm:"index:idx_asset_name_broker"`
	Broker           string          `json:"broker" gorm:"index:idx_asset_name_broker"`
	AssetTypeId      int             `json:"assetTypeId" gorm:"index:idx_asset_asset_type_id"`
	AssetType        AssetType       `json:"-" gorm:"foreignKey:AssetTypeId;references:Id"`
	DefaultIncrement decimal.Decimal `json:"defaultIncrement" gorm:"type:decimal(14,2);default:0"`
	Sequence         int             `json:"sequence" gorm:""`
	IsActive         bool            `json:"isActive" gorm:"type:boolean"`
}

type AssetType struct {
	Id          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"index:idx_asset_type_name"`
	IsCash      bool    `json:"isCash" gorm:"type:boolean"`
	IsLiability bool    `json:"isLiability" gorm:"type:boolean"`
	Sequence    int     `json:"sequence" gorm:""`
	IsActive    bool    `json:"isActive" gorm:"type:boolean"`
	Assets      []Asset `json:"assets" gorm:"foreignKey:AssetTypeId;references:Id"`
}

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

type DateResult struct {
	PastDate    []date.Date
	CurrentDate date.Date
	FutureDate  []date.Date
}
