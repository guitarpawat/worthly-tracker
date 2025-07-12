package view

import (
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type GetRecordTable struct {
	AssetTypes []RecordTableAssetType
	Dates      DateResult
}

type DateResult struct {
	PastDate    []date.Date
	CurrentDate date.Date
	FutureDate  []date.Date
}

type EditRecordTable struct {
	AssetTypes []RecordTableAssetType
	Date       date.Date
	Action     constant.EditRecordTableAction
}

type RecordTableAssetType struct {
	Id      int
	Name    string
	IsCash  bool
	Records []RecordTableData
}

type RecordTableData struct {
	Id               int
	AssetId          int
	Name             string
	Broker           string
	BoughtValue      decimal.Decimal
	CurrentValue     decimal.Decimal
	RealizedValue    decimal.Decimal
	DefaultIncrement decimal.Decimal
	IsCash           bool
	Note             string
}
