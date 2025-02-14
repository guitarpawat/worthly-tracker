package cmd

import (
	"github.com/guitarpawat/worthly-tracker/config"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/http"
	"github.com/guitarpawat/worthly-tracker/internal/service"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "worthly-tracker",
	Short: "Personal Net Worth Tracker",
	RunE:  run,
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cfg, err := config.Init(cfgFile)
	if err != nil {
		return err
	}
	logs.Init(cfg.Logger)

	sqlite, err := db.NewSqlite(ctx, cfg.Datasource.Sqlite)
	if err != nil {
		return err
	}

	repos := db.NewRepositories(sqlite)

	// service
	recordService := service.NewRecords(repos.Record)

	// handler
	recordHandler := http.NewRecordHandler(recordService)

	router := http.NewRouter(http.RouterConfig{Port: cfg.Server.Port}, recordHandler)

	err = router.Start()
	logs.Log().Errorf("server stopped with error: %w", err)
	return err
}
