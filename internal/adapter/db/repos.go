package db

import "gorm.io/gorm"

type Repositories struct {
	Asset       *AssetRepository
	AssetType   *AssetTypeRepository
	Record      *RecordRepository
	ValueOffset *ValueOffsetRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Asset:       &AssetRepository{db: db},
		AssetType:   &AssetTypeRepository{db: db},
		Record:      &RecordRepository{db: db},
		ValueOffset: &ValueOffsetRepository{db: db},
	}
}
