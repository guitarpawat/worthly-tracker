package cmd

import (
	"github.com/guitarpawat/worthly-tracker/config"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
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

	_ = db.NewRepositories(sqlite)

	return nil
}
