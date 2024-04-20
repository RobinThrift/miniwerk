package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/segmentio/ksuid"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := requestIDWithCtx(r.Context(), ksuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}

		wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func(ctx context.Context) {
			dur := float64(time.Since(start)) / float64(time.Millisecond)
			var log = slog.InfoContext
			if wrapped.Status() >= 400 {
				log = slog.ErrorContext
			}

			log(ctx, url, "method", r.Method, "status", wrapped.Status(), "response_time_ms", dur, "bytes_written", wrapped.BytesWritten())
		}(r.Context())

		next.ServeHTTP(wrapped, r)
	})
}
