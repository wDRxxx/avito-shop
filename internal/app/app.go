package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/wDRxxx/avito-shop/internal/closer"
	"github.com/wDRxxx/avito-shop/internal/config"
)

type App struct {
	wg *sync.WaitGroup

	serviceProvider *serviceProvider

	httpServer *http.Server
}

func NewApp(ctx context.Context, wg *sync.WaitGroup, envPath string) (*App, error) {
	err := config.Load(envPath)
	if err != nil {
		return nil, err
	}

	app := &App{wg: wg}

	err = app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initDeps(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()

	a.initHTTPServer(ctx)

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) {
	s := a.serviceProvider.HTTPServer(ctx, a.wg)
	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HttpConfig().Address(),
		Handler:           s.Handler(),
		ReadHeaderTimeout: a.serviceProvider.HttpConfig().ReadHeaderTimeout(),
	}
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	a.wg.Add(1)
	go func() {
		closer.Add(1, func() error {
			a.wg.Done()
			return nil
		})

		err := a.runHttpServer()
		if err != nil {
			log.Fatalf("error running http server: %v", err)
		}
	}()

	a.wg.Wait()

	return nil
}

func (a *App) runHttpServer() error {
	slog.Info("starting http server...")

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
