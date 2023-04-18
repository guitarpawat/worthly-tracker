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
	Id               *int                `json:"id,omitempty"`
	Name             *string             `json:"name,omitempty"`
	Broker           *string             `json:"broker,omitempty"`
	TypeId           *int                `json:"typeId,omitempty"`
	TypeName         *string             `json:"typeName,omitempty"`
	Category         *string             `json:"category,omitempty"`
	DefaultIncrement decimal.NullDecimal `json:"defaultIncrement,omitempty"`
	Sequence         *int                `json:"sequence,omitempty"`
	IsActive         *bool               `json:"isActive,omitempty"`
}

type OffsetDetail struct {
	Id            *int                `json:"id,omitempty"`
	AssetId       *int                `json:"assetId,omitempty"`
	AssetName     *string             `json:"assetName,omitempty"`
	Broker        *string             `json:"broker,omitempty"`
	EffectiveDate *string             `json:"effectiveDate,omitempty"`
	OffsetPrice   decimal.NullDecimal `json:"offsetPrice,omitempty"`
	Note          *string             `json:"note,omitempty"`
}
