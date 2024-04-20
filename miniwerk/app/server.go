package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/RobinThrift/miniwerk/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	srv *http.Server
}

func newServer(config *Config, mux *chi.Mux) *server {
	srv := &server{
		srv: &http.Server{
			Addr:    config.Addr,
			Handler: mux,
		},
	}

	mux.Use(
		requestIDMiddleware,
		logReqMiddleware,
		middleware.Compress(5),
	)

	mux.Get("/health", http.HandlerFunc(srv.handleHealth))
	mux.Handle("/static/*", ui.Files("/static/"))

	return srv
}

func (s server) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "starting http server on "+s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping http server")
	return s.srv.Shutdown(ctx)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		slog.ErrorContext(r.Context(), "error writing http health response", "error", err)
	}
}
