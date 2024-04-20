package ui

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type ErrorPageErr struct {
	Err     error
	Code    int
	Title   string
	Details string
}

func (err ErrorPageErr) Error() string {
	return err.Err.Error()
}

func (r *Renderer) RenderErrorPage(ctx context.Context, err error, output io.Writer) error {
	var errPageErr ErrorPageErr

	if !errors.As(err, &errPageErr) {
		errPageErr = ErrorPageErr{
			Code:    http.StatusInternalServerError,
			Title:   "Unknown Error",
			Details: err.Error(),
		}
	}

	if errPageErr.Title == "" && errPageErr.Err != nil {
		errPageErr.Title = errPageErr.Err.Error()
	}

	renderErr := r.Render(ctx, RenderInput{Page: "Error", Data: errPageErr}, output)
	if renderErr != nil {
		slog.ErrorContext(ctx, "error rendering error page", "error", renderErr)
	}

	return nil
}

type IndexPageProps struct {
}

func (r *Renderer) RenderIndexPage(ctx context.Context, data IndexPageProps, output io.Writer) error {
	return r.Render(ctx, RenderInput{Page: "Index", Data: data}, output)
}
