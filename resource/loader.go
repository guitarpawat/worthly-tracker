package resource

import (
	"embed"
	"io/fs"
	"worthly-tracker/logs"
)

//go:embed *
var resources embed.FS

func Loader() *embed.FS {
	return &resources
}

func MustLoadDirFS(subDir string) fs.FS {
	dir, err := fs.Sub(resources, subDir)
	if err != nil {
		logs.Log().Panicf("Cannot load resource subdirectory: %v", subDir)
	}
	return dir
}

func MustLoadDirEntry(subDir string) []fs.DirEntry {
	dir, err := fs.ReadDir(resources, subDir)
	if err != nil {
		logs.Log().Panicf("Cannot load resource subdirectory: %v", subDir)
	}
	return dir
}

func MustLoadFile(resource fs.FS, name string) string {
	if resource == nil {
		resource = resources
	}

	data, err := fs.ReadFile(resource, name)
	if err != nil {
		logs.Log().Panicf("Cannot load resource file: %v", err)
	}

	return string(data)
}
