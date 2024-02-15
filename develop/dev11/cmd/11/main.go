package main

import (
	"context"
	"github.com/swmh/wbl2/develop/dev11/internal/app"
	"github.com/swmh/wbl2/develop/dev11/internal/repo/sqlite"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := sqlite.New()
	if err != nil {
		logger.Error("Cannot connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}

	s := service.New(db)

	appCfg := app.Config{
		Addr:    ":8099",
		Service: s,
		Logger:  logger,
	}

	a := app.New(appCfg)

	go func() {
		logger.Info("Server started")
		err := a.Run()
		if err != nil {
			logger.Error("Server closed", slog.String("err", err.Error()))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	logger.Info("Shutdown...")
	a.Shutdown(ctx)
	logger.Info("Server closed")

}
