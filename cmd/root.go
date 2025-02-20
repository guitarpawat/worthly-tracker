package cmd

import (
	"errors"
	"github.com/guitarpawat/worthly-tracker/config"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/db"
	"github.com/guitarpawat/worthly-tracker/internal/adapter/http"
	"github.com/guitarpawat/worthly-tracker/internal/service"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"github.com/spf13/cobra"
	net "net/http"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
		<-gracefulStop
		logs.Log().Info("server is shutting down")
		err := router.Close()
		logs.Log().Errorf("cannot stop server: %v", err)
	}()

	err = router.Start()
	if errors.Is(err, net.ErrServerClosed) {
		return nil
	}
	logs.Log().Errorf("server stopped with error: %v", err)
	return err
}
