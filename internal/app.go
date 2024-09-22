package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const httpTimeoutClose = 5 * time.Second

type Option func(app *App)

type App struct {
	httpSrv *http.Server
}

func New(opts ...Option) *App {
	app := &App{}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) Run() error {
	var err error

	a.httpSrv, err = NewHTTPServer(":8882")
	if err != nil {
		return fmt.Errorf("http serve init: %w", err)
	}

	go func() {
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http serve error: %v", err)
		}
	}()

	return nil
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutClose)
	defer cancel()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("http closing error: %w", err)
	}

	return nil
}
