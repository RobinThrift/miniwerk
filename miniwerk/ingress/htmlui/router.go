package htmlui

import (
	"log/slog"
	"net/http"

	"github.com/RobinThrift/miniwerk/ui"
	"github.com/go-chi/chi/v5"
)

type HTMLUIRouter struct {
	Renderer *ui.Renderer
}

func (h *HTMLUIRouter) Setup(mux *chi.Mux) {
	mux.Get("/", h.handle(h.indexPage))
}

func (h *HTMLUIRouter) handle(hf HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hf(w, r)
		if err != nil {
			slog.ErrorContext(r.Context(), r.URL.Path, "error", err)
			h.Renderer.RenderErrorPage(r.Context(), err, w)
		}
	})
}
