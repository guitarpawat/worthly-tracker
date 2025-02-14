package model

type GetRecordResponse struct {
	AssetTypes []GetRecordAssetType
	Dates      DateResult
}

type GetRecordAssetType struct {
	Name    string
	Records []Record
}
