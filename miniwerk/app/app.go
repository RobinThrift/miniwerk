package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/RobinThrift/miniwerk/miniwerk/ingress/htmlui"
	"github.com/RobinThrift/miniwerk/ui"
	"github.com/go-chi/chi/v5"
)

type App struct {
	srv      *server
	renderer *ui.Renderer
}

func New(ctx context.Context, config *Config) (*App, error) {
	setupLogger(config)

	slog.DebugContext(ctx, "setting up WASM runtime")
	renderer, err := ui.NewRenderer(ctx)
	if err != nil {
		return nil, err
	}

	mux := chi.NewMux()

	slog.DebugContext(ctx, "setting up HTTP server")
	srv := newServer(config, mux)

	slog.DebugContext(ctx, "setting up HTTP ingress routes")
	(&htmlui.HTMLUIRouter{
		Renderer: renderer,
	}).Setup(mux)

	return &App{
		srv:      srv,
		renderer: renderer,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.srv.Start(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return errors.Join(
		app.renderer.Close(ctx),
		app.srv.Stop(ctx),
	)
}
