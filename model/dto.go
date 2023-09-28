package model

import "github.com/shopspring/decimal"

type AssetRecord struct {
	Id               *int             `json:"id,omitempty"`
	AssetId          *int             `json:"assetId,omitempty"`
	Name             *string          `json:"name,omitempty"`
	Broker           *string          `json:"broker,omitempty"`
	DefaultIncrement *decimal.Decimal `json:"defaultIncrement,omitempty"`
	BoughtValue      *decimal.Decimal `json:"boughtValue,omitempty"`
	CurrentValue     *decimal.Decimal `json:"currentValue,omitempty"`
	RealizedValue    *decimal.Decimal `json:"realizedValue,omitempty"`
	Note             *string          `json:"note,omitempty"`
}

type AssetTypeRecord struct {
	Id          *int          `json:"id,omitempty"`
	Name        *string       `json:"name,omitempty"`
	IsCash      *bool         `json:"isCash,omitempty"`
	IsLiability *bool         `json:"isLiability,omitempty"`
	Assets      []AssetRecord `json:"assets,omitempty"`
}

type DateList struct {
	Current *Date  `json:"current,omitempty"`
	Prev    []Date `json:"prev,omitempty"`
	Next    []Date `json:"next,omitempty"`
}

type AssetTypeDetail struct {
	Id          *int    `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	IsCash      *bool   `json:"isCash,omitempty"`
	IsLiability *bool   `json:"isLiability,omitempty"`
	Sequence    *int    `json:"sequence,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
}

type AssetDetail struct {
	Id               *int             `json:"id,omitempty"`
	Name             *string          `json:"name,omitempty"`
	Broker           *string          `json:"broker,omitempty"`
	TypeId           *int             `json:"typeId,omitempty" db:"type_id"`
	TypeName         *string          `json:"typeName,omitempty" db:"type_name"`
	Category         *string          `json:"category,omitempty"`
	DefaultIncrement *decimal.Decimal `json:"defaultIncrement,omitempty" db:"default_increment"`
	Sequence         *int             `json:"sequence,omitempty"`
	IsActive         *bool            `json:"isActive,omitempty" db:"is_active"`
}

type OffsetDetail struct {
	Id            *int            `json:"id,omitempty"`
	AssetId       int             `json:"assetId,omitempty"`
	EffectiveDate Date            `json:"effectiveDate,omitempty"`
	OffsetPrice   decimal.Decimal `json:"offsetPrice,omitempty"`
	Note          *string         `json:"note,omitempty"`
}

type SequenceDetail struct {
	Id       int `json:"id"`
	Sequence int `json:"sequence"`
}

type NameDetail struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
