package db

type Repositories struct {
	Asset     *AssetsRepository
	AssetType *AssetTypesRepository
	Record    *RecordsRepository
}

func NewRepositories(conn Conn) *Repositories {
	return &Repositories{
		Asset:     NewAssetsRepository(conn),
		AssetType: NewAssetTypesRepository(conn),
		Record:    NewRecordsRepository(conn),
	}
}
