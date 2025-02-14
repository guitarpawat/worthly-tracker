package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/spf13/viper"
	"strings"
)

//go:embed config.yaml
var defaultConfig []byte

type Config struct {
	Server     *Server `validate:"required"`
	Logger     logs.Config
	Datasource *Datasource `validate:"required"`
}

type Datasource struct {
	Sqlite *db.SqliteConfig `validate:"required"`
}

type Server struct {
	Port int `validate:"required"`
}

func setDefaultConfig(v *viper.Viper) {
	v.SetDefault("datasource.sqlite.uri", "file::memory:?cache=shared&")
}

func Init(filePath string) (Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaultConfig(v)

	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewReader(defaultConfig))
	if err != nil {
		return Config{}, err
	}

	// read additional config
	if filePath != "" {
		v.SetConfigFile(filePath)
		err = v.MergeInConfig()
		if err != nil {
			return Config{}, err
		}
	}

	cfg := Config{
		Datasource: &Datasource{
			Sqlite: &db.SqliteConfig{
				Uri:         v.GetString("datasource.sqlite.uri"),
				AutoMigrate: v.GetBool("datasource.sqlite.autoMigrate"),
			},
		},
		Logger: logs.Config{
			LogLevel: viper.GetString("log.level"),
		},
		Server: &Server{
			Port: v.GetInt("server.port"),
		},
	}

	// validate config
	err = validator.New(validator.WithRequiredStructEnabled()).Struct(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("fail to validate config: %w", err)
	}

	return cfg, err
}
