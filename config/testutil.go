//go:build test || unit || integration

package config

import (
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"worthly-tracker/logs"
	"worthly-tracker/resource"
)

func InitTest() {
	var file fs.File
	var err error
	viper.SetConfigType("yaml")
	if os.Getenv("CI") == "true" {
		logs.Log().Debug("Using default CI config")
		file, err = resource.Loader().Open("ci_test.yaml")
	} else {
		logs.Log().Debug("Using default test config")
		file, err = resource.Loader().Open("test.yaml")
	}
	if err != nil {
		logs.Log().Panicf("Cannot read config file: %v", err)
	}
	err = viper.ReadConfig(file)
	if err != nil {
		logs.Log().Panicf("Cannot read config file: %v", err)
	}
	logs.Log().Info("Configuration file load successfully")
}
