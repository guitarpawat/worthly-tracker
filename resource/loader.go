package resource

import (
	"embed"
)

//go:embed static/*
var StaticContent embed.FS

//go:embed migration/*
var DbMigration embed.FS
