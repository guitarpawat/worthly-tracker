package config

import (
	"github.com/spf13/viper"
	"os"
	"worthly-tracker/logs"
	"worthly-tracker/resource"
)

func Init() {
	viper.SetConfigType("yaml")
	args := os.Args
	filePath := ""
	if len(args) >= 2 {
		filePath = args[1]
	}
	if len(filePath) > 0 {
		logs.Log().Debugf("Using conf filepath %v", filePath)
		viper.SetConfigFile(filePath)
		err := viper.ReadInConfig()
		if err != nil {
			logs.Log().Panicf("Cannot read config file: %v", err)
		}
	} else {
		logs.Log().Debug("Using default dev config")
		file, err := resource.Loader().Open("dev.yaml")
		if err != nil {
			logs.Log().Panicf("Cannot read config file: %v", err)
		}
		err = viper.ReadConfig(file)
		if err != nil {
			logs.Log().Panicf("Cannot read config file: %v", err)
		}
	}
	logs.Log().Info("Configuration file load successfully")
}
