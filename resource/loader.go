package resource

import (
	"embed"
)

//go:embed static/*
var StaticContent embed.FS
