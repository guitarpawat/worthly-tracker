//go:build test || unit || integration

package config

import (
	"github.com/spf13/viper"
	"worthly-tracker/logs"
	"worthly-tracker/resource"
)

func InitTest() {
	viper.SetConfigType("yaml")
	logs.Log().Debug("Using default test config")
	file, err := resource.Loader().Open("test.yaml")
	if err != nil {
		logs.Log().Panicf("Cannot read config file: %v", err)
	}
	err = viper.ReadConfig(file)
	if err != nil {
		logs.Log().Panicf("Cannot read config file: %v", err)
	}
	logs.Log().Info("Configuration file load successfully")
}
