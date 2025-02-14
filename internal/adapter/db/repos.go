package db

import "gorm.io/gorm"

type Repositories struct {
	Asset     *AssetsRepository
	AssetType *AssetTypesRepository
	Record    *RecordsRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Asset:     &AssetsRepository{db: db},
		AssetType: &AssetTypesRepository{db: db},
		Record:    &RecordsRepository{db: db},
	}
}
