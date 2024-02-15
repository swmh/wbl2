package app

import (
	"context"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	transport "github.com/swmh/wbl2/develop/dev11/internal/transport/http"
	"log/slog"
)

type Config struct {
	Addr    string
	Service *service.Service
	Logger  *slog.Logger
}

type App struct {
	server transport.Server
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)

}

func New(c Config) *App {
	trCfg := transport.Config{
		Addr:    c.Addr,
		Service: c.Service,
		Logger:  c.Logger,
	}

	tr := transport.New(trCfg)
	return &App{
		server: tr,
	}
}
