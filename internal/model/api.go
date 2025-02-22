package model

import (
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type CreateRecordRequest struct {
	Date    date.Date                   `validate:"required"`
	Records []CreateRecordRequestRecord `validate:"required"`
}

type CreateRecordRequestRecord struct {
	AssetId       string `validate:"number"`
	BoughtValue   decimal.Decimal
	CurrentValue  decimal.Decimal
	RealizedValue decimal.Decimal
	Note          string
}
