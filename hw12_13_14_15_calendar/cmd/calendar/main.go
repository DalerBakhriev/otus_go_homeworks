package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store"
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store/memorystore"
	"github.com/DalerBakhriev/otus_go_homeworks/hw12_13_14_15_calendar/internal/store/sqlstore"
	"github.com/jmoiron/sqlx"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.File)
	defer logg.Sync()

	var store store.Store
	if config.DB.inMemory {
		store = memorystore.New()
	} else {
		db, err := sqlx.Connect("postgres", config.DB.url)
		if err != nil {
			panic(fmt.Errorf("connect to db: %w", err))
		}
		store = sqlstore.New(db)
	}

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
