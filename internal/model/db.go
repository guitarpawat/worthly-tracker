package model

import (
	"database/sql"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

// AssetDetail represents an asset with its details
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

// AssetNameDetail represents an asset with just id and name
type AssetNameDetail struct {
	Id   int    `json:"id" example:"1"`
	Name string `json:"name" example:"BTP"`
}

// AssetSequenceDetail contains an asset's sequence information
type AssetSequenceDetail struct {
	Id       int `json:"id" example:"1"`
	Sequence int `json:"sequence" example:"1"`
}

// AssetTypeDetail represents an asset type with its details
type AssetTypeDetail struct {
	Id          sql.NullInt64 `json:"id,omitempty" example:"1"`
	Name        string        `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      bool          `json:"isCash,omitempty" example:"false" db:"is_cash"`
	IsLiability bool          `json:"isLiability,omitempty" example:"false" db:"is_liability"`
	Sequence    sql.NullInt64 `json:"sequence,omitempty" example:"1"`
	IsActive    bool          `json:"isActive,omitempty" example:"true" db:"is_active"`
}

// AssetTypeNameDetail represents an asset type with just id and name
type AssetTypeNameDetail struct {
	Id   int    `json:"id" example:"1"`
	Name string `json:"name" example:"BTP"`
}

// AssetTypeSequenceDetail contains an asset type's sequence information
type AssetTypeSequenceDetail struct {
	Id       int `json:"id" example:"1"`
	Sequence int `json:"sequence" example:"1"`
}

// DateList represents a list of dates with context (current, previous, next)
type DateList struct {
	Current date.Date   `json:"current,omitempty" format:"date"` // Selected date
	Prev    []date.Date `json:"prev,omitempty" format:"date"`    // Prev 12 days from selected date
	Next    []date.Date `json:"next,omitempty" format:"date"`    // Next 12 days from selected date
}

// AssetTypeRecord represents an asset type with its associated assets
type AssetTypeRecord struct {
	Id          sql.NullInt64 `json:"id,omitempty" example:"1"`
	Name        string        `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      bool          `json:"isCash,omitempty" example:"false"`
	IsLiability bool          `json:"isLiability,omitempty" example:"false"`
	Assets      []AssetRecord `json:"assets,omitempty"`
}

// AssetRecord represents a record for a specific asset
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

// JoinedRecord represents a record joined with asset and asset type information
type JoinedRecord struct {
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
