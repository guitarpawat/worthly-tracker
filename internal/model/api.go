package model

import "github.com/shopspring/decimal"

type GetRecordResponse struct {
	AssetTypes []GetRecordAssetType
	Dates      DateResult
}

type GetRecordAssetType struct {
	Id      int
	Name    string
	IsCash  bool
	Records []GetRecordData
}

type GetRecordData struct {
	Id            int
	Name          string
	Broker        string
	BoughtValue   decimal.Decimal
	CurrentValue  decimal.Decimal
	RealizedValue decimal.Decimal
	IsCash        bool
	Note          string
}
