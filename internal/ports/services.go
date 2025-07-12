package ports

import (
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/internal/view"
	"github.com/rickb777/date/v2"
)

// RecordsService defines the interface for record service operations
type RecordsService interface {
	GetByDate(d date.Date) (view.GetRecordTable, error)
	GetByDateForEdit(d date.Date) (view.EditRecordTable, error)
	GetDraft() (view.EditRecordTable, error)
	CreateRecords(req model.CreateRecordRequest) error
	UpdateRecords(req model.UpdateRecordRequest) error
	DeleteRecords(d date.Date) error
}
