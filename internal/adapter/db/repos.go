package db

import "github.com/guitarpawat/worthly-tracker/internal/ports"

// Repositories provides access to all repository implementations
type Repositories struct {
	Asset     ports.AssetsRepository
	AssetType ports.AssetTypesRepository
	Record    ports.RecordsRepository
}

// NewRepositories creates a new Repositories instance with all repository implementations
func NewRepositories(conn ports.Conn) *Repositories {
	return &Repositories{
		Asset:     NewAssetsRepository(conn),
		AssetType: NewAssetTypesRepository(conn),
		Record:    NewRecordsRepository(conn),
	}
}
