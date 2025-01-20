package entity

import "github.com/shopspring/decimal"

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
