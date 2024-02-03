package model

import "github.com/shopspring/decimal"

type AssetRecord struct {
	Id               *int             `json:"id,omitempty" example:"1"`
	AssetId          *int             `json:"assetId,omitempty" example:"1"`
	Name             *string          `json:"name,omitempty" example:"BTP"`
	Broker           *string          `json:"broker,omitempty" example:"SCBAM"`
	DefaultIncrement *decimal.Decimal `json:"defaultIncrement,omitempty" example:"0.00"`
	BoughtValue      *decimal.Decimal `json:"boughtValue,omitempty" example:"100.00"`
	CurrentValue     *decimal.Decimal `json:"currentValue,omitempty" example:"101.50"`
	RealizedValue    *decimal.Decimal `json:"realizedValue,omitempty" example:"0.00"`
	Note             *string          `json:"note,omitempty" example:"Something worth mention"`
}

type AssetTypeRecord struct {
	Id          *int          `json:"id,omitempty" example:"1"`
	Name        *string       `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      *bool         `json:"isCash,omitempty" example:"false"`
	IsLiability *bool         `json:"isLiability,omitempty" example:"false"`
	Assets      []AssetRecord `json:"assets,omitempty"`
}

type DateList struct {
	Current *Date  `json:"current,omitempty" format:"date"` // Selected date
	Prev    []Date `json:"prev,omitempty" format:"date"`    // Prev 12 days from selected date
	Next    []Date `json:"next,omitempty" format:"date"`    // Next 12 days from selected date
}

type AssetTypeDetail struct {
	Id          *int    `json:"id,omitempty" example:"1"`
	Name        *string `json:"name,omitempty" example:"Mutual Funds"`
	IsCash      *bool   `json:"isCash,omitempty" example:"false"`
	IsLiability *bool   `json:"isLiability,omitempty" example:"false"`
	Sequence    *int    `json:"sequence,omitempty" example:"1"`
	IsActive    *bool   `json:"isActive,omitempty" example:"true"`
}

type AssetDetail struct {
	Id               *int             `json:"id,omitempty" example:"1"`
	Name             *string          `json:"name,omitempty" example:"BTP"`
	Broker           *string          `json:"broker,omitempty" example:"SCBAM"`
	TypeId           *int             `json:"typeId,omitempty" db:"type_id" example:"1"`
	TypeName         *string          `json:"typeName,omitempty" db:"type_name" example:"Mutual Fund"`
	DefaultIncrement *decimal.Decimal `json:"defaultIncrement,omitempty" db:"default_increment" example:"1000.00"`
	Sequence         *int             `json:"sequence,omitempty" example:"1"`
	IsActive         *bool            `json:"isActive,omitempty" db:"is_active" example:"true"`
}

type OffsetDetail struct {
	Id            *int            `json:"id,omitempty" example:"1"`
	AssetId       int             `json:"assetId,omitempty" example:"1"`
	EffectiveDate Date            `json:"effectiveDate,omitempty" format:"date"`
	OffsetPrice   decimal.Decimal `json:"offsetPrice,omitempty" example:"-500.00"`
	Note          *string         `json:"note,omitempty" example:"Something worth mention"`
}

type SequenceDetail struct {
	Id       int `json:"id" example:"1"`
	Sequence int `json:"sequence" example:"1"`
}

type NameDetail struct {
	Id   int    `json:"id" example:"1"`
	Name string `json:"name" example:"BTP"`
}
