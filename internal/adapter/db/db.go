package db

import (
	"github.com/guitarpawat/worthly-tracker/internal/ports"
)

type tx interface {
	ports.DB
	ports.TX
}

// Ensure that our types implement the ports interfaces
var _ ports.TX = (tx)(nil)
