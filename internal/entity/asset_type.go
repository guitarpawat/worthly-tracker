package entity

type AssetType struct {
	Id          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"index:idx_asset_type_name"`
	IsCash      bool    `json:"isCash" gorm:"type:boolean"`
	IsLiability bool    `json:"isLiability" gorm:"type:boolean"`
	Sequence    int     `json:"sequence" gorm:""`
	IsActive    bool    `json:"isActive" gorm:"type:boolean"`
	Assets      []Asset `json:"assets" gorm:"foreignKey:AssetTypeId;references:Id"`
}
