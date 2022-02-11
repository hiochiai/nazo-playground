package static

import (
	"net/http"
)

const PathPrefix = "/static"

type Handler struct {
	fileServer http.Handler
}

func NewHandler(staticResourceDir string) *Handler {

	httpDir := http.Dir(staticResourceDir)
	fileServer := http.StripPrefix(PathPrefix+"/", http.FileServer(httpDir))

	return &Handler{
		fileServer: fileServer,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fileServer.ServeHTTP(w, r)
}
