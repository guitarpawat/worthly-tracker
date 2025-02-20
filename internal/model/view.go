package model

import (
	"github.com/guitarpawat/worthly-tracker/internal/constant"
	"github.com/rickb777/date/v2"
	"github.com/shopspring/decimal"
)

type GetRecordTableView struct {
	AssetTypes []RecordTableViewAssetType
	Dates      DateResult
}

type EditRecordTableView struct {
	AssetTypes []RecordTableViewAssetType
	Date       date.Date
	Action     constant.EditRecordTableAction
}

type RecordTableViewAssetType struct {
	Id      int
	Name    string
	IsCash  bool
	Records []RecordTableViewData
}

type RecordTableViewData struct {
	Id               int
	Name             string
	Broker           string
	BoughtValue      decimal.Decimal
	CurrentValue     decimal.Decimal
	RealizedValue    decimal.Decimal
	DefaultIncrement decimal.Decimal
	IsCash           bool
	Note             string
}
