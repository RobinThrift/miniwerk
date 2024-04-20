package htmlui

import (
	"net/http"

	"github.com/RobinThrift/miniwerk/ui"
)

func (h *HTMLUIRouter) indexPage(w http.ResponseWriter, r *http.Request) error {
	return h.Renderer.RenderIndexPage(r.Context(), ui.IndexPageProps{}, w)
}
